package abstractor;

import java.io.*;

import org.junit.jupiter.api.*;

import abstractor.app.App;
import abstractor.app.Config;
import abstractor.core.Tester;
import abstractor.core.json.JsonFormat;
import abstractor.core.json.JsonNode;

public class AppTests {

    @Test
    public void test0001() throws Exception { runApp("test0001"); }

    @Test
    public void test0002() throws Exception { runApp("test0002"); }

    @Test
    public void test1001() throws Exception { testClass("test1001", "Foo"); }

    @Test
    public void test1002() throws Exception { testClass("test1002", "Foo"); }

    @Test
    public void test1003() throws Exception { testClass("test1003", "Foo"); }

    @Test
    public void test1004() throws Exception { testClass("test1004", "Foo"); }

    @Test
    public void test1005() throws Exception { testClass("test1005", "Foo"); }

    @Test
    public void test1006() throws Exception { testClass("test1006", "Foo"); }

    @Test
    public void test1007() throws Exception { testClass("test1007", "Foo"); }

    static private void runApp(String testName) throws Exception {
        final String testPath = "../testData/java/" + testName;
        final String absFile  = testPath + "/abstraction.yaml";
        final String diffFile = testPath + "/abstraction.diff";
        final String logFile  = testPath + "/abstraction.log";

        final ByteArrayOutputStream buffer = new ByteArrayOutputStream();
        final ByteArrayOutputStream logBuf = new ByteArrayOutputStream();
        final JsonFormat format = JsonFormat.Relaxed();

        final Config cfg  = new Config();
        cfg.input         = testPath;
        cfg.verbose       = true;
        cfg.writeIndices  = false;
        cfg.writeKinds    = false;
        cfg.writeRefs     = false;
        cfg.defaultOut    = new PrintStream(buffer);
        cfg.logOut        = new PrintStream(logBuf);
        cfg.logErr        = new PrintStream(logBuf);
        cfg.defaultFormat = format;

        try {
            final boolean success = App.run(cfg);
            if (!success) {
                Tester.printLogs(logBuf, logFile);
                Assertions.fail("App.run returned false.");
            }
        } catch(Exception ex) {
            Tester.printLogs(logBuf, logFile);
            Assertions.fail(ex);
        }

        final JsonNode expJson = JsonNode.parseFile(absFile);
        final String exp = format.format(expJson);
        final String result = buffer.toString().trim();
        if (!exp.equals(result)) {
            Tester.printLogs(logBuf, logFile);
            Tester.printDiff(exp, result, diffFile);
            Assertions.fail("unexpected lines (see diff)");
        }
    }

    static private void testClass(String testName, String className) throws Exception {
        final String testPath = "../testData/java/" + testName;
        final String absFile  = testPath + "/abstraction.yaml";
        final String diffFile = testPath + "/abstraction.diff";
        final String logFile  = testPath + "/abstraction.log";

        final Tester t = new Tester(4);
        t.addClassFromFile(testPath+"/"+className+".java");
        t.checkProjectWithFile(absFile, diffFile, logFile);
    }
}
