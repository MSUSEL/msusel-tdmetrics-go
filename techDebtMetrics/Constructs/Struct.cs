using Constructs.Data;
using Constructs.Tooling;

namespace Constructs;

public class Struct : ITypeDesc, IInitializable {
   
    void IInitializable.Initialize(Project project, Node node) { }

    public void ToStub(Journal j) { }
}
