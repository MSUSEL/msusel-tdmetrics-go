import json.Parser;
import janis.Janis;
import json.JsonMap;
import json.JsonObj;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

public class JanisTests {

    static private void checkJson(JsonObj obj, String... expLines) {
        JsonObj expObj = null;
        try {
            expObj = Parser.parse(String.join("\n", expLines).replace("'", "\""));
        } catch (Exception e) {
            fail("Expected JSON is invalid: " + e.getMessage());
        }

        try {
            obj.removeOmitted();
            expObj.assertCompare(obj);
        } catch (Exception e) {
            System.out.println(obj.toString());
            fail(e.getMessage());
        }
    }

    @Test
    public void test001() {
        JsonMap data = Janis.read("./testProjects/test001/src/");
        checkJson(data,
            "{",
            "   'methods': [",
            "      {",
            "         'name':       'main',",
            "         'parameters': [[ 'String' ]],",
            "         'receiver':   'Main',",
            "         'returns':    'void'",
            "      }, {",
            "         'name':     'getName',",
            "         'receiver': 'Person',",
            "         'returns':  'String'",
            "      }",
            "   ],",
            "   'packages': [",
            "       { 'name': '<root>' }",
            "   ],",
            "   'types': [",
            "       { 'name': 'Main',   'package': '<root>' },",
            "       { 'name': 'Person', 'package': '<root>' }",
            "   ]",
            "}");
    }
}
