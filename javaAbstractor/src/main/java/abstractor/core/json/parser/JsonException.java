package abstractor.core.json.parser;

public class JsonException extends Exception {
    public JsonException(String msg) { super(msg); }
    public JsonException(String msg, Exception ex) { super(msg, ex); }
}
