namespace Constructs.Tooling;

internal interface IInitializable {

    void Initialize(TypeGetter getter, Data.Node node);
}
