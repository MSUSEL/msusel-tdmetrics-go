import json.Parser;
import janis.Janis;
import json.JsonMap;
import json.JsonObj;
import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;

public class JanisTests {

    static private void checkJson(JsonObj obj, String... expLines) {
        try {
            JsonObj expObj = Parser.parse(String.join("\n", expLines));
            assertEquals(expObj, obj);
        } catch (Exception e) {
            fail("Expected JSON is invalid: " + e.getMessage());
        }
    }

    @Test
    public void test001() {
        JsonMap data = Janis.read("./testProjects/test001/src/");
        checkJson(data,
            "{" +
            "   \"methods\": [" +
            "       { \"name\": \"main\", \"receiver\": \"Main\"}," +
            "       { \"name\": \"getName\", \"receiver\": \"Person\"}" +
            "   ]," +
            "   \"packages\": [" +
            "       { \"name\": \"<root>\"}" +
            "   ]," +
            "   \"types\": [" +
            "       { \"name\": \"Main\", \"package\": \"<root>\" }," +
            "       { \"name\": \"Person\", \"package\": \"<root>\" }" +
            "   ]" +
            "}");
    }
}