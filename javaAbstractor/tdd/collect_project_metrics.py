"""Collect class-level metrics for one Apache project at a specific commit.

Purpose: cross-validation dataset for the local Java abstractor's own AST
computed metrics (WMC, TCC, ATFD, LOC, NOM, and related). This script drives
two independent tools per project so that later analysis can compare their
outputs against the abstractor's:

  1. CK  (https://github.com/mauricioaniche/ck) - primary metric source.
     Produces class.csv / method.csv / field.csv with raw metric columns.

  2. PMD 7 (https://pmd.github.io/) - secondary/cross-check source. Uses the
     checked-in ruleset `pmd_metrics.xml` in this folder, which references
     PMD's built-in metric-based design rules with thresholds pushed to
     their minimum so they fire on every class/method. Their standard
     violation messages embed the metric value, so parsing the CSV report
     yields per-class/method values.

Workflow per invocation:

  1. Resolve target repo dir as <repo-base>/<org>/<repo>.
  2. If missing -> `git clone`. If present -> assume it's a valid clone
     (user's stated convention: rename/delete manually to force re-clone).
  3. Fetch if needed, then `git checkout <sha>`. Skipped if HEAD already
     matches. Never runs anything destructive.
  4. Run CK if `--ck-jar` exists; skip with warning otherwise.
  5. Run PMD if `--pmd-bin` exists; skip with warning otherwise.
  6. Write run.json alongside outputs with full provenance.

This script never deletes or force-modifies any file outside `--out`
(and the created-if-missing clone). Interrupts, tool failures, and
missing binaries are handled non-destructively.

Usage:
    ./collect_project_metrics.py commons-cli \\
        apache/commons-cli \\
        92f1def0bb3c0345295012e36b7150cfd1d7b6ab

Note: This script was created by Opus 4.7 High
"""

from __future__ import annotations

import argparse
import json
import os
import re
import shutil
import subprocess
import sys
import time
import zipfile
from datetime import datetime, timezone
from pathlib import Path

# ---------------------------------------------------------------------------
# Small helpers
# ---------------------------------------------------------------------------


def _log(msg: str) -> None:
    print(f"[collect] {msg}", file=sys.stderr, flush=True)


def _run(cmd: list[str], cwd: Path | None = None, check: bool = True,
         capture: bool = True, env: dict | None = None,
         timeout: float | None = None) -> subprocess.CompletedProcess:
    """Thin wrapper around subprocess.run with consistent logging."""
    _log("$ " + " ".join(str(c) for c in cmd) + (f"    (cwd={cwd})" if cwd else ""))
    return subprocess.run(
        cmd,
        cwd=str(cwd) if cwd else None,
        check=check,
        text=True,
        capture_output=capture,
        env=env,
        timeout=timeout,
    )


def _short_sha(sha: str) -> str:
    return sha[:10] if len(sha) >= 10 else sha


# ---------------------------------------------------------------------------
# Git checkout stage
# ---------------------------------------------------------------------------


def ensure_checkout(repo_base: Path, git_slug: str, commit_sha: str) -> Path:
    """Ensure `<repo_base>/<org>/<repo>` exists and points at `commit_sha`.

    Returns the absolute path to the checkout.
    """
    org, _, name = git_slug.partition("/")
    if not org or not name:
        raise ValueError(f"expected git slug of form 'org/repo', got {git_slug!r}")
    repo_dir = repo_base / org / name

    if not repo_dir.exists():
        _log(f"cloning {git_slug} -> {repo_dir}")
        (repo_base / org).mkdir(parents=True, exist_ok=True)
        _run([
            "git", "clone",
            "--no-single-branch",
            f"https://github.com/{git_slug}.git",
            str(repo_dir),
        ])
    else:
        _log(f"reusing existing clone at {repo_dir}")
        if not (repo_dir / ".git").exists():
            raise RuntimeError(
                f"{repo_dir} exists but is not a git checkout. "
                "Rename or remove it manually to force a re-clone."
            )

    # Skip checkout if we are already on the target commit.
    head = _run(["git", "rev-parse", "HEAD"], cwd=repo_dir).stdout.strip()
    if head == commit_sha:
        _log(f"HEAD already at {commit_sha}, skipping checkout")
        return repo_dir

    # Try direct checkout first.
    try:
        _run(["git", "-c", "advice.detachedHead=false",
              "checkout", commit_sha], cwd=repo_dir)
    except subprocess.CalledProcessError:
        _log(f"checkout failed, fetching all refs and retrying")
        _run(["git", "fetch", "--all", "--tags", "--quiet"], cwd=repo_dir)
        _run(["git", "-c", "advice.detachedHead=false",
              "checkout", commit_sha], cwd=repo_dir)

    head = _run(["git", "rev-parse", "HEAD"], cwd=repo_dir).stdout.strip()
    if head != commit_sha:
        raise RuntimeError(
            f"checkout landed at HEAD={head}, expected {commit_sha}"
        )
    return repo_dir


