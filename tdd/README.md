# Technical Debt Dataset

The TDD (Technical Debt Dataset) is found at
<https://github.com/clowee/The-Technical-Debt-Dataset>

## Getting started

If this folder doesn't contain a `td_V2.db` file, copy the `td_V2.db` from the
[V2.0](https://github.com/clowee/The-Technical-Debt-Dataset/releases/tag/2.0.1)
release into this folder.

## Agent Notes

I asked an AI agent to read the `td_V2.db` file and provide me nodes on what
is contained and how to read it. The following is the agent's notes.

### Database Overview

The `td_V2.db` file is a SQLite database containing technical debt analysis data
for 31 Apache Java projects. The data was collected using SonarQube and includes
code smells, git history, Jira issues, and refactoring information.

### Tables

| Table | Description |
|-------|-------------|
| `PROJECTS` | 31 Apache projects with Git/Jira/Sonar links |
| `SONAR_ANALYSIS` | Analysis snapshots (project, date, revision) |
| `SONAR_ISSUES` | Code smells, bugs, vulnerabilities detected |
| `SONAR_MEASURES` | Project-level metrics per analysis |
| `SONAR_RULES` | Rule definitions and descriptions |
| `GIT_COMMITS` | Commit history metadata |
| `GIT_COMMITS_CHANGES` | Per-file changes in commits |
| `JIRA_ISSUES` | Bug/feature tracking data |
| `REFACTORING_MINER` | Detected refactorings per commit |
| `SZZ_FAULT_INDUCING_COMMITS` | SZZ algorithm fault-inducing commits |

### Projects and Newest Analysis Dates

| Project | Git Repository | Analysis Date | Commit SHA |
|---------|----------------|---------------|------------|
| santuario | apache/santuario-java | 2018-10-15 | `be4e2331f77adb1e479406ebf973e516bbf5e32b` |
| commons-beanutils | apache/commons-beanutils | 2018-10-12 | `c4da598872233b59af41a221bd2bdcefbbca1259` |
| commons-validator | apache/commons-validator | 2018-10-07 | `a3771313c9f1833abf32c7c294ad1de4810e532d` |
| commons-net | apache/commons-net | 2018-10-05 | `fb7aae4c64f7d2bf6dced00c49c3ffc428b2d572` |
| commons-configuration | apache/commons-configuration | 2018-09-27 | `15b4031ba94a60f20b854e6ce2c7964d77086387` |
| commons-vfs | apache/commons-vfs | 2018-09-27 | `d72192f18bfaed730b4f37a2f94853e1503ffd74` |
| commons-daemon | apache/commons-daemon | 2018-09-19 | `1ffa799cb3ddf5a4a918e59e46cd9868ee766b19` |
| commons-bcel | apache/commons-bcel | 2018-07-11 | `6ed18c5bef0f5b93b54783a8e8fb2b9042da26ac` |
| commons-codec | apache/commons-codec | 2018-07-01 | `db51a1cb41e9155ca028a73b0637b32a2c37c43a` |
| commons-ognl | apache/commons-ognl | 2018-06-08 | `6ec1a1a4588b82c0972ca2ff35b85d9b50cc4604` |
| commons-jxpath | apache/commons-jxpath | 2018-05-15 | `eff47ab8ca52fdbc91d1313cc224324465dd043e` |
| commons-exec | apache/commons-exec | 2018-05-15 | `2da60ab3eefaaa2f8a434ded1eebe1ce17efd34a` |
| commons-jexl | apache/commons-jexl | 2018-01-22 | `d3e702149a3db297d6db2c0b7671807f5c7b98fc` |
| commons-dbcp | apache/commons-dbcp | 2017-12-18 | `d8dd39b32bbb04a28ea86eb826c56aa6783f3faf` |
| commons-io | apache/commons-io | 2017-12-05 | `65c4a9c0ec651dd99f28b9fae40378728d071985` |
| commons-fileupload | apache/commons-fileupload | 2017-12-05 | `cae90facebc54803232a0593003914ca77193a73` |
| commons-jelly | apache/commons-jelly | 2017-09-27 | `48c008cc2328402e0976295625b32c5197ba2324` |
| commons-digester | apache/commons-digester | 2017-08-25 | `c1d0e563339faec040eb036ae97a7b7bf07ba865` |
| commons-collections | apache/commons-collections | 2017-06-12 | `f0f364fd9d946483f947011a3557c1e6f2e5d8ee` |
| commons-cli | apache/commons-cli | 2017-06-05 | `92f1def0bb3c0345295012e36b7150cfd1d7b6ab` |
| commons-dbutils | apache/commons-dbutils | 2017-05-30 | `2f48485a82697d9aed060ba36f6d5beb3a58ed8b` |
| httpcomponents-client | apache/httpcomponents-client | 2017-05-10 | `8a1b96bfa75382c0b94d70f6914fbb9bfeb0451e` |
| httpcomponents-core | apache/httpcomponents-core | 2017-05-09 | `3a677d47cb872b6ede20b28e93d3206f08b349ac` |
| zookeeper | apache/zookeeper | 2016-12-21 | `eac693cc76a34f96b9116ef33d1e92af7129416d` |
| hive | apache/hive | 2015-03-03 | `a4d91eaf2925239aa29342f7e5b0f8680c842390` |
| thrift | apache/thrift | 2012-11-16 | `a2123693838410c1e78170419e9bb91cb01151b4` |
| archiva | apache/archiva | 2012-02-24 | `374fc983abc92df8aa4f8ef30caee94b34312ad2` |
| felix | apache/felix | 2009-07-17 | `bdb6cb5cac0d81e9cd3fda666065e0e577eb9c41` |
| cayenne | apache/cayenne | 2008-07-07 | `b9988a83e364b9b470873dff8996dcf401d08dc4` |
| cocoon | apache/cocoon | 2007-02-05 | `a80f73b27592a2794c9133ee03d2e402bf12ecc1` |
| batik | apache/batik | 2002-08-13 | `2bb3a6ea5a6258ff6372e2493b81d7768d6bb494` |

**To clone a project at the exact analyzed version:**

```bash
# Example: Clone commons-io at the analyzed commit
git clone https://github.com/apache/commons-io.git
cd commons-io
git checkout 65c4a9c0ec651dd99f28b9fae40378728d071985
```

### WMC, TCC, and ATFD Metrics

**Important**: The raw WMC (Weighted Method Count), TCC (Tight Class Cohesion),
and ATFD (Access to Foreign Data) values are **not stored directly** in this
database. These metrics are used internally by SonarQube's code smell detection
plugins to identify God Classes and other anti-patterns.

What IS available:

1. **`code_smells:complex_class`** - Classes flagged as "God Classes" using the
   detection strategy from Lanza & Marinescu's "Object-Oriented Metrics in Practice"
   which uses WMC, TCC, and ATFD thresholds internally.

2. **`SONAR_MEASURES`** - Project-level metrics including:
   - `COMPLEXITY` - Total cyclomatic complexity
   - `CLASS_COMPLEXITY` - Average complexity per class
   - `COGNITIVE_COMPLEXITY` - Cognitive complexity metric
   - `AFFERENT_COUPLINGS` / `EFFERENT_COUPLINGS` - Coupling metrics
   - `CODE_SMELLS` - Total code smell count

3. **`squid:ClassCyclomaticComplexity`** / **`squid:MethodCyclomaticComplexity`** -
   Issues for classes/methods exceeding complexity thresholds.

### Code Smell Rules Available

The `code_smells` plugin provides these detection rules:

- `complex_class` - God Class (high cyclomatic complexity)
- `blob_class` - Large class monopolizing processing
- `large_class` - Class too large in LOC
- `long_method` - Methods too long
- `lazy_class` - Classes with few fields/methods
- `spaghetti_code` - Unstructured long methods
- `swiss_army_knife` - Complex class with many interfaces
- `long_parameter_list` - Methods with too many parameters
- `refused_parent_bequest` - Broken polymorphism
- `many_field_attributes_not_complex` - Data classes (related to ATFD)

### Example Queries

**List all projects:**

```sql
SELECT PROJECT_KEY, GIT_LINK, JIRA_LINK 
FROM PROJECTS;
```

**Get newest analysis for each project:**

```sql
SELECT p.PROJECT_KEY, a.DATE, a.ANALYSIS_KEY
FROM PROJECTS p
JOIN SONAR_ANALYSIS a ON p.PROJECT_ID = a.PROJECT_ID
WHERE a.DATE = (
    SELECT MAX(a2.DATE) 
    FROM SONAR_ANALYSIS a2 
    WHERE a2.PROJECT_ID = p.PROJECT_ID
)
ORDER BY a.DATE DESC;
```

**Find God Classes (complex_class) in a project at newest version:**

```sql
-- Get all God Class issues for commons-io at its newest analysis
SELECT i.COMPONENT, i.CREATION_DATE, i.STATUS
FROM SONAR_ISSUES i
JOIN SONAR_ANALYSIS a ON i.PROJECT_ID = a.PROJECT_ID 
    AND i.CREATION_ANALYSIS_KEY = a.ANALYSIS_KEY
WHERE i.PROJECT_ID = 'org.apache:commons-io'
    AND i.RULE = 'code_smells:complex_class'
    AND a.DATE = (
        SELECT MAX(DATE) FROM SONAR_ANALYSIS 
        WHERE PROJECT_ID = 'org.apache:commons-io'
    );
```

**Get project-level complexity metrics for newest analysis:**

```sql
SELECT m.PROJECT_ID, m.COMPLEXITY, m.CLASS_COMPLEXITY, 
       m.COGNITIVE_COMPLEXITY, m.CODE_SMELLS, m.CLASSES
FROM SONAR_MEASURES m
JOIN SONAR_ANALYSIS a ON m.PROJECT_ID = a.PROJECT_ID 
    AND m.ANALYSIS_KEY = a.ANALYSIS_KEY
WHERE m.PROJECT_ID = 'org.apache:commons-io'
    AND a.DATE = (
        SELECT MAX(DATE) FROM SONAR_ANALYSIS 
        WHERE PROJECT_ID = 'org.apache:commons-io'
    );
```

**Count code smells by type for a project:**

```sql
SELECT RULE, COUNT(*) as count
FROM SONAR_ISSUES
WHERE PROJECT_ID = 'org.apache:commons-io'
    AND TYPE = 'CODE_SMELL'
GROUP BY RULE
ORDER BY count DESC
LIMIT 20;
```

### Python Example

```python
import sqlite3

db_path = "td_V2.db"
conn = sqlite3.connect(db_path)
conn.row_factory = sqlite3.Row

# Get God Class issues for commons-io at newest version
query = """
SELECT i.COMPONENT, i.CREATION_DATE, i.STATUS, i.MESSAGE
FROM SONAR_ISSUES i
JOIN SONAR_ANALYSIS a ON i.PROJECT_ID = a.PROJECT_ID
WHERE i.PROJECT_ID = 'org.apache:commons-io'
    AND i.RULE = 'code_smells:complex_class'
    AND a.DATE = (SELECT MAX(DATE) FROM SONAR_ANALYSIS 
                  WHERE PROJECT_ID = 'org.apache:commons-io')
GROUP BY i.COMPONENT
"""

for row in conn.execute(query):
    # Extract class name from component path
    component = row['COMPONENT']
    class_path = component.split(':')[-1]  # e.g., src/main/java/...
    print(f"God Class: {class_path}")

conn.close()
```

### Relationship to This Project

The Java Abstractor in this repository processes these same Apache projects to
extract structural information (classes, methods, fields, interfaces, etc.)
conforming to `docs/genFeatureDef.md`. The TDD database provides pre-computed
SonarQube metrics for validation and correlation with our own technical debt
analysis.

**Use case**: Compare our abstraction's computed metrics against the TDD's
SonarQube-detected code smells to validate the abstractor's accuracy.

### Computing WMC, TCC, and ATFD Directly

Since the TDD database doesn't store raw metric values, you can compute them
using these tools:

#### Option 1: SourceMeter (Recommended for Research)

[SourceMeter](https://www.sourcemeter.com/) is a static analysis tool that
exports class-level metrics including WMC, TCC, and ATFD directly to CSV.

```bash
# Download SourceMeter for Java from sourcemeter.com
# Run analysis on a Maven project:
SourceMeterJava \
    -projectName=commons-io \
    -projectBaseDir=/path/to/commons-io \
    -resultsDir=./results \
    -runFB=false \
    -runPMD=false

# Output includes CSV files with class metrics:
# results/commons-io/java/<timestamp>/commons-io-Class.csv
# Columns include: WMC, TCC, ATFD, LOC, NOM, etc.
```

**SourceMeter metric definitions:**
- **WMC** (Weighted Methods per Class): Sum of cyclomatic complexity of all methods
- **TCC** (Tight Class Cohesion): Ratio of directly connected method pairs
- **ATFD** (Access to Foreign Data): Number of external attributes accessed

#### Option 2: PMD with Custom Rules

[PMD](https://pmd.github.io/) can compute some metrics but requires custom
rules for TCC and ATFD. WMC is available via the metrics framework.

```bash
# Install PMD (https://pmd.github.io/)
# Create a ruleset that exports metrics:

pmd check -d /path/to/source \
    -R category/java/design.xml \
    -f csv \
    -r metrics-report.csv
```

**PMD metrics available:**
- `WMC` - Available via `CyclomaticComplexity` rule aggregated per class
- `TCC` - Not built-in; requires custom rule or post-processing
- `ATFD` - Not built-in; requires custom rule

For WMC specifically, use PMD's designer to create a custom XPath rule:

```xml
<rule name="ClassWMC"
      language="java"
      message="WMC = {0}"
      class="net.sourceforge.pmd.lang.rule.xpath.XPathRule">
    <properties>
        <property name="xpath">
            <value>
//ClassOrInterfaceDeclaration[
    @Interface = false()
]
            </value>
        </property>
    </properties>
</rule>
```

#### Option 3: SonarQube with API Export

If you have a SonarQube server with projects analyzed, you can export metrics
via the Web API:

```bash
# Get component metrics (requires SonarQube server)
curl -u admin:password \
    "http://localhost:9000/api/measures/component_tree?\
component=org.apache:commons-io&\
metricKeys=complexity,class_complexity,cognitive_complexity&\
qualifiers=FIL"
```

**Note**: SonarQube computes but doesn't directly expose TCC and ATFD in its
standard API. The `code_smells` plugin uses them internally for God Class
detection but doesn't publish them as separate metrics.

#### Option 4: JHawk / Understand

Commercial tools like [JHawk](https://www.virtualmachinery.com/jhawk.htm) and
[SciTools Understand](https://scitools.com/) provide comprehensive OO metrics
including WMC, TCC (as "Cohesion"), and coupling metrics.

#### Option 5: Compute from Java Abstractor Output

This project's Java Abstractor extracts the structural information needed to
compute these metrics. After running the abstractor:

```python
# Pseudo-code for computing metrics from abstractor JSON output
import json

with open('abstraction.json') as f:
    project = json.load(f)

for obj in project.get('objects', []):
    methods = obj.get('methods', [])
    
    # WMC: Sum of complexity of all methods
    wmc = sum(m.get('complexity', 1) for m in methods)
    
    # TCC: Requires analyzing method bodies for shared field access
    # (field access data available in method metrics)
    
    # ATFD: Count of foreign attribute accesses
    # (available from 'reads' and 'writes' in method metrics)
    
    print(f"{obj['name']}: WMC={wmc}")
```

The `techDebtMetrics` .NET component is designed to compute these metrics from
the abstractor's JSON output once fully implemented.

### Metric Thresholds (God Class Detection)

Per Lanza & Marinescu's detection strategy, a class is a "God Class" if:

| Metric | Threshold | Meaning |
|--------|-----------|---------|
| WMC | > 47 | High total method complexity |
| TCC | < 0.33 | Low cohesion (< 1/3 of method pairs connected) |
| ATFD | > 5 | Accesses many foreign attributes |

A class is flagged as a God Class if **all three** conditions are met.

### References

- Dataset: <https://github.com/clowee/The-Technical-Debt-Dataset>
- Lanza, M. & Marinescu, R. (2006). *Object-Oriented Metrics in Practice*. Springer.
- God Class detection strategy: pp. 80-81 of the above reference
- SourceMeter: <https://www.sourcemeter.com/>
- PMD: <https://pmd.github.io/>
- SonarQube API: <https://docs.sonarqube.org/latest/extend/web-api/>

### Running python tools

The following is a short summary of how to run the python tools. All three
are stdlib-only Python 3 scripts and live alongside this README. Run them
from this folder unless a `--` flag says otherwise.

The commit SHA of every project is hardcoded in `build_tdd_flat.py` as the
`README_PROJECTS` constant, and mirrors the "Projects and Newest Analysis
Dates" table above. `run_metrics_pipeline.py` imports that constant, so
neither script requires you to type SHAs on the command line — you refer
to projects by short `project_key` (e.g. `commons-io`).

### build_tdd_flat.py

Rebuilds `tdd_flat.json` — a flat, human-readable snapshot of the TDD
database restricted to the 31 pinned commits and to the code-smell rules
called out above.

```bash
# Rebuild tdd_flat.json using ./td_V2.db.
python3 build_tdd_flat.py

# Point at a db / output path in a different location.
python3 build_tdd_flat.py --db /some/where/td_V2.db --out /tmp/tdd.json
```

### collect_project_metrics.py

Runs the metric collectors (CK and/or PMD) against **one** project at a
specific commit. Clones the repo into
`~/go/src/github.com/<org>/<repo>/` on first run, then reuses the clone
on subsequent runs (delete or rename the folder to force a fresh clone).

```bash
# Run both CK and PMD for one project (positional args: key, slug, sha).
python3 collect_project_metrics.py \
    commons-io apache/commons-io 65c4a9c0ec651dd99f28b9fae40378728d071985

# PMD only (useful before ck.jar is installed).
python3 collect_project_metrics.py --tools pmd \
    commons-io apache/commons-io 65c4a9c0ec651dd99f28b9fae40378728d071985
```

Defaults it will look for:

- `--ck-jar   ~/tools/ck.jar`
- `--pmd-bin  ~/pmd-bin-7.26.0/bin/pmd`
- `--repo-base ~/go/src/github.com`
- `--out      ./metrics_output`

Missing CK jar or PMD binary is non-fatal: the corresponding step is
skipped with a clear message so the other tool still runs.

### run_metrics_pipeline.py

Batch driver that loops `collect_project_metrics.py` over a subset (or
all) of the pinned projects. Uses `README_PROJECTS` from
`build_tdd_flat.py` for the SHA of each requested `project_key`.

Four ways to select targets:

```bash
# 1. Default: run a curated 4-project test set
#    (commons-cli, commons-exec, commons-codec, commons-io).
python3 run_metrics_pipeline.py

# 2. Run all 31 projects.
python3 run_metrics_pipeline.py --all

# 3. Explicit subset by short project_key.
python3 run_metrics_pipeline.py --targets commons-io hive zookeeper

# 4. Same as any of the above, restricted to one tool.
python3 run_metrics_pipeline.py --all --tools pmd
```

Per-project results land in `metrics_output/<key>-<short_sha>/`; a
pipeline-level summary lands in
`metrics_output/pipeline_run_<UTC>.json`. Failures don't stop the run
unless you pass `--stop-on-error`.

### run_abstractor_pipeline.py

Batch driver that runs the Java abstractor against pinned Apache
projects and drops one `abstractions/<project_key>.json` (and matching
`.log` capturing stdout/stderr) per project. Same target-selection UX
as `run_metrics_pipeline.py` (test set / `--all` / `--targets`).

By default the abstractor jar is rebuilt once at the start via
`mvn clean compile assembly:single` in `../javaAbstractor/`. Each
target is then checked out at its pinned commit and the abstractor is
invoked with `-i <repo>/ -o abstractions/<key>.json -v`. Existing
`<key>.json` / `<key>.log` are overwritten every run. Each project has
a per-run timeout (in whole minutes) controlled by `--timeout`
(default `10`; use `--timeout 0` to disable). On timeout the child is
killed, a `# TIMED OUT ...` line is appended to the log, and the
project is recorded as failed.

Repos are expected to already be cloned under `~/go/src/github.com/`
(run `collect_project_metrics.py` first if a repo is missing -- this
script never clones).

```bash
# Default: run the curated 4-project test set.
python3 run_abstractor_pipeline.py

# All 31 pinned projects.
python3 run_abstractor_pipeline.py --all

# Explicit subset.
python3 run_abstractor_pipeline.py --targets commons-io commons-bcel

# Reuse the existing jar (skip the mvn build).
python3 run_abstractor_pipeline.py --skip-build --targets commons-io

# Stop on first failure instead of continuing through the list.
python3 run_abstractor_pipeline.py --all --stop-on-error

# Override the per-project timeout (whole minutes; 0 = no timeout).
python3 run_abstractor_pipeline.py --all --timeout 20
```

Console prints one `starting <key>...` line and one
`<key> successful (Ns)` / `<key> FAILED (exit N, Ns)` line per
project, plus a final summary table. A pipeline-level summary is
written to `abstractions/abstraction_run_<UTC>.json`.

### build_per_project.py

Merges the TDD portion of `tdd_flat.json` with the per-class CK and PMD
outputs produced by `run_metrics_pipeline.py`, writing one file per
project at `per_project/<project_key>.json`. Each output file is a
single top-level object mirroring an entry in `tdd_flat.json`'s
`projects` list plus a `classes` array containing the merged CK + PMD
metrics (with `methods_ck` and `methods_pmd` nested per class).

Repo-absolute paths from CK/PMD and Sonar-encoded paths from
`code_smell_issues.component` are all trimmed to repo-relative form so
paths across sources are directly comparable. CK's `NaN` values (e.g.
`tcc` / `lcc` on 0-or-1-method classes) are emitted as the JSON string
`"NaN"` so the output is spec-valid JSON.

```bash
# Merge all 31 projects.
python3 build_per_project.py

# Only rebuild a subset.
python3 build_per_project.py --targets commons-io commons-cli

# Read / write to non-default locations.
python3 build_per_project.py \
    --metrics /some/where/metrics_output \
    --out     /some/where/per_project
```

Requires `tdd_flat.json` (from `build_tdd_flat.py`) and populated
`metrics_output/<key>-<sha>/{ck,pmd}/` directories (from
`run_metrics_pipeline.py`).
