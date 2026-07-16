"""Build one merged JSON snapshot per project.

Reads:
  - tdd_flat.json       (SonarQube TDD extract, project measures + smells)
  - metrics_output/<key>-<sha>/ck/{class,method}.csv    (CK output, per class/method)
  - metrics_output/<key>-<sha>/pmd/report.csv           (PMD violation CSV)

Writes:
  - per_project/<project_key>.json

Each output file is a single top-level object with no meta wrapper, mirroring
the shape of one entry in tdd_flat.json's `projects` list plus:

    "classes": [
      {
        "file":  "src/main/java/.../Foo.java",   # repo-relative
        "class": "org.apache.example.Foo",
        "type":  "class",
        "ck":  { ...class-level CK columns... },
        "pmd": { ...class-level PMD rollup... },
        "methods_ck":  [ ...CK method rows... ],
        "methods_pmd": [ ...PMD per-method rollups... ]
      }
    ],
    "pmd_orphans": [
      { "rule": "...", "file": "...", "line": ..., "description": "..." }
    ]

Rationale for keeping ck/pmd side-by-side rather than inner-joining methods:
  CK's method identifier is `foo/1[String]` while PMD's message uses
  `foo(String)`. A cross-tool join by signature is fragile. Line numbers
  work for most methods but not for overloads on the same line or lambdas.
  Both sides live in the same class entry; downstream analysis can pick
  its own join strategy.
"""

from __future__ import annotations

import argparse
import csv
import json
import math
import re
import sys
from collections import defaultdict
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from build_tdd_flat import README_PROJECTS  # noqa: E402

# ---------------------------------------------------------------------------
# CK column plumbing.
# ---------------------------------------------------------------------------

# CK class.csv "identifier" columns hoisted to the outer class entry.
_CK_CLASS_ID_COLS = {"file", "class", "type"}
# CK method.csv identifier columns hoisted / stripped.
_CK_METHOD_STRIP  = {"file", "class"}


def _num(v):
    """Coerce CK CSV string values to int/float when possible.

    Strips thousand-separator commas (e.g., PMD's "WMC=1,022").
    NaN and Infinity are emitted as strings ("NaN", "Infinity",
    "-Infinity") so the output is spec-valid JSON. CK produces NaN for
    e.g. TCC/LCC on 0-or-1-method classes.
    """
    if v is None or v == "":
        return None
    if isinstance(v, str):
        if v in ("NaN", "nan", "NAN"):
            return "NaN"
        stripped = v.replace(",", "") if _COMMA_INT.match(v) else v
    else:
        stripped = v
    try:
        return int(stripped)
    except (TypeError, ValueError):
        pass
    try:
        f = float(stripped)
    except (TypeError, ValueError):
        return v
    if math.isnan(f):
        return "NaN"
    if math.isinf(f):
        return "Infinity" if f > 0 else "-Infinity"
    return f


_COMMA_INT = re.compile(r"^-?\d{1,3}(?:,\d{3})+$")


def _row_to_dict(row: dict, strip: set[str]) -> dict:
    out: dict = {}
    for k, v in row.items():
        if k in strip:
            continue
        out[k] = _num(v)
    return out


# ---------------------------------------------------------------------------
# PMD message parsing.
# ---------------------------------------------------------------------------

# Regex library keyed by (rule_name, is_class_scope). Each pattern must
# capture the metric value in group 'v' (and optionally the class/method
# name in group 'name'). PMD messages are stable across 7.x point releases
# but if they ever change these regexes are the only thing to update.

# PMD emits messages for classes, enums, interfaces, records, and (rarely)
# annotation types. All share the shape `The <kind> 'Name' has ...`. Anonymous
# / synthetic declarations sometimes come through as `The class '' has ...`.
_KIND        = r"(?:class|enum|interface|record|annotation(?: type)?|@interface)"
_TYPE_QUOTED = r"'(?P<name>[^']*)'"
_METH_QUOTED = r"'(?P<name>[^']+)'"