# ---------------------------------------------------------------------------
# CK stage
# ---------------------------------------------------------------------------


def _detect_ck_version(ck_jar: Path) -> str | None:
    try:
        with zipfile.ZipFile(ck_jar) as z:
            with z.open("META-INF/MANIFEST.MF") as f:
                for line in f.read().decode("utf-8", errors="ignore").splitlines():
                    if line.lower().startswith("implementation-version"):
                        return line.split(":", 1)[1].strip()
    except Exception:
        pass
    return None


def run_ck(ck_jar: Path, repo_dir: Path, out_dir: Path,
           use_jars: bool = True) -> dict:
    """Run CK against `repo_dir`. Output CSVs land in `out_dir`.

    CK CLI args:
        java -jar ck.jar <project_dir> <use_jars> <max_files> <fields_var> <out>
    where:
        use_jars      - resolve dependency JARs on classpath ("true"/"false")
        max_files     - 0 = auto partitioning
        fields_var    - emit variable/field-level metrics too ("true"/"false")
        out           - output directory prefix (CK appends "class.csv" etc.)
    """
    out_dir.mkdir(parents=True, exist_ok=True)
    stdout_path = out_dir / "ck.stdout.log"
    stderr_path = out_dir / "ck.stderr.log"
    started = time.time()
    cmd = [
        "java", "-jar", str(ck_jar),
        str(repo_dir),
        "true" if use_jars else "false",
        "0",
        "false",              # skip field/variable detail (huge, unused here)
        str(out_dir) + os.sep,
    ]
    _log("$ " + " ".join(cmd))
    with open(stdout_path, "w") as so, open(stderr_path, "w") as se:
        proc = subprocess.run(cmd, stdout=so, stderr=se, text=True)
    elapsed = time.time() - started

    result: dict = {
        "cmd":         cmd,
        "exit_code":   proc.returncode,
        "elapsed_sec": round(elapsed, 2),
        "stdout_log":  str(stdout_path.name),
        "stderr_log":  str(stderr_path.name),
        "outputs":     {},
        "version":     _detect_ck_version(ck_jar),
    }
    for expected in ("class.csv", "method.csv", "field.csv", "variable.csv"):
        p = out_dir / expected
        if p.exists():
            result["outputs"][expected] = {
                "size_bytes": p.stat().st_size,
                "row_count":  max(0, sum(1 for _ in p.open()) - 1),
            }
    return result


# ---------------------------------------------------------------------------
# PMD stage
# ---------------------------------------------------------------------------


_JAVA_SOURCE_RE = re.compile(
    r"<(?:maven\.compiler\.source|maven\.compiler\.release|source)>"
    r"\s*(?:\$\{[^}]+\}|(\d+(?:\.\d+)?))\s*"
    r"</(?:maven\.compiler\.source|maven\.compiler\.release|source)>"
)


def _detect_pmd_java_version(repo_dir: Path, fallback: str) -> tuple[str, str]:
    """Best-effort Java source version inference from a top-level pom.xml.

    Returns (version_string_for_pmd, reason). E.g. ("17", "fallback") or
    ("8",  "pom.xml maven.compiler.source").
    """
    pom = repo_dir / "pom.xml"
    if pom.exists():
        try:
            text = pom.read_text(encoding="utf-8", errors="ignore")
            m = _JAVA_SOURCE_RE.search(text)
            if m and m.group(1):
                raw = m.group(1)
                # "1.8" -> "8", "17" -> "17"
                if raw.startswith("1.") and raw[2:].isdigit():
                    raw = raw[2:]
                return raw, "pom.xml maven.compiler.source"
        except Exception:
            pass
    return fallback, "fallback default"


def _pmd_version(pmd_bin: Path) -> str | None:
    try:
        p = subprocess.run(
            [str(pmd_bin), "--version"],
            check=False, text=True, capture_output=True, timeout=10,
        )
        out = (p.stdout or "") + (p.stderr or "")
        m = re.search(r"PMD\s+(\d+\.\d+\.\d+)", out)
        return m.group(1) if m else None
    except Exception:
        return None


