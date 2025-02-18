package abstractor.core;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertLinesMatch;

import java.util.Arrays;

import org.junit.jupiter.api.Test;

import abstractor.core.constructs.Project;
import abstractor.core.json.JsonHelper;
import abstractor.core.json.JsonNode;
import abstractor.core.log.Logger;

public class ClassTests {

    static private Project classFromSource(String ...lines) {
        final Project proj = new Project();
        final Logger log = new Logger(true);
        final Abstractor ab = new Abstractor(log, proj);
        ab.addClassFromSource(lines);
        assertEquals(0, log.errorCount());
        return proj;
    }

    static private void checkJson(Project proj, String ...lines) {
        assertDoesNotThrow(() -> {
            final JsonHelper jh = new JsonHelper();
            final JsonNode result = proj.toJson(jh);
            final JsonNode exp = JsonNode.parse(lines);
            assertLinesMatch(
                Arrays.asList(exp.toString(false).split("\n")),
                Arrays.asList(result.toString(false).split("\n")));
        });
    }

    @Test
    public void ClassWithOneBasicField() {
        final Project proj = classFromSource(
            "public class Foo {",
            "  public int bar;",
            "}");

        checkJson(proj,
            "{",
            "  basics: [ int ],",
            "  fields: [",
            "    { name: bar, type: basic1 }",
            "  ],",
            "  language: java,",
            "  locs: { 1: unknown },",
            "  methods: [",
            "    { name: clone,     package: 1, receiver: 1 },",
            "    { name: equals,    package: 1, receiver: 1 },",
            "    { name: finalize,  package: 1, receiver: 1 },",
            "    { name: getClass,  package: 1, receiver: 1 },",
            "    { name: hashCode,  package: 1, receiver: 1 },",
            "    { name: notify,    package: 1, receiver: 1 },",
            "    { name: notifyAll, package: 1, receiver: 1 },",
            "    { name: toString,  package: 1, receiver: 1 },",
            "    { name: wait,      package: 1, receiver: 1 }",
            "  ],",
            "  objects: [",
            "    {",
            "      name: Foo, package: 1, data: 1, loc: 1,",
            "      methods: [ 1, 2, 3, 4, 5, 6, 7, 8, 9 ]",
            "    }",
            "  ],",
            "  packages: [",
            "    {",
            "      name: \"\",",
            "      methods: [ 1, 2, 3, 4, 5, 6, 7, 8, 9 ],",
            "      objects: [ 1 ]",
            "    }",
            "  ],",
            "  structDescs: [",
            "    { fields: [ 1 ] }",
            "  ]",
            "}");
    }
}
