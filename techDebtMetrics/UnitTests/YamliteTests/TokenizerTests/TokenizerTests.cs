using Yamlite.Tokenizer;

namespace UnitTests.YamliteTests.TokenizerTests;

public class TokenizerTests {

    [Test]
    public void SimpleObject() {
        Tokenizer t = Tokenizer.Yamlite();
        checkTokens(t.Tokenize("{ hello: world }"),
            new Token("OpenObject", "{"),
            new Token("Value", "hello"),
            new Token("Colon", ":"),
            new Token("Value", "world "),
            new Token("CloseObject", "}"));
    }

    static private void checkTokens(IEnumerable<Token> tokens, params Token[] expToken) =>
        Assert.That(expToken, Is.EqualTo(tokens));
}
