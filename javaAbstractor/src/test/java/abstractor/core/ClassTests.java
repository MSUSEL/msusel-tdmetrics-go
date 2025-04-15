package abstractor.core;

import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

public class ClassTests {

    @Test
    @Disabled
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
            "    { name: $get,      signature:  9 },", //  1. $get(index int) int
            "    { name: $get,      signature: 10 },", //  2. $get(index int) T
            "    { name: $len,      signature:  2 },", //  3. $len() int
            "    { name: $set,      signature: 11 },", //  4. $set(index int, value int)
            "    { name: $set,      signature: 12 },", //  5. $set(index int, value T)
            "    { name: bar,       signature: 13 },", //  6. bar( ??? ) ???
            "    { name: baz,       signature:  1 },", //  7. baz( ??? ) ???
            "    { name: cat,       signature: 14 },", //  8. cat( ??? ) ???
            "    { name: clone,     signature:  5 },", //  9. clone( ??? ) ???
            "    { name: equals,    signature:  8 },", // 10. equals( ??? ) ???
            "    { name: finalize,  signature:  1 },", // 11. finalize( ??? ) ???
            "    { name: getClass,  signature:  3 },", // 12. getClass( ??? ) ???
            "    { name: hashCode,  signature:  2 },", // 13. hashCode( ??? ) ???
            "    { name: notify,    signature:  1 },", // 14. notify( ??? ) ???
            "    { name: notifyAll, signature:  1 },", // 15. notifyAll( ??? ) ???
            "    { name: toString,  signature:  3 },", // 16. toString( ??? ) ???
            "    { name: wait,      signature:  1 },", // 17. wait( ??? ) ???
            "    { name: wait,      signature:  6 },", // 18. wait( ??? ) ???
            "    { name: wait,      signature:  7 }",  // 19. wait( ??? ) ???
            "  ],",
            "  arguments: [",
            "    {              type: basic1         },", //  1. <unnamed> int
            "    {              type: basic2         },", //  2. <unnamed> ???
            "    {              type: interfaceDesc5 },", //  3. <unnamed> ???
            "    {              type: object1        },", //  4. <unnamed> ???
            "    {              type: object2        },", //  4. <unnamed> ???
            "    {              type: typeParam1     },", //  5. <unnamed> T
            "    { name: arg0,  type: basic3         },", //  6. arg0 ???
            "    { name: arg0,  type: object2        },", //  7. arg0 ???
            "    { name: arg1,  type: basic2         },", //  8. arg1 ???
            "    { name: index, type: basic2         },", //  9. index int
            "    { name: t,     type: interfaceInst1 },", // 10. t $Array<T>
            "    { name: value, type: basic2         },", // 11. value int
            "    { name: value, type: typeParam1     },", // 12. value T
            "    { name: x,     type: basic2         },", // 13. x int
            "    { name: y,     type: basic2         }",  // 14. y int
            "  ],",
            "  basics: [ boolean, int, long ],",
            "  interfaceDecls: [",
            "    {", // 1. $Array<T>
            "      name: $Array, interface: 2,",
            "      typeParams: [ 1 ],",
            "    },",
            "  ],",
            "  interfaceDescs: [",
            "    { abstracts: [ 1, 3, 4 ] },", // 1. interface { $get(index int) int, $len() int, $set(index int, value int) }
            "    { abstracts: [ 2, 3, 5 ] },", // 2. interface { $get(index int) T,   $len() int, $set(index int, value T) }
            "    {", // 3. ???
            "      abstracts: [ 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19 ],",
            "      pin: object1,",
            "    },",
            "    {", // 4. ???
            "      abstracts: [ 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19 ],",
            "      pin: object2,",
            "    },",
            "    { abstracts: [ ] }", // 5. Object
            "  ],",
            "  interfaceInsts: [",
            "    { generic: 1, instanceTypes: [ basic2 ], resolved: 1 }", // 1. $Array<int>
            "  ],",
            "  language: java,",
            "  locs: { 1: unknown },",
            "  methods: [",
            "    { name: Foo,       package: 1, receiver: 1, signature:  4, vis: public        },", //  1. Foo ???
            "    { name: Object,    package: 2, receiver: 2, signature:  5, vis: public        },", //  2. Object ???
            "    { name: bar,       package: 1, receiver: 1, signature: 13, loc: 2, metrics: 1 },", //  3. bar(x int, y int) int
            "    { name: baz,       package: 1, receiver: 1, signature:  1, loc: 5, metrics: 2 },", //  4. baz()
            "    { name: cat,       package: 1, receiver: 1, signature: 14, loc: 8, metrics: 3 },", //  5. cat(t ...$Array<T>)
            "    { name: clone,     package: 2, receiver: 2, signature:  5, vis: protected     },", //  6. clone ???
            "    { name: equals,    package: 2, receiver: 2, signature:  8, vis: public        },", //  7. equals ???
            "    { name: finalize,  package: 2, receiver: 2, signature:  1, vis: protected     },", //  8. finalize ???
            "    { name: getClass,  package: 2, receiver: 2, signature:  3, vis: public        },", //  9. getClass ???
            "    { name: hashCode,  package: 2, receiver: 2, signature:  2, vis: public        },", // 10. hashCode ???
            "    { name: notify,    package: 2, receiver: 2, signature:  1, vis: public        },", // 11. notify ???
            "    { name: notifyAll, package: 2, receiver: 2, signature:  1, vis: public        },", // 12. notifyAll ???
            "    { name: toString,  package: 2, receiver: 2, signature:  3, vis: public        },", // 13. toString ???
            "    { name: wait,      package: 2, receiver: 2, signature:  1, vis: public        },", // 14. wait ???
            "    { name: wait,      package: 2, receiver: 2, signature:  6, vis: public        },", // 15. wait ???
            "    { name: wait,      package: 2, receiver: 2, signature:  7, vis: public        }",  // 16. wait ???
            "  ],",
            "  metrics: [",
            "    { loc: 2, codeCount: 3, complexity: 1, indents: 3, lineCount: 3 },", // 1. metrics for bar
            "    { loc: 5, codeCount: 3, complexity: 1, indents: 3, lineCount: 3 },", // 2. metrics for baz
            "    { loc: 8, codeCount: 3, complexity: 1, indents: 3, lineCount: 3 }",  // 3. metrics for cat
            "  ],",
            "  objects: [",
            "    {", // 1. Foo struct{}
            "      name: Foo, package: 1, data: 1, interface: 3,",
            "      vis: public, loc: 1,",
            "      methods: [ 1, 3, 4, 5 ]",
            "    },",
            "    {", // 2. Object
            "      name: Object, package: 2, data: 1, interface: 4,",
            "      vis: public,",
            "      methods: [ 2, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16 ]",
            "    }",
            "  ],",
            "  packages: [",
            "    {", // 1. main package
            "      methods: [ 1, 3, 4, 5 ],",
            "      objects: [ 1 ]",
            "    },",
            "    {", // 2. build-in package
            "      name: java.lang,",
            "      methods: [ 2, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16 ],",
            "      objects: [ 2 ],",
            "    },",
            "  ],",
            "  signatures: [",
            "    {                                    },", //  1. func()
            "    {                     results: [ 2 ] },", //  2. func() int
            "    {                     results: [ 3 ] },", //  3. func( ??? ) ???
            "    {                     results: [ 4 ] },", //  4. func( ??? ) ???
            "    {                     results: [ 5 ] },", //  5. func( ??? ) ???
            "    { params: [ 7 ]                      },", //  6. func( ??? ) ???
            "    { params: [ 7, 9 ]                   },", //  7. func( ??? ) ???  
            "    { params: [ 8 ],      results: [ 1 ] },", //  8. func(index int) int
            "    { params: [ 10 ],     results: [ 2 ] },", //  9. func(index int) T
            "    { params: [ 10 ],     results: [ 6 ] },", // 10. func( ??? ) ???
            "    { params: [ 10, 12 ]                 },", // 11. func(index int, value int)
            "    { params: [ 10, 13 ]                 },", // 12. func(index int, value T)
            "    { params: [ 14, 15 ], results: [ 2 ] },", // 13. func(x int, y int) int
            "    { params: [ 11 ], variadic: true     }",  // 14. func(t ...$Array<T>)
            "  ],",
            "  structDescs: [",
            "    { }", // 1. struct{}
            "  ],",
            "  typeParams: [",
            "    { name: T, type: interfaceDesc5 }", // 1. T object
            "  ]",
            "}");
    }

    @Test
    @Disabled
    public void ClassTestWithTypeParams() {
        final Tester t = Tester.classFromSource(
            "public class Foo<T extends Object> {",
            "  public T value",
            "}");
        t.checkProject(
            "{",
            "  abstracts: [",
            "    { name: clone,     signature: 5 },", //  1. clone() Object
            "    { name: equals,    signature: 8 },", //  2. equals(arg0 Object) boolean
            "    { name: finalize,  signature: 1 },", //  3. finalize()
            "    { name: getClass,  signature: 3 },", //  4. getClass() any // TODO: Is this really "any"?
            "    { name: hashCode,  signature: 2 },", //  5. hashCode() int
            "    { name: notify,    signature: 1 },", //  6. notify()
            "    { name: notifyAll, signature: 1 },", //  7. notifyAll()
            "    { name: toString,  signature: 3 },", //  8. toString() any // TODO: Shouldn't this be "String"?
            "    { name: wait,      signature: 1 },", //  9. wait()
            "    { name: wait,      signature: 6 },", // 10. wait(arg0 long)
            "    { name: wait,      signature: 7 }",  // 11. wait(arg0 long, arg1 int)
            "  ],",
            "  arguments: [",
            "    {             type: basic1         },", // 1. <unnamed> boolean
            "    {             type: basic2         },", // 2. <unnamed> int
            "    {             type: interfaceDesc3 },", // 3. <unnamed> any
            "    {             type: object1        },", // 4. <unnamed> Foo
            "    {             type: object2        },", // 5. <unnamed> Object
            "    { name: arg0, type: basic3         },", // 6. arg0 long
            "    { name: arg0, type: object2        },", // 7. arg0 Object
            "    { name: arg1, type: basic2         }",  // 8. arg1 int
            "  ],",
            "  basics: [ boolean, int, long ],",
            "  fields: [",
            "    { name: value, type: typeParam1, vis: public }", // 1. value T
            "  ],",
            "  interfaceDescs: [",
            "    {", // 1. interface for Foo 
            "      abstracts: [ 1, 10, 11, 2, 3, 4, 5, 6, 7, 8, 9 ],",
            "      pin: object1",
            "    },",
            "    {", // 2.interface for Object
            "      abstracts: [ 1, 10, 11, 2, 3, 4, 5, 6, 7, 8, 9 ],",
            "      pin: object2",
            "    },",
            "    { abstracts: [ ] }", // 3. any
            "  ],",
            "  language: java,",
            "  locs: { 1: unknown },",
            "  methods: [",
            "    { name: Foo,       package: 1, receiver: 1, signature: 4, vis: public    },", //  1. ???
            "    { name: Object,    package: 2, receiver: 2, signature: 5, vis: public    },", //  2. ???
            "    { name: clone,     package: 2, receiver: 2, signature: 5, vis: protected },", //  3. clone() Object
            "    { name: equals,    package: 2, receiver: 2, signature: 8, vis: public    },", //  4. equals(arg0 Object) boolean
            "    { name: finalize,  package: 2, receiver: 2, signature: 1, vis: protected },", //  5. finalize()
            "    { name: getClass,  package: 2, receiver: 2, signature: 3, vis: public    },", //  6. getClass() any // TODO: Is this really "any"?
            "    { name: hashCode,  package: 2, receiver: 2, signature: 2, vis: public    },", //  7. hashCode() int
            "    { name: notify,    package: 2, receiver: 2, signature: 1, vis: public    },", //  8. notify()
            "    { name: notifyAll, package: 2, receiver: 2, signature: 1, vis: public    },", //  9. notifyAll()
            "    { name: toString,  package: 2, receiver: 2, signature: 3, vis: public    },", // 10. toString() any // TODO: Shouldn't this be "String"?
            "    { name: wait,      package: 2, receiver: 2, signature: 1, vis: public    },", // 11. wait()
            "    { name: wait,      package: 2, receiver: 2, signature: 6, vis: public    },", // 12. wait(arg0 long)
            "    { name: wait,      package: 2, receiver: 2, signature: 7, vis: public    }",  // 13. wait(arg0 long, arg1 int)
            "  ],",
            "  objects: [",
            "    {", // 1. Foo<T> { value T }
            "      name: Foo, package: 1, data: 2, interface: 1,",
            "      vis: public, loc: 1,",
            "      methods: [ 1 ],",
            "      typeParams: [ 1 ]",
            "    },",
            "    {", // 2. Object { }
            "      name: Object, package: 2, data: 1, interface: 2,",
            "      methods: [ 10, 11, 12, 13, 2, 3, 4, 5, 6, 7, 8, 9 ],",
            "      vis: public",
            "    }",
            "  ],",
            "  packages: [",
            "    {", // 1. main package
            "      methods: [ 1 ],",
            "      objects: [ 1 ]",
            "    },",
            "    {", // 2. built-in package
            "      methods: [ 10, 11, 12, 13, 2, 3, 4, 5, 6, 7, 8, 9 ],",
            "      name: java.lang,",
            "      objects: [ 2 ]",
            "    }",
            "  ],",
            "  signatures: [",
            "    {                               },", // 1. func()
            "    {                results: [ 2 ] },", // 2. func() int
            "    {                results: [ 3 ] },", // 3. func() any
            "    {                results: [ 4 ] },", // 4. func() Foo
            "    {                results: [ 5 ] },", // 5. func() Object
            "    { params: [ 6 ]                 },", // 6. func(arg0 long)
            "    { params: [ 6, 8 ]              },", // 7. func(arg0 long, arg1 int)
            "    { params: [ 7 ], results: [ 1 ] }",  // 8. func(arg0 Object) boolean
            "  ],",
            "  structDescs: [",
            "    { },",              // 1. struct { }
            "    { fields: [ 1 ] }", // 2. struct { value T }
            "  ],",
            "  typeParams: [",
            "    { name: T, type: interfaceDesc3 }", // T any
            "  ]",
            "}");
    }
}
