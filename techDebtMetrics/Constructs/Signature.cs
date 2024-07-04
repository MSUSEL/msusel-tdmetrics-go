using Constructs.Exceptions;
using Constructs.Extensions;
using Constructs.Tooling;
using System.Collections.Generic;
using System.Text.Json.Nodes;

namespace Constructs;

public class Signature : ITypeDesc, IInitializable {
    public bool Variadic { get; private set; }

    public IReadOnlyList<Named> Params => this.inParams.AsReadOnly();
    private readonly List<Named> inParams = [];

    public IReadOnlyList<Named> TypeParams => this.inTypeParams.AsReadOnly();
    private readonly List<Named> inTypeParams = [];

    public ITypeDesc ReturnType => this.inReturnType ??
        throw new UninitializedException("returnType");
    private ITypeDesc? inReturnType;

    void IInitializable.Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();

        if (obj.ContainsKey("variadic"))
            this.Variadic = obj.ReadValue<bool>("variadic");

        obj.ReadIndexTypeList("params", getter, this.inParams);
        obj.ReadIndexTypeList("typeParams", getter, this.inTypeParams);

        if (obj.ContainsKey("return"))
            this.inReturnType = obj.ReadIndexType<ITypeDesc>("return", getter);
    }

    public string ToStub() => throw new System.NotImplementedException(); // TODO: Implement
}
