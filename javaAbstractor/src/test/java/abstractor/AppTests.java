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
import abstractor.core.json.JsonFormat;
import abstractor.core.json.JsonNode;

public class AppTests {
    
    @Test
    public void test0001() { run("test0001"); }

    static private void run(String testName) {
        final String testPath = "../testData/java/" + testName;
        final String absFile  = testPath + "/abstraction.yaml";

        final ByteArrayOutputStream buffer = new ByteArrayOutputStream();
        final JsonFormat format = JsonFormat.Relaxed();

        final Config cfg  = new Config();
        cfg.input         = testPath;
        cfg.verbose       = true;
        cfg.writeIndices  = false;
        cfg.writeKinds    = false;
        cfg.defaultOut    = new PrintStream(buffer);
        cfg.defaultFormat = format;
        assertDoesNotThrow(() -> {
            assertTrue(App.run(cfg), "App.run returned false.");
        }, "Error running App.run");

        final JsonNode expJson = assertDoesNotThrow(() -> JsonNode.parseFile(absFile));
        final String exp = format.format(expJson);
        assertLines(exp, buffer.toString().trim());
    }

    static private void assertLines(String exp, String result) {
        if (!exp.equals(result)) {
            final String diff = String.join("\n\t", new Diff().PlusMinusByLine(exp, result));
            System.out.println("Error: unexpected lines\n\t" + diff);
            fail("unexpected lines (see diff above)");
        }
    }
}
