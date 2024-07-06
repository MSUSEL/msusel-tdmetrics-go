using Yamlite.Tokenizer;

namespace UnitTests.YamliteTests.TokenizerTests;

public class TokenizerTests {

    [Test]
    public void SimpleObject() =>
        checkTokens("{ hello: world }",
            new Token("OpenObject", "{", 0, 0, 1),
            new Token("Value", "hello", 6, 6, 1),
            new Token("Colon", ":", 7, 7, 1),
            new Token("Value", "world", 13, 13, 1),
            new Token("CloseObject", "}", 15, 15, 1));

    [Test]
    public void SimpleArray() =>
        checkTokens("[ hello, world ]",
            new Token("OpenArray", "[", 0, 0, 1),
            new Token("Value", "hello", 6, 6, 1),
            new Token("Comma", ",", 7, 7, 1),
            new Token("Value", "world", 13, 13, 1),
            new Token("CloseArray", "]", 15, 15, 1));

    [Test]
    public void QuoteLiterals() =>
        checkTokens("\"[\\\"]\" , ab (c\") de , 0.0 , -10 , 'ain''t' ",
            new Token("DoubleValue", "\"[\\\"]\"", 5, 5, 1),
            new Token("Comma", ",", 7, 7, 1),
            new Token("Value", "ab (c\") de", 18, 18, 1),
            new Token("Comma", ",", 20, 20, 1),
            new Token("Value", "0.0", 24, 24, 1),
            new Token("Comma", ",", 26, 26, 1),
            new Token("Value", "-10", 30, 30, 1),
            new Token("Comma", ",", 32, 32, 1),
            new Token("SingleValue", "'ain''t'", 41, 41, 1));

    [Test]
    public void Comments() =>
        checkTokens("hello\nsmall # comment\n\t blue\t \nplanet",
            new Token("Value", "hello", 4, 4, 1),
            new Token("Value", "small", 10, 4, 2),
            new Token("Value", "blue", 27, 5, 3),
            new Token("Value", "planet", 36, 5, 4));

    static private void checkTokens(string text, params Token[] expToken) =>
        Assert.That(expToken, Is.EqualTo(Tokenizer.Yamlite.Tokenize(text)));
}
