"""Build a flat JSON snapshot of the Technical Debt Dataset.

Extracts, for each of the 31 Apache projects listed in README.md, the
SonarQube analysis pinned to a specific commit (README's commit map),
plus:
  - Project-level SONAR_MEASURES (size, complexity, coupling, duplication,
    quality, violations by severity, SQALE technical debt).
  - Per-rule counts across all `code_smells:*` rules.
  - Full per-issue records for `code_smells:*` issues that were OPEN at the
    time of that analysis (created on or before, and either never closed or
    closed after).

Output: tdd_flat.json in the same folder.

Usage:
    python3 build_tdd_flat.py [--db path/to/td_V2.db] [--out path/to/tdd_flat.json]

Only depends on the Python stdlib (sqlite3, json, argparse, datetime).
"""

from __future__ import annotations

import argparse
import json
import sqlite3
import sys
from collections import defaultdict
from datetime import datetime, timezone
from pathlib import Path

# ---------------------------------------------------------------------------
# README-pinned commit map. project_key -> (git_link, commit_sha).
# Order preserved from README.md so the JSON output mirrors README ordering.
# ---------------------------------------------------------------------------
README_PROJECTS: list[tuple[str, str, str]] = [
    ("santuario",              "apache/santuario-java",         "be4e2331f77adb1e479406ebf973e516bbf5e32b"),
    ("commons-beanutils",      "apache/commons-beanutils",      "c4da598872233b59af41a221bd2bdcefbbca1259"),
    ("commons-validator",      "apache/commons-validator",      "a3771313c9f1833abf32c7c294ad1de4810e532d"),
    ("commons-net",            "apache/commons-net",            "fb7aae4c64f7d2bf6dced00c49c3ffc428b2d572"),
    ("commons-configuration",  "apache/commons-configuration",  "15b4031ba94a60f20b854e6ce2c7964d77086387"),
    ("commons-vfs",            "apache/commons-vfs",            "d72192f18bfaed730b4f37a2f94853e1503ffd74"),
    ("commons-daemon",         "apache/commons-daemon",         "1ffa799cb3ddf5a4a918e59e46cd9868ee766b19"),
    ("commons-bcel",           "apache/commons-bcel",           "6ed18c5bef0f5b93b54783a8e8fb2b9042da26ac"),
    ("commons-codec",          "apache/commons-codec",          "db51a1cb41e9155ca028a73b0637b32a2c37c43a"),
    ("commons-ognl",           "apache/commons-ognl",           "6ec1a1a4588b82c0972ca2ff35b85d9b50cc4604"),
    ("commons-jxpath",         "apache/commons-jxpath",         "eff47ab8ca52fdbc91d1313cc224324465dd043e"),
    ("commons-exec",           "apache/commons-exec",           "2da60ab3eefaaa2f8a434ded1eebe1ce17efd34a"),
    ("commons-jexl",           "apache/commons-jexl",           "d3e702149a3db297d6db2c0b7671807f5c7b98fc"),
    ("commons-dbcp",           "apache/commons-dbcp",           "d8dd39b32bbb04a28ea86eb826c56aa6783f3faf"),
    ("commons-io",             "apache/commons-io",             "65c4a9c0ec651dd99f28b9fae40378728d071985"),
    ("commons-fileupload",     "apache/commons-fileupload",     "cae90facebc54803232a0593003914ca77193a73"),
    ("commons-jelly",          "apache/commons-jelly",          "48c008cc2328402e0976295625b32c5197ba2324"),
    ("commons-digester",       "apache/commons-digester",       "c1d0e563339faec040eb036ae97a7b7bf07ba865"),
    ("commons-collections",    "apache/commons-collections",    "f0f364fd9d946483f947011a3557c1e6f2e5d8ee"),
    ("commons-cli",            "apache/commons-cli",            "92f1def0bb3c0345295012e36b7150cfd1d7b6ab"),
    ("commons-dbutils",        "apache/commons-dbutils",        "2f48485a82697d9aed060ba36f6d5beb3a58ed8b"),
    ("httpcomponents-client",  "apache/httpcomponents-client",  "8a1b96bfa75382c0b94d70f6914fbb9bfeb0451e"),
    ("httpcomponents-core",    "apache/httpcomponents-core",    "3a677d47cb872b6ede20b28e93d3206f08b349ac"),
    ("zookeeper",              "apache/zookeeper",              "eac693cc76a34f96b9116ef33d1e92af7129416d"),
    ("hive",                   "apache/hive",                   "a4d91eaf2925239aa29342f7e5b0f8680c842390"),
    ("thrift",                 "apache/thrift",                 "a2123693838410c1e78170419e9bb91cb01151b4"),
    ("archiva",                "apache/archiva",                "374fc983abc92df8aa4f8ef30caee94b34312ad2"),
    ("felix",                  "apache/felix",                  "bdb6cb5cac0d81e9cd3fda666065e0e577eb9c41"),
    ("cayenne",                "apache/cayenne",                "b9988a83e364b9b470873dff8996dcf401d08dc4"),
    ("cocoon",                 "apache/cocoon",                 "a80f73b27592a2794c9133ee03d2e402bf12ecc1"),
    ("batik",                  "apache/batik",                  "2bb3a6ea5a6258ff6372e2493b81d7768d6bb494"),
]

