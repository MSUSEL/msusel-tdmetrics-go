using Constructs;
using System;
using System.IO;

namespace UnitTests.ConstructTests;

public class ConstructTests {

    [Test]
    public void StubTest0001() => runStubTest(testPath("go", 1));

    [Test]
    public void StubTest0002() => runStubTest(testPath("go", 2));

    [Test]
    public void StubTest0003() => runStubTest(testPath("go", 3));

    [Test]
    public void StubTest0004() => runStubTest(testPath("go", 4));

    [Test]
    public void StubTest0005() => runStubTest(testPath("go", 5));

    // Skip /go/test0006, or make a stub test able to pick one package.

    [Test]
    public void StubTest0007() => runStubTest(testPath("go", 7));

    [Test]
    public void StubTest0008() => runStubTest(testPath("go", 8));

    [Test]
    public void StubTest0009() => runStubTest(testPath("go", 9));

    [Test]
    public void StubTest0010() => runStubTest(testPath("go", 10));

    [Test]
    public void StubTest0011() => runStubTest(testPath("go", 11));

    [Test]
    public void StubTest0012() => runStubTest(testPath("go", 12));

    [Test]
    public void StubTest0013() => runStubTest(testPath("go", 13));

    [Test]
    public void StubTest0014() => runStubTest(testPath("go", 14));
    
    [Test]
    public void StubTest0015() => runStubTest(testPath("go", 15));
    
    [Test]
    public void StubTest0016() => runStubTest(testPath("go", 16));

    #region Test Tools...

    static private readonly string repoDir;

    static private readonly string testDataDir;

    static ConstructTests() {
        const string repoName = "msusel-tdmetrics-go";
        string curDir = Environment.CurrentDirectory;
        int index = curDir.LastIndexOf(repoName);
        if (index == -1) throw new Exception("Failed to find root directory of the repo from " + curDir);
        index += repoName.Length;
        repoDir = curDir[0..index];
        testDataDir = repoDir + "/testData";
    }

    static private string testPath(string sourceLang, int testNum) =>
        string.Format("{0}/{1}/test{2:D4}", testDataDir, sourceLang, testNum);

    static private Project readTestPackage(string testPath, string fileName = "abstraction.yaml") =>
        Project.FromFile(testPath + "/" + fileName);

    static private string readExpectedStub(string testPath, string fileName = "expStub.txt") =>
        File.ReadAllText(testPath + "/" + fileName).Trim();

    static private void runStubTest(string testPath) {
        Project proj = readTestPackage(testPath);
        string got = proj.ToString();
        string exp = readExpectedStub(testPath).ReplaceLineEndings("\n");
        if (got != exp) {
            Console.WriteLine(got);
            Assert.That(got, Is.EqualTo(exp));
        }
    }

    #endregion
}