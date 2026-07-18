using Yaml = Commons.Data.Yaml;

namespace GroundTruth;

public class ClassMetrics(Yaml.Object obj) {
    public readonly Yaml.Object Object = obj;
    
    public string File  => this.Object.ReadString("file");
    public string Class => this.Object.ReadString("class");
    public string Name  => this.Class.Split('.').Last();

    public Yaml.Object Ck => this.Object.ReadNode("ck").AsObject();
    public int    Wmc    => this.Ck.ReadInt("wmc");
    public int    FanIn  => this.Ck.ReadInt("fanin");
    public int    FanOut => this.Ck.ReadInt("fanout");
    public int    Loc    => this.Ck.ReadInt("loc");
    public double Tcc    => this.Ck.ReadDouble("tcc");
    public double Lcc    => this.Ck.ReadDouble("lcc");

    public Yaml.Object Pmd  => this.Object.ReadNode("pmd").AsObject();
    public int Coupling     => this.Pmd.ReadInt("coupling");
    public int Ncss         => this.Pmd.ReadInt("ncss");
    public int NcssHighest  => this.Pmd.ReadInt("ncss_highest");
    public int CycloTotal   => this.Pmd.ReadInt("cyclo_total");
    public int CycloHighest => this.Pmd.ReadInt("cyclo_highest");
    public int PublicCount  => this.Pmd.ReadInt("public_count");
    

    // TODO: Collect Methods
    

}
