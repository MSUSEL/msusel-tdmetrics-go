// Ignore Spelling: Yamlite

using System.Collections.Generic;

namespace Yamlite;

public class Yamlite {
    static public INode? Parse(string text) {
        IEnumerable<Tokenizer.Token> tokens = Tokenizer.Tokenizer.Yamlite.Tokenize(text);
        Yamlite ym = new();
        foreach (Tokenizer.Token token in tokens) ym.addToken(token);
        return ym.finish();
    }

    private readonly Stack<INode> nodes = new();

    private Yamlite() { }

    private void addToken(Tokenizer.Token token) {
        switch (token.Name) {
            case "OpenObject": this.openObject(); break;
            case "CloseObject": this.closeObject(); break;
            case "OpenArray": this.openArray(); break;
            case "CloseArray": this.closeArray(); break;
            case "Colon": this.colon(); break;
            case "Comma": this.comma(); break;
            case "SingleValue": this.singleValue(token.Value); break;
            case "DoubleValue": this.doubleValue(token.Value); break;
            case "Value": this.basicvalue(token.Value); break;
            default: throw new System.ArgumentException("Unexpected token " + token + ".");
        }
    }

    private void openObject() { }
    private void closeObject() { }
    private void openArray() { }
    private void closeArray() { }
    private void colon() { }
    private void comma() { }

    private void singleValue(string value) { }
    private void doubleValue(string value) { }
    private void basicvalue(string value) { }

    private void value(string value) { }

    private INode finish() {

        return this.nodes.Pop();
    }
}
