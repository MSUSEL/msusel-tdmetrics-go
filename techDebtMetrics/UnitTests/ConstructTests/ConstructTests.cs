using Constructs;
using System;
using System.IO;

namespace UnitTests.ConstructTests;

public class ConstructTests {

    [Test]
    public void StubTest0001() => runStubTest(1);

    [Test]
    public void StubTest0002() => runStubTest(2);

    [Test]
    public void StubTest0003() => runStubTest(3);

    [Test]
    public void StubTest0004() => runStubTest(4);

    [Test]
    public void StubTest0005() => runStubTest(5);

    #region Test Tools...

    static private string testDataDir;

    static ConstructTests() {
        string curDir = Environment.CurrentDirectory;
        int index = curDir.LastIndexOf("UnitTests");
        if (index == -1) throw new Exception("Failed to find test data folder from " + curDir);
        testDataDir = curDir[0..index]+"UnitTests/TestData";
    }

    static private string getTestPath(int testNum, string fileName) =>
        string.Format("{0}/Test{1:D4}/{2}", testDataDir, testNum, fileName);

    static private Project readTestPackage(int testNum, string fileName = "abstraction.yaml") =>
        Project.FromFile(getTestPath(testNum, fileName));

    static private string readExpectedStub(int testNum, string fileName = "expStub.txt") =>
        File.ReadAllText(getTestPath(testNum, fileName)).Trim();

    static private void runStubTest(int testNum) {
        Project proj = readTestPackage(testNum);
        string got = proj.ToString();
        string exp = readExpectedStub(testNum).ReplaceLineEndings("\n");
        if (got != exp) {
            Console.WriteLine(got);
            Assert.That(got, Is.EqualTo(exp));
        }
    }

    #endregion
}