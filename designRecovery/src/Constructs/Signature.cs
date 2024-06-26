﻿using designRecovery.src.Extensions;
using System.Text.Json.Nodes;

namespace designRecovery.src.Constructs;

internal class Signature : ITypeDesc {
    public bool Variadic { get; private set; }

    private readonly List<Named> inParams = [];
    public IReadOnlyList<Named> Params => this.inParams.AsReadOnly();

    private readonly List<Named> inTypeParams = [];
    public IReadOnlyList<Named> TypeParams => this.inTypeParams.AsReadOnly();

    private ITypeDesc? inReturnType;
    public ITypeDesc ReturnType => this.inReturnType ??
        throw new UninitializedException("returnType");

    public void Initialize(ITypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();

        if (obj.ContainsKey("variadic"))
            this.Variadic = obj.ReadValue<bool>("variadic");

        obj.ReadIndexTypeList("params", getter, this.inParams);
        obj.ReadIndexTypeList("typeParams", getter, this.inTypeParams);

        if (obj.ContainsKey("return"))
            this.inReturnType = obj.ReadIndexType<ITypeDesc>("return", getter);
    }
}
