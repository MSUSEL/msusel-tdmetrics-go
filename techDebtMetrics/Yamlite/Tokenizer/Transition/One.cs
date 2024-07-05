namespace Yamlite.Tokenizer.Transition;

internal record One(char Inner) : ITransition {
    public bool Accept(char c) => this.Inner == c;
}
