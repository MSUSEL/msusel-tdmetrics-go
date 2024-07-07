using Constructs.Exceptions;
using Constructs.Tooling;

namespace Constructs;

public class ValueDef : IConstruct, IInitializable {
    public string Name { get; private set; } = "";
    public bool IsConst { get; private set; } = false;

    public ITypeDesc Type => this.inType ??
        throw new UninitializedException("type");
    private ITypeDesc? inType;

    void IInitializable.Initialize(TypeGetter getter, Data.Node node) {
        Data.Object obj = node.AsObject();
        this.Name = obj.ReadString("name");
        if (obj.Contains("isConst"))
            this.IsConst = obj.ReadBool("isConst");
        this.inType = obj.ReadIndexType<ITypeDesc>("type", getter);
    }

    public void ToStub(Journal j) {
        if (this.IsConst) j.Write("const ");
        j.Write(this.Name).Write(" ").Write(this.Type);
    }
}
