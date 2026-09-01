"""Batch driver for collect_project_metrics.py.

Iterates a curated list of Apache projects (pinned to the same commits as
README.md / tdd_flat.json) and invokes collect_project_metrics.py for each.

Design notes:
  - Reuses collect_project_metrics.py as a subprocess so the per-project
    logic stays in one place and both scripts can be run standalone.
  - Reads the canonical project list from build_tdd_flat.README_PROJECTS
    (single source of truth). Local --targets overrides which subset runs.
  - Sequential by default. Metrics collection is I/O-heavy on git and
    CPU-heavy on PMD/CK, and parallel runs would fight for the same
    per-repo working tree if the same project were listed twice. Kept
    linear for now; can add a --jobs flag later once the pipeline is
    proven.
  - Non-destructive: never deletes clones, never rewrites tracked files,
    never stages anything. On per-target failure, records the failure
    and moves on unless --stop-on-error is passed.

Usage:
    ./run_metrics_pipeline.py                            # test set (4 projects)
    ./run_metrics_pipeline.py --all                      # all 31 projects
    ./run_metrics_pipeline.py --targets commons-io hive  # explicit subset
    ./run_metrics_pipeline.py --tools pmd                # PMD only, skip CK

Note: This script was created by Opus 4.7 High
"""

from __future__ import annotations

import argparse
import json
import subprocess
import sys
import time
from datetime import datetime, timezone
from pathlib import Path

# Same-folder import: build_tdd_flat.py holds the canonical (project_key,
# git_slug, commit_sha) list used to build tdd_flat.json.
sys.path.insert(0, str(Path(__file__).resolve().parent))
from build_tdd_flat import README_PROJECTS  # noqa: E402

# ---------------------------------------------------------------------------
# Curated test subset.
#
# Chosen to cover different sizes/ages while staying fast:
#   commons-cli    - already cloned at target SHA (checkout is a no-op)
#   commons-exec   - very small codebase, ~20 code-smell issues in tdd_flat
#   commons-codec  - small, ~44 issues
#   commons-io     - mid-size, ~87 issues, well-known reference project
# ---------------------------------------------------------------------------
DEFAULT_TEST_TARGETS: list[str] = [
    "commons-cli",
    "commons-exec",
    "commons-codec",
    "commons-io",
]


def _log(msg: str) -> None:
    print(f"[pipeline] {msg}", file=sys.stderr, flush=True)


def _resolve_targets(all_projects: bool, targets: list[str] | None
                     ) -> list[tuple[str, str, str]]:
    if all_projects:
        return list(README_PROJECTS)

    wanted = targets if targets else DEFAULT_TEST_TARGETS
    by_key = {p[0]: p for p in README_PROJECTS}
    missing = [t for t in wanted if t not in by_key]
    if missing:
        raise SystemExit(
            f"error: unknown project_key(s): {missing}\n"
            f"       known keys: {sorted(by_key)}"
        )
    return [by_key[t] for t in wanted]


def _run_one(collect_script: Path, project_key: str, git_slug: str,
             commit_sha: str, passthrough: list[str]) -> dict:
    cmd = [
        sys.executable, str(collect_script),
        project_key, git_slug, commit_sha,
        *passthrough,
    ]
    _log(f"--- {project_key} ({git_slug} @ {commit_sha[:10]}) ---")
    _log("$ " + " ".join(cmd))
    started = time.time()
    # Inherit stdout/stderr so live output streams to the terminal.
    proc = subprocess.run(cmd, check=False)
    elapsed = time.time() - started
    return {
        "project_key":  project_key,
        "git_slug":     git_slug,
        "commit_sha":   commit_sha,
        "cmd":          cmd,
        "exit_code":    proc.returncode,
        "elapsed_sec":  round(elapsed, 2),
    }


