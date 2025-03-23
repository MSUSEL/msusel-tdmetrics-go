package abstractor.core;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Disabled;

public class ClassTests {

    @Test
    public void ClassTestWithFields() {
        final Tester t = Tester.classFromSource(
            "public class Foo {",
            "  public int bar;",
            // TODO: add once we have interfaceDesc. "  public int[] baz;",
            "}");
        t.checkProject(
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
        final Tester t = Tester.classFromSource(
            "public class Foo<T extends Object> {",
            "  public T value",
            "}");
        t.checkProject(
            "{",
            "}");
    }

    @Test
    public void ClassTestWithMethods() {
        final Tester t = Tester.classFromSource(
            "public class Foo {",
            "  int bar(int x, int y) {",
            "    return x + y*2;",
            "  }",
            "  void baz() {",
            "    System.out.println(\"Baz\");",
            "  }",
            "  void cat(int ...t) {",
            "    System.out.println(t)",
            "  }",
            "}");
        t.checkProject(
            "{",
            "  abstracts: [",
            "    { name: $get, signature: 3 },", // 1.
            "    { name: $get, signature: 4 },", // 2.
            "    { name: $len, signature: 2 },", // 3.
            "    { name: $set, signature: 5 },", // 4.
            "    { name: $set, signature: 6 }",  // 5.
            "  ],",
            "  arguments: [",
            "    {              type: basic1         },", // 1.
            "    {              type: typeParam1     },", // 2.
            "    { name: index, type: basic0         },", // 3.
            "    { name: t,     type: interfaceInst1 },", // 4.
            "    { name: value, type: basic1         },", // 5.
            "    { name: value, type: typeParam1     },", // 6.
            "    { name: x,     type: basic1         },", // 7.
            "    { name: y,     type: basic1         }",  // 8.
            "  ],",
            "  basics: [ int ],",
            "  interfaceDecls: [",
            "    {", // 1.
            "      name: $Array, interface: 2,",
            "      typeParams: [ 1 ],",
            "    },",
            "  ],",
            "  interfaceDescs: [",
            "    { abstracts: [ 1, 3, 4 ] },", // 1.
            "    { abstracts: [ 2, 3, 5 ] },", // 2.
            "    { abstracts: [ ] },",         // 3.
            "  ],",
            "  interfaceInsts: [",
            "    { generic: 1, instanceTypes: [ 1 ], resolved: 1 },", // 1.
            "  ],",
            "  language: java,",
            "  locs: { 1: unknown },",
            "  methods: [",
            "    { loc: 2, metrics: 1, name: bar, package: 1, receiver: 1, signature: 7 },", // 1.
            "    { loc: 5, metrics: 2, name: baz, package: 1, receiver: 1, signature: 1 },", // 2.
            "    { loc: 8, metrics: 3, name: cat, package: 1, receiver: 1, signature: 8 }",  // 3.
            "  ],",
            "  metrics: [",
            "    { loc: 2, codeCount: 3, complexity: 1, indents: 3, lineCount: 3 },", // 1.
            "    { loc: 5, codeCount: 3, complexity: 1, indents: 3, lineCount: 3 },", // 2.
            "    { loc: 8, codeCount: 3, complexity: 1, indents: 3, lineCount: 3 }",  // 3.
            "  ],",
            "  objects: [",
            "    {",
            "      data: 1,",
            "      interface: null,",
            "      loc: 1,",
            "      methods: [ 1, 2, 3 ],",
            "      name: Foo,",
            "      package: 1,",
            "      visibility: public",
            "    }",
            "  ],",
            "  packages: [",
            "    {",
            "      methods: [ 1, 2, 3 ],",
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