_PMD_PATTERNS: dict[str, list[dict]] = {
    "CyclomaticComplexity": [
        {"scope": "class",
         "re": re.compile(
             r"The " + _KIND + r" " + _TYPE_QUOTED +
             r" has a total cyclomatic complexity of (?P<v>-?\d+)"
             r"(?: \(highest (?P<hi>-?\d+)\))?"
         ),
         "field": "cyclo_total", "extra": {"hi": "cyclo_highest"}},
        {"scope": "method",
         "re": re.compile(
             r"The (?:method|constructor) " + _METH_QUOTED +
             r" has a(?:n\b)? cyclomatic complexity of (?P<v>-?\d+)"
         ),
         "field": "cyclo"},
    ],
    "NcssCount": [
        {"scope": "class",
         "re": re.compile(
             r"The " + _KIND + r" " + _TYPE_QUOTED +
             r" has a NCSS line count of (?P<v>-?\d+)"
             r"(?: \(Highest = (?P<hi>-?\d+)\))?"
         ),
         "field": "ncss", "extra": {"hi": "ncss_highest"}},
        {"scope": "method",
         "re": re.compile(
             r"The (?:method|constructor) " + _METH_QUOTED +
             r" has a NCSS line count of (?P<v>-?\d+)"
         ),
         "field": "ncss"},
    ],
    "CognitiveComplexity": [
        {"scope": "method",
         "re": re.compile(
             r"The (?:method|constructor) " + _METH_QUOTED +
             r" has a cognitive complexity of (?P<v>-?\d+)"
         ),
         "field": "cognitive"},
    ],
    "NPathComplexity": [
        {"scope": "method",
         "re": re.compile(
             r"The (?:method|constructor) " + _METH_QUOTED +
             r" has an NPath complexity of (?P<v>-?\d+)"
         ),
         "field": "npath"},
    ],
    "CouplingBetweenObjects": [
        # No class name in the message; use the file basename.
        {"scope": "class",
         "re": re.compile(
             r"A value of (?P<v>-?\d+) may denote a high amount of coupling"
         ),
         "field": "coupling"},
    ],
    "ExcessiveParameterList": [
        {"scope": "method",
         "re": re.compile(
             r"Avoid long parameter lists \((?P<v>-?\d+) parameters"
         ),
         "field": "params"},
    ],
    "ExcessivePublicCount": [
        # PMD emits "This class has N public methods..." for any type decl.
        {"scope": "class",
         "re": re.compile(
             r"This class has (?P<v>-?\d+) public methods and attributes"
         ),
         "field": "public_count"},
    ],
    "GodClass": [
        # PMD formats very large numbers with thousand-separator commas
        # (e.g. WMC=1,022 seen in hive), so digit runs may contain commas.
        {"scope": "class",
         "re": re.compile(
             r"Possible God Class \(WMC=(?P<wmc>-?[\d,]+),"
             r" ATFD=(?P<atfd>-?[\d,]+), TCC=(?P<tcc>-?[\d.,]+)%\)"
         ),
         "field": "god_class"},
    ],
}


# ---------------------------------------------------------------------------
# Path trimming.
# ---------------------------------------------------------------------------


def _trim_ck_path(abs_path: str, repo_root: str) -> str:
    """Strip the repo_root prefix from a CK/PMD absolute path."""
    if not abs_path:
        return abs_path
    if abs_path.startswith(repo_root):
        stripped = abs_path[len(repo_root):]
        if stripped.startswith("/"):
            stripped = stripped[1:]
        return stripped
    return abs_path


def _trim_sonar_component(component: str) -> str:
    """org.apache:foo:src/main/java/.../X.java -> src/main/java/.../X.java."""
    if not component:
        return component
    parts = component.split(":", 2)
    if len(parts) == 3:
        return parts[2]
    return component


# ---------------------------------------------------------------------------
# Per-project build.
# ---------------------------------------------------------------------------


