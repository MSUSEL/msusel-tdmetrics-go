using Yaml = Commons.Data.Yaml;

namespace GroundTruth;

internal class MethodMetrics(Yaml.Object ckObj, Yaml.Object pmdObj) {
    public readonly Yaml.Object Ck = ckObj;
    public bool Constructor => this.Ck.ReadBool("constructor");
    public bool Line        => this.Ck.ReadBool("line");
    public bool FanIn       => this.Ck.ReadBool("fanin");
    public bool FanOut      => this.Ck.ReadBool("fanout");
    public bool Wmc         => this.Ck.ReadBool("wmc");
    public bool Loc         => this.Ck.ReadBool("loc");

    public readonly Yaml.Object Pmd = pmdObj;
    public string Signature => this.Pmd.ReadString("signature");
    public int Ncss   => this.Pmd.ReadInt("ncss");
    public int Cyclo  => this.Pmd.ReadInt("cyclo");
    public int NPath  => this.Pmd.ReadInt("npath");
    public int Params => this.Pmd.ReadInt("params");
}
