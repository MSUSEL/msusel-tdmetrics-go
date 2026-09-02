using Commons.Data.Repo;
using Constructs;
using System.Collections.Generic;
using System.Linq;
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
            where c.Type != GT.DeclType.Anonymous
            where c.Type != GT.DeclType.Interface
            select new KeyValuePair<string, GT.DeclMetrics>(c.FullName, c)
        );

        foreach (KeyValuePair<string, ObjectDecl> p in projObjects) {
            GT.DeclMetrics gtObj = gtClasses[p.Key];
            ObjectDecl projObj = p.Value;
            Assert.AreEqual(gtObj.FullName, projObj.FullName, "the object full names should match");

            List<GT.MethodMetrics> gtMethods = [..
                from m in gtObj.Methods
                where !m.Name.StartsWith("(initializer ")
                where !isPmdProblem(m)
                select m
            ];

            SortedSet<int> gtLines     = [.. from c in gtMethods select c.Line];
            SortedSet<int> projLines   = [.. from c in projObj.Methods select c.Location.LineNo];
            SortedSet<int> found       = [.. gtLines.Intersect(projLines)];
            SortedSet<int> gtMissing   = [.. gtLines.Except(found)];
            SortedSet<int> projMissing = [.. projLines.Except(found)];

            using (Assert.EnterMultipleScope()) {
                Assert.AreEqual(gtMethods.Count, gtLines.Count, "Duplicate method lines in ground truth class " + gtObj.FullName);
                Assert.AreEqual(projObj.Methods.Count, projLines.Count, "Duplicate method lines in object declaration " + projObj.FullName);
                Assert.Zero(gtMissing.Count, "Ground truth, " + gtObj.FullName + ", missing count:\n  " + string.Join("\n  ", gtMissing));
                Assert.Zero(projMissing.Count, "Object declaration, " + projObj.FullName + ", missing count:\n  " + string.Join("\n  ", projMissing));
            }
        }
    }

    static private bool isPmdProblem(GT.MethodMetrics m) {
        // PMD was the only one that read this method.
        if (!m.HasPmd || m.HasCk) return false;

        // PMD put these interface abstract in the wrong location.
        return m.FullName switch {
            "org.apache.bcel.util.ClassPath#getSize" => true,
            "org.apache.bcel.util.InstructionFinder#checkCode" => true,
            _ => false,
        };
    }

    [Test] public void GroundTruthArchiva()              => checkGroundTruth(JavaTarget.Archiva);
    [Test] public void GroundTruthBatik()                => checkGroundTruth(JavaTarget.Batik);
    [Test] public void GroundTruthCayenne()              => checkGroundTruth(JavaTarget.Cayenne);
    [Test] public void GroundTruthCocoon()               => checkGroundTruth(JavaTarget.Cocoon);
    [Test] public void GroundTruthCommonsBcel()          => checkGroundTruth(JavaTarget.CommonsBcel);
    [Test] public void GroundTruthCommonsBeanutils()     => checkGroundTruth(JavaTarget.CommonsBeanutils);
    [Test] public void GroundTruthCommonsCli()           => checkGroundTruth(JavaTarget.CommonsCli);
    [Test] public void GroundTruthCommonsCodec()         => checkGroundTruth(JavaTarget.CommonsCodec);
    [Test] public void GroundTruthCommonsCollections()   => checkGroundTruth(JavaTarget.CommonsCollections);
    [Test] public void GroundTruthCommonsConfiguration() => checkGroundTruth(JavaTarget.CommonsConfiguration);
    [Test] public void GroundTruthCommonsDaemon()        => checkGroundTruth(JavaTarget.CommonsDaemon);
    [Test] public void GroundTruthCommonsDbcp()          => checkGroundTruth(JavaTarget.CommonsDbcp);
    [Test] public void GroundTruthCommonsDbutils()       => checkGroundTruth(JavaTarget.CommonsDbutils);
    [Test] public void GroundTruthCommonsDigester()      => checkGroundTruth(JavaTarget.CommonsDigester);
    [Test] public void GroundTruthCommonsExec()          => checkGroundTruth(JavaTarget.CommonsExec);
    [Test] public void GroundTruthCommonsFileUpload()    => checkGroundTruth(JavaTarget.CommonsFileUpload);
    [Test] public void GroundTruthCommonsIo()            => checkGroundTruth(JavaTarget.CommonsIo);
    [Test] public void GroundTruthCommonsJelly()         => checkGroundTruth(JavaTarget.CommonsJelly);
    [Test] public void GroundTruthCommonsJexl()          => checkGroundTruth(JavaTarget.CommonsJexl);
    [Test] public void GroundTruthCommonsJxpath()        => checkGroundTruth(JavaTarget.CommonsJxpath);
    [Test] public void GroundTruthCommonsNet()           => checkGroundTruth(JavaTarget.CommonsNet);
    [Test] public void GroundTruthCommonsOgnl()          => checkGroundTruth(JavaTarget.CommonsOgnl);
    [Test] public void GroundTruthCommonsValidator()     => checkGroundTruth(JavaTarget.CommonsValidator);
    [Test] public void GroundTruthCommonsVfs()           => checkGroundTruth(JavaTarget.CommonsVfs);
    [Test] public void GroundTruthFelix()                => checkGroundTruth(JavaTarget.Felix);
    [Test] public void GroundTruthHive()                 => checkGroundTruth(JavaTarget.Hive);
    [Test] public void GroundTruthHttpComponentsClient() => checkGroundTruth(JavaTarget.HttpComponentsClient);
    [Test] public void GroundTruthHttpComponentsCore()   => checkGroundTruth(JavaTarget.HttpComponentsCore);
    [Test] public void GroundTruthSantuario()            => checkGroundTruth(JavaTarget.Santuario);
    [Test] public void GroundTruthThrift()               => checkGroundTruth(JavaTarget.Thrift);
    [Test] public void GroundTruthZookeeper()            => checkGroundTruth(JavaTarget.Zookeeper);

    static private void checkGroundTruth(JavaTarget target) {
        GT.GroundTruth gt = GT.GroundTruth.FromZip(Repo.MetricsZip, target);
        Project proj;
        try {
            proj = Project.FromFile(Repo.AbstractedJava(target));
        } catch (System.IO.FileNotFoundException) {
            Assert.Fail("Missing the abstraction file: " + target.ProjectKey + ".json");
            return;
        }

        string groupId = target.GroupId;
        //string groupId = proj.GroupId;
        //Assert.IsNotEmpty(groupId, "The groupId needs to not be empty");
        //Assert.True(groupId.StartsWith("org.apache"), "Expected the groupId to start with \"org.apache\" but it was \"" + groupId + "\"");
        
        Assert.AreEqual(proj.CommitHash, target.CommitSha, "Commit hash should match for " + groupId);

        Dictionary<string, ObjectDecl> projObjects = new(
            from c in proj.ObjectDecls
            where c.Package.Name.StartsWith(groupId)
            select new KeyValuePair<string, ObjectDecl>(c.FullName, c)
        );
        Assert.IsNotEmpty(projObjects, "Must check at least one class for " + groupId + "\n" +
            "   NOTICE: If projObjects is empty, check that the groupId is the the root package name "+
            "(e.g. \"commons-codec\" needs to be \"org.apache.commons.codec\")");

        Dictionary<string, GT.DeclMetrics> gtClasses = new(
            from c in gt.Declarations
            where c.FullName.StartsWith(groupId)
            where !c.InTestPath
            where c.Type != GT.DeclType.Anonymous
            where c.Type != GT.DeclType.Interface
            select new KeyValuePair<string, GT.DeclMetrics>(c.FullName, c)
        );
        if (projObjects.Count != gtClasses.Count) {
            List<string> missing = [..
                from key in gtClasses.Keys
                where !projObjects.ContainsKey(key)
                orderby key
                select key
            ];
            string missingMsg = missing.Count <= 0 ? "" :
                "\n  Missing (in ground truth but not in objects):\n    " + string.Join("\n    ", missing);

            List<string> extra = [..
                from key in projObjects.Keys
                where !gtClasses.ContainsKey(key)
                orderby key
                select key
            ];
            string extraMsg = extra.Count <= 0 ? "" :
                "\n  Extra (in objects but not in ground truth):\n    " + string.Join("\n    ", extra);

            Assert.AreEqual(projObjects.Count, gtClasses.Count,
                "The number of classes are expected to match" +
                missingMsg + extraMsg);
        }

        int methods = 0;
        foreach (KeyValuePair<string, ObjectDecl> p in projObjects) {
            methods += checkGroundTruth(gtClasses[p.Key], p.Value);
        }
        Assert.NotZero(methods, "Must check at least one method");
    }

    static private int checkGroundTruth(GT.DeclMetrics gtObj, ObjectDecl projObj) {
        Assert.AreEqual(gtObj.FullName, projObj.FullName, "The objects should match");

        Dictionary<int, GT.MethodMetrics> gtMetByLine = new(
            from m in gtObj.Methods
            where !m.Name.StartsWith("(initializer ")
            where !isPmdProblem(m)
            select new KeyValuePair<int, GT.MethodMetrics>(m.Line, m)
        );
        SortedSet<int> projLines = [.. from m in projObj.Methods select m.Location.LineNo];
        List<string> missing = [..
            from m in gtMetByLine.Values
            where !projLines.Contains(m.Line)
            orderby m.Line
            select m.Line + ": " + m.Name
        ];
        List<string> extra = [..
            from m in projObj.Methods
            where !gtMetByLine.ContainsKey(m.Location.LineNo)
            orderby m.Location.LineNo
            select m.Location.LineNo + ": " + m.Name
        ];
        List<string> gtMethodLines = [..
            from m in gtObj.Methods
            where !m.Name.StartsWith("(initializer ")
            orderby m.Line
            select m.Line + ": " + m.Name
        ];
        List<string> projMethodLines = [..
            from m in projObj.Methods
            orderby m.Location.LineNo
            select m.Location.LineNo + ": " + m.Name
        ];

        string missingMsg = missing.Count <= 0 ? "" :
            "\n  Missing (in ground truth but not in objects):\n    " + string.Join("\n    ", missing);
        string extraMsg = extra.Count <= 0 ? "" :
            "\n  Extra (in objects but not in ground truth):\n    " + string.Join("\n    ", extra);
        string getMethodsMsg = gtMethodLines.Count <= 0 ? "":
            "\n  Methods in ground truth were:\n    " + string.Join("\n    ", gtMethodLines);
        string projMethodsMsg = projMethodLines.Count <= 0 ? "":
            "\n  Methods in objects were:\n    " + string.Join("\n    ", projMethodLines);

        if ((missing.Count > 0) || (extra.Count > 0)) {
            Assert.AreEqual(projObj.Methods.Count, gtMetByLine.Count,
                "The number of methods in " + projObj.FullName + " are expected to match" + 
                missingMsg + extraMsg + getMethodsMsg + projMethodsMsg);
        }

        Assert.Zero(missing.Count, "missing.Count");
        Assert.Zero(extra.Count, "extra.Count");

        int methods = 0;
        foreach (MethodDecl projMet in projObj.Methods) {
            if (gtMetByLine.TryGetValue(projMet.Location.LineNo, out GT.MethodMetrics? gtMet)) {
                Assert.AreEqual(gtMet.Cyclo, projMet.Metrics?.PmdCyclo ?? 1, "Cyclomatic for " + gtMet.FullName + ":" + projMet);
                methods++;
                continue;
            }
            Assert.Fail("Failed to find " + projMet.Location.LineNo + ", for " + projMet + ", in ground truth methods by line number" + 
                missingMsg + extraMsg + getMethodsMsg + projMethodsMsg);
        }
        Assert.AreEqual(gtObj.CycloTotal, projObj.PmdWmc, "WMC for " + projObj.FullName + " @ "  + projObj.Location);
        return methods;
    }
}
