using designRecovery.src.Extensions;
using System.Text.Json.Nodes;

namespace designRecovery.src.Constructs;

internal class Struct {    
    private readonly List<Named> inFields = [];
    public IReadOnlyList<Named> Fields => this.inFields.AsReadOnly();

	public void Initialize(ITypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
		obj.ReadIndexTypeList("fields", getter, this.inFields);
	}
}
