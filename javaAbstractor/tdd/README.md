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

| Project | Git Repository | Newest Analysis |
|---------|----------------|-----------------|
| santuario | apache/santuario-java | 2018-10-15 |
| commons-beanutils | apache/commons-beanutils | 2018-10-12 |
| commons-validator | apache/commons-validator | 2018-10-07 |
| commons-net | apache/commons-net | 2018-10-05 |
| commons-configuration | apache/commons-configuration | 2018-09-27 |
| commons-vfs | apache/commons-vfs | 2018-09-27 |
| commons-daemon | apache/commons-daemon | 2018-09-19 |
| commons-bcel | apache/commons-bcel | 2018-07-11 |
| commons-codec | apache/commons-codec | 2018-07-01 |
| commons-ognl | apache/commons-ognl | 2018-06-08 |
| commons-jxpath | apache/commons-jxpath | 2018-05-15 |
| commons-exec | apache/commons-exec | 2018-05-15 |
| commons-jexl | apache/commons-jexl | 2018-01-22 |
| commons-dbcp | apache/commons-dbcp | 2017-12-18 |
| commons-io | apache/commons-io | 2017-12-05 |
| commons-fileupload | apache/commons-fileupload | 2017-12-05 |
| commons-jelly | apache/commons-jelly | 2017-09-27 |
| commons-digester | apache/commons-digester | 2017-08-25 |
| commons-collections | apache/commons-collections | 2017-06-12 |
| commons-cli | apache/commons-cli | 2017-06-05 |
| commons-dbutils | apache/commons-dbutils | 2017-05-30 |
| httpcomponents-client | apache/httpcomponents-client | 2017-05-10 |
| httpcomponents-core | apache/httpcomponents-core | 2017-05-09 |
| zookeeper | apache/zookeeper | 2016-12-21 |
| hive | apache/hive | 2015-03-03 |
| thrift | apache/thrift | 2012-11-16 |
| archiva | apache/archiva | 2012-02-24 |
| felix | apache/felix | 2009-07-17 |
| cayenne | apache/cayenne | 2008-07-07 |
| cocoon | apache/cocoon | 2007-02-05 |
| batik | apache/batik | 2002-08-13 |

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
