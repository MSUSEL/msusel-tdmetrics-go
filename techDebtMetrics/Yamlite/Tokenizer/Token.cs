namespace Yamlite.Tokenizer;

internal record Token(string Name, string Value, int Offset, int Column, int Line) { }
