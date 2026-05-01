# Research: The Technical Debt Dataset (TDD)

## Source

- Repository: https://github.com/clowee/The-Technical-Debt-Dataset
- Local DB: `javaAbstractor/tdd/td_V2.db` (SQLite, ~1.5 GB)

## Database Schema

### PROJECTS Table

```sql
CREATE TABLE PROJECTS (
  PROJECT_KEY TEXT,
  GIT_LINK TEXT,
  JIRA_LINK TEXT,
  SONAR_PROJECT_KEY TEXT,
  PROJECT_ID TEXT
);
```

### Other Tables

- `GIT_COMMITS` — commit history per project (PROJECT_ID, COMMIT_HASH, dates, etc.)
- `GIT_COMMITS_CHANGES` — file-level changes per commit
- `JIRA_ISSUES` — issue tracker data
- `SONAR_ANALYSIS` — SonarQube analysis snapshots
- `SONAR_ISSUES` — SonarQube detected issues
- `SONAR_MEASURES` — SonarQube metrics per analysis
- `SONAR_RULES` — SonarQube rule definitions
- `REFACTORING_MINER` — detected refactorings
- `SZZ_FAULT_INDUCING_COMMITS` — fault-inducing commit identification

## Projects (31 total)

All are Apache Foundation Java projects:

| PROJECT_KEY | GIT_LINK |
|-------------|----------|
| batik | https://github.com/apache/batik |
| commons-bcel | https://github.com/apache/commons-bcel |
| commons-beanutils | https://github.com/apache/commons-beanutils |
| cocoon | https://github.com/apache/cocoon |
| commons-codec | https://github.com/apache/commons-codec |
| commons-collections | https://github.com/apache/commons-collections |
| commons-cli | https://github.com/apache/commons-cli |
| commons-exec | https://github.com/apache/commons-exec |
| commons-fileupload | https://github.com/apache/commons-fileupload |
| commons-io | https://github.com/apache/commons-io |
| commons-jelly | https://github.com/apache/commons-jelly |
| commons-jexl | https://github.com/apache/commons-jexl |
| commons-configuration | https://github.com/apache/commons-configuration |
| commons-daemon | https://github.com/apache/commons-daemon |
| commons-dbcp | https://github.com/apache/commons-dbcp |
| commons-dbutils | https://github.com/apache/commons-dbutils |
| commons-digester | https://github.com/apache/commons-digester |
| felix | https://github.com/apache/felix |
| httpcomponents-client | https://github.com/apache/httpcomponents-client |
| httpcomponents-core | https://github.com/apache/httpcomponents-core |
| commons-jxpath | https://github.com/apache/commons-jxpath |
| commons-net | https://github.com/apache/commons-net |
| commons-ognl | https://github.com/apache/commons-ognl |
| santuario | https://github.com/apache/santuario-java |
| commons-validator | https://github.com/apache/commons-validator |
| commons-vfs | https://github.com/apache/commons-vfs |
| zookeeper | https://github.com/apache/zookeeper |
| archiva | https://github.com/apache/archiva |
| cayenne | https://github.com/apache/cayenne |
| hive | https://github.com/apache/hive |
| thrift | https://github.com/apache/thrift |

## Relevance to Java Abstractor

The abstractor needs to:

1. Clone each project from its GIT_LINK (at specific commits, potentially from
   GIT_COMMITS table).
2. Run the abstractor on the project source code.
3. Produce valid JSON/YAML output conforming to `genFeatureDef.md`.

### Key Challenges

- These are real-world, large Apache projects with complex Java features:
  generics, nested classes, enums, annotations, inheritance hierarchies,
  anonymous classes, lambdas, etc.
- Some projects (hive, zookeeper, batik) are very large.
- The abstractor currently requires a `pom.xml` (Maven project). All 31 projects
  should be Maven-based (Apache standard).
- Projects may use different Java versions; the abstractor targets Java 17.