# All `code_smells:*` rules we emit counts for. Kept static so the schema is
# stable across projects (missing rules -> count 0). `blob_class` is called
# out in README.md but has zero rows in this DB; still emitted for schema
# stability.
CODE_SMELL_RULES: list[str] = [
    "antisingleton",
    "baseclass_abstract",
    "blob_class",
    "class_data_private",
    "complex_class",
    "large_class",
    "lazy_class",
    "long_method",
    "long_parameter_list",
    "many_field_attributes_not_complex",
    "refused_parent_bequest",
    "spaghetti_code",
    "speculative_generality",
    "swiss_army_knife",
]

# SONAR_MEASURES columns to include, grouped for readability. Each maps to a
# JSON key. Numeric where possible; strings kept as-is otherwise.
MEASURE_COLUMNS: list[tuple[str, str]] = [
    ("lines",                          "lines"),
    ("ncloc",                          "ncloc"),
    ("classes",                        "classes"),
    ("files",                          "files"),
    ("directories",                    "directories"),
    ("functions",                      "functions"),
    ("statements",                     "statements"),
    ("comment_lines",                  "COMMENT_LINES"),
    ("comment_lines_density",          "COMMENT_LINES_DENSITY"),
    ("complexity",                     "COMPLEXITY"),
    ("class_complexity",               "CLASS_COMPLEXITY"),
    ("function_complexity",            "FUNCTION_COMPLEXITY"),
    ("file_complexity",                "FILE_COMPLEXITY"),
    ("cognitive_complexity",           "COGNITIVE_COMPLEXITY"),
    ("afferent_couplings",             "AFFERENT_COUPLINGS"),
    ("efferent_couplings",             "EFFERENT_COUPLINGS"),
    ("package_dependency_cycles",      "PACKAGE_DEPENDENCY_CYCLES"),
    ("number_of_classes_and_interfaces","NUMBER_OF_CLASSES_AND_INTERFACES"),
    ("duplicated_lines",               "DUPLICATED_LINES"),
    ("duplicated_blocks",              "DUPLICATED_BLOCKS"),
    ("duplicated_files",               "DUPLICATED_FILES"),
    ("duplicated_lines_density",       "DUPLICATED_LINES_DENSITY"),
    ("code_smells",                    "CODE_SMELLS"),
    ("bugs",                           "BUGS"),
    ("vulnerabilities",                "VULNERABILITIES"),
    ("coverage",                       "COVERAGE"),
    ("violations",                     "VIOLATIONS"),
    ("blocker_violations",             "BLOCKER_VIOLATIONS"),
    ("critical_violations",            "CRITICAL_VIOLATIONS"),
    ("major_violations",               "MAJOR_VIOLATIONS"),
    ("minor_violations",               "MINOR_VIOLATIONS"),
    ("info_violations",                "INFO_VIOLATIONS"),
    ("sqale_index",                    "SQALE_INDEX"),
    ("sqale_rating",                   "SQALE_RATING"),
    ("sqale_debt_ratio",               "SQALE_DEBT_RATIO"),
    ("development_cost",               "DEVELOPMENT_COST"),
    ("reliability_rating",             "RELIABILITY_RATING"),
    ("security_rating",                "SECURITY_RATING"),
    ("reliability_remediation_effort", "RELIABILITY_REMEDIATION_EFFORT"),
    ("security_remediation_effort",    "SECURITY_REMEDIATION_EFFORT"),
]

# Per-issue SONAR_ISSUES columns to emit. Kept in a stable order.
ISSUE_COLUMNS: list[tuple[str, str]] = [
    ("issue_key",     "ISSUE_KEY"),
    ("rule",          "RULE"),
    ("type",          "TYPE"),
    ("severity",      "SEVERITY"),
    ("status",        "STATUS"),
    ("resolution",    "RESOLUTION"),
    ("component",     "COMPONENT"),
    ("start_line",    "START_LINE"),
    ("end_line",      "END_LINE"),
    ("start_offset",  "START_OFFSET"),
    ("end_offset",    "END_OFFSET"),
    ("effort",        "EFFORT"),
    ("debt",          "DEBT"),
    ("tags",          "TAGS"),
    ("creation_date", "CREATION_DATE"),
    ("close_date",    "CLOSE_DATE"),
    ("message",       "MESSAGE"),
]


