using Constructs.Exceptions;
using Constructs.Extensions;
using Constructs.Tooling;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Constructs;

public class TypeDef : ITypeDesc, IInitializable {
    public string Name { get; private set; } = "";

    public ITypeDesc Type => this.inType ??
        throw new UninitializedException("type");
    private ITypeDesc? inType;

    public IReadOnlyList<Named> TypeParams => this.inTypeParams.AsReadOnly();
    private readonly List<Named> inTypeParams = [];

    public IReadOnlyList<Method> Methods => this.inMethods.AsReadOnly();
    private readonly List<Method> inMethods = [];

    public Interface Interface => this.inInterface ??
        throw new UninitializedException("interface");
    private Interface? inInterface;

    void IInitializable.Initialize(TypeGetter getter, Data.Node node) {
        Data.Object obj = node.AsObject();
        this.Name = obj.ReadString("name");
        this.inType = obj.ReadIndexType<ITypeDesc>("type", getter);
        this.inInterface = obj.ReadIndexType<Interface>("interface", getter);
        obj.ReadIndexTypeList("typeParams", getter, this.inTypeParams);

        if (obj.Contains("methods")) {
            Data.Array methodsArr = obj["methods"].AsArray();
            for (int i = 0; i < methodsArr.Count; i++) {
                Method m = new();
                (m as IInitializable).Initialize(getter, methodsArr[i]);
                this.inMethods.Add(m);
            }
        }
    }

    public string ToStub() {
        StringBuilder sb = new();
        if (this.Type == this.Interface) {
            sb.Append("interface ");
            sb.Append(this.Name);

            if (this.TypeParams.Count > 0) {
                sb.Append('<');
                sb.Append(this.TypeParams.Select(tp => tp.ToStub().Indent()).Join());
                sb.Append('>');
            }

            sb.Append(": ");
            sb.Append(this.Type.ToStub());
            return sb.ToString();
        }

        sb.Append("class ");
        sb.Append(this.Name);

        if (this.TypeParams.Count > 0) {
            sb.Append('<');
            sb.Append(this.TypeParams.Select(tp => tp.ToStub().Indent()).Join());
            sb.Append('>');
        }

        sb.AppendLine("{");
        sb.Append("   Data ");
        sb.Append(this.Type.ToStub().Indent());
        
        if (this.Methods.Count > 0) {
            sb.AppendLine();
            foreach (Method m in this.Methods) {
                sb.Append("   ");
                sb.Append(m.ToStub().Indent());
                sb.AppendLine(";");
            }
        }
        sb.Append('}');
        return sb.ToString();
    }
}
