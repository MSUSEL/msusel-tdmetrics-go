using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>A field is named types in a structure.</summary>
/// <see cref="../../docs/genFeatureDef.md#field"/>
public class Field : IConstruct, IInitializable {

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; private set; } = 0;

    /// <summary>The name of the field.</summary>
    public string Name { get; private set; } = "";

    /// <summary>The type of the field.</summary>
    public ITypeDesc Type => this.inType ??
        throw new UninitializedException(nameof(this.Type));
    private ITypeDesc? inType;

    /// <summary>Enumerates all the constructs that are directly part of this construct.</summary>
    public IEnumerable<IConstruct> SubConstructs => [this.Type];

    void IInitializable.Initialize(Project project, int index, Node node) {
        this.Index = index;
        Object obj = node.AsObject();
        this.Name = obj.ReadString("name");
        this.inType = obj.ReadKey<ITypeDesc>("type", project);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        if (j.Long)
            j.AsShort.Write(this.Type, suffix: " ");
        j.Write(this.Name);
    }
}
