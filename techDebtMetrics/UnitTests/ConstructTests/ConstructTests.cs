using Constructs;
using System;
using System.IO;
using Commons.Data.Repo;

namespace UnitTests.ConstructTests;

public class ConstructTests {

    [Test] public void StubTest0001() => runStubTest("go", 1);
    [Test] public void StubTest0002() => runStubTest("go", 2);
    [Test] public void StubTest0003() => runStubTest("go", 3);
    [Test] public void StubTest0004() => runStubTest("go", 4);
    [Test] public void StubTest0005() => runStubTest("go", 5);
    // Skip /go/test0006, or make a stub test able to pick one package.
    [Test] public void StubTest0007() => runStubTest("go", 7);
    [Test] public void StubTest0008() => runStubTest("go", 8);
    [Test] public void StubTest0009() => runStubTest("go", 9);
    [Test] public void StubTest0010() => runStubTest("go", 10);
    [Test] public void StubTest0011() => runStubTest("go", 11);
    [Test] public void StubTest0012() => runStubTest("go", 12);
    [Test] public void StubTest0013() => runStubTest("go", 13);
    [Test] public void StubTest0014() => runStubTest("go", 14);
    [Test] public void StubTest0015() => runStubTest("go", 15);
    [Test] public void StubTest0016() => runStubTest("go", 16);
    [Test] public void StubTest0017() => runStubTest("go", 17);
    [Test] public void StubTest0018() => runStubTest("go", 18);

    #region Test Tools...

    static private Project readTestPackage(string testPath, string fileName = "abstraction.yaml") =>
        Project.FromFile(testPath + "/" + fileName);

    static private string readExpectedStub(string testPath, string fileName = "expStub.txt") =>
        File.ReadAllText(testPath + "/" + fileName).Trim();

    static private void runStubTest(string sourceLang, int testNum) =>
        runStubTest(Repo.TestPath(sourceLang, testNum));

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