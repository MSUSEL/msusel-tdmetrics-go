using Constructs.Data;
using Constructs.Tooling;

namespace Constructs;

public class Selection : IConstruct, IInitializable {

    // TODO: Implements

    void IInitializable.Initialize(Project project, Node node) { }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) { }
}
