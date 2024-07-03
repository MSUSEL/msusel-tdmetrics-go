using DesignRecovery.Extensions;
using System.Collections.Generic;
using System.Text.Json.Nodes;

namespace DesignRecovery.Constructs;

public class TypeDef : ITypeDesc {
    public string Name { get; private set; } = "";

    private ITypeDesc? inType;
    public ITypeDesc Type => this.inType ??
        throw new UninitializedException("type");

    private readonly List<Named> inTypeParams = [];
    public IReadOnlyList<Named> TypeParams => this.inTypeParams.AsReadOnly();

    private readonly List<Method> inMethods = [];
    public IReadOnlyList<Method> Methods => this.inMethods.AsReadOnly();

    private Interface? inInterface;
    public Interface Interface => this.inInterface ??
        throw new UninitializedException("interface");

    public void Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
        this.Name = obj.ReadValue<string>("name");
        this.inType = obj.ReadIndexType<ITypeDesc>("type", getter);
        this.inInterface = obj.ReadIndexType<Interface>("interface", getter);
        obj.ReadIndexTypeList("typeParams", getter, this.inTypeParams);

        JsonArray? methodsArr = obj["methods"]?.AsArray();
        if (methodsArr is not null) {
            for (int i = 0; i < methodsArr.Count; i++) {
                JsonNode methodNode = methodsArr[i] ??
                    throw new MissingDataException("methods[" + i + "]");

                Method m = new();
                m.Initialize(getter, methodNode);
                this.inMethods.Add(m);
            }
        }
    }
}
