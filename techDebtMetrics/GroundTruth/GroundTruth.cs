using System.IO.Compression;
using Yaml = Commons.Data.Yaml;
using Commons.Data.Repo;

namespace GroundTruth;

public class GroundTruth(Yaml.Node node) {
    static public GroundTruth FromZip(string zipFile, JavaTarget target) =>
        FromZip(zipFile, target.ProjectKey + ".json");

    static public GroundTruth FromZip(string zipFile, string targetFile) {
        using ZipArchive zip = ZipFile.Open(zipFile, ZipArchiveMode.Read);
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

    public Yaml.Object Measures => this.Root.TryReadNode("measures")?.AsObject() ?? new Yaml.Object();

    public int    Lines               => this.Measures.ReadInt("lines");
    public int    Classes             => this.Measures.ReadInt("classes");
    public int    Files               => this.Measures.ReadInt("files");
    public int    Functions           => this.Measures.ReadInt("functions");
    public int    Complexity          => this.Measures.ReadInt("complexity");
    public int    CognitiveComplexity => this.Measures.ReadInt("cognitive_complexity");
    public double ClassComplexity     => this.Measures.ReadDouble("class_complexity");
    public double FunctionComplexity  => this.Measures.ReadDouble("function_complexity");
    public double FileComplexity      => this.Measures.ReadDouble("file_complexity");

    public List<DeclMetrics> Declarations {
        get {
            if (field is not null) return field;
            List<DeclMetrics> list = [];
            Yaml.Node? node = this.Root.TryReadNode("classes");
            if (node is not null) {
                foreach (Yaml.Node item in node.AsArray().Items)
                    list.Add(new(item.AsObject()));
            }
            field = list;
            return field;
        }
    }
}
