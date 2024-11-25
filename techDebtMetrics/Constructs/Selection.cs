using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;

namespace Constructs;

public class Selection : IConstruct, IInitializable {
    public string Name { get; private set; } = "";

    public IConstruct Origin => this.inOrigin ??
        throw new UninitializedException(nameof(this.Origin));
    private IConstruct? inOrigin;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.Name     = obj.ReadString("name");
        this.inOrigin = obj.ReadKey<IConstruct>("origin", project);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) =>
        j.AsShort.Write(this.Origin).Write(".").Write(this.Name);
}
