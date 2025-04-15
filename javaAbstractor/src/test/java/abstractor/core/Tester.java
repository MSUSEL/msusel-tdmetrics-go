package abstractor.core;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.fail;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Paths;

import abstractor.core.constructs.Construct;
import abstractor.core.constructs.Project;
import abstractor.core.diff.Diff;
import abstractor.core.json.JsonFormat;
import abstractor.core.json.JsonHelper;
import abstractor.core.json.JsonNode;
import abstractor.core.json.Jsonable;
import abstractor.core.log.Logger;

public class Tester {
    
    static public Tester classFromSource(String ...lines) {
        final Tester t = new Tester(4);
        t.addClassFromSource(lines);
        return t;
    }

    static private String getTestMethodName(int depth) {
        final StackTraceElement[] ste = Thread.currentThread().getStackTrace();
        return  ste[depth].getClassName() + "." + ste[depth].getMethodName();
    }

    private final ByteArrayOutputStream buffer;
    private final Logger log;
    private final Project proj;
    private final Abstractor ab;

    public Tester() { this(3); }

    public Tester(int depth) {
        this.buffer = new ByteArrayOutputStream();
        PrintStream ps = new PrintStream(buffer);
        this.log = new Logger(true, ps, ps);

        this.log.log("");
        this.log.log("Testing " + getTestMethodName(depth) + "...");

        this.proj = new Project();
        this.ab = new Abstractor(this.log, this.proj);
    }

    public void printLogs() {
        System.out.println(this.buffer.toString());
        this.buffer.reset();
    }

    public void addClassFromFile(String path) {
        assertDoesNotThrow(() -> { 
            final byte[] data =  Files.readAllBytes(Paths.get(path));
            this.addClassFromSource(new String(data, StandardCharsets.UTF_8));
        });
    }

    public void addClassFromSource(String ...lines) {
        try {
            this.ab.addClassFromSource(lines);
            this.ab.finish();
        } catch (Exception ex) {
            this.printLogs();
            fail(ex);
        }

        if (this.log.errorCount() > 0) {
            this.printLogs();
            fail("expected zero errors");
        }
    }

    public void checkJsonWithFile(Jsonable j, String path) {
        final JsonHelper jh = new JsonHelper();
        final String result = j == null ? "null" : JsonFormat.Relaxed().format(j.toJson(jh));
        final String exp = this.formatJsonFromFile(path);
        this.assertLines(exp, result);
    }
    
    public void checkProjectWithFile(String path) {
        this.checkJsonWithFile(this.proj, path);
    }

    public void checkJson(Jsonable j, String ...lines) {
        final JsonHelper jh = new JsonHelper();
        final String result = j == null ? "null" : JsonFormat.Relaxed().format(j.toJson(jh));
        final String exp = this.formatJson(lines);
        this.assertLines(exp, result);
    }
    
    public void checkProject(String ...lines) {
        this.checkJson(this.proj, lines);
    }

    public void checkConstruct(String key, String ...lines) {
        final Construct con = this.proj.getConstructWithKey(key);
        if (con == null) {
            this.printLogs();
            fail("unable to find "+key+" in given project");
            return;
        }
        this.checkJson(con, lines);
    }

    public String formatJson(String ...lines) {
        try {
            return JsonFormat.Relaxed().format(JsonNode.parse(lines));
        } catch(Exception ex) {
            this.printLogs();
            fail(ex);
            return "This should be unreachable.";
        }
    }

    public String formatJsonFromFile(String path) {
        try {
            return JsonFormat.Relaxed().format(JsonNode.parseFile(path));
        } catch(Exception ex) {
            this.printLogs();
            fail(ex);
            return "This should be unreachable.";
        }
    }

    public void assertLines(String exp, String result) {
        if (!exp.equals(result)) {
            this.printLogs();
            final String diff = String.join("\n\t", new Diff().PlusMinusByLine(exp, result));
            System.out.println("Error: unexpected lines\n\t" + diff);
            fail("unexpected lines (see diff above)");
        }
    }
}
