namespace Yamlite.Tokenizer.Transition;

internal interface ITransition {
    public bool Accept(char c);
}
