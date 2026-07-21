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
    public void GroundTruthClassMatching() {
        JavaTarget target = JavaTarget.CommonsBcel;
        GT.GroundTruth gt = GT.GroundTruth.FromZip(Repo.MetricsZip, target);
        Project proj = Project.FromFile(Repo.AbstractedJava(target));

        Assert.AreEqual(proj.CommitHash, target.CommitSha, "commit hash should match");

        List<ObjectDecl> projObjects = [..
            from c in proj.ObjectDecls
            where c.Package.Name.StartsWith("org.apache.bcel")
            select c    
        ];

        List<GT.DeclMetrics> gtClasses = [..
            from c in gt.Declarations
            where c.FullName.StartsWith("org.apache.bcel")
            where !c.InTestPath
            where c.Type != GT.DeclType.Anonymous // TODO: Probably need to fold this into the nest to be counted instead of skipping it.
            where c.Type != GT.DeclType.Interface
            select c
        ];

        SortedSet<string> gtNames     = [.. from c in gtClasses select c.FullName];
        SortedSet<string> projNames   = [.. from c in projObjects select c.FullName];
        SortedSet<string> found       = [.. gtNames.Intersect(projNames)];
        SortedSet<string> gtMissing   = [.. gtNames.Except(found)];
        SortedSet<string> projMissing = [.. projNames.Except(found)];

        Assert.Multiple(() => {
            Assert.AreEqual(gtClasses.Count,   gtNames.Count,   "Duplicate class names in ground truth data");
            Assert.AreEqual(projObjects.Count, projNames.Count, "Duplicate class names in object declarations");
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
        foreach (GT.DeclMetrics gtObj in gt.Declarations) {
            ObjectDecl? obj = proj.ObjectDecls.FirstOrDefault(c => c.FullName == gtObj.FullName);
            if (obj is null) {
                System.Console.WriteLine("Failed to find class/object " + gtObj.FullName);
                continue;
            }
            this.checkGroundTruth(gtObj, obj);
        }
    }

    private void checkGroundTruth(GT.DeclMetrics gtObj, ObjectDecl obj) {
        foreach (GT.MethodMetrics gtMet in gtObj.Methods) {
            MethodDecl? met = obj.Methods.FirstOrDefault(m => m.Location.LineNo == gtMet.Line);
            if (met is null) {
                System.Console.WriteLine("Failed to find method " + gtMet.Name + " in " + gtObj.FullName);
                continue;
            }
            this.checkGroundTruth(gtObj, obj, gtMet, met);
        }
    }

    private void checkGroundTruth(GT.DeclMetrics gtObj, ObjectDecl obj, GT.MethodMetrics gtMet, MethodDecl met) {
        // TODO: Add more.
        Assert.AreEqual(gtMet.Loc,   met.Metrics?.LineCount ?? 0,  "Lines of code for " + gtMet.Name + " in " + gtObj.FullName);
        Assert.AreEqual(gtMet.Cyclo, met.Metrics?.Complexity ?? 0, "Complexity for " + gtMet.Name + " in " + gtObj.FullName);
    }
}
