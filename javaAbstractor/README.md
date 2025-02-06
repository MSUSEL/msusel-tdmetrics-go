# Java Abstractor

- [Java Abstractor](#java-abstractor)
  - [Setup Abstractor](#setup-abstractor)
  - [Running Abstractor](#running-abstractor)
  - [Running Tests](#running-tests)
  - [Resources](#resources)

## Setup Abstractor

1. Install OpenJDK 17.0.14 for Java 17.
   1. This project was built using
      [OpenLogic](https://www.openlogic.com/openjdk-downloads).
   2. See the [Java 17 Almanac](https://javaalmanac.io/jdk/17/)  

2. Install Maven.
   1. This project was built using
      [Apache Maven v3.9.9](https://maven.apache.org/download.cgi)

3. Check the install by running the Maven version check `mvn -v`.

   The result should look similar to following except with your own
   install paths:

    ```Plain
    msusel-tdmetrics-go\javaAbstractor> mvn -v
    Apache Maven 3.9.9 (8e8579a9e76f7d015ee5ec7bfcdc97d260186937)
    Maven home: C:\Program Files\Apache\Maven\3.9.9
    Java version: 17.0.14, vendor: OpenLogic, runtime: C:\Program Files (x86)\OpenJDK\17.0.14
    Default locale: en_US, platform encoding: Cp1252
    OS name: "windows 11", version: "10.0", arch: "amd64", family: "windows"
    ```

## Running Abstractor

1. Compile with `mvn clean compile assembly:single`.

2. Run with `java -jar .\target\abstractor-0.1-jar-with-dependencies.jar <options>`.

3. For help with the `<options>` use `-help`.

## Running Tests

Run the tests with `mvn test`.

## Resources

- Spoon
  - [Spoon Forge](https://spoon.gforge.inria.fr/)
    - Has the Spoon BibTeX reference to [Paper](https://inria.hal.science/hal-01078532/document)
  - [Maven Repository](https://central.sonatype.com/artifact/fr.inria.gforge.spoon/spoon-core)
  - [Examples](https://github.com/SpoonLabs/spoon-examples/tree/master)

- Eclipse JDT API
  - [Project](https://projects.eclipse.org/projects/eclipse.jdt)
  - [Github](https://github.com/eclipse-jdt/eclipse.jdt.core)
  - [Programmer Guide](https://github.com/eclipse-jdt/eclipse.jdt.core/wiki/Programmer-Guide)
  - [Vogella Plugin Example](https://www.vogella.com/tutorials/EclipseJDT/article.html)
  - [Maven Repository](https://mvnrepository.com/artifact/org.eclipse.jdt/org.eclipse.jdt.core/3.40.0)

- [How to build a CLI app with Maven](https://www.sohamkamani.com/java/cli-app-with-maven/).
  Includes how to setup tests with junit.

- [SpotBugs for Maven to Lint Java](https://spotbugs.readthedocs.io/en/latest/maven.html)
