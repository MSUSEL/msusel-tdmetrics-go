package abstractor.core.json;

import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Paths;

import abstractor.core.json.parser.*;

public interface JsonNode {

    static public JsonNode parseFile(String path) throws Exception {
        final byte[] data = Files.readAllBytes(Paths.get(path));
        return parse(new String(data, StandardCharsets.UTF_8));
    }

    static public JsonNode parse(String ...lines) throws Exception {
        final String input = String.join("\n", lines);
        final TokenCleaner tok = new TokenCleaner(new Tokenizer(new Scanner(input)));
        final Parser src = new Parser(tok);
        return src.read();
    }

    boolean isEmpty();
}
