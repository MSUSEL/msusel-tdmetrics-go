// Ignore Spelling: Yamlite

using System.Collections.Generic;
using System.Linq;

namespace Yamlite;

public class Yamlite {
    static public INode? Parse(string text) {
        IEnumerable<Tokenizer.Token> tokens = Tokenizer.Tokenizer.Yamlite.
            Tokenize(text).Select(cleanToken);

        Yamlite ym = new();
        //foreach (Tokenizer.Token token in tokens) ym.addToken(token);
        return ym.finish();
    }

    private readonly Stack<INode> nodes = new();

    private Yamlite() { }

    static private Tokenizer.Token cleanToken(Tokenizer.Token token) {
        return token.Name switch {
            "SingleValue" => cleanSingleValue(token),
            "DoubleValue" => cleanDoubleValue(token)
            _             => token,
        };
    }

    static private Tokenizer.Token cleanSingleValue(Tokenizer.Token token) {
        string value = token.Value.Trim()[1..-2].Replace("''", "'");
        return new Tokenizer.Token("Value", value, token.Pos);
    }

    static private Tokenizer.Token cleanDoubleValue(Tokenizer.Token token) {
        string value = token.Value.Trim()[1..-2];
        // TODO: Remove escapes
        return new Tokenizer.Token("Value", value, token.Pos);
    }

    private INode finish() {

        return this.nodes.Pop();
    }
}
