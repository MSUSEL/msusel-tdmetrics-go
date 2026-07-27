using Commons.Data.Repo;
using Constructs;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using GT = GroundTruth;

namespace Runner.Commands;

/// <summary>
/// Loads the ground-truth data and the abstracted project for a target and
/// writes a CSV comparing per-object WMC values from each source. Uses the
/// same class-filtering rules as <c>GroundTruthTests.checkGroundTruth</c>.
/// </summary>
public sealed class GtCheckCommand : ICommand {

    public string Name => "gtcheck";

    public string Description =>
        "gtcheck <projectKey> [outputCsvPath]  " +
        "Writes GodWmc/PmdWmc for each matched object to a CSV. " +
        "Default output: <projectKey>-gtcheck.csv in the current directory.";

    public void Run(string[] args) {
        if (args.Length < 1) {
            Console.Error.WriteLine("gtcheck: missing <projectKey> argument.");
            Console.Error.WriteLine("usage: " + this.Description);
            Environment.Exit(2);
            return;
        }

        string projectKey = args[0];
        string outputPath = args.Length >= 2
            ? args[1]
            : projectKey + "-gtcheck.csv";

        JavaTarget? maybeTarget = findTarget(projectKey);
        if (maybeTarget is null) {
            Console.Error.WriteLine("gtcheck: unknown project key '" + projectKey + "'.");
            Console.Error.WriteLine("Known keys: " +
                string.Join(", ", from t in JavaTarget.Targets select t.ProjectKey));
            Environment.Exit(2);
            return;
        }
        JavaTarget target = maybeTarget.Value;

        GT.GroundTruth gt = GT.GroundTruth.FromZip(Repo.MetricsZip, target);
        Project proj = Project.FromFile(Repo.AbstractedJava(target));
        string groupId = proj.GroupId;

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

        int written = 0;
        int skipped = 0;
        using (StreamWriter writer = new(outputPath, false, Encoding.UTF8)) {
            writer.WriteLine("god_wmc, cyclo_total, pmd_wmc, full_name, location");
            foreach (KeyValuePair<string, ObjectDecl> p in projObjects.OrderBy(kv => kv.Key)) {
                ObjectDecl projObj = p.Value;
                if (!gtClasses.TryGetValue(p.Key, out GT.DeclMetrics? gtObj) || gtObj is null) {
                    Console.Error.WriteLine("gtcheck: no ground truth for '" + p.Key + "', skipping.");
                    skipped++;
                    continue;
                }
                writer.WriteLine(
                    gtObj.GodWmc.ToString() + ", " +
                    gtObj.CycloTotal.ToString() + ", " +
                    projObj.PmdWmc.ToString() + ", " +
                    csvField(projObj.FullName) + ", " +
                    csvField(projObj.Location.ToString()));
                written++;
            }
        }

        Console.WriteLine("gtcheck: wrote " + written + " rows to " + outputPath +
            (skipped > 0 ? " (" + skipped + " skipped)" : ""));
    }

    static private JavaTarget? findTarget(string projectKey) {
        foreach (JavaTarget t in JavaTarget.Targets)
            if (t.ProjectKey == projectKey) return t;
        return null;
    }

    static private string csvField(string value) {
        if (value.IndexOfAny(['"', ',', '\r', '\n']) < 0) return value;
        return "\"" + value.Replace("\"", "\"\"") + "\"";
    }
}