def _load_ck_classes(ck_dir: Path, repo_root: str) -> tuple[list[dict], dict]:
    """Return (classes_list, index) where index maps
    (repo_rel_file, class_key) -> class_entry.

    A single class is indexed under all of:
      - its fully-qualified name          (`org.foo.Outer`)
      - its simple name after last `.`    (`Outer`)
      - its simple name after last `$`    (`Inner` for `org.foo.Outer$Inner`)

    The `$` variant lets PMD messages (which use the simple name `'Inner'`)
    resolve to CK's inner-class entry (which is named `Outer$Inner`).
    """
    class_csv = ck_dir / "class.csv"
    classes:  list[dict] = []
    index:    dict       = {}
    if not class_csv.exists():
        return classes, index
    for row in csv.DictReader(class_csv.open()):
        rel = _trim_ck_path(row.get("file", ""), repo_root)
        fq  = row.get("class", "")
        entry = {
            "file":        rel,
            "class":       fq,
            "type":        row.get("type"),
            "ck":          _row_to_dict(row, _CK_CLASS_ID_COLS),
            "pmd":         {},
            "methods_ck":  [],
            "methods_pmd": [],
        }
        classes.append(entry)
        index[(rel, fq)] = entry
        dot_simple    = fq.rsplit(".", 1)[-1] if "." in fq else fq
        dollar_simple = dot_simple.rsplit("$", 1)[-1] if "$" in dot_simple else dot_simple
        index.setdefault((rel, dot_simple),    entry)
        index.setdefault((rel, dollar_simple), entry)
    return classes, index


def _attach_ck_methods(ck_dir: Path, index: dict,
                       repo_root: str) -> tuple[int, dict]:
    """Attach CK methods to their class entries.

    Returns (n_added, line_map). line_map is keyed by (file, line_int) and
    maps to the CK class's fully-qualified name. Used later to attribute
    PMD method-scope rows to the correct class when a file has nested
    classes.
    """
    method_csv = ck_dir / "method.csv"
    line_map: dict = {}
    if not method_csv.exists():
        return 0, line_map
    added = 0
    for row in csv.DictReader(method_csv.open()):
        rel = _trim_ck_path(row.get("file", ""), repo_root)
        fq  = row.get("class", "")
        entry = index.get((rel, fq))
        if entry is None:
            simple = fq.rsplit(".", 1)[-1] if "." in fq else fq
            entry = index.get((rel, simple))
        if entry is None:
            continue
        entry["methods_ck"].append(_row_to_dict(row, _CK_METHOD_STRIP))
        added += 1
        try:
            ln = int(row.get("line", "") or 0)
            if ln:
                line_map[(rel, ln)] = fq
        except ValueError:
            pass
    return added, line_map


def _class_simple_from_file(rel: str) -> str:
    return Path(rel).stem


def _synth_class(index: dict, classes: list[dict], rel: str,
                 cls_name: str) -> dict:
    """Create and register a PMD-only class entry when no CK entry exists.

    Some codebases (e.g., felix) have many files declaring the same
    fully-qualified class across separate modules; CK deduplicates but PMD
    scans every file. Synthesised entries have empty `ck` / `methods_ck`
    so downstream consumers can tell a PMD-only entry from a full one.
    """
    entry = {
        "file":        rel,
        "class":       cls_name,
        "type":        "unknown",
        "ck":          {},
        "pmd":         {},
        "methods_ck":  [],
        "methods_pmd": [],
        "synthesised": True,
    }
    classes.append(entry)
    index[(rel, cls_name)] = entry
    simple = cls_name.rsplit(".", 1)[-1] if "." in cls_name else cls_name
    index.setdefault((rel, simple), entry)
    return entry


