using Constructs.Exceptions;
using Constructs.Extensions;
using Constructs.Tooling;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Text.Json.Nodes;

namespace Constructs;

public class Interface : ITypeDesc, IInitializable {
    public IReadOnlyList<Named> TypeParams => this.inTypeParams.AsReadOnly();
    private readonly List<Named> inTypeParams = [];

    public IReadOnlyList<Interface> Interfaces => this.inInterfaces.AsReadOnly();
    private readonly List<Interface> inInterfaces = [];

    public IReadOnlyDictionary<string, ITypeDesc> Methods => this.inMethods.AsReadOnly();
    private readonly Dictionary<string, ITypeDesc> inMethods = [];

    public Union? Union { get; private set; }

    void IInitializable.Initialize(TypeGetter getter, Data.Node node) {
        Data.Object obj = node.AsObject();

        obj.ReadIndexTypeList("typeParams", getter, this.inTypeParams);
        obj.ReadIndexTypeList("interfaces", getter, this.inInterfaces);

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

    public string ToStub() {
        StringBuilder sb = new();
        sb.Append("interface");
        
        if (this.TypeParams.Count > 0) {
            sb.Append('<');
            sb.Append(this.TypeParams.Select(tp => tp.ToStub()).Join());
            sb.Append('>');
        }
        
        if (this.Interfaces.Count > 0) {
            sb.Append(':');
            sb.Append(this.Interfaces.Select(tp => tp.ToStub()).Join(", "));
        }

        sb.Append(" {");
        if (this.Methods.Count > 0) {
            sb.AppendLine();
            foreach (KeyValuePair<string, ITypeDesc> pair in this.Methods) {
                sb.Append("   ");
                sb.Append(pair.Key);
                sb.Append(' ');
                sb.Append(pair.Value.ToStub());
                sb.AppendLine();
            }
        }
        sb.Append('}');
        return sb.ToString();
    }
}
