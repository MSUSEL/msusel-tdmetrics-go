"""Batch driver that runs the Java abstractor against pinned Apache projects.

For each target `(project_key, git_slug, commit_sha)` taken from
`build_tdd_flat.README_PROJECTS`:

  1. Ensure the checkout at `<repo-base>/<git_slug>` exists and its HEAD is
     at the pinned `commit_sha` (fetch + checkout only if needed; never
     clones -- run collect_project_metrics.py first if the repo is missing).
  2. Invoke the Java abstractor jar with `-i <repo> -o abstractions/<key>.json -v`.
     Combined stdout+stderr is captured to `abstractions/<key>.log`
     (overwritten each run). No per-project timeout.
  3. Print a short "starting / successful / FAILED" line per project.

Before the per-project loop, the abstractor jar is rebuilt once via
`mvn clean compile assembly:single` in `../javaAbstractor/`. A build
failure aborts the pipeline before any project runs.

Design mirrors run_metrics_pipeline.py: sequential (metrics collection
is I/O + CPU heavy and clones can't be time-shared), non-destructive,
per-target failures are recorded and the pipeline continues unless
`--stop-on-error` is passed.

Usage:
    ./run_abstractor_pipeline.py                            # test set (4 projects)
    ./run_abstractor_pipeline.py --all                      # all 31 projects
    ./run_abstractor_pipeline.py --targets commons-io hive  # explicit subset
    ./run_abstractor_pipeline.py --skip-build               # reuse existing jar

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
# Curated test subset -- same 4 projects as run_metrics_pipeline.py so a
# quick pipeline smoke test covers the same code paths across tools.
# ---------------------------------------------------------------------------
DEFAULT_TEST_TARGETS: list[str] = [
    "commons-cli",
    "commons-exec",
    "commons-codec",
    "commons-io",
]

HERE          = Path(__file__).resolve().parent
REPO_ROOT     = HERE.parent
JAVA_ABS_DIR  = REPO_ROOT / "javaAbstractor"
DEFAULT_JAR   = JAVA_ABS_DIR / "target" / "abstractor-0.1-jar-with-dependencies.jar"
DEFAULT_ABS   = HERE / "abstractions"
DEFAULT_BASE  = Path("~/go/src/github.com").expanduser()


def _log(msg: str) -> None:
    print(f"[abstractor] {msg}", file=sys.stderr, flush=True)


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


def _build_jar(jar: Path) -> bool:
    _log(f"building javaAbstractor jar in {JAVA_ABS_DIR}")
    print("building javaAbstractor jar...", flush=True)
    cmd = ["mvn", "clean", "compile", "assembly:single"]
    started = time.time()
    proc = subprocess.run(cmd, cwd=str(JAVA_ABS_DIR))
    elapsed = time.time() - started
    if proc.returncode != 0:
        print(f"jar build FAILED (exit {proc.returncode}, {elapsed:.1f}s)",
              flush=True)
        return False
    if not jar.exists():
        print(f"jar build reported success but {jar} not found", flush=True)
        return False
    print(f"jar build OK ({elapsed:.1f}s)", flush=True)
    return True


def _ensure_checkout(repo_base: Path, git_slug: str, commit_sha: str) -> Path:
    """Verify the repo exists and HEAD matches commit_sha (fetch/checkout if not).

    Never clones. If the repo directory is missing, raises with a hint to
    run collect_project_metrics.py first.
    """
    org, _, name = git_slug.partition("/")
    if not org or not name:
        raise ValueError(f"expected git slug of form 'org/repo', got {git_slug!r}")
    repo_dir = repo_base / org / name

    if not repo_dir.exists():
        raise RuntimeError(
            f"{repo_dir} does not exist. Run collect_project_metrics.py "
            f"first (it clones on demand) or clone {git_slug} manually."
        )
    if not (repo_dir / ".git").exists():
        raise RuntimeError(
            f"{repo_dir} exists but is not a git checkout."
        )

    head = subprocess.run(
        ["git", "rev-parse", "HEAD"], cwd=str(repo_dir),
        check=True, text=True, capture_output=True,
    ).stdout.strip()
    if head == commit_sha:
        return repo_dir

    # Try a direct checkout, then fetch + retry.
    r = subprocess.run(
        ["git", "-c", "advice.detachedHead=false", "checkout", commit_sha],
        cwd=str(repo_dir), text=True, capture_output=True,
    )
    if r.returncode != 0:
        subprocess.run(["git", "fetch", "--all", "--tags", "--quiet"],
                       cwd=str(repo_dir), check=True)
        subprocess.run(
            ["git", "-c", "advice.detachedHead=false", "checkout", commit_sha],
            cwd=str(repo_dir), check=True,
        )

    head = subprocess.run(
        ["git", "rev-parse", "HEAD"], cwd=str(repo_dir),
        check=True, text=True, capture_output=True,
    ).stdout.strip()
    if head != commit_sha:
        raise RuntimeError(
            f"checkout landed at HEAD={head}, expected {commit_sha}"
        )
    return repo_dir


def _run_one(jar: Path, abs_dir: Path, repo_base: Path,
             project_key: str, git_slug: str, commit_sha: str) -> dict:
    json_path = abs_dir / f"{project_key}.json"
    log_path  = abs_dir / f"{project_key}.log"

    print(f"starting {project_key}...", flush=True)
    started = time.time()

    result: dict = {
        "project_key": project_key,
        "git_slug":    git_slug,
        "commit_sha":  commit_sha,
        "json_path":   str(json_path),
        "log_path":    str(log_path),
    }

    try:
        repo_dir = _ensure_checkout(repo_base, git_slug, commit_sha)
    except Exception as e:
        elapsed = time.time() - started
        log_path.write_text(f"checkout failed: {e}\n")
        result.update({
            "exit_code":   -1,
            "elapsed_sec": round(elapsed, 2),
            "error":       f"checkout: {e}",
        })
        print(f"{project_key} FAILED (checkout error, {elapsed:.1f}s) "
              f"-- see {log_path.relative_to(HERE)}", flush=True)
        return result

    cmd = [
        "java", "-jar", str(jar),
        "-i", str(repo_dir) + "/",
        "-o", str(json_path),
        "-v",
    ]

    # Combined stdout+stderr to the log file, overwritten each run.
    with log_path.open("w") as lf:
        lf.write(f"# cmd: {' '.join(cmd)}\n")
        lf.write(f"# cwd: {HERE}\n")
        lf.write(f"# repo: {repo_dir} @ {commit_sha}\n")
        lf.write(f"# started: {datetime.now(timezone.utc).isoformat()}\n")
        lf.write("# ---\n")
        lf.flush()
        proc = subprocess.run(cmd, cwd=str(HERE), stdout=lf,
                              stderr=subprocess.STDOUT)

    elapsed = time.time() - started
    result.update({
        "cmd":         cmd,
        "exit_code":   proc.returncode,
        "elapsed_sec": round(elapsed, 2),
    })

    if proc.returncode == 0:
        print(f"{project_key} successful ({elapsed:.1f}s)", flush=True)
    else:
        print(f"{project_key} FAILED (exit {proc.returncode}, {elapsed:.1f}s) "
              f"-- see {log_path.relative_to(HERE)}", flush=True)
    return result


def main() -> int:
    ap = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    ap.add_argument("--all", action="store_true",
                    help="run the full README project list (31 projects). "
                         "Overrides --targets.")
    ap.add_argument("--targets", nargs="+", metavar="PROJECT_KEY",
                    help=f"explicit project_keys to run "
                         f"(default: {DEFAULT_TEST_TARGETS})")
    ap.add_argument("--jar", default=str(DEFAULT_JAR),
                    help=f"path to abstractor jar "
                         f"(default: {DEFAULT_JAR})")
    ap.add_argument("--abs-out", default=str(DEFAULT_ABS),
                    help=f"output dir for <key>.json / <key>.log and summary "
                         f"(default: {DEFAULT_ABS})")
    ap.add_argument("--repo-base", default=str(DEFAULT_BASE),
                    help=f"base dir holding <org>/<repo> checkouts "
                         f"(default: {DEFAULT_BASE})")
    ap.add_argument("--skip-build", action="store_true",
                    help="skip 'mvn clean compile assembly:single'; reuse "
                         "the existing jar as-is")
    ap.add_argument("--stop-on-error", action="store_true",
                    help="stop the pipeline on the first non-zero exit "
                         "(default: continue and record failures)")
    args = ap.parse_args()

    jar       = Path(args.jar).expanduser().resolve()
    abs_dir   = Path(args.abs_out).expanduser().resolve()
    repo_base = Path(args.repo_base).expanduser().resolve()

    targets = _resolve_targets(args.all, args.targets)
    _log(f"targets ({len(targets)}): {[t[0] for t in targets]}")

    abs_dir.mkdir(parents=True, exist_ok=True)

    if not args.skip_build:
        if not _build_jar(jar):
            return 2
    else:
        _log(f"--skip-build set; using existing jar at {jar}")
        if not jar.exists():
            _log(f"error: jar not found: {jar}")
            return 2

    started_at = datetime.now(timezone.utc).isoformat()
    results: list[dict] = []
    for project_key, git_slug, commit_sha in targets:
        try:
            r = _run_one(jar, abs_dir, repo_base,
                         project_key, git_slug, commit_sha)
        except KeyboardInterrupt:
            _log("interrupted; recording partial results")
            results.append({
                "project_key": project_key, "git_slug": git_slug,
                "commit_sha":  commit_sha,  "exit_code": 130,
                "interrupted": True,
            })
            print(f"{project_key} FAILED (interrupted)", flush=True)
            break
        results.append(r)
        if r["exit_code"] != 0 and args.stop_on_error:
            _log(f"stopping: {project_key} exited {r['exit_code']}")
            break

    finished_at = datetime.now(timezone.utc).isoformat()

    stamp = started_at.replace(":", "").replace("-", "").split(".")[0]
    summary_path = abs_dir / f"abstraction_run_{stamp}Z.json"

    ok   = [r for r in results if r.get("exit_code") == 0]
    fail = [r for r in results if r.get("exit_code") not in (0, None)]

    summary = {
        "started_at":  started_at,
        "finished_at": finished_at,
        "jar":         str(jar),
        "repo_base":   str(repo_base),
        "abs_out":     str(abs_dir),
        "targets_requested": [t[0] for t in targets],
        "counts": {
            "total":   len(results),
            "success": len(ok),
            "failed":  len(fail),
        },
        "results": results,
    }
    summary_path.write_text(json.dumps(summary, indent=2, ensure_ascii=False))

    print("", flush=True)
    print(f"summary: {len(ok)} ok / {len(fail)} failed / {len(results)} total",
          flush=True)
    for r in results:
        status = "OK  " if r.get("exit_code") == 0 else f"FAIL({r.get('exit_code')})"
        print(f"  {status}  {r['project_key']:<24} "
              f"{r.get('elapsed_sec', '?'):>6}s", flush=True)
    print(f"wrote {summary_path}", flush=True)

    return 0 if not fail else 1


if __name__ == "__main__":
    raise SystemExit(main())
