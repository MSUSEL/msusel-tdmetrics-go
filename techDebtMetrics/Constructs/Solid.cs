using Constructs.Exceptions;
using Constructs.Extensions;
using Constructs.Tooling;
using System.Collections.Generic;
using System.Linq;
using System.Text.Json.Nodes;

namespace Constructs;

public class Solid : ITypeDesc, IInitializable {
    public ITypeDesc Target => this.inTarget ??
        throw new UninitializedException("target");
    private ITypeDesc? inTarget;

    public IReadOnlyList<ITypeDesc> TypeParams => this.inTypeParams.AsReadOnly();
    private readonly List<ITypeDesc> inTypeParams = [];

    void IInitializable.Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
        this.inTarget = obj.ReadIndexType<ITypeDesc>("target", getter);
        obj.ReadIndexTypeList("typeParams", getter, this.inTypeParams);
    }

    public string ToStub() => this.Target.ToStub() +
        "<" + this.TypeParams.Select(tp => tp.ToStub()).Join() + ">";
}
