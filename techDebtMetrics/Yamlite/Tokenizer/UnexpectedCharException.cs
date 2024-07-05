namespace Yamlite.Tokenizer;

internal class UnexpectedCharException(Scanner s) :
    System.Exception("Unexpected character '" + s.Current + "' (" + (int)s.Current + ") at " +
        "offset: " + s.CurrentOffset + ", " +
        "column: " + s.CurrentColumn + ", " +
        "line: " + s.CurrentLine + ".") { }
