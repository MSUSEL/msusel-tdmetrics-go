using System;
using System.IO;
using Constructs;

namespace UnitTests.Constructs;

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

    static private string getTestPath(int testNum, string fileName) =>
        string.Format("../../../TestData/Test{0:D4}/{1}", testNum, fileName);

    static private Project readTestPackage(int testNum, string fileName = "abstraction.yaml") =>
        Project.FromFile(getTestPath(testNum, fileName));

    static private string readExpectedStub(int testNum, string fileName = "expStub.txt") =>
        File.ReadAllText(getTestPath(testNum, fileName)).Trim();

    static private void runStubTest(int testNum) {
        Project proj = readTestPackage(testNum);
        string got = proj.ToStub();
        string exp = readExpectedStub(testNum).ReplaceLineEndings("\n");
        if (got != exp) {
            Console.WriteLine(got);
            Assert.That(got, Is.EqualTo(exp));
        }
    }

    #endregion
}