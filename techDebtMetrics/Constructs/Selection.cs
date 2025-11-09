using Commons.Data.Yaml;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>A selection of information from another, e.g. selection of a field from an object.</summary>
/// <see cref="../../docs/genFeatureDef.md#selection"/>
public class Selection : IConstruct, IInitializable<Project> {

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; private set; } = 0;

    /// <summary>The name of the field, abstract, etc that is being selected from the origin.</summary>
    public string Name { get; private set; } = "";

    /// <summary>The origin object, interface, etc that is being selected out of.</summary>
    public IConstruct Origin => this.inOrigin ??
        throw new UninitializedException(nameof(this.Origin));
    private IConstruct? inOrigin;

    /// <summary>Enumerates all the constructs that are directly part of this construct.</summary>
    public IEnumerable<IConstruct> SubConstructs => [this.Origin];

    void IInitializable<Project>.Initialize(Project project, int index, Node node) {
        this.Index = index;
        Object obj = node.AsObject();
        this.Name = obj.ReadString("name");
        this.inOrigin = obj.ReadKey<IConstruct>(project, "origin");
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) =>
        j.AsShort.Write(this.Origin).Write(".").Write(this.Name);
}
