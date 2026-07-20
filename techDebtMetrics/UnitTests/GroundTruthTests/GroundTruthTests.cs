using GT = GroundTruth;
using Commons.Data.Repo;
using Constructs;
using System.Linq;
using System.Collections.Generic;

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
    public void GroundTruthClasses() {
        JavaTarget target = JavaTarget.CommonsBcel;
        GT.GroundTruth gt = GT.GroundTruth.FromZip(Repo.MetricsZip, target);
        Project proj = Project.FromFile(Repo.AbstractedJava(target));
        SortedSet<string> gtNames     = [.. gt.Classes.Select(c => c.FullName)];
        SortedSet<string> projNames   = [.. proj.ObjectDecls.Select(c => c.FullName)];
        SortedSet<string> found       = [.. gtNames.Intersect(projNames)];
        SortedSet<string> gtMissing   = [.. gtNames.Except(found)];
        SortedSet<string> projMissing = [.. projNames.Except(found)];

        Assert.Multiple(() => {
            Assert.AreEqual(gtNames.Count,   gt.Classes.Count,       "Duplicate class names in ground truth data");
            Assert.AreEqual(projNames.Count, proj.ObjectDecls.Count, "Duplicate class names in object declarations");
            Assert.Zero(gtMissing.Count,   "Ground truth missing count:\n  " + string.Join("\n  ", gtMissing));
            Assert.Zero(projMissing.Count, "Object declarations missing count:\n  " + string.Join("\n  ", projMissing));
        });
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
        foreach (GT.MethodMetrics gtMet in gtObj.Methods) {
            MethodDecl? met = obj.Methods.FirstOrDefault(m => m.Location.LineNo == gtMet.Line);
            if (met is null) {
                System.Console.WriteLine("Failed to find method " + gtMet.Name + " in " + gtObj.FullName);
                continue;
            }
            this.checkGroundTruth(gtObj, obj, gtMet, met);
        }
    }

    private void checkGroundTruth(GT.ClassMetrics gtObj, ObjectDecl obj, GT.MethodMetrics gtMet, MethodDecl met) {
        // TODO: Add more.
        Assert.AreEqual(gtMet.Loc,   met.Metrics?.LineCount ?? 0,  "Lines of code for " + gtMet.Name + " in " + gtObj.FullName);
        Assert.AreEqual(gtMet.Cyclo, met.Metrics?.Complexity ?? 0, "Complexity for " + gtMet.Name + " in " + gtObj.FullName);
    }
}
