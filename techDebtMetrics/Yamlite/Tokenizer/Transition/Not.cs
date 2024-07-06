namespace Yamlite.Tokenizer.Transition;

internal record Not(string Inner) : ITransition {
    public bool Accept(char c) => !this.Inner.Contains(c);
}
