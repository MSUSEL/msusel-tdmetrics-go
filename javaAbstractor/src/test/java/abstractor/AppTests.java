package abstractor;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.fail;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;

import org.junit.jupiter.api.Test;

import abstractor.app.App;
import abstractor.app.Config;
import abstractor.core.Tester;
import abstractor.core.diff.Diff;
import abstractor.core.json.JsonFormat;
import abstractor.core.json.JsonNode;

public class AppTests {

    @Test
    public void test0001() { runApp("test0001"); }

    @Test
    public void test0002() { runApp("test0002"); }

    @Test
    public void test1002() { testClass("test1001", "Foo"); }

    @Test
    public void test1003() { testClass("test1002", "Foo"); }

    @Test
    public void test1004() { testClass("test1003", "Foo"); }

    static private void runApp(String testName) {
        final String testPath = "../testData/java/" + testName;
        final String absFile  = testPath + "/abstraction.yaml";

        final ByteArrayOutputStream buffer = new ByteArrayOutputStream();
        final ByteArrayOutputStream logBuf = new ByteArrayOutputStream();
        final JsonFormat format = JsonFormat.Relaxed();

        final Config cfg  = new Config();
        cfg.input         = testPath;
        cfg.verbose       = true;
        cfg.writeIndices  = false;
        cfg.writeKinds    = false;
        cfg.defaultOut    = new PrintStream(buffer);
        cfg.logOut        = new PrintStream(logBuf);
        cfg.logErr        = new PrintStream(logBuf);
        cfg.defaultFormat = format;

        try {
            final boolean success = App.run(cfg);
            if (!success) {
                System.out.println(logBuf.toString());
                fail("App.run returned false.");
            }
        } catch(Exception ex) {
            System.out.println(logBuf.toString());
            fail(ex);
        }

        final JsonNode expJson = assertDoesNotThrow(() -> JsonNode.parseFile(absFile));
        final String exp = format.format(expJson);
        final String result = buffer.toString().trim();
        if (!exp.equals(result)) {
            System.out.println(logBuf.toString());
            final String diff = String.join("\n\t", new Diff().PlusMinusByLine(exp, result));
            System.out.println("Error: unexpected lines\n\t" + diff);
            fail("unexpected lines (see diff above)");
        }
    }

    static private void testClass(String testName, String className) {
        final String testPath = "../testData/java/" + testName;
        final String absFile  = testPath + "/abstraction.yaml";

        final Tester t = new Tester(4);
        t.addClassFromFile(testPath+"/"+className+".java");
        t.checkProjectWithFile(absFile);
    }
}
