package abstractor.core.json;

import java.io.PrintStream;

import abstractor.core.json.parser.*;

public interface JsonNode {
    static public JsonNode parse(String ...lines) throws Exception {
        final String input = String.join("\n", lines);
        final PushBackScanner scan = new PushBackScanner(new Scanner(input));
        final TokenCleaner tok = new TokenCleaner(new Tokenizer(scan));
        final Parser src = new Parser(tok);
        return src.read();
    }

    boolean isEmpty();
    void toString(PrintStream sb, boolean minimize, String indent);
    String toString(boolean minimize);
}
