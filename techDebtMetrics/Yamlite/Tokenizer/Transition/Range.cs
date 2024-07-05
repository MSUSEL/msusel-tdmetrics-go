using System.Linq;

namespace Yamlite.Tokenizer.Transition;

internal record Range(char Low, char High) : ITransition {
    public bool Accept(char c) => this.Low >= c && this.High <= c;
}