def run_pmd(pmd_bin: Path, ruleset: Path, repo_dir: Path,
            out_dir: Path, java_version_default: str) -> dict:
    out_dir.mkdir(parents=True, exist_ok=True)
    report_path = out_dir / "report.csv"
    stdout_path = out_dir / "pmd.stdout.log"
    stderr_path = out_dir / "pmd.stderr.log"

    # Locate the source root(s) to analyse. Prefer src/ if present, else
    # feed PMD the whole repo (with binaries excluded via --exclude).
    src_root = repo_dir / "src"
    target = src_root if src_root.exists() else repo_dir

    java_ver, java_ver_reason = _detect_pmd_java_version(
        repo_dir, fallback=java_version_default
    )

    started = time.time()
    cmd = [
        str(pmd_bin), "check",
        "-d",    str(target),
        "-R",    str(ruleset),
        "-f",    "csv",
        "-r",    str(report_path),
        "--use-version", f"java-{java_ver}",
        "--no-cache",
        "--no-progress",
        "--no-fail-on-violation",
        "--no-fail-on-error",
    ]
    _log("$ " + " ".join(cmd))
    with open(stdout_path, "w") as so, open(stderr_path, "w") as se:
        proc = subprocess.run(cmd, stdout=so, stderr=se, text=True)
    elapsed = time.time() - started

    result: dict = {
        "cmd":          cmd,
        "exit_code":    proc.returncode,
        "elapsed_sec":  round(elapsed, 2),
        "stdout_log":   str(stdout_path.name),
        "stderr_log":   str(stderr_path.name),
        "target_dir":   str(target),
        "java_version": java_ver,
        "java_version_reason": java_ver_reason,
        "version":      _pmd_version(pmd_bin),
    }
    if report_path.exists():
        result["report"] = {
            "path":       report_path.name,
            "size_bytes": report_path.stat().st_size,
            "row_count":  max(0, sum(1 for _ in report_path.open()) - 1),
        }
    return result


# ---------------------------------------------------------------------------
# Entry point
# ---------------------------------------------------------------------------