def main() -> int:
    here = Path(__file__).resolve().parent
    default_collect = here / "collect_project_metrics.py"
    default_out     = here / "metrics_output"

    ap = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    ap.add_argument("--all", action="store_true",
                    help="run the full README project list (31 projects). "
                         "Overrides --targets.")
    ap.add_argument("--targets", nargs="+", metavar="PROJECT_KEY",
                    help=f"explicit project_keys to run "
                         f"(default: {DEFAULT_TEST_TARGETS})")
    ap.add_argument("--collect-script", default=str(default_collect),
                    help=f"path to collect_project_metrics.py "
                         f"(default: {default_collect})")
    ap.add_argument("--out", default=str(default_out),
                    help=f"output root dir for individual runs and the "
                         f"aggregate summary (default: {default_out})")
    ap.add_argument("--tools", default=None,
                    help="pass-through --tools value for collect script "
                         "(e.g. 'ck', 'pmd', 'ck,pmd'). Omitted -> collect "
                         "script default.")
    ap.add_argument("--ck-jar",   default=None)
    ap.add_argument("--pmd-bin",  default=None)
    ap.add_argument("--repo-base", default=None)
    ap.add_argument("--stop-on-error", action="store_true",
                    help="stop the pipeline on the first non-zero exit "
                         "(default: continue and record failures)")
    args = ap.parse_args()

    collect_script = Path(args.collect_script).expanduser().resolve()
    if not collect_script.exists():
        _log(f"error: collect script not found: {collect_script}")
        return 2

    targets = _resolve_targets(args.all, args.targets)
    _log(f"targets ({len(targets)}): {[t[0] for t in targets]}")

    # Assemble pass-through args for the per-project script.
    passthrough: list[str] = []
    if args.tools     is not None: passthrough += ["--tools",     args.tools]
    if args.ck_jar    is not None: passthrough += ["--ck-jar",    args.ck_jar]
    if args.pmd_bin   is not None: passthrough += ["--pmd-bin",   args.pmd_bin]
    if args.repo_base is not None: passthrough += ["--repo-base", args.repo_base]
    if args.out       is not None: passthrough += ["--out",       args.out]

    started_at = datetime.now(timezone.utc).isoformat()
    results: list[dict] = []
    for project_key, git_slug, commit_sha in targets:
        try:
            r = _run_one(collect_script, project_key, git_slug, commit_sha,
                         passthrough)
        except KeyboardInterrupt:
            _log("interrupted; recording partial results")
            results.append({
                "project_key": project_key, "git_slug": git_slug,
                "commit_sha":  commit_sha,  "exit_code": 130,
                "interrupted": True,
            })
            break
        results.append(r)
        if r["exit_code"] != 0 and args.stop_on_error:
            _log(f"stopping: {project_key} exited {r['exit_code']}")
            break

    finished_at = datetime.now(timezone.utc).isoformat()

    out_root = Path(args.out).expanduser().resolve()
    out_root.mkdir(parents=True, exist_ok=True)
    stamp = started_at.replace(":", "").replace("-", "").split(".")[0]
    summary_path = out_root / f"pipeline_run_{stamp}Z.json"

    ok    = [r for r in results if r.get("exit_code") == 0]
    fail  = [r for r in results if r.get("exit_code") not in (0, None)]

    summary = {
        "started_at":  started_at,
        "finished_at": finished_at,
        "collect_script": str(collect_script),
        "passthrough_args": passthrough,
        "targets_requested": [t[0] for t in targets],
        "counts": {
            "total":    len(results),
            "success":  len(ok),
            "failed":   len(fail),
        },
        "results": results,
    }
    summary_path.write_text(json.dumps(summary, indent=2, ensure_ascii=False))

    _log("")
    _log(f"summary: {len(ok)} ok / {len(fail)} failed / {len(results)} total")
    for r in results:
        status = "OK " if r.get("exit_code") == 0 else f"FAIL({r.get('exit_code')}) "
        _log(f"  {status}  {r['project_key']:<24} {r.get('elapsed_sec', '?'):>6}s")
    _log(f"wrote {summary_path}")

    return 0 if not fail else 1


if __name__ == "__main__":
    raise SystemExit(main())
