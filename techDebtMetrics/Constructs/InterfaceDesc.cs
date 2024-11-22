using Constructs.Data;
using Constructs.Tooling;

namespace Constructs;

public class InterfaceDesc : ITypeDesc, IInitializable {

    // TODO: Implements

    void IInitializable.Initialize(Package package, Node node) { }

    public void ToStub(Journal j) { }
}
