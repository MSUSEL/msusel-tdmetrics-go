using System.IO.Compression;
using Yaml = Commons.Data.Yaml;
using Commons.Data.Repo;

namespace GroundTruth;

public class GroundTruth(Yaml.Node node) {
    static public GroundTruth FromZip(string zipfile, JavaTarget target) =>
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

    public readonly Yaml.Object Root = node.AsObject();

    public string ProjectKey => this.Root.TryReadString("project_key");
    public string ProjectId  => this.Root.TryReadString("project_id");
    public string GitLink    => this.Root.TryReadString("git_link");
    public string CommitSha  => this.Root.TryReadString("commit_sha");
    public Measures Measures => new(this.Root.TryReadNode("measures")?.AsObject() ?? new Yaml.Object());

    public List<DeclMetrics> Declarations {
        get {
            if (field is not null) return field;
            List<DeclMetrics> list = [];
            Yaml.Node? node = this.Root.TryReadNode("classes");
            if (node is not null) {
                foreach (Yaml.Node item in node.AsArray().Items) {
                    

                    list.Add(new(item.AsObject()));
                }
            }
            field = list;
            return field;
        }
    }
}
