using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;

namespace Constructs;

public class Field : IConstruct, IInitializable {
    public string Name { get; private set; } = "";

    public ITypeDesc Type => this.inType ??
        throw new UninitializedException(nameof(this.Type));
    private ITypeDesc? inType;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.Name   = obj.ReadString("name");
        this.inType = obj.ReadKey<ITypeDesc>("type", project);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        j.Write(this.Name);
        if (!j.Short)
            j.Write(" ").AsShort.Write(this.Type).Write(";");
    }
}
