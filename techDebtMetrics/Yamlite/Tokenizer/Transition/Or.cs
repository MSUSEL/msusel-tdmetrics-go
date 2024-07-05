using System.Linq;

namespace Yamlite.Tokenizer.Transition;

internal record Or(params ITransition[] Inner) : ITransition {
    public bool Accept(char c) => this.Inner.Any(i => i.Accept(c));
}
