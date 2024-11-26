using Constructs.Tooling;

namespace Constructs;

public interface IDeclaration : IConstruct {
    public string Name { get; }
    public Location Location { get; }
    public Package Package { get; }
}
