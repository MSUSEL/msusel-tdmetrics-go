package json;

import java.util.ArrayDeque;
import java.util.Queue;

public class Parser {

    private record Gram(
        int offset,
        int line,
        int column,
        char value) {

        @Override
        public String toString() {
            return String.valueOf(this.value);
        }
    }

    static private class Scanner {
        private final String json;
        private int offset;
        private int line;
        private int column;
        private final Queue<Gram> queue;

        public Scanner(String json) {
            this.json   = json;
            this.offset = -1;
            this.line   = 1;
            this.column = 0;
            this.queue  = new ArrayDeque<>();
        }

        public void pushback(Gram gram) {
            this.queue.add(gram);
        }

        public Gram step() {
            if (!this.queue.isEmpty()) return this.queue.poll();
            if (this.offset >= json.length()) return null;

            ++this.offset;
            if (this.offset >= json.length()) return null;

            char value = json.charAt(this.offset);
            if (value == '\n') {
                this.column = 0;
                ++this.line;
            } else ++this.column;
            return new Gram(offset, line, column, value);
        }
    }

    private enum TokenType {
        OpenMap,
        CloseMap,
        OpenArray,
        CloseArray,
        Comma,
        Colon,
        Str,
        Int,
        Num,
        True,
        False,
        Null
    }

    private record Token(
       TokenType type,
       int offset,
       int line,
       int column,
       String value) { }

    static private class Tokenizer {
        private final Scanner scanner;
        private final Queue<Token> queue;

        public Tokenizer(String json) {
            this.scanner = new Scanner(json);
            this.queue = new ArrayDeque<>();
        }

        public void pushback(Token token) {
            this.queue.add(token);
        }

        public Token step() throws Exception {
            if (!this.queue.isEmpty()) return this.queue.poll();
            while (true) {
                Gram gram = this.scanner.step();
                if (gram == null) return null;
                switch (gram.value) {
                    case '{': return this.simpleToken(TokenType.OpenMap, gram);
                    case '}': return this.simpleToken(TokenType.CloseMap, gram);
                    case '[': return this.simpleToken(TokenType.OpenArray, gram);
                    case ']': return this.simpleToken(TokenType.CloseArray, gram);
                    case ':': return this.simpleToken(TokenType.Colon, gram);
                    case ',': return this.simpleToken(TokenType.Comma, gram);
                    case '-': return this.readNegDigitStart(gram, new StringBuilder("-"));
                    case '0': return this.readOptionalNum(gram, new StringBuilder("0"));
                    case '1', '2', '3', '4', '5', '6', '7', '8', '9':
                        return this.readDigitTail(gram, new StringBuilder(gram.toString()));
                    case '"': return this.readStr(gram);
                    case 't': return this.readId(gram, "true", TokenType.True);
                    case 'f': return this.readId(gram, "false", TokenType.False);
                    case 'n': return this.readId(gram, "null", TokenType.Null);
                    case ' ', '\n', '\r', '\t': break;
                    default: this.unexpectedGram(gram);
                }
            }
        }

        private Token newToken(TokenType type, Gram start, String value) {
            return new Token(type, start.offset, start.line, start.column, value);
        }

        private Token simpleToken(TokenType type, Gram gram) {
            return newToken(type, gram, gram.toString());
        }

        private void unexpectedGram(Gram gram) throws Exception {
            throw new Exception("Unexpected character: (offset: " + gram.offset +
                ", line: " + gram.line + ", column: " + gram.column + ", value: " + gram.value + ")");
        }

        private void unexpectedEnd() throws Exception {
            throw new Exception("Unexpected end of JSON characters.");
        }

        private Token readNegDigitStart(Gram start, StringBuilder builder) throws Exception {
            Gram gram = this.scanner.step();
            if (gram == null) this.unexpectedEnd();
            switch (gram.value) {
                case '0' -> {
                    builder.append(gram.value);
                    return this.readOptionalNum(start, builder);
                }
                case '1', '2', '3', '4', '5', '6', '7', '8', '9' -> {
                    builder.append(gram.value);
                    return this.readDigitTail(gram, builder);
                }
                default -> {
                    this.unexpectedGram(gram);
                    return null;
                }
            }
        }

