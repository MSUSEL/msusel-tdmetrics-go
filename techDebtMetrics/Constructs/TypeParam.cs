using Commons.Data.Reader;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>A type parameter name and type description.</summary>
/// <see cref="../../docs/genFeatureDef.md#type-parameter"/>
public class TypeParam : ITypeDesc, IInitializable<Project> {

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; private set; } = 0;

    /// <summary>The name of the type parameter.</summary>
    public string Name { get; private set; } = "";

    /// <summary>The type description of the type parameter.</summary>
    public ITypeDesc Type => this.inType ??
        throw new UninitializedException(nameof(this.Type));
    private ITypeDesc? inType;

    /// <summary>Enumerates all the constructs that are directly part of this construct.</summary>
    public IEnumerable<IConstruct> SubConstructs => [this.Type];

    void IInitializable<Project>.Initialize(Project project, int index, Node node) {
        this.Index = index;
        Object obj = node.AsObject();
        this.Name = obj.TryReadString("name");
        this.inType = obj.ReadKey<ITypeDesc>(project, "type");
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        if (j.Long)
            j.AsShort.Write(this.Type, suffix: " ");
        j.Write(this.Name);
    }
}
