using System.Collections.Generic;

namespace Constructs.Tooling;

internal class TypeGetter {
    public TypeGetter(Project p) {
        List<ITypeDesc> types = [
            .. p.Types.Basics,
            .. p.Types.Interfaces,
            .. p.Types.Named,
            .. p.Types.Signatures,
            .. p.Types.Solids,
            .. p.Types.Structs,
            .. p.Types.Unions];
        foreach (Package pkg in p.Packages)
            types.AddRange(pkg.Types);

        this.types = types;
        this.packages = p.Packages;
    }

    private readonly IReadOnlyList<ITypeDesc> types;
    private readonly IReadOnlyList<Package> packages;

    public T GetTypeAtIndex<T>(uint index) where T : ITypeDesc {
        if (index <= 0 || index > this.types.Count)
            throw new System.ArgumentOutOfRangeException(nameof(index),
                "Read type [1.." + this.types.Count + "] at " + index + ".");

        ITypeDesc type = this.types[(int)index-1];
        return type is T t ? t :
            throw new System.InvalidCastException("type at index " + index +
                " was expected to be " + typeof(T).FullName + " but got " + type.GetType().FullName);
    }

    public Package GetPackageAtIndex(uint index) => this.packages[(int)index];
}
