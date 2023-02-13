# Experiment 001 - Participation

- Documentation:
  - [Getting the java AST](https://www.geeksforgeeks.org/abstract-syntax-tree-ast-in-java/)
  - [Checkstyle via command line](https://checkstyle.org/cmdline.html)
- Running experiments:
  - First navigate to this folder: `cd .\javaexps\exp001\`
  - To run checkstyle to analyze some code with a command similar to:
      `java -jar checkstyle-8.43-all.jar -c /google_checks.xml YourFile.java`
    - e.g. `java -jar C:\Data\Code\checkstyle-10.6.0-all.jar -c sun_checks.xml data001\main.java`
  - To print AST with a command similar to:
      `java -jar checkstyle-8.43-all.jar -t YourFile.java`
    - e.g. `java -jar C:\Data\Code\checkstyle-10.6.0-all.jar -t data001\main.java`
