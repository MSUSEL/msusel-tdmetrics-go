using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class Interface : ITypeDesc, IInitializable {
    public IReadOnlyList<Named> TypeParams => this.inTypeParams.AsReadOnly();
    private readonly List<Named> inTypeParams = [];

    public IReadOnlyList<Interface> Inherits => this.inInherits.AsReadOnly();
    private readonly List<Interface> inInherits = [];

    public IReadOnlyDictionary<string, ITypeDesc> Methods => this.inMethods.AsReadOnly();
    private readonly Dictionary<string, ITypeDesc> inMethods = [];

    public Union? Union { get; private set; }

    public bool IsAny => this.TypeParams.Count <= 0 && this.Methods.Count <= 0 && this.Inherits.Count <= 0;

    void IInitializable.Initialize(TypeGetter getter, Data.Node node) {
        Data.Object obj = node.AsObject();

        obj.ReadIndexTypeList("typeParams", getter, this.inTypeParams);
        obj.ReadIndexTypeList("inherits", getter, this.inInherits);

        if (obj.Contains("methods")) {
            Data.Object methodObj = obj["methods"].AsObject();
            foreach (KeyValuePair<string, Data.Node> pair in methodObj.Children) {
                string key = pair.Key;
                uint typeIndex = pair.Value.AsUint();
                ITypeDesc value = getter.GetTypeAtIndex<ITypeDesc>(typeIndex);
                this.inMethods[key] = value;
            }
        }

        if (obj.Contains("union"))
            this.Union = obj.ReadIndexType<Union>("union", getter);
    }

    public void ToStub(Journal j) {
        if (this.IsAny) {
            j.Write("any");
            return;
        }

        j.Write("interface");
        Journal j2 = j.Indent;
        j2.AsLong.Write(this.TypeParams, "<", ">");
        j2.AsShort.Write(this.Inherits, ":");
        j.Write(" {");
        if (this.Methods.Count > 0) {
            j2.WriteLine();
            foreach (KeyValuePair<string, ITypeDesc> pair in this.Methods) {
                j2.Write(pair.Key);
                if (pair.Value is not Signature)
                    j2.Write(" ");
                j2.AsShort.Write(pair.Value);
                j2.WriteLine(";");
            }
        }
        j.Write("}");
    }
}
