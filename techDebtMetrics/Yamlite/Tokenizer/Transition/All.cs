namespace Yamlite.Tokenizer.Transition;

internal record All() : ITransition {
    public bool Accept(char c) => true;
}
