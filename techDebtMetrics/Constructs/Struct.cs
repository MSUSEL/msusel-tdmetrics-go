using Constructs.Extensions;
using Constructs.Tooling;
using System.Collections.Generic;
using System.Text.Json.Nodes;

namespace Constructs;

public class Struct : ITypeDesc, IInitializable {
    public IReadOnlyList<Named> Fields => this.inFields.AsReadOnly();
    private readonly List<Named> inFields = [];

    void IInitializable.Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
        obj.ReadIndexTypeList("fields", getter, this.inFields);
    }

    public string ToStub() => throw new System.NotImplementedException(); // TODO: Implement
}
