using DesignRecovery.Extensions;
using System.Collections.Generic;
using System.Text.Json.Nodes;

namespace DesignRecovery.Constructs;

public class Struct : ITypeDesc, IInitializable {
    private readonly List<Named> inFields = [];
    public IReadOnlyList<Named> Fields => this.inFields.AsReadOnly();

    public void Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
        obj.ReadIndexTypeList("fields", getter, this.inFields);
    }
}
