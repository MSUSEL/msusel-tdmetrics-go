# Java Abstractor

- [Java Abstractor](#java-abstractor)
  - [Resources](#resources)
  - [Setup Abstractor](#setup-abstractor)
  - [Running Abstractor](#running-abstractor)

## Resources

- [Eclipse JDT API](https://www.vogella.com/tutorials/EclipseJDT/article.html)

## Setup Abstractor

1. Install an OpenJDK. This project was built using
   [RedHat OpenJDK v1.8.0](https://developers.redhat.com/products/openjdk/download)

2. Install Maven. This project was built using
   [Apache Maven v3.9.9](https://maven.apache.org/download.cgi)

3. Check the install by running the Maven version check `mvn -v`.

   The result should look similar to following except with your own
   install paths:

    ```Plain
    msusel-tdmetrics-go\javaAbstractor> mvn -v
    Apache Maven 3.9.9 (8e8579a9e76f7d015ee5ec7bfcdc97d260186937)
    Maven home: C:\Program Files\Apache\Maven\3.9.9
    Java version: 1.8.0_392, vendor: Red Hat, Inc., runtime: C:\Program Files (x86)\RedHat\OpenJDK\1.8.0\jre
    Default locale: en_US, platform encoding: Cp1252
    OS name: "windows 11", version: "10.0", arch: "amd64", family: "windows"
    ```

4. From the same folder as this README.md, run `mvn compile`.

   This will finish installing by downloading
   all the needed dependencies for this project.

## Running Abstractor
