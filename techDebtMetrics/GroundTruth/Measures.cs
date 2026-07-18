using Yaml = Commons.Data.Yaml;

namespace GroundTruth;

public class Measures(Yaml.Object obj) {
    public Yaml.Object Object => obj;

    public int    Lines               => this.Object.ReadInt("lines");
    public int    Classes             => this.Object.ReadInt("classes");
    public int    Files               => this.Object.ReadInt("files");
    public int    Functions           => this.Object.ReadInt("functions");
    public int    Complexity          => this.Object.ReadInt("complexity");
    public int    CognitiveComplexity => this.Object.ReadInt("cognitive_complexity");
    public double ClassComplexity     => this.Object.ReadDouble("class_complexity");
    public double FunctionComplexity  => this.Object.ReadDouble("function_complexity");
    public double FileComplexity      => this.Object.ReadDouble("file_complexity");
}