def _blank_to_none(v):
    if v is None:
        return None
    if isinstance(v, str) and v.strip() == "":
        return None
    return v


def _to_number(v):
    """Coerce a DB value to int/float when it clearly is one, else return as-is."""
    v = _blank_to_none(v)
    if v is None or not isinstance(v, str):
        return v
    try:
        if "." in v or "e" in v or "E" in v:
            return float(v)
        return int(v)
    except ValueError:
        return v


def load_project_index(conn: sqlite3.Connection) -> dict[str, dict]:
    """project_key (README name) -> {project_id, sonar_project_key, git_link_db, jira_link}."""
    out: dict[str, dict] = {}
    for row in conn.execute(
        "SELECT PROJECT_KEY, PROJECT_ID, SONAR_PROJECT_KEY, GIT_LINK, JIRA_LINK FROM PROJECTS"
    ):
        out[row[0]] = {
            "project_id":         row[1],
            "sonar_project_key":  row[2],
            "git_link_db":        row[3],
            "jira_link":          row[4],
        }
    return out


def pick_analysis(conn: sqlite3.Connection, project_id: str, commit_sha: str):
    """Return (analysis_key, date, revision, matched_by) for the analysis pinned
    to `commit_sha`. Falls back to MAX(DATE) if no exact SHA match.
    """
    cur = conn.execute(
        "SELECT ANALYSIS_KEY, DATE, REVISION FROM SONAR_ANALYSIS "
        "WHERE PROJECT_ID = ? AND REVISION = ?",
        (project_id, commit_sha),
    )
    rows = cur.fetchall()
    if rows:
        # If multiple analyses share the same REVISION, take the newest.
        rows.sort(key=lambda r: r[1] or "", reverse=True)
        ak, dt, rev = rows[0]
        return ak, dt, rev, "commit_sha"

    # Fallback: newest analysis for the project.
    cur = conn.execute(
        "SELECT ANALYSIS_KEY, DATE, REVISION FROM SONAR_ANALYSIS "
        "WHERE PROJECT_ID = ? ORDER BY DATE DESC LIMIT 1",
        (project_id,),
    )
    row = cur.fetchone()
    if row is None:
        return None, None, None, "no_analysis"
    return row[0], row[1], row[2], "max_date_fallback"


def fetch_measures(conn: sqlite3.Connection, project_id: str, analysis_key: str) -> dict:
    col_sql = ", ".join(f'"{c}"' for _, c in MEASURE_COLUMNS)
    cur = conn.execute(
        f"SELECT {col_sql} FROM SONAR_MEASURES "
        f"WHERE PROJECT_ID = ? AND ANALYSIS_KEY = ?",
        (project_id, analysis_key),
    )
    row = cur.fetchone()
    if row is None:
        return {out: None for out, _ in MEASURE_COLUMNS}
    result: dict = {}
    for (out_key, _), val in zip(MEASURE_COLUMNS, row):
        # Values may be stored as strings in some cols; coerce numerics.
        result[out_key] = _to_number(val)
    return result


def fetch_code_smell_issues(
    conn: sqlite3.Connection, project_id: str, analysis_date: str
) -> tuple[list[dict], dict[str, int]]:
    """Return (issues_list, counts_by_short_rule_name).

    Only `code_smells:*` issues alive at `analysis_date` are returned:
      CREATION_DATE <= analysis_date
      AND (CLOSE_DATE IS NULL OR CLOSE_DATE = '' OR CLOSE_DATE > analysis_date)
    """
    col_sql = ", ".join(f'"{c}"' for _, c in ISSUE_COLUMNS)
    cur = conn.execute(
        f"SELECT {col_sql} FROM SONAR_ISSUES "
        f"WHERE PROJECT_ID = ? "
        f"  AND RULE LIKE 'code_smells:%' "
        f"  AND CREATION_DATE IS NOT NULL "
        f"  AND CREATION_DATE <= ? "
        f"  AND (CLOSE_DATE IS NULL OR CLOSE_DATE = '' OR CLOSE_DATE > ?)",
        (project_id, analysis_date, analysis_date),
    )
    issues: list[dict] = []
    counts: dict[str, int] = defaultdict(int)
    for row in cur:
        rec: dict = {}
        for (out_key, _), val in zip(ISSUE_COLUMNS, row):
            if out_key in ("start_line", "end_line", "start_offset", "end_offset"):
                v = _blank_to_none(val)
                rec[out_key] = int(v) if isinstance(v, (int, float)) else _to_number(v)
            else:
                rec[out_key] = _blank_to_none(val)
        issues.append(rec)
        rule = rec.get("rule") or ""
        short = rule.split(":", 1)[1] if ":" in rule else rule
        counts[short] += 1

    # Ensure every configured rule appears in counts (with 0 if absent).
    counts_full = {r: int(counts.get(r, 0)) for r in CODE_SMELL_RULES}
    # Include any unexpected rules we hit, too (shouldn't happen given the LIKE
    # filter above, but keeps the output honest if the DB grows new ones).
    for r, n in counts.items():
        if r not in counts_full:
            counts_full[r] = int(n)
    return issues, counts_full


