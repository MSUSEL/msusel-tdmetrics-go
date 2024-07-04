using Constructs.Exceptions;
using Constructs.Extensions;
using Constructs.Tooling;
using System.Collections.Generic;
using System.Text.Json.Nodes;

namespace Constructs;

public class TypeDef : ITypeDesc, IInitializable {
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

    void IInitializable.Initialize(TypeGetter getter, JsonNode node) {
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
                (m as IInitializable).Initialize(getter, methodNode);
                this.inMethods.Add(m);
            }
        }
    }

    public string ToStub() => throw new System.NotImplementedException(); // TODO: Implement
}
