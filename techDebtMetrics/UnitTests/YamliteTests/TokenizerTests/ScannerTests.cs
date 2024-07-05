using Yamlite.Tokenizer;

namespace UnitTests.YamliteTests.TokenizerTests;

public class ScannerTests {
    
    [Test]
    public void StringWithOneNewline() {
        Scanner s = new("Hello\nWorld".GetEnumerator());
        checkCurrent(s, 0, 0, 1, '\0', "");
        checkStart(s, 0, 0, 1);

        checkNext(s, 0, 0, 1, 'H', "H");
        checkNext(s, 1, 1, 1, 'e', "He");
        checkNext(s, 2, 2, 1, 'l', "Hel");
        checkNext(s, 3, 3, 1, 'l', "Hell");
        checkNext(s, 4, 4, 1, 'o', "Hello");
        
        checkNext(s, 5, 5, 1, '\n', "Hello\n");
        checkNext(s, 6, 0, 2, 'W', "Hello\nW");
        checkNext(s, 7, 1, 2, 'o', "Hello\nWo");
        checkNext(s, 8, 2, 2, 'r', "Hello\nWor");
        checkNext(s, 9, 3, 2, 'l', "Hello\nWorl");
        checkNext(s, 10, 4, 2, 'd', "Hello\nWorld");
        
        Assert.That(s.MoveNext(), Is.False);
        checkCurrent(s, 10, 4, 2, 'd', "Hello\nWorld");
        checkStart(s, 0, 0, 1);

        Assert.That(s.Take(3), Is.EqualTo("Hel"));
        checkCurrent(s, 2, 2, 1, 'l', "");
        checkStart(s, 2, 2, 1);
        
        checkNext(s, 3, 3, 1, 'l', "l");
        checkNext(s, 4, 4, 1, 'o', "lo");
        
        checkNext(s, 5, 5, 1, '\n', "lo\n");
        checkNext(s, 6, 0, 2, 'W', "lo\nW");
        checkNext(s, 7, 1, 2, 'o', "lo\nWo");

        Assert.That(s.Take(4), Is.EqualTo("lo\nW"));
        checkCurrent(s, 6, 0, 2, 'W', "");
        checkStart(s, 6, 0, 2);
        
        checkNext(s, 7, 1, 2, 'o', "o");
        checkNext(s, 8, 2, 2, 'r', "or");
        checkNext(s, 9, 3, 2, 'l', "orl");
        checkNext(s, 10, 4, 2, 'd', "orld");
        
        Assert.That(s.MoveNext(), Is.False);
        checkCurrent(s, 10, 4, 2, 'd', "orld");
        checkStart(s, 6, 0, 2);
        
        Assert.That(s.Take(4), Is.EqualTo("orld"));
        checkCurrent(s, 10, 4, 2, 'd', "");
        checkStart(s, 10, 4, 2);
        
        Assert.That(s.MoveNext(), Is.False);
        checkCurrent(s, 10, 4, 2, 'd', "");
        checkStart(s, 10, 4, 2);
    }

    [Test]
    public void ResetAndTakeNone() {
        Scanner s = new("Hello".GetEnumerator());
        checkNext(s, 0, 0, 1, 'H', "H");
        checkNext(s, 1, 1, 1, 'e', "He");
        checkNext(s, 2, 2, 1, 'l', "Hel");
        checkNext(s, 3, 3, 1, 'l', "Hell");
        checkNext(s, 4, 4, 1, 'o', "Hello");
        
        Assert.That(s.Take(2), Is.EqualTo("He"));
        checkCurrent(s, 1, 1, 1, 'e', "");
        checkNext(s, 2, 2, 1, 'l', "l");
        checkNext(s, 3, 3, 1, 'l', "ll");
        
        Assert.That(s.Take(0), Is.EqualTo(""));
        checkCurrent(s, 1, 1, 1, '\0', "");
        checkNext(s, 2, 2, 1, 'l', "l");
        checkNext(s, 3, 3, 1, 'l', "ll");
        checkNext(s, 4, 4, 1, 'o', "llo");
        
        Assert.That(s.Take(0), Is.EqualTo(""));
        checkCurrent(s, 1, 1, 1, '\0', "");
        checkNext(s, 2, 2, 1, 'l', "l");

        s.Reset();
        checkNext(s, 0, 0, 1, 'H', "H");
        checkNext(s, 1, 1, 1, 'e', "He");
        checkNext(s, 2, 2, 1, 'l', "Hel");
        checkNext(s, 3, 3, 1, 'l', "Hell");
        checkNext(s, 4, 4, 1, 'o', "Hello");

        s.Dispose();
        checkCurrent(s, 0, 0, 1, '\0', "");
        checkStart(s, 0, 0, 1);
    }

    [Test]
    public void AllNewlines() {
        Scanner s = new("\n\n\n".GetEnumerator());

        checkNext(s, 0, 0, 1, '\n', "\n");
        checkNext(s, 1, 0, 2, '\n', "\n\n");
        checkNext(s, 2, 0, 3, '\n', "\n\n\n");
        
        Assert.That(s.Take(2), Is.EqualTo("\n\n"));
        checkCurrent(s, 1, 0, 2, '\n', "");
        
        checkNext(s, 2, 0, 3, '\n', "\n");
        Assert.That(s.MoveNext(), Is.False);
    }

    [Test]
    public void TakeFullBuffer() {
        Scanner s = new("Hello".GetEnumerator());
        checkNext(s, 0, 0, 1, 'H', "H");
        checkNext(s, 1, 1, 1, 'e', "He");
        checkNext(s, 2, 2, 1, 'l', "Hel");
        checkNext(s, 3, 3, 1, 'l', "Hell");
        checkNext(s, 4, 4, 1, 'o', "Hello");
        
        Assert.That(s.Take(5), Is.EqualTo("Hello"));
        checkCurrent(s, 4, 4, 1, 'o', "");

        Assert.That(s.MoveNext(), Is.False);
        checkCurrent(s, 4, 4, 1, 'o', "");

        Assert.That(s.Take(0), Is.EqualTo(""));
        checkCurrent(s, 4, 4, 1, '\0', "");
    }
    
    static private void checkNext(Scanner s, int offset, int column, int line, char c, string str) =>
        Assert.Multiple(() => {
            Assert.That(s.MoveNext(), Is.True);
            checkCurrent(s, offset, column, line, c, str);
        });

    static private void checkCurrent(Scanner s, int offset, int column, int line, char c, string str) =>
        Assert.Multiple(() => {
            Assert.That(s.CurrentOffset, Is.EqualTo(offset), "current offset");
            Assert.That(s.CurrentColumn, Is.EqualTo(column), "current column");
            Assert.That(s.CurrentLine, Is.EqualTo(line), "current line");
            Assert.That(s.Current, Is.EqualTo(c), "current character");
            Assert.That(s.ToString(), Is.EqualTo(str), "current buffer text");
            Assert.That(s.Count, Is.EqualTo(str.Length), "current buffer count");
        });

    static private void checkStart(Scanner s, int offset, int column, int line) =>
        Assert.Multiple(() => {
            Assert.That(s.StartOffset, Is.EqualTo(offset), "current offset");
            Assert.That(s.StartColumn, Is.EqualTo(column), "current column");
            Assert.That(s.StartLine, Is.EqualTo(line), "current line");
        });
}
