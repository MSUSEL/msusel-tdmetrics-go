package abstractor.core;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.fail;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Disabled;

import abstractor.core.constructs.Project;
import abstractor.core.diff.Diff;
import abstractor.core.json.JsonHelper;
import abstractor.core.json.JsonNode;
import abstractor.core.log.Logger;

public class ClassTests {

    static private Project classFromSource(String ...lines) {
        final Project proj = new Project();
        final Logger log = new Logger(true);
        final Abstractor ab = new Abstractor(log, proj);
        ab.addClassFromSource(lines);
        assertEquals(0, log.errorCount(), "expected zero errors");
        return proj;
    }

    static private void checkJson(Project proj, String ...lines) {
        assertDoesNotThrow(() -> {
            final JsonHelper jh = new JsonHelper();
            final String result = proj.toJson(jh).toString(false);
            final String exp = JsonNode.parse(lines).toString(false);
            assertLines(exp, result);
        });
    }

    static private void assertLines(String exp, String result) {
        if (!exp.equals(result)) {
            final String diff = String.join("\n\t", new Diff().PlusMinusByLine(exp, result));
            System.out.println("Error: unexpected lines\n\t" + diff);
            fail("unexpected lines (see diff above)");
        }
    }

    @Test
    public void ClassTestWithFields() {
        final Project proj = classFromSource(
            "public class Foo {",
            "  public int bar;",
            // TODO: add once we have interfaceDesc. "  public int[] baz;",
            "}");

        checkJson(proj,
            "{",
            "  basics: [ int ],",
            "  fields: [",
            "    { name: bar, type: basic1, visibility: public }",
            "  ],",
            "  language: java,",
            "  locs: { 1: unknown },",
            "  objects: [",
            "    {",
            "      name: Foo, package: 1, data: 1, interface: null,", // TODO: Fix null
            "      loc: 1, visibility: public,",
            "    }",
            "  ],",
            "  packages: [",
            "    {",
            "      objects: [ 1 ]",
            "    }",
            "  ],",
            "  structDescs: [",
            "    { fields: [ 1 ] }",
            "  ]",
            "}");
    }

    @Test
    @Disabled // TODO: Enable once type params work
    public void ClassTestWithTypeParams() {
        final Project proj = classFromSource(
            "public class Foo<T extends Object> {",
            "  public T value",
            "}");

        checkJson(proj,
            "{",
            "}");
    }

    @Test
    public void ClassTestWithMethods() {
        final Project proj = classFromSource(
            "public class Foo {",
            "  int bar(int x, int y) {",
            "    return x + y*2;",
            "  }",
            "  void baz() {",
            "    System.out.println(\"Baz\");",
            "  }",
            "}");

        checkJson(proj,
            "{",
            "  arguments: [",
            "    { type: basic1 },",
            "    { name: x, type: basic1 },",
            "    { name: y, type: basic1 }",
            "  ],",
            "  basics: [ int ],",
            "  language: java,",
            "  locs: { 1: unknown },",
            "  methods: [",
            "    { name: bar, package: 1, receiver: 1, signature: 2 },",
            "    { name: baz, package: 1, receiver: 1, signature: 1 }",
            "  ],", // TODO: Add metrics
            "  objects: [",
            "    {",
            "      data: 1,",
            "      interface: null,",
            "      loc: 1,",
            "      methods: [ 1, 2 ],",
            "      name: Foo,",
            "      package: 1,",
            "      visibility: public",
            "    }",
            "  ],",
            "  packages: [",
            "    {",
            "      methods: [ 1, 2 ],",
            "      objects: [ 1 ]",
            "    }",
            "  ],",
            "  signatures: [",
            "    { },",
            "    {",
            "      params: [ 2, 3 ],",
            "      results: [ 1 ]",
            "    }",
            "  ],",
            "  structDescs: [",
            "    { }",
            "  ]",
            "}");
    }
}
