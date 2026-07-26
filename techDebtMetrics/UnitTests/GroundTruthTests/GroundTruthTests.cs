using Commons.Data.Repo;
using Constructs;
using System.Collections;
using System.Collections.Generic;
using System.Linq;
using System.Text.RegularExpressions;
using GT = GroundTruth;

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
        string groupId = proj.GroupId;
        Assert.AreEqual(proj.CommitHash, target.CommitSha, "commit hash should match for " + groupId);

        List<ObjectDecl> projObjects = [..
            from c in proj.ObjectDecls
            where c.Package.Name.StartsWith(groupId)
            select c    
        ];

        List<GT.DeclMetrics> gtClasses = [..
            from c in gt.Declarations
            where c.FullName.StartsWith(groupId)
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

        using (Assert.EnterMultipleScope()) {
            Assert.AreEqual(gtClasses.Count,   gtNames.Count,   "Duplicate class names in ground truth data");
            Assert.AreEqual(projObjects.Count, projNames.Count, "Duplicate class names in object declarations");
            Assert.Zero(gtMissing.Count,   "Ground truth missing count:\n  " + string.Join("\n  ", gtMissing));
            Assert.Zero(projMissing.Count, "Object declarations missing count:\n  " + string.Join("\n  ", projMissing));
        }
    }
    
    [Test]
    public void GroundTruthMethodMatching() {
        JavaTarget target = JavaTarget.CommonsBcel;
        GT.GroundTruth gt = GT.GroundTruth.FromZip(Repo.MetricsZip, target);
        Project proj = Project.FromFile(Repo.AbstractedJava(target));
        string groupId = proj.GroupId;
        Assert.AreEqual(proj.CommitHash, target.CommitSha, "commit hash should match for " + groupId);

        Dictionary<string, ObjectDecl> projObjects = new(
            from c in proj.ObjectDecls
            where c.Package.Name.StartsWith(groupId)
            select new KeyValuePair<string, ObjectDecl>(c.FullName, c)
        );

        Dictionary<string, GT.DeclMetrics> gtClasses = new(
            from c in gt.Declarations
            where c.FullName.StartsWith(groupId)
            where !c.InTestPath
            where c.Type != GT.DeclType.Anonymous // TODO: Probably need to fold this into the nest to be counted instead of skipping it.
            where c.Type != GT.DeclType.Interface
            select new KeyValuePair<string, GT.DeclMetrics>(c.FullName, c)
        );

        foreach (KeyValuePair<string, ObjectDecl> p in projObjects) {
            GT.DeclMetrics gtObj = gtClasses[p.Key];
            ObjectDecl projObj = p.Value;
            Assert.AreEqual(gtObj.FullName, projObj.FullName, "the object full names should match");

            string debugTarget = "org.apache.bcel.generic.PUSH"; // TODO: REMOVE!!
            if (projObj.FullName == debugTarget) { // TODO: REMOVE!!
                System.Console.WriteLine("====================================");
                System.Console.WriteLine(debugTarget);
                System.Console.WriteLine("---[ gt ]---");
                foreach (GT.MethodMetrics m in gtObj.Methods) {
                    string ckStr = m.HasCk ? " [CK: " + m.CkLine + "]" : "";
                    string pmdStr = m.HasPmd ? " [PMD: " + m.PmdLine + "]" : "";
                    System.Console.WriteLine(m.Name + ckStr + pmdStr);
                }
                System.Console.WriteLine("---[ proj ]---");
                foreach (MethodDecl c in projObj.Methods) {
                    System.Console.WriteLine(c.Name + " [proj: " + c.Location.LineNo + "]");
                }
            }

            List<GT.MethodMetrics> gtMethods = [..
                from m in gtObj.Methods
                where !m.Name.StartsWith("(initializer ")
                where m.FullName != "org.apache.bcel.util.ClassPath#getSize" // PMD was the only one that read this, it seems wrong.
                where m.FullName != "org.apache.bcel.util.InstructionFinder#checkCode" // PMD put this interface abstract in the wrong location too.
                //where !m.Modifiers.Abstract
                select m
            ];

            SortedSet<int> gtLines     = [.. from c in gtMethods select c.Line];
            SortedSet<int> projLines   = [.. from c in projObj.Methods select c.Location.LineNo];
            SortedSet<int> found       = [.. gtLines.Intersect(projLines)];
            SortedSet<int> gtMissing   = [.. gtLines.Except(found)];
            SortedSet<int> projMissing = [.. projLines.Except(found)];

            if (projObj.FullName == debugTarget) { // TODO: REMOVE!!
                System.Console.WriteLine("gtLines:     [" + string.Join(", ", gtLines) + "]"); // TODO: REMOVE
                System.Console.WriteLine("projLines:   [" + string.Join(", ", projLines) + "]"); // TODO: REMOVE
                System.Console.WriteLine("found:       [" + string.Join(", ", found) + "]"); // TODO: REMOVE
                System.Console.WriteLine("gtMissing:   [" + string.Join(", ", gtMissing) + "]"); // TODO: REMOVE
                System.Console.WriteLine("projMissing: [" + string.Join(", ", projMissing) + "]"); // TODO: REMOVE
            }

            using (Assert.EnterMultipleScope()) {
                Assert.AreEqual(gtMethods.Count, gtLines.Count, "Duplicate method lines in ground truth class " + gtObj.FullName);
                Assert.AreEqual(projObj.Methods.Count, projLines.Count, "Duplicate method lines in object declaration " + projObj.FullName);
                Assert.Zero(gtMissing.Count, "Ground truth, " + gtObj.FullName + ", missing count:\n  " + string.Join("\n  ", gtMissing));
                Assert.Zero(projMissing.Count, "Object declaration, " + projObj.FullName + ", missing count:\n  " + string.Join("\n  ", projMissing));
            }
        }
    }

    [Test]
    public void GroundTruthCommonBcel() =>
        this.checkGroundTruth(JavaTarget.CommonsBcel);

    private void checkGroundTruth(JavaTarget target) {
        GT.GroundTruth gt = GT.GroundTruth.FromZip(Repo.MetricsZip, target);
        Project proj = Project.FromFile(Repo.AbstractedJava(target));
        string groupId = proj.GroupId;
        Assert.AreEqual(proj.CommitHash, target.CommitSha, "commit hash should match for " + groupId);

        Dictionary<string, ObjectDecl> projObjects = new(
            from c in proj.ObjectDecls
            where c.Package.Name.StartsWith(groupId)
            select new KeyValuePair<string, ObjectDecl>(c.FullName, c)
        );

        Dictionary<string, GT.DeclMetrics> gtClasses = new(
            from c in gt.Declarations
            where c.FullName.StartsWith(groupId)
            where !c.InTestPath
            where c.Type != GT.DeclType.Anonymous // TODO: Probably need to fold this into the nest to be counted instead of skipping it.
            where c.Type != GT.DeclType.Interface
            select new KeyValuePair<string, GT.DeclMetrics>(c.FullName, c)
        );

        foreach (KeyValuePair<string, ObjectDecl> p in projObjects)
            this.checkGroundTruth(gtClasses[p.Key], p.Value);
    }

    private void checkGroundTruth(GT.DeclMetrics gtObj, ObjectDecl obj) {
        Assert.AreEqual(gtObj.FullName, obj.FullName, "the objects should match");
        //Assert.AreEqual(gtObj.Wmc, obj.Wmc, "WMC for " + obj.FullName);
        
        // TODO: FIX

        //foreach (GT.MethodMetrics gtMet in gtObj.Methods) {
        //    MethodDecl? met = obj.Methods.FirstOrDefault(m => m.Location.LineNo == gtMet.Line);
        //    if (met is null) {
        //        System.Console.WriteLine("Failed to find method " + gtMet.Name + " in " + gtObj.FullName);
        //        continue;
        //    }
        //    this.checkGroundTruth(gtObj, obj, gtMet, met);
        //}
    }

    //private void checkGroundTruth(GT.DeclMetrics gtObj, ObjectDecl obj, GT.MethodMetrics gtMet, MethodDecl met) {
    //    // TODO: Add more.
    //    Assert.AreEqual(gtMet.Loc,   met.Metrics?.LineCount ?? 0,  "Lines of code for " + gtMet.Name + " in " + gtObj.FullName);
    //    Assert.AreEqual(gtMet.Cyclo, met.Metrics?.Complexity ?? 0, "Complexity for " + gtMet.Name + " in " + gtObj.FullName);
    //}
}
