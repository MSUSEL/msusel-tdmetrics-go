namespace Yamlite.Tokenizer.Transition;

internal record Any(string Inner) : ITransition {
    public bool Accept(char c) => this.Inner.Contains(c);
}
