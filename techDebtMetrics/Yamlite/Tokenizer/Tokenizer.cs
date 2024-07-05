using System.Collections.Generic;

namespace Yamlite.Tokenizer;

internal class Tokenizer(State start) {
    public IReadOnlyList<State> States => this.inStates.AsReadOnly();
    public readonly List<State> inStates = [start];

    public State Start { get; init; } = start;

    public void Add(params State[] states) => this.inStates.AddRange(states);

    public IEnumerable<Token> Tokenize(IEnumerable<char> text) {
        State current = this.Start;
        int acceptIndex = 0;
        State? accept = null;
        Scanner source = new(text.GetEnumerator());

        while (true) {
            if (source.MoveNext()) {
                State? next = current.Next(source.Current);
                if (next is null) {
                    if (accept is null)
                        throw new UnexpectedCharException(source);

                    string value = source.Take(acceptIndex);
                    yield return new Token(accept.TokenName, value);
                    continue;
                }
                
                if (next.IsAccept) {
                    accept = next;
                    acceptIndex = source.Count;
                }
                continue;
            }

            if (source.Count > 0) {
                if (accept is null)
                    throw new UnacceptedException(source);

                string value = source.Take(acceptIndex);
                yield return new Token(accept.TokenName, value);
                continue;
            }

            break;
        }
    }
}
