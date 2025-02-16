package abstractor.core;

import static org.junit.jupiter.api.Assertions.assertEquals;

import org.junit.jupiter.api.Test;

import abstractor.core.constructs.Project;
import abstractor.core.json.JsonHelper;
import abstractor.core.log.Logger;

public class ClassTests {

    @Test
    public void ClassWithOneBasicField() {
        Project proj = new Project();
        Logger log = new Logger(true);
        Abstractor ab = new Abstractor(log, proj);
        ab.addClassFromSource(
            "public class Foo {",
            "  public int bar;",
            "}");

        JsonHelper jh = new JsonHelper();
        String result = proj.toJson(jh).toString();
        String exp = "";

        assertEquals(exp, result);
    }
}
