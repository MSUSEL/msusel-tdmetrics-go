package abstractor.core;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.fail;

import abstractor.core.constructs.Construct;
import abstractor.core.constructs.Project;
import abstractor.core.diff.Diff;
import abstractor.core.json.JsonHelper;
import abstractor.core.json.JsonNode;
import abstractor.core.json.Jsonable;
import abstractor.core.log.Logger;

public class Testing {
    private Testing() { }
    
    static public Project classFromSource(String ...lines) {
        final Project proj = new Project();
        final Logger log = new Logger(true);
        final Abstractor ab = new Abstractor(log, proj);
        ab.addClassFromSource(lines);
        assertEquals(0, log.errorCount(), "expected zero errors");
        return proj;
    }

    static public void checkJson(Jsonable j, String ...lines) {
        final JsonHelper jh = new JsonHelper();
        final String result = j == null ? "null" : j.toJson(jh).toString(false);
        final String exp = formatJson(lines);
        assertLines(exp, result);
    }

    static public void checkConstruct(Project proj, String key, String ...lines) {
        assertNotNull(proj);
        final Construct con = proj.getConstruct(key);
        if (con == null) {
            fail("unable to find "+key+" in given project");
            return;
        }
        checkJson(con, lines);
    }

    static public String formatJson(String ...lines) {
        try {
            return JsonNode.parse(lines).toString(false);
        } catch(Exception ex) {
            fail(ex);
            return "Error";
        }
    }

    static public void assertLines(String exp, String result) {
        if (!exp.equals(result)) {
            final String diff = String.join("\n\t", new Diff().PlusMinusByLine(exp, result));
            System.out.println("Error: unexpected lines\n\t" + diff);
            fail("unexpected lines (see diff above)");
        }
    }

}
