using GT = GroundTruth;
using Commons.Data.Repo;
using Constructs;
using System.Linq;

namespace UnitTests.GroundTruthTests;

public class GroundTruthTests {

    [Test]
    public void GroundTruthReadZip() {
        GT.GroundTruth gt = GT.GroundTruth.FromZip(Repo.MetricsZip, JavaTarget.CommonsBcel);
        Assert.AreEqual(gt.ProjectKey, "commons-bcel");
        Assert.AreEqual(gt.ProjectId,  "org.apache:bcel");
        Assert.AreEqual(gt.GitLink,    "apache/commons-bcel");
        Assert.AreEqual(gt.CommitSha,  "6ed18c5bef0f5b93b54783a8e8fb2b9042da26ac");
    }

    [Test]
    public void GroundTruthCommonBcel() =>
        this.checkGroundTruth(JavaTarget.CommonsBcel);

    private void checkGroundTruth(JavaTarget target) {
        GT.GroundTruth gt = GT.GroundTruth.FromZip(Repo.MetricsZip, target);
        Project proj = Project.FromFile(Repo.AbstractedJava(target));
        foreach (GT.ClassMetrics gtObj in gt.Classes) {
            ObjectDecl? obj = proj.ObjectDecls.FirstOrDefault(c => c.FullName == gtObj.FullName);
            if (obj is null) {
                System.Console.WriteLine("Failed to find class/object " + gtObj.FullName);
                continue;
            }
            this.checkGroundTruth(gtObj, obj);
        }
    }

    private void checkGroundTruth(GT.ClassMetrics gtObj, ObjectDecl obj) {
        Assert.Multiple(() => {
            foreach (GT.MethodMetrics gtMet in gtObj.Methods) {
                MethodDecl? met = obj.Methods.FirstOrDefault(m => m.Location.LineNo == gtMet.Line);
                if (met is null) {
                    System.Console.WriteLine("Failed to find method " + gtMet.Name + " in " + gtObj.FullName);
                    continue;
                }
                this.checkGroundTruth(gtObj, obj, gtMet, met);
            }
        });
    }

    private void checkGroundTruth(GT.ClassMetrics gtObj, ObjectDecl obj, GT.MethodMetrics gtMet, MethodDecl met) {
        // TODO: Add more.
        Assert.AreEqual(gtMet.Loc,   met.Metrics?.LineCount ?? 0,  "Lines of code for " + gtMet.Name + " in " + gtObj.FullName);
        Assert.AreEqual(gtMet.Cyclo, met.Metrics?.Complexity ?? 0, "Complexity for " + gtMet.Name + " in " + gtObj.FullName);
    }
}