        private Token readDigitTail(Gram start, StringBuilder builder) throws Exception {
            while (true) {
                Gram gram = this.scanner.step();
                if (gram == null) return this.newToken(TokenType.Int, start, builder.toString());
                switch (gram.value) {
                    case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9' ->
                        builder.append(gram.value);
                    case '.' -> {
                        builder.append(gram.value);
                        return this.readFactionTail(start, builder);
                    }
                    case 'e', 'E' -> {
                        builder.append(gram.value);
                        return this.readExpStart(start, builder);
                    }
                    default -> {
                        this.scanner.pushback(gram);
                        return this.newToken(TokenType.Int, start, builder.toString());
                    }
                }
            }
        }

        private Token readOptionalNum(Gram start, StringBuilder builder) throws Exception {
            Gram gram = this.scanner.step();
            if (gram == null) return newToken(TokenType.Int, start, builder.toString());
            switch (gram.value) {
                case '.' -> {
                    builder.append(gram.value);
                    return this.readFactionStart(start, builder);
                }
                case 'e', 'E' -> {
                    builder.append(gram.value);
                    return this.readExpStart(start, builder);
                }
                default -> {
                    this.scanner.pushback(gram);
                    return this.newToken(TokenType.Int, start, builder.toString());
                }
            }
        }

        private Token readFactionStart(Gram start, StringBuilder builder) throws Exception {
            Gram gram = this.scanner.step();
            if (gram == null) this.unexpectedEnd();
            if (gram.value >= '0' && gram.value <= '9') {
                builder.append(gram.value);
                return this.readExpTail(start, builder);
            }
            this.unexpectedGram(gram);
            return null;
        }

        private Token readFactionTail(Gram start, StringBuilder builder) throws Exception {
            while (true) {
                Gram gram = this.scanner.step();
                if (gram == null) return this.newToken(TokenType.Num, start, builder.toString());
                switch (gram.value) {
                    case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9' ->
                        builder.append(gram.value);
                    case 'e', 'E' -> {
                        builder.append(gram.value);
                        return this.readExpStart(start, builder);
                    }
                    default -> {
                        this.scanner.pushback(gram);
                        return this.newToken(TokenType.Num, start, builder.toString());
                    }
                }
            }
        }

        private Token readExpStart(Gram start, StringBuilder builder) throws Exception {
            Gram gram = this.scanner.step();
            if (gram == null) this.unexpectedEnd();
            switch (gram.value) {
                case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9' -> {
                    builder.append(gram.value);
                    return this.readExpTail(start, builder);
                }
                case '-', '+' -> {
                    builder.append(gram.value);
                    return this.readExpSign(start, builder);
                }
                default -> {
                    this.unexpectedGram(gram);
                    return null;
                }
            }
        }

        private Token readExpSign(Gram start, StringBuilder builder) throws Exception {
            Gram gram = this.scanner.step();
            if (gram == null) this.unexpectedEnd();
            if (gram.value >= '0' && gram.value <= '9') {
                builder.append(gram.value);
                return this.readExpTail(start, builder);
            }
            this.unexpectedGram(gram);
            return null;
        }

        private Token readExpTail(Gram start, StringBuilder builder) {
            while (true) {
                Gram gram = this.scanner.step();
                if (gram == null) return this.newToken(TokenType.Num, start, builder.toString());
                if (gram.value >= '0' && gram.value <= '9') builder.append(gram.value);
                else {
                    this.scanner.pushback(gram);
                    return this.newToken(TokenType.Num, start, builder.toString());
                }
            }
        }

        private int readHexNibble() throws Exception {
            Gram gram = this.scanner.step();
            if (gram == null) this.unexpectedEnd();
            if (gram.value >= 'a' && gram.value <= 'f') return (int)gram.value - (int)'a' + 10;
            if (gram.value >= 'A' && gram.value <= 'F') return (int)gram.value - (int)'A' + 10;
            if (gram.value >= '0' && gram.value <= '9') return (int)gram.value - (int)'0';
            this.unexpectedGram(gram);
            return 0;
        }

