namespace Yamlite.Tokenizer;

internal class UnacceptedException(Scanner s) :
    System.Exception("Unaccepted text \"" + s.ToString() + "\" starting at " +
        "offset: " + s.StartOffset + ", " +
        "column: " + s.StartColumn + ", " +
        "line: " + s.StartLine + " and with " +
        "length: " + s.Count + ".") { }
