using System.Collections.Generic;
using System.Linq;
using Yamlite.Tokenizer.Transition;

namespace Yamlite.Tokenizer;

internal class State(string tokenName = "") {

    public string TokenName { get; init; } = tokenName;

    public bool IsAccept => string.IsNullOrEmpty(this.TokenName);

    public IReadOnlyDictionary<ITransition, State> Children => this.inChildren.AsReadOnly();
    private readonly Dictionary<ITransition, State> inChildren = [];

    public void Add(ITransition transition, State child) => this.inChildren[transition] = child;

    public State? Next(char c) => this.inChildren.FirstOrDefault(p => p.Key.Accept(c)).Value;
}