def _attach_pmd(pmd_dir: Path, index: dict, classes: list[dict],
                repo_root: str, line_map: dict) -> list[dict]:
    """Parse PMD report.csv, attach per-class rollups and per-method entries
    to the CK-derived index, synthesise class entries when PMD sees a class
    CK doesn't, and return truly-unparseable rows as orphans.
    """
    report = pmd_dir / "report.csv"
    orphans: list[dict] = []
    if not report.exists():
        return orphans
    method_bags: dict[tuple, dict] = {}
    for row in csv.DictReader(report.open()):
        rule = row.get("Rule", "")
        desc = row.get("Description", "")
        rel  = _trim_ck_path(row.get("File", ""), repo_root)
        line = row.get("Line")
        try:
            line_i = int(line) if line else None
        except ValueError:
            line_i = None

        pats = _PMD_PATTERNS.get(rule)
        if not pats:
            orphans.append({"rule": rule, "file": rel, "line": line_i,
                            "description": desc, "reason": "unknown_rule"})
            continue

        matched = False
        for spec in pats:
            m = spec["re"].search(desc)
            if not m:
                continue
            matched = True
            gd = m.groupdict()

            if spec["scope"] == "class":
                # Regex-captured name is authoritative; only fall back to
                # the filename when the message name is empty (e.g., PMD's
                # "The class '' has ..." on synthetic/anonymous decls).
                # Never fall back to the file's outer class for a named
                # inner class - that would silently overwrite outer-class
                # metrics with inner-class metrics.
                raw_name = gd.get("name") or ""
                cls_name = raw_name or _class_simple_from_file(rel)
                entry = index.get((rel, cls_name))
                if entry is None and "." in cls_name:
                    entry = index.get((rel, cls_name.rsplit(".", 1)[-1]))
                if entry is None and "$" in cls_name:
                    entry = index.get((rel, cls_name.rsplit("$", 1)[-1]))
                if entry is None:
                    entry = _synth_class(index, classes, rel, cls_name)

                if rule == "GodClass":
                    entry["pmd"]["god_class"]   = True
                    entry["pmd"]["god_wmc"]     = _num(gd["wmc"])
                    entry["pmd"]["god_atfd"]    = _num(gd["atfd"])
                    entry["pmd"]["god_tcc_pct"] = _num(gd["tcc"])
                else:
                    entry["pmd"][spec["field"]] = _num(gd["v"])
                    for src, dst in spec.get("extra", {}).items():
                        if gd.get(src) is not None:
                            entry["pmd"][dst] = _num(gd[src])
                break

            # Method-scope: route to the class that actually contains
            # this method. CK's method.csv gives us (file, line) -> class,
            # which correctly handles nested classes; without it we would
            # always attach to the file's outer class.
            #
            # Some PMD rules emit no method name in their message (e.g.
            # ExcessiveParameterList: "Avoid long parameter lists (N ...)").
            # Dedup rule: prefer merging into an existing record at the
            # same line; use the sig-bearing key only when a name was
            # captured, otherwise a line-only key.
            sig      = gd.get("name") or ""
            # PMD sometimes reports the method's Javadoc/annotation line
            # while CK reports the declaration keyword line, so probe a
            # small forward window (Javadoc rarely exceeds ~6 lines).
            resolved = None
            if line_i is not None:
                for probe in range(line_i, line_i + 8):
                    resolved = line_map.get((rel, probe))
                    if resolved is not None:
                        break
            cls_hint  = resolved or _class_simple_from_file(rel)
            entry     = index.get((rel, cls_hint))
            if entry is None and resolved:
                simple = resolved.rsplit(".", 1)[-1] if "." in resolved else resolved
                simple = simple.rsplit("$", 1)[-1] if "$" in simple else simple
                entry  = index.get((rel, simple))
            if entry is None:
                entry = _synth_class(index, classes, rel, cls_hint)

            if sig:
                key = (rel, cls_hint, line_i, sig)
            else:
                key = None
                for mrec_ in entry["methods_pmd"]:
                    if mrec_.get("line") == line_i:
                        key = ("_existing_", id(mrec_))
                        method_bags[key] = mrec_
                        break
                if key is None:
                    key = (rel, cls_hint, line_i, "")

            mrec = method_bags.get(key)
            if mrec is None:
                mrec = {"signature": sig, "line": line_i}
                method_bags[key] = mrec
                entry["methods_pmd"].append(mrec)
            elif sig and not mrec.get("signature"):
                mrec["signature"] = sig
            mrec[spec["field"]] = _num(gd["v"])
            break

        if not matched:
            orphans.append({"rule": rule, "file": rel, "line": line_i,
                            "description": desc, "reason": "unparsed_message"})
    return orphans


