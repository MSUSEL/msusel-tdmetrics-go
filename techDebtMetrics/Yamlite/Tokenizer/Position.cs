namespace Yamlite.Tokenizer;

internal record Position(int Offset, int Column, int Line) {
    public Position Step(char c) => c == '\n' ?
        new Position(this.Offset+1, 0, this.Line+1) :
        new Position(this.Offset+1, this.Column+1, this.Line);
}
