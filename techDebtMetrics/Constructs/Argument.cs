using Commons.Data.Reader;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>An argument is a parameter or result value in a signature.</summary>
/// <see cref="../../docs/genFeatureDef.md#argument"/>
public class Argument : IConstruct, IInitializable<Project> {

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; private set; } = 0;

    /// <summary>Optional name of argument or an empty string.</summary>
    public string Name { get; private set; } = "";

    /// <summary>The type of the argument.</summary>
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

    public void ToStub(Journal j) =>
        j.AsShort.Write(this.Type).Write(this.Name, prefix: " ");
}
