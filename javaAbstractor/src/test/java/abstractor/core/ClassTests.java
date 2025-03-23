package abstractor.core;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Disabled;

public class ClassTests {

    @Test
    public void ClassTestWithFields() {
        final Tester t = Tester.classFromSource(
            "public class Foo {",
            "  public int bar;",
            "  public int[] baz;",
            "}");
        t.checkProject(
            "{",
            "  abstracts: [",
            "    { name: $get, signature: 2 },", // 1. $get(index int) int
            "    { name: $get, signature: 3 },", // 2. $get(index int) T
            "    { name: $len, signature: 1 },", // 3. $len() int
            "    { name: $set, signature: 4 },", // 4. $set(index int, value int)
            "    { name: $set, signature: 5 }",  // 5. $set(index int, value T)
            "  ],",
            "  arguments: [",
            "    {              type: basic1     },", // 1. <unnamed> int
            "    {              type: typeParam1 },", // 2. <unnamed> T
            "    { name: index, type: basic1     },", // 3. index int
            "    { name: value, type: basic1     },", // 4. value int
            "    { name: value, type: typeParam1 }",  // 5. value T
            "  ],",
            "  basics: [ int ],",
            "  fields: [",
            "    { name: bar, type: basic1,         visibility: public },", // 1. bar int
            "    { name: baz, type: interfaceInst1, visibility: public }",  // 2. baz $Array<int>
            "  ],",
            "  interfaceDecls: [",
            "    { name: $Array, interface: 2, typeParams: [ 1 ] }", // 1. $Array<T>
            "  ],",
            "  interfaceDescs: [",
            "    { abstracts: [ 1, 3, 4 ] },", // 1. interface { $get(index int) int, $len() int, $set(index int, value int) }
            "    { abstracts: [ 2, 3, 5 ] },", // 2. interface { $get(index int) T,   $len() int, $set(index int, value T) }
            "    { abstracts: [ ] }",          // 3. interface Object
            "  ],",
            "  interfaceInsts: [",
            "    { generic: 1, instanceTypes: [ basic1 ], resolved: 1 },", // 1. $Array<int>
            "  ],",
            "  language: java,",
            "  locs: { 1: unknown },",
            "  objects: [",
            "    {", // 1. Foo { bar int, baz $Array<int> }
            "      name: Foo, package: 1, data: 1, interface: null,", // TODO: Fix null
            "      loc: 1, visibility: public,",
            "    }",
            "  ],",
            "  packages: [",
            "    { objects: [ 1 ] }", // 1. main package
            "  ],",
            "  signatures: [",
            "    { results: [ 1 ] },",                // 1. func() int
            "    { params: [ 3 ], results: [ 1 ] },", // 2. func(index int) int
            "    { params: [ 3 ], results: [ 2 ] },", // 3. func(index int) T
            "    { params: [ 3, 4 ] },",              // 4. func(index int, value int)
            "    { params: [ 3, 5 ] }",               // 5. func(index int, value T)
            "  ],",
            "  structDescs: [",
            "    { fields: [ 1, 2 ] }", // 1. struct { bar int, baz $Array<int> }
            "  ],",
            "  typeParams: [",
            "    { name: T, type: interfaceDesc3 }", // 1. T object
            "  ]",
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
            "    { name: $get, signature: 3 },", // 1. $get(index int) int
            "    { name: $get, signature: 4 },", // 2. $get(index int) T
            "    { name: $len, signature: 2 },", // 3. $len() int
            "    { name: $set, signature: 5 },", // 4. $set(index int, value int)
            "    { name: $set, signature: 6 }",  // 5. $set(index int, value T)
            "  ],",
            "  arguments: [",
            "    {              type: basic1         },", // 1. <unnamed> int
            "    {              type: typeParam1     },", // 2. <unnamed> T
            "    { name: index, type: basic1         },", // 3. index int
            "    { name: t,     type: interfaceInst1 },", // 4. t $Array<T>
            "    { name: value, type: basic1         },", // 5. value int
            "    { name: value, type: typeParam1     },", // 6. value T
            "    { name: x,     type: basic1         },", // 7. x int
            "    { name: y,     type: basic1         }",  // 8. y int
            "  ],",
            "  basics: [ int ],",
            "  interfaceDecls: [",
            "    {", // 1. $Array<T>
            "      name: $Array, interface: 2,",
            "      typeParams: [ 1 ],",
            "    },",
            "  ],",
            "  interfaceDescs: [",
            "    { abstracts: [ 1, 3, 4 ] },", // 1. interface { $get(index int) int, $len() int, $set(index int, value int) }
            "    { abstracts: [ 2, 3, 5 ] },", // 2. interface { $get(index int) T,   $len() int, $set(index int, value T) }
            "    { abstracts: [ ] }",         // 3. Object
            "  ],",
            "  interfaceInsts: [",
            "    { generic: 1, instanceTypes: [ basic1 ], resolved: 1 }", // 1. $Array<int>
            "  ],",
            "  language: java,",
            "  locs: { 1: unknown },",
            "  methods: [",
            "    { loc: 2, metrics: 1, name: bar, package: 1, receiver: 1, signature: 7 },", // 1. bar func(x int, y int) int { ... }
            "    { loc: 5, metrics: 2, name: baz, package: 1, receiver: 1, signature: 1 },", // 2. baz func()() { ... }
            "    { loc: 8, metrics: 3, name: cat, package: 1, receiver: 1, signature: 8 }",  // 3. cat func(t ...$Array<T>) { ... }
            "  ],",
            "  metrics: [",
            "    { loc: 2, codeCount: 3, complexity: 1, indents: 3, lineCount: 3 },", // 1. metrics for bar
            "    { loc: 5, codeCount: 3, complexity: 1, indents: 3, lineCount: 3 },", // 2. metrics for baz
            "    { loc: 8, codeCount: 3, complexity: 1, indents: 3, lineCount: 3 }",  // 3. metrics for cat
            "  ],",
            "  objects: [",
            "    {", // 1. Foo struct{}
            "      name: Foo, package: 1,",
            "      data: 1, interface: null,", // TODO: Fix interface
            "      visibility: public, loc: 1,",
            "      methods: [ 1, 2, 3 ]",
            "    }",
            "  ],",
            "  packages: [",
            "    {", // 1. main package
            "      methods: [ 1, 2, 3 ],",
            "      objects: [ 1 ]",
            "    }",
            "  ],",
            "  signatures: [",
            "    { },",                                  // 1. func()()
            "    { results: [ 1 ] },",                   // 2. func() int
            "    { params: [ 3 ], results: [ 1 ] },",    // 3. func(index int) int
            "    { params: [ 3 ], results: [ 2 ] },",    // 4. func(index int) T
            "    { params: [ 3, 5 ] },",                 // 5. func(index int, value int)
            "    { params: [ 3, 6 ] },",                 // 6. func(index int, value T)
            "    { params: [ 7, 8 ], results: [ 1 ] },", // 7. func(x int, y int) int
            "    { params: [ 4 ], variadic: true }",     // 8. func(t ...$Array<T>)
            "  ],",
            "  structDescs: [",
            "    { }", // 1. struct{}
            "  ],",
            "  typeParams: [",
            "    { name: T, type: interfaceDesc3 }", // 1. T object
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
}
