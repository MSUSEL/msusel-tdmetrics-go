namespace Yamlite.Tokenizer;

internal class UnexpectedCharException(Scanner s) :
    System.Exception("Unexpected character '" + s.Current + "' (" + (int)s.Current + ") at " +
        "offset: " + s.CurrentPos.Offset + ", " +
        "column: " + s.CurrentPos.Column + ", " +
        "line: " + s.CurrentPos.Line + ".") { }
