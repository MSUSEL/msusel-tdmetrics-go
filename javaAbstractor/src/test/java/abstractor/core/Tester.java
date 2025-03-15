package abstractor.core;

import static org.junit.jupiter.api.Assertions.fail;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;

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
        Tester t = new Tester();
        t.addClassFromSource(lines);
        return t;
    }

    static private String getTestMethodName() {
        final StackTraceElement[] ste = Thread.currentThread().getStackTrace();
        final int depth = 3;
        return  ste[depth].getClassName() + "." + ste[depth].getMethodName();
    }

    private final ByteArrayOutputStream buffer;
    private final Logger log;
    private final Project proj;
    private final Abstractor ab;

    public Tester() {
        this.buffer = new ByteArrayOutputStream();
        PrintStream ps = new PrintStream(buffer);
        this.log = new Logger(true, ps, ps);

        this.log.log("");
        this.log.log("Testing " + getTestMethodName() + "...");

        this.proj = new Project();
        this.ab = new Abstractor(this.log, this.proj);
    }

    private void dumpBuffer() {
        System.out.println(this.buffer.toString());
        this.buffer.reset();
    }

    public void addClassFromSource(String ...lines) {
        this.ab.addClassFromSource(lines);
        if (this.log.errorCount() > 0) {
            this.dumpBuffer();
            fail("expected zero errors");
        }
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
        final Construct con = this.proj.getConstruct(key);
        if (con == null) {
            this.dumpBuffer();
            fail("unable to find "+key+" in given project");
            return;
        }
        this.checkJson(con, lines);
    }

    public String formatJson(String ...lines) {
        try {
            return JsonFormat.Relaxed().format(JsonNode.parse(lines));
        } catch(Exception ex) {
            this.dumpBuffer();
            fail(ex);
            return "This should be unreachable.";
        }
    }

    public void assertLines(String exp, String result) {
        if (!exp.equals(result)) {
            this.dumpBuffer();
            final String diff = String.join("\n\t", new Diff().PlusMinusByLine(exp, result));
            System.out.println("Error: unexpected lines\n\t" + diff);
            fail("unexpected lines (see diff above)");
        }
    }
}