def _build_one(project_key: str, git_slug: str, commit_sha: str,
               tdd_entry: dict, metrics_root: Path, repo_base: Path,
               out_root: Path) -> tuple[Path, dict]:
    org, _, repo = git_slug.partition("/")
    repo_root = str((repo_base / org / repo).resolve())
    short_sha = commit_sha[:10]
    project_metrics = metrics_root / f"{project_key}-{short_sha}"
    ck_dir  = project_metrics / "ck"
    pmd_dir = project_metrics / "pmd"

    classes, index         = _load_ck_classes(ck_dir, repo_root)
    n_methods, line_map    = _attach_ck_methods(ck_dir, index, repo_root)
    orphans                = _attach_pmd(pmd_dir, index, classes,
                                         repo_root, line_map)

    # Base fields from the tdd_flat entry, with component paths trimmed.
    base = dict(tdd_entry)
    smells = base.get("code_smell_issues") or []
    trimmed_smells = []
    for s in smells:
        s2 = dict(s)
        if s2.get("component"):
            s2["component"] = _trim_sonar_component(s2["component"])
        trimmed_smells.append(s2)
    base["code_smell_issues"] = trimmed_smells

    base["classes"]     = classes
    base["pmd_orphans"] = orphans

    out_path = out_root / f"{project_key}.json"
    out_path.write_text(json.dumps(
        base, indent=2, ensure_ascii=False, allow_nan=False,
    ))

    stats = {
        "classes":     len(classes),
        "methods_ck":  n_methods,
        "orphans":     len(orphans),
        "size_bytes":  out_path.stat().st_size,
    }
    return out_path, stats


def main() -> int:
    here = Path(__file__).resolve().parent
    ap = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    ap.add_argument("--tdd-flat",    default=str(here / "tdd_flat.json"))
    ap.add_argument("--metrics",     default=str(here / "metrics_output"))
    ap.add_argument("--repo-base",   default=str(Path.home() / "go" / "src" / "github.com"))
    ap.add_argument("--out",         default=str(here / "per_project"))
    ap.add_argument("--targets", nargs="+",
                    help="only build these project_keys (default: all 31)")
    args = ap.parse_args()

    tdd = json.loads(Path(args.tdd_flat).read_text())
    tdd_by_key = {p["project_key"]: p for p in tdd["projects"]}

    metrics_root = Path(args.metrics).resolve()
    repo_base    = Path(args.repo_base).resolve()
    out_root     = Path(args.out).resolve()
    out_root.mkdir(parents=True, exist_ok=True)

    keys = args.targets if args.targets else [p[0] for p in README_PROJECTS]
    by_key_readme = {p[0]: p for p in README_PROJECTS}

    grand = defaultdict(int)
    for project_key in keys:
        if project_key not in by_key_readme:
            print(f"[skip] unknown project_key: {project_key}", file=sys.stderr)
            continue
        _, git_slug, commit_sha = by_key_readme[project_key]
        tdd_entry = tdd_by_key.get(project_key, {
            "project_key":  project_key,
            "git_link":     git_slug,
            "commit_sha":   commit_sha,
        })
        path, stats = _build_one(
            project_key, git_slug, commit_sha,
            tdd_entry, metrics_root, repo_base, out_root,
        )
        print(f"  {project_key:<24} classes={stats['classes']:>5}  "
              f"methods_ck={stats['methods_ck']:>6}  "
              f"orphans={stats['orphans']:>5}  "
              f"{stats['size_bytes']//1024:>6} KB  -> {path.name}",
              file=sys.stderr)
        for k, v in stats.items():
            grand[k] += v

    print(f"totals: {dict(grand)}", file=sys.stderr)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
