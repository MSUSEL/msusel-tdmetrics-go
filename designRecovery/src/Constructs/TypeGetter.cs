// Ignore Spelling: proj

namespace designRecovery.src.Constructs;

internal class TypeGetter {
    public TypeGetter(Project proj) {
        List<ITypeDesc> types = [
            .. proj.Types.Basics,
            .. proj.Types.Interfaces,
            .. proj.Types.Named,
            .. proj.Types.Signatures,
            .. proj.Types.Solids,
            .. proj.Types.Structs,
            .. proj.Types.Unions];
        foreach (Package pkg in proj.Packages)
            types.AddRange(pkg.Types);

        this.types = types;
        this.packages = proj.Packages;
    }

    private readonly IReadOnlyList<ITypeDesc> types;
    private readonly IReadOnlyList<Package> packages;

    public T GetTypeAtIndex<T>(uint index) where T : ITypeDesc {
        ITypeDesc type = this.types[(int)index];
        return type is T t ? t :
            throw new InvalidCastException("type at index "+index+" was expected to be "+typeof(T).FullName+" but got "+type.GetType().FullName);
    }

    public Package GetPackageAtIndex(uint index) => this.packages[(int)index];
}
