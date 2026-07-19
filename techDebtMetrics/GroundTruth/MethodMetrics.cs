using Yaml = Commons.Data.Yaml;

namespace GroundTruth;

public class MethodMetrics(Yaml.Object ckObj, Yaml.Object pmdObj) {
    public readonly Yaml.Object Ck = ckObj;
    public bool Constructor => this.Ck.ReadBool("constructor");
    public int  Line        => this.Ck.ReadInt("line");
    public int  FanIn       => this.Ck.ReadInt("fanin");
    public int  FanOut      => this.Ck.ReadInt("fanout");
    public int  Wmc         => this.Ck.ReadInt("wmc");
    public int  Loc         => this.Ck.ReadInt("loc");

    public readonly Yaml.Object Pmd = pmdObj;
    public string Signature => this.Pmd.ReadString("signature");
    public int    Ncss      => this.Pmd.ReadInt("ncss");
    public int    Cyclo     => this.Pmd.ReadInt("cyclo");
    public int    NPath     => this.Pmd.ReadInt("npath");
    public int    Params    => this.Pmd.ReadInt("params");
}
