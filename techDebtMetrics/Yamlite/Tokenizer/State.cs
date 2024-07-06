using System.Collections.Generic;
using Yamlite.Tokenizer.Transition;

namespace Yamlite.Tokenizer;

internal class State(string tokenName = "", bool consume = false) {
    public string TokenName { get; init; } = tokenName;
    public bool IsConsume { get; init; } = consume;
    public bool IsAccept { get; init; } = consume || !string.IsNullOrEmpty(tokenName);
    
    private readonly List<State> children = [];
    private readonly List<ITransition> trans = [];

    public void Add(State child, ITransition transition) {
        this.children.Add(child);
        this.trans.Add(transition);
    }

    public State? Next(char c) {
        int index = this.trans.FindIndex(t => t.Accept(c));
        return index < 0 ? null : this.children[index];
    }
}
