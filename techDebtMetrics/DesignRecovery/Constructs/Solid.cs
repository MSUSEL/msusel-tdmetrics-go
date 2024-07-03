using DesignRecovery.Extensions;
using System.Collections.Generic;
using System.Text.Json.Nodes;

namespace DesignRecovery.Constructs;

public class Solid : ITypeDesc {
    private ITypeDesc? inTarget;
    public ITypeDesc Target => this.inTarget ??
        throw new UninitializedException("target");

    private readonly List<ITypeDesc> inTypeParams = [];
    public IReadOnlyList<ITypeDesc> TypeParams => this.inTypeParams.AsReadOnly();

    public void Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
        this.inTarget = obj.ReadIndexType<ITypeDesc>("target", getter);
        obj.ReadIndexTypeList("typeParams", getter, this.inTypeParams);
    }
}