        private Token readStr(Gram start) throws Exception {
            StringBuilder builder = new StringBuilder();
            while (true) {
                Gram gram = this.scanner.step();
                if (gram == null) this.unexpectedEnd();
                if (gram.value == '\\') {
                    gram = this.scanner.step();
                    if (gram == null) this.unexpectedEnd();
                    switch (gram.value) {
                        case '"', '\\', '/' -> builder.append(gram.value);
                        case 'b' -> builder.append('\b');
                        case 'f' -> builder.append('\f');
                        case 'n' -> builder.append('\n');
                        case 'r' -> builder.append('\r');
                        case 't' -> builder.append('\t');
                        case 'u' -> {
                            int code = (this.readHexNibble() << 12) |
                                (this.readHexNibble() << 8) |
                                (this.readHexNibble() << 4) |
                                this.readHexNibble();
                            builder.append((char) code);
                        }
                        default -> throw new Exception("Unexpected escape sequence: (offset: " + gram.offset +
                            ", line: " + gram.line + ", column: " + gram.column + ", value: " + gram.value + ")");
                    }
                } else if (gram.value == '"')
                    return newToken(TokenType.Str, start, builder.toString());
                else builder.append(gram.value);
            }
        }

        private Token readId(Gram start, String expected, TokenType type) throws Exception {
            for (int i = 1; i < expected.length(); ++i) {
                Gram gram = this.scanner.step();
                if (gram == null) this.unexpectedEnd();
                if (gram.value != expected.charAt(i)) this.unexpectedGram(gram);
            }
            return newToken(type, start, expected);
        }
    }

    static public JsonObj parse(String json) throws Exception {
        Tokenizer tokenizer = new Tokenizer(json);
        JsonObj obj = readObj(tokenizer);
        Token tail = tokenizer.step();
        if (tail != null) unexpectedToken(tail);
        return obj;
    }

    static private void unexpectedToken(Token token) throws Exception {
        throw new Exception("Unexpected token: (type: " + token.type + ", offset: " + token.offset +
            ", line: " + token.line + ", column: " + token.column + ", value: " + token.value + ")");
    }

    static private void unexpectedEnd() throws Exception {
        throw new Exception("Unexpected end of JSON tokens.");
    }

    static private JsonObj readObj(Tokenizer tokenizer) throws Exception {
        Token token = tokenizer.step();
        if (token == null) return null;
        switch (token.type) {
            case OpenMap   -> { return readMap(tokenizer); }
            case OpenArray -> { return readArray(tokenizer); }
            case Str   -> { return new JsonStr(token.value); }
            case Int   -> { return new JsonInt(Integer.parseInt(token.value)); }
            case Num   -> { return new JsonNum(Double.parseDouble(token.value)); }
            case True  -> { return new JsonBool(true); }
            case False -> { return new JsonBool(false); }
            case Null  -> { return null; }
            default -> {
                unexpectedToken(token);
                return null;
            }
        }
    }

    static private JsonMap readMap(Tokenizer tokenizer) throws Exception {
        JsonMap map = new JsonMap();
        Token token = tokenizer.step();
        if (token == null) unexpectedEnd();
        if (token.type == TokenType.CloseMap) return map;
        tokenizer.pushback(token);

        while (true) {
            token = tokenizer.step();
            if (token == null) unexpectedEnd();
            if (token.type != TokenType.Str)
                throw new Exception("Keys must be strings: (type: " + token.type + ", offset: " + token.offset +
                    ", line: " + token.line + ", column: " + token.column + ", value: " + token.value + ")");
            String key = token.value;

            token = tokenizer.step();
            if (token == null) unexpectedEnd();
            if (token.type != TokenType.Colon) unexpectedToken(token);

            JsonObj value = readObj(tokenizer);
            map.put(key, value);

            token = tokenizer.step();
            if (token == null) unexpectedEnd();
            if (token.type == TokenType.CloseMap) return map;
            if (token.type != TokenType.Comma) unexpectedToken(token);
        }
    }

    static private JsonList readArray(Tokenizer tokenizer) throws Exception {
        JsonList list = new JsonList();
        Token token = tokenizer.step();
        if (token == null) unexpectedEnd();
        if (token.type == TokenType.CloseArray) return list;
        tokenizer.pushback(token);

        while (true) {
            JsonObj value = readObj(tokenizer);
            list.add(value);

            token = tokenizer.step();
            if (token == null) unexpectedEnd();
            if (token.type == TokenType.CloseArray) return list;
            if (token.type != TokenType.Comma) unexpectedToken(token);
        }
    }
}
