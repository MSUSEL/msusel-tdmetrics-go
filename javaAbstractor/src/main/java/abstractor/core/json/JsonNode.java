package abstractor.core.json;

import abstractor.core.json.parser.*;

public interface JsonNode {
    static public JsonNode parse(String ...lines) throws Exception {
        final String input = String.join("\n", lines);
        final TokenCleaner tok = new TokenCleaner(new Tokenizer(new Scanner(input)));
        final Parser src = new Parser(tok);
        return src.read();
    }

    boolean isEmpty();
}
