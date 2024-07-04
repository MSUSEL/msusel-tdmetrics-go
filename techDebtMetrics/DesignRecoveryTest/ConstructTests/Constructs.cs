using DesignRecovery.Constructs;

namespace DesignRecoveryTest.Constructs;

public class ConstructTests {

    [Test]
    public void StubTest0001() => runStubTest(1);

    #region Test Tools...

    static private string getTestPath(int testNum, string fileName) =>
        string.Format("../TestData/Test{0:D6}/{1}", testNum, fileName);

    static private Project readTestPackage(int testNum, string fileName = "absraction.yaml") =>
        Project.FromJsonFile(getTestPath(testNum, fileName));

    static private string readExpectedStub(int testNum, string fileName = "expStub.txt") =>
        File.ReadAllText(getTestPath(testNum, fileName));

    static private void runStubTest(int testNum) {
        Project proj = readTestPackage(testNum);
        string got = proj.ToStub();
        string exp = readExpectedStub(testNum);
        Assert.That(got, Is.EqualTo(exp));
    }

    #endregion
}