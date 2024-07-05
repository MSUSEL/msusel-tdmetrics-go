using System.Linq;

namespace Yamlite.Tokenizer.Transition;

internal record Any(params char[] Inner) : ITransition {
    public bool Accept(char c) => this.Inner.Contains(c);
}