def main() -> int:
    here = Path(__file__).resolve().parent
    default_repo_base = Path.home() / "go" / "src" / "github.com"
    default_ck_jar    = Path.home() / "tools" / "ck.jar"
    default_pmd_bin   = Path.home() / "pmd-bin-7.26.0" / "bin" / "pmd"
    default_ruleset   = here / "pmd_metrics.xml"
    default_out_root  = here / "metrics_output"

    ap = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    ap.add_argument("project_key",
                    help="short project name used in output dir name")
    ap.add_argument("git_slug",
                    help="github <org>/<repo>, e.g. apache/commons-cli")
    ap.add_argument("commit_sha",
                    help="full commit SHA to check out")
    ap.add_argument("--repo-base",    default=str(default_repo_base),
                    help=f"clone base dir (default: {default_repo_base})")
    ap.add_argument("--ck-jar",       default=str(default_ck_jar),
                    help=f"path to ck.jar (default: {default_ck_jar})")
    ap.add_argument("--pmd-bin",      default=str(default_pmd_bin),
                    help=f"path to pmd binary (default: {default_pmd_bin})")
    ap.add_argument("--pmd-ruleset",  default=str(default_ruleset),
                    help=f"pmd ruleset xml (default: {default_ruleset})")
    ap.add_argument("--out",          default=str(default_out_root),
                    help=f"output root dir (default: {default_out_root})")
    ap.add_argument("--tools", default="ck,pmd",
                    help="comma-separated subset of {ck,pmd} to run "
                         "(default: ck,pmd)")
    ap.add_argument("--pmd-java-version", default="17",
                    help="fallback PMD Java source version (default: 17)")
    ap.add_argument("--ck-no-jars", action="store_true",
                    help="disable CK's use_jars (skip dep resolution, "
                         "faster but reduces coupling accuracy)")
    args = ap.parse_args()

    tools = {t.strip() for t in args.tools.split(",") if t.strip()}
    unknown = tools - {"ck", "pmd"}
    if unknown:
        _log(f"unknown tools requested: {sorted(unknown)}")
        return 2

    repo_base   = Path(args.repo_base).expanduser().resolve()
    ck_jar      = Path(args.ck_jar).expanduser().resolve()
    pmd_bin     = Path(args.pmd_bin).expanduser().resolve()
    pmd_ruleset = Path(args.pmd_ruleset).expanduser().resolve()
    out_root    = Path(args.out).expanduser().resolve()

    if not shutil.which("git"):
        _log("error: git not found on PATH")
        return 2
    if not shutil.which("java"):
        _log("error: java not found on PATH (needed for CK; PMD launcher too)")
        return 2

    out_dir = out_root / f"{args.project_key}-{_short_sha(args.commit_sha)}"
    out_dir.mkdir(parents=True, exist_ok=True)

    started_at = datetime.now(timezone.utc).isoformat()
    _log(f"output dir: {out_dir}")

    # -- git ---------------------------------------------------------------
    repo_base.mkdir(parents=True, exist_ok=True)
    repo_dir = ensure_checkout(repo_base, args.git_slug, args.commit_sha)
    resolved_head = _run(["git", "rev-parse", "HEAD"], cwd=repo_dir).stdout.strip()

    provenance: dict = {
        "generated_at":     started_at,
        "project_key":      args.project_key,
        "git_slug":         args.git_slug,
        "target_commit":    args.commit_sha,
        "resolved_head":    resolved_head,
        "repo_dir":         str(repo_dir),
        "out_dir":          str(out_dir),
        "tools_requested":  sorted(tools),
        "steps":            {},
        "errors":           [],
    }

    # -- CK ----------------------------------------------------------------
    if "ck" in tools:
        ck_step: dict = {"skipped": False}
        if not ck_jar.exists():
            msg = (
                f"CK jar not found at {ck_jar}. Skipping CK. "
                f"Install with e.g.:\n"
                f"  mkdir -p {ck_jar.parent} && \\\n"
                f"  curl -L -o {ck_jar} "
                f"'https://github.com/mauricioaniche/ck/releases/download/"
                f"0.7.1-SNAPSHOT/ck-0.7.1-SNAPSHOT-jar-with-dependencies.jar'"
            )
            _log(msg)
            ck_step["skipped"] = True
            ck_step["reason"]  = msg.split("\n", 1)[0]
            provenance["errors"].append(msg.split("\n", 1)[0])
        else:
            ck_out = out_dir / "ck"
            ck_step.update(run_ck(
                ck_jar, repo_dir, ck_out,
                use_jars=not args.ck_no_jars,
            ))
            if ck_step.get("exit_code", 0) != 0:
                provenance["errors"].append(
                    f"CK exited with {ck_step['exit_code']}"
                )
        provenance["steps"]["ck"] = ck_step

    # -- PMD ---------------------------------------------------------------
    if "pmd" in tools:
        pmd_step: dict = {"skipped": False}
        if not pmd_bin.exists():
            msg = (f"PMD binary not found at {pmd_bin}. Skipping PMD. "
                   f"Pass --pmd-bin <path> to override.")
            _log(msg)
            pmd_step["skipped"] = True
            pmd_step["reason"]  = msg
            provenance["errors"].append(msg)
        elif not pmd_ruleset.exists():
            msg = f"PMD ruleset not found at {pmd_ruleset}. Skipping PMD."
            _log(msg)
            pmd_step["skipped"] = True
            pmd_step["reason"]  = msg
            provenance["errors"].append(msg)
        else:
            pmd_out = out_dir / "pmd"
            pmd_step.update(run_pmd(
                pmd_bin, pmd_ruleset, repo_dir, pmd_out,
                java_version_default=args.pmd_java_version,
            ))
            # PMD exit code: 0 = no violations, 4 = violations found (with
            # --no-fail-on-violation this should stay 0). Anything else is
            # a real error.
            if pmd_step.get("exit_code", 0) not in (0, 4):
                provenance["errors"].append(
                    f"PMD exited with {pmd_step['exit_code']}"
                )
        provenance["steps"]["pmd"] = pmd_step

    provenance["finished_at"] = datetime.now(timezone.utc).isoformat()

    (out_dir / "run.json").write_text(
        json.dumps(provenance, indent=2, ensure_ascii=False)
    )
    _log(f"wrote {out_dir / 'run.json'}")

    if provenance["errors"]:
        _log(f"completed with {len(provenance['errors'])} warning(s)/error(s)")
    else:
        _log("completed cleanly")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except KeyboardInterrupt:
        _log("interrupted")
        raise SystemExit(130)
