using System.IO.Compression;
using Yaml = Commons.Data.Yaml;

namespace GroundTruth;

public class GroundTruth {
    static public GroundTruth FromZip(string zipfile, Target target) =>
        FromZip(zipfile, target.ProjectKey + ".json");

    static public GroundTruth FromZip(string zipfile, string targetFile) {
        using ZipArchive zip = ZipFile.Open(zipfile, ZipArchiveMode.Read);
        foreach (ZipArchiveEntry entry in zip.Entries) {
            if (entry.Name == targetFile)
                return FromJson(entry.Open());
        }
        throw new Exception("Failed to find target: " + targetFile);
    }

    static public GroundTruth FromJson(string jsonFile) {
        using StreamReader jsonStream = new(jsonFile);
        return new(Yaml.Node.Parse(jsonStream));
    }

    static public GroundTruth FromJson(Stream jsonStream) => new(Yaml.Node.Parse(jsonStream));

    public GroundTruth(Yaml.Node node) {
        this.Root = node.AsObject();

        List<ClassMetrics> classes = [];
        Yaml.Array array = this.Root.ReadNode("classes").AsArray();
        foreach (Yaml.Node item in array.Items)
            classes.Add(new(item.AsObject()));
        this.Classes = classes.AsReadOnly();
    }

    public readonly Yaml.Object Root;
    public readonly IReadOnlyList<ClassMetrics> Classes;

    public string ProjectKey => this.Root.ReadString("project_key");
    public string ProjectId  => this.Root.ReadString("project_id");
    public string GitLink    => this.Root.ReadString("git_link");
    public string CommitSha  => this.Root.ReadString("commit_sha");

    public Measures Measures => new(this.Root.ReadNode("measures").AsObject());
}
