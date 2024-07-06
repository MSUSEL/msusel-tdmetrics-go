// Ignore Spelling: Yamlite

using System.Collections.Generic;
using Yamlite.Tokenizer.Transition;

namespace Yamlite.Tokenizer;

internal class Tokenizer(State start) {
    private readonly State start = start;

    private class RunnerState(State start, Scanner source) {
        private readonly State start = start;
        private readonly Scanner source = source;
        private State current = start;
        private int acceptIndex = 0;
        private State? accept = null;

        public bool Running = true;

        private Token? createToken() {
            if (this.accept is null)
                throw new UnexpectedCharException(source);

            string value = this.source.Take(this.acceptIndex);
            Token? token = this.accept.IsConsume ? null :
                new Token(this.accept.TokenName, value, this.source.CurrentPos);

            this.current = this.start;
            this.acceptIndex = 0;
            this.accept = null;
            return token;
        }

        public Token? Step() {
            if (this.source.MoveNext()) {
                State? next = this.current.Next(this.source.Current);
                if (next is null) return this.createToken();

                if (next.IsAccept) {
                    this.accept = next;
                    this.acceptIndex = this.source.Count;
                }
                this.current = next;
                return null;
            }

            if (this.source.Count > 0) return this.createToken();

            this.Running = false;
            return null;
        }
    }


    public IEnumerable<Token> Tokenize(IEnumerable<char> text) {
        RunnerState r = new(this.start, new(text.GetEnumerator()));
        while (r.Running) {
            Token? token = r.Step();
            if (token is not null) yield return token;
        }
    }

    static public Tokenizer Yamlite {
        get {
            State start = new();
            start.Add(new("OpenObject"), new Any("{"));
            start.Add(new("CloseObject"), new Any("}"));
            start.Add(new("OpenArray"), new Any("["));
            start.Add(new("CloseArray"), new Any("]"));
            start.Add(new("Colon"), new Any(":"));
            start.Add(new("Comma"), new Any(","));

            State whitespace = new(consume: true);
            start.Add(whitespace, new Any(" \n\r\t"));
            whitespace.Add(whitespace, new Any(" \n\r\t"));

            State commment = new(consume: true);
            start.Add(commment, new Any("#"));
            commment.Add(commment, new Not("\n\r"));

            State innerSingle = new();
            State endSingle = new("SingleValue");
            start.Add(innerSingle, new Any("'"));
            innerSingle.Add(endSingle, new Any("'"));
            innerSingle.Add(innerSingle, new Not("\n\r"));
            endSingle.Add(innerSingle, new Any("'"));

            State innerDouble = new();
            State escapeDouble = new();
            State endDouble = new("DoubleValue");
            start.Add(innerDouble, new Any("\""));
            innerDouble.Add(escapeDouble, new Any("\\"));
            innerDouble.Add(endDouble, new Any("\""));
            innerDouble.Add(innerDouble, new Not("\n\r"));
            escapeDouble.Add(innerDouble, new Not("\n\r"));

            State valState = new("Value");
            State valTail = new();
            start.Add(valState, new Not(" \n\r\t{}[]:,\"\'#&*?|<>=!@\\"));
            valState.Add(valTail, new Any(" \t"));
            valTail.Add(valTail, new Any(" \t"));
            valTail.Add(valState, new Not("\n\r}]:,#"));
            valState.Add(valState, new Not("\n\r}]:,#"));

            return new Tokenizer(start);
        }
    }
}
