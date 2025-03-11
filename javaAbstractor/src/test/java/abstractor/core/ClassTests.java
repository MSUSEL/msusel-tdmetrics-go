package abstractor.core;

import static abstractor.core.Testing.*;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Disabled;

import abstractor.core.constructs.Project;

public class ClassTests {

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
            //"  void cat(int ...t) {",
            //"    System.out.println(t)",
            //"  }",
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
            "    { loc: 2, metrics: 1, name: bar, package: 1, receiver: 1, signature: 2 },",
            "    { loc: 5, metrics: 2, name: baz, package: 1, receiver: 1, signature: 1 }",
            "  ],",
            "  metrics: [",
            "    { loc: 2, codeCount: 3, complexity: 1, indents: 3, lineCount: 3 },",
            "    { loc: 5, codeCount: 3, complexity: 1, indents: 3, lineCount: 3 }",
            "  ],",
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
