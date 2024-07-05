namespace Yamlite.Tokenizer.Transition;

internal record Not(ITransition Inner) : ITransition {
    public bool Accept(char c) => !this.Inner.Accept(c);
}
