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

    public void ToStub(Journal j) {
        if (this.Type == this.Interface) {
            j.Write("interface ").Write(this.Name).Write(this.TypeParams, "<", ">").Write(": ").Write(this.Type);
            return;
        }

        j.Write("class ").Write(this.Name);
        j.Write(this.TypeParams, "<", ">");
        j.WriteLine("{").Write("   Data ").Indent.AsShort.Write(this.Type);
        j.Indent.Write(this.Methods, "\n", "", ";\n");
        j.Write("}");
    }
}
