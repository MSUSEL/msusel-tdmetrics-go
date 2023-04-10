import json.JsonList;
import json.JsonMap;
import json.JsonObj;
import json.Parser;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

public class JsonTests {

    private void checkParse(String json, Object expObj) {
        try {
            JsonObj obj = Parser.parse(json);
            JsonObj exp = JsonObj.convert(expObj);
            assertEquals(exp, obj);
        } catch (Exception e) {
            fail(e);
        }
    }

    private void checkParseError(String json, String expErr) {
        try {
            JsonObj obj = Parser.parse(json);
            fail("Expected an error but got "+obj);
        } catch (Exception e) {
            assertEquals(expErr, e.getMessage());
        }
    }

    @Test
    public void parseIdentifiersTests() {
        checkParse("", null);
        checkParse("true", true);
        checkParse("false", false);
        checkParse("null", null);
    }

    @Test
    public void parseIntegerTests() {
        checkParse("0", 0);
        checkParse("1", 1);
        checkParse("-1", -1);
        checkParse("123", 123);
        checkParse("-123", -123);
    }

    @Test
    public void parseNumberTests() {
        checkParse("0.0", 0.0);
        checkParse("1.0", 1.0);
        checkParse("-1.0", -1.0);
        checkParse("1.23", 1.23);
        checkParse("-1.23", -1.23);

        checkParse("0e1", 0e1);
        checkParse("1e1", 1e1);
        checkParse("-1e1", -1e1);
        checkParse("1e12", 1e12);
        checkParse("-1e12", -1e12);

        checkParse("0e+1", 0e+1);
        checkParse("1e+1", 1e+1);
        checkParse("-1e+1", -1e+1);
        checkParse("1e+12", 1e+12);
        checkParse("-1e+12", -1e+12);

        checkParse("0e-1", 0e-1);
        checkParse("1e-1", 1e-1);
        checkParse("-1e-1", -1e-1);
        checkParse("1e-12", 1e-12);
        checkParse("-1e-12", -1e-12);

        checkParse("1.23e4", 1.23e4);
        checkParse("1.23e-4", 1.23e-4);
        checkParse("-1.23e4", -1.23e4);
        checkParse("-1.23e-4", -1.23e-4);
    }

    @Test
    public void parseStringTests() {
        checkParse("\"\"", "");
        checkParse("\"Hello\"", "Hello");
        checkParse("\"Hello\\nWorld\"", "Hello\nWorld");
        checkParse("\"A\\bB\\fC\\nD\\rE\\tF\"", "A\bB\fC\nD\rE\tF");
        checkParse("\"\\u0020\\u0042\\u007E\\u007e\\u039b\"", " B~~Λ");
        checkParse("\"\\\"'\\\"/\\/\"", "\"'\"//");
    }

    @Test
    public void parseArrayTests() {
        checkParse("[]", new JsonList());
        checkParse(" [  ] ", new JsonList());
        checkParse("\t[\n]\t", new JsonList());
        checkParse("[1]", new JsonList().with(1));
        checkParse("[1, 2]", new JsonList().with(1).with(2));
        checkParse("[true, false, null, \"Hello\", 42, 3.14]",
            new JsonList().with(true).with(false).with(null).with("Hello").with(42).with(3.14));
        checkParse("[[1, 2], [3, 4]]",
            new JsonList().with(new JsonList().with(1).with(2)).with(new JsonList().with(3).with(4)));
    }

    @Test
    public void parseMapTests() {
        checkParse("{}", new JsonMap());
        checkParse(" {  } ", new JsonMap());
        checkParse("{\"Hello\": \"World\"}", new JsonMap().with("Hello", "World"));
        checkParse("{\"Hello\": 1}", new JsonMap().with("Hello", 1));
        checkParse("{\"A\": true, \"B\": null}", new JsonMap().with("A", true).with("B", null));
    }

    @Test
    public void parseErrorsTests() {
        checkParseError("\"", "Unexpected end of JSON characters.");
        checkParseError("x", "Unexpected character: (offset: 0, line: 1, column: 1, value: x)");
        checkParseError("{x}", "Unexpected character: (offset: 1, line: 1, column: 2, value: x)");
        checkParseError("{\nx}", "Unexpected character: (offset: 2, line: 2, column: 1, value: x)");
        checkParseError("}", "Unexpected token: (type: CloseMap, offset: 0, line: 1, column: 1, value: })");
        checkParseError("{1:1}", "Keys must be strings: (type: Int, offset: 1, line: 1, column: 2, value: 1)");
        checkParseError("\"\\q\"", "Unexpected escape sequence: (offset: 2, line: 1, column: 3, value: q)");
        checkParseError(".0", "Unexpected character: (offset: 0, line: 1, column: 1, value: .)");
        checkParseError("01", "Unexpected token: (type: Int, offset: 1, line: 1, column: 2, value: 1)");
        checkParseError("0.", "Unexpected end of JSON characters.");
        checkParseError("0.0.0", "Unexpected character: (offset: 3, line: 1, column: 4, value: .)");
        checkParseError("1e1.1", "Unexpected character: (offset: 3, line: 1, column: 4, value: .)");
        checkParseError("1e1e1", "Unexpected character: (offset: 3, line: 1, column: 4, value: e)");
        checkParseError("{\"A\":", "Unexpected end of JSON tokens.");
        checkParseError("{0.}", "Unexpected character: (offset: 3, line: 1, column: 4, value: })");
        checkParseError("0e", "Unexpected end of JSON characters.");
    }
}
