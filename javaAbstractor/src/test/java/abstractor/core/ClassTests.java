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

        // TODO: Improve
        String exp = String.join("\n",
            "{",
            "  \"basics\": [ \"int\" ],",
            "  \"fields\": [",
            "    { \"name\": \"bar\", \"type\": \"basic1\" }",
            "  ],",
            "  \"language\": \"java\",",
            "  \"locs\": { \"1\": \"unknown\" },",
            "  \"methods\": [",
            "    { \"name\": \"clone\", \"package\": 1, \"receiver\": 1 },",
            "    { \"name\": \"equals\", \"package\": 1, \"receiver\": 1 },",
            "    { \"name\": \"finalize\", \"package\": 1, \"receiver\": 1 },",
            "    { \"name\": \"getClass\", \"package\": 1, \"receiver\": 1 },",
            "    { \"name\": \"hashCode\", \"package\": 1, \"receiver\": 1 },",
            "    { \"name\": \"notify\", \"package\": 1, \"receiver\": 1 },",
            "    { \"name\": \"notifyAll\", \"package\": 1, \"receiver\": 1 },",
            "    { \"name\": \"toString\", \"package\": 1, \"receiver\": 1 },",
            "    { \"name\": \"wait\", \"package\": 1, \"receiver\": 1 }",
            "  ],",
            "  \"objects\": [",
            "    {",
            "      \"data\": 1,",
            "      \"loc\": 1,",
            "      \"methods\": [ 1, 2, 3, 4, 5, 6, 7, 8, 9 ],",
            "      \"name\": \"Foo\",",
            "      \"package\": 1",
            "    }",
            "  ],",
            "  \"packages\": [",
            "    {",
            "      \"methods\": [ 1, 2, 3, 4, 5, 6, 7, 8, 9 ],",
            "      \"name\": \"\",",
            "      \"objects\": [ 1 ]",
            "    }",
            "  ],",
            "  \"structDescs\": [",
            "    {",
            "      \"fields\": [ 1 ]",
            "    }",
            "  ]",
            "}");

        assertEquals(exp, result);
    }
}