def build(db_path: Path, out_path: Path) -> None:
    conn = sqlite3.connect(str(db_path))
    conn.row_factory = None
    conn.text_factory = str

    projects_index = load_project_index(conn)

    projects_out: list[dict] = []
    warnings: list[str] = []

    for project_key, git_link, commit_sha in README_PROJECTS:
        meta = projects_index.get(project_key)
        if meta is None:
            warnings.append(f"{project_key}: not found in PROJECTS table")
            projects_out.append({
                "project_key":  project_key,
                "git_link":     git_link,
                "commit_sha":   commit_sha,
                "error":        "PROJECT not found in DB",
            })
            continue

        project_id = meta["project_id"]
        analysis_key, analysis_date, revision, matched_by = pick_analysis(
            conn, project_id, commit_sha
        )
        if analysis_key is None:
            warnings.append(f"{project_key}: no SONAR_ANALYSIS rows for {project_id}")
            projects_out.append({
                "project_key":  project_key,
                "project_id":   project_id,
                "git_link":     git_link,
                "commit_sha":   commit_sha,
                "error":        "no analyses in DB",
            })
            continue

        if matched_by == "max_date_fallback":
            warnings.append(
                f"{project_key}: commit_sha {commit_sha!r} not found in "
                f"SONAR_ANALYSIS; fell back to newest analysis "
                f"(rev={revision!r}, date={analysis_date!r})"
            )

        measures = fetch_measures(conn, project_id, analysis_key)
        issues, counts = fetch_code_smell_issues(conn, project_id, analysis_date)

        entry: dict = {
            "project_key":        project_key,
            "project_id":         project_id,
            "sonar_project_key":  meta["sonar_project_key"],
            "git_link":           git_link,
            "git_link_db":        meta["git_link_db"],
            "jira_link":          meta["jira_link"],
            "commit_sha":         commit_sha,
            "analysis_key":       analysis_key,
            "analysis_date":      analysis_date,
            "analysis_revision":  revision,
            "analysis_matched_by": matched_by,
            "measures":           measures,
            "code_smell_counts":  counts,
            "code_smell_issues":  issues,
        }
        projects_out.append(entry)

        print(
            f"  {project_key:<24} issues={len(issues):>5}  "
            f"matched_by={matched_by}",
            file=sys.stderr,
        )

    payload = {
        "_meta": {
            "generated_at":       datetime.now(timezone.utc).isoformat(),
            "source_db":          db_path.name,
            "project_count":      len(projects_out),
            "issue_filter":       (
                "RULE LIKE 'code_smells:%' "
                "AND CREATION_DATE <= analysis_date "
                "AND (CLOSE_DATE IS NULL OR CLOSE_DATE > analysis_date)"
            ),
            "analysis_selection": (
                "SONAR_ANALYSIS row where REVISION == README commit_sha; "
                "falls back to MAX(DATE) if that revision is not present."
            ),
            "code_smell_rules":   CODE_SMELL_RULES,
            "measure_columns":    [k for k, _ in MEASURE_COLUMNS],
            "issue_columns":      [k for k, _ in ISSUE_COLUMNS],
            "warnings":           warnings,
        },
        "projects": projects_out,
    }

    out_path.write_text(json.dumps(payload, indent=2, ensure_ascii=False))
    print(f"wrote {out_path} ({out_path.stat().st_size:,} bytes)", file=sys.stderr)
    if warnings:
        print(f"{len(warnings)} warning(s):", file=sys.stderr)
        for w in warnings:
            print(f"  - {w}", file=sys.stderr)


def main() -> int:
    here = Path(__file__).resolve().parent
    ap = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    ap.add_argument("--db",  default=str(here / "td_V2.db"), help="path to td_V2.db")
    ap.add_argument("--out", default=str(here / "tdd_flat.json"), help="output JSON path")
    args = ap.parse_args()

    db_path  = Path(args.db).resolve()
    out_path = Path(args.out).resolve()

    if not db_path.exists():
        print(f"error: db not found: {db_path}", file=sys.stderr)
        return 2

    build(db_path, out_path)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
