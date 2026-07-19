using YamlDotNet.RepresentationModel;
using Yaml = Commons.Data.Yaml;

namespace GroundTruth;

public class ClassMetrics(Yaml.Object obj) {
    public readonly Yaml.Object Object = obj;
    
    public string File  => this.Object.TryReadString("file");
    public string Class => this.Object.TryReadString("class");
    public string Name  => this.Class.Split('.').Last();

    public Yaml.Object Ck => this.Object.TryReadNode("ck")?.AsObject() ?? new Yaml.Object();
    public int    Wmc    => this.Ck.TryReadInt("wmc");
    public int    FanIn  => this.Ck.TryReadInt("fanin");
    public int    FanOut => this.Ck.TryReadInt("fanout");
    public int    Loc    => this.Ck.TryReadInt("loc");
    public double Tcc    => this.Ck.TryReadDouble("tcc");
    public double Lcc    => this.Ck.TryReadDouble("lcc");

    public Yaml.Object Pmd  => this.Object.TryReadNode("pmd")?.AsObject() ?? new Yaml.Object();
    public int Coupling     => this.Pmd.TryReadInt("coupling");
    public int Ncss         => this.Pmd.TryReadInt("ncss");
    public int NcssHighest  => this.Pmd.TryReadInt("ncss_highest");
    public int CycloTotal   => this.Pmd.TryReadInt("cyclo_total");
    public int CycloHighest => this.Pmd.TryReadInt("cyclo_highest");
    public int PublicCount  => this.Pmd.TryReadInt("public_count");

    public List<MethodMetrics> Methods {
        get {
            if (field is not null) return field;

            Dictionary<int, Yaml.Object> ck = [];
            Yaml.Array ckArray = this.Object.TryReadNode("methods_ck")?.AsArray() ?? new Yaml.Array();
            foreach (Yaml.Node ckNode in ckArray.Items) {
                Yaml.Object ckMethod = ckNode.AsObject();
                int line = ckMethod.TryReadInt("line");
                if (ck.ContainsKey(line))
                    throw new Exception("Multiple methods on the same line " + line + " in class " + this.Class);
                ck[line] = ckMethod;
            }

            List<MethodMetrics> list = [];
            Yaml.Array pmdArray = this.Object.TryReadNode("methods_pmd")?.AsArray() ?? new Yaml.Array();
            foreach (Yaml.Node pmdNode in pmdArray.Items) {
                Yaml.Object pmdMethod = pmdNode.AsObject();
                int line = pmdMethod.TryReadInt("line");

                if (ck.TryGetValue(line, out Yaml.Object? ckMethod)) {
                    list.Add(new MethodMetrics(ckMethod, pmdMethod));
                    ck.Remove(line);
                    continue;
                }

                list.Add(new MethodMetrics(new Yaml.Object(), pmdMethod));
            }

            foreach (KeyValuePair<int, Yaml.Object> pair in ck) {
                list.Add(new MethodMetrics(pair.Value, new Yaml.Object()));
            }

            list.Sort((a, b) => a.Line - b.Line);
            field = list;
            return field;
        }
    }

}
