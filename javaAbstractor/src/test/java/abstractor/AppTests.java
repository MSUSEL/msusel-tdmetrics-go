package abstractor;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.junit.jupiter.api.Assertions.fail;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;

import org.junit.jupiter.api.Test;

import abstractor.app.App;
import abstractor.app.Config;
import abstractor.core.diff.Diff;
import abstractor.core.json.JsonNode;

public class AppTests {
    
    @Test
    public void test0001() { run("test0001"); }

    static private void run(String testName) {
        final String path = "../testData/java/" + testName;

        final Config cfg = new Config();
        cfg.input   = path;
        cfg.verbose = true;
        final ByteArrayOutputStream buffer = new ByteArrayOutputStream();
        cfg.defaultOut = new PrintStream(buffer);
        assertTrue(App.run(cfg), "App.run returned false.");

        assertDoesNotThrow(() -> {
            final JsonNode expJson = JsonNode.parseFile(path+"/abstraction.yaml");
            final String exp = expJson.toString();
            assertLines(exp, buffer.toString().trim());
        });
    }

    static public void assertLines(String exp, String result) {
        if (!exp.equals(result)) {
            final String diff = String.join("\n\t", new Diff().PlusMinusByLine(exp, result));
            System.out.println("Error: unexpected lines\n\t" + diff);
            fail("unexpected lines (see diff above)");
        }
    }
}
