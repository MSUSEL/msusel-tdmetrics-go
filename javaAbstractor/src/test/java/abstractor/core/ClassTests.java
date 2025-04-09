package abstractor.core;

import org.junit.jupiter.api.Test;

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
            "    { name: $get,      signature:  8 },", //  1. $get(index int) int
            "    { name: $get,      signature:  9 },", //  2. $get(index int) T
            "    { name: $len,      signature:  2 },", //  3. $len() int
            "    { name: $set,      signature: 10 },", //  4. $set(index int, value int)
            "    { name: $set,      signature: 11 },", //  5. $set(index int, value T)
            "    { name: clone,     signature:  4 },", //  6. clone() Object
            "    { name: equals,    signature:  7 },", //  7. equals(arg0 Object) boolean
            "    { name: finalize,  signature:  1 },", //  8. finalize()
            "    { name: getClass,  signature:  3 },", //  9. getClass() any
            "    { name: hashCode,  signature:  2 },", // 10. hashCode() int
            "    { name: notify,    signature:  1 },", // 11. notify()
            "    { name: notifyAll, signature:  1 },", // 12. notifyAll()
            "    { name: toString,  signature:  3 },", // 13. toString() any
            "    { name: wait,      signature:  1 },", // 14. wait()
            "    { name: wait,      signature:  5 },", // 15. wait(arg0 long)
            "    { name: wait,      signature:  6 }",  // 16. wait(arg0 long, arg1 int)
            "  ],",
            "  arguments: [",
            "    {              type: basic1         },", //  1. <unnamed> boolean
            "    {              type: basic2         },", //  2. <unnamed> int
            "    {              type: interfaceDesc5 },", //  3. <unnamed> any
            "    {              type: object2        },", //  4. <unnamed> Object
            "    {              type: typeParam1     },", //  5. <unnamed> T
            "    { name: arg0,  type: basic3         },", //  6. arg0 long
            "    { name: arg0,  type: object2        },", //  7. arg0 Object
            "    { name: arg1,  type: basic2         },", //  8. arg1 int
            "    { name: index, type: basic2         },", //  9. index int
            "    { name: value, type: basic2         },", // 10. value int
            "    { name: value, type: typeParam1     }",  // 11. value T
            "  ],",
            "  basics: [ boolean, int, long ],",
            "  fields: [",
            "    { name: bar, type: basic2,         vis: public },", // 1. bar int
            "    { name: baz, type: interfaceInst1, vis: public }",  // 2. baz $Array<int>
            "  ],",
            "  interfaceDecls: [",
            "    { name: $Array, interface: 2, typeParams: [ 1 ] }", // 1. $Array<T>
            "  ],",
            "  interfaceDescs: [",
            "    {", // 1. interface { $get(index int) int, $len() int, $set(index int, value int) }
            "      abstracts: [ 1, 3, 4 ]",
            "    },",
            "    {", // 2. interface { $get(index int) T, $len() int, $set(index int, value T) }
            "      abstracts: [ 2, 3, 5 ]",
            "    },",
            "    {",  // 3. interface for Foo // TODO: fix with inheritance
            "      abstracts: [ 10, 11, 12, 13, 14, 15, 16, 6, 7, 8, 9 ],",
            "      pin: object1,",
            "    },",
            "    {", // 4. interface for Object
            "      abstracts: [ 10, 11, 12, 13, 14, 15, 16, 6, 7, 8, 9 ],",
            "      pin: object2,",
            "    },",
            "    {", // 5. any
            "      abstracts: [ ]",
            "    }",
            "  ],",
            "  interfaceInsts: [",
            "    { generic: 1, instanceTypes: [ basic2 ], resolved: 1 },", // 1. $Array<int>
            "  ],",
            "  language: java,",
            "  locs: { 1: unknown },",
            "  methods: [",
            "    { name: clone,     package: 2, receiver: 2, signature: 4, vis: protected },",
            "    { name: equals,    package: 2, receiver: 2, signature: 7, vis: public },",
            "    { name: finalize,  package: 2, receiver: 2, signature: 1, vis: protected },",
            "    { name: getClass,  package: 2, receiver: 2, signature: 3, vis: public },",
            "    { name: hashCode,  package: 2, receiver: 2, signature: 2, vis: public },",
            "    { name: notify,    package: 2, receiver: 2, signature: 1, vis: public },",
            "    { name: notifyAll, package: 2, receiver: 2, signature: 1, vis: public },",
            "    { name: toString,  package: 2, receiver: 2, signature: 3, vis: public },",
            "    { name: wait,      package: 2, receiver: 2, signature: 1, vis: public },",
            "    { name: wait,      package: 2, receiver: 2, signature: 5, vis: public },",
            "    { name: wait,      package: 2, receiver: 2, signature: 6, vis: public },",
            "  ],",
            "  objects: [",
            "    {", // 1. Foo { bar int, baz $Array<int> }
            "      name: Foo, package: 1, data: 2, interface: 3,",
            "      loc: 1, vis: public,",
            "    },",
            "    {", // 2. Object { }
            "      data: 1, interface: 4, name: Object, package: 2, vis: public,",
            "      methods: [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11 ],",
            "    },",
            "  ],",
            "  packages: [",
            "    { objects: [ 1 ] },", // 1. main package
            "    {", // 2. builtin package
            "      methods: [ 1, 10, 11, 2, 3, 4, 5, 6, 7, 8, 9 ],",
            "      name: java.lang,",
            "      objects: [ 2 ],",
            "    },",
            "  ],",
            "  signatures: [",
            "    {                                  },", //  1. func()
            "    {                   results: [ 2 ] },", //  2. func() int
            "    {                   results: [ 3 ] },", //  3. func() any
            "    {                   results: [ 4 ] },", //  4. func() Object
            "    { params: [ 6 ]                    },", //  5. func(arg0 long)
            "    { params: [ 6, 8 ]                 },", //  6. func(arg0 long, arg1 int)
            "    { params: [ 7 ],    results: [ 1 ] },", //  7. func(arg0 Object) boolean
            "    { params: [ 9 ],    results: [ 2 ] },", //  8. func(index int) int
            "    { params: [ 9 ],    results: [ 5 ] },", //  9. func(index int) T
            "    { params: [ 9, 10 ]                },", // 10. func(index int, value int)
            "    { params: [ 9, 11 ]                }",  // 11. func(index int, value T)
            "  ],",
            "  structDescs: [",
            "    { },",                 // 1. struct{ }
            "    { fields: [ 1, 2 ] }", // 3. struct{ bar int, baz $Array<int> }
            "  ],",
            "  typeParams: [",
            "    { name: T, type: interfaceDesc5 }", // 1. T any
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
            "    { name: $get,      signature:  8 },", //  1. $get(index int) int
            "    { name: $get,      signature:  9 },", //  2. $get(index int) T
            "    { name: $len,      signature:  2 },", //  3. $len() int
            "    { name: $set,      signature: 10 },", //  4. $set(index int, value int)
            "    { name: $set,      signature: 11 },", //  5. $set(index int, value T)
            "    { name: bar,       signature: 12 },", //  6. bar( ??? ) ???
            "    { name: baz,       signature:  1 },", //  7. baz( ??? ) ???
            "    { name: cat,       signature: 13 },", //  8. cat( ??? ) ???
            "    { name: clone,     signature:  4 },", //  9. clone( ??? ) ???
            "    { name: equals,    signature:  7 },", // 10. equals( ??? ) ???
            "    { name: finalize,  signature:  1 },", // 11. finalize( ??? ) ???
            "    { name: getClass,  signature:  3 },", // 12. getClass( ??? ) ???
            "    { name: hashCode,  signature:  2 },", // 13. hashCode( ??? ) ???
            "    { name: notify,    signature:  1 },", // 14. notify( ??? ) ???
            "    { name: notifyAll, signature:  1 },", // 15. notifyAll( ??? ) ???
            "    { name: toString,  signature:  3 },", // 16. toString( ??? ) ???
            "    { name: wait,      signature:  1 },", // 17. wait( ??? ) ???
            "    { name: wait,      signature:  5 },", // 18. wait( ??? ) ???
            "    { name: wait,      signature:  6 }",  // 19. wait( ??? ) ???
            "  ],",
            "  arguments: [",
            "    {              type: basic1         },", //  1. <unnamed> int
            "    {              type: basic2         },", //  2. <unnamed> ???
            "    {              type: interfaceDesc5 },", //  3. <unnamed> ???
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
            "    {", // 3.
            "      abstracts: [ 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 6, 7, 8, 9 ],",
            "      pin: object1,",
            "    },",
            "    {", // 4.
            "      abstracts: [ 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 9 ],",
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
            "    { name: bar,       package: 1, receiver: 1, signature: 12, loc: 2, metrics: 1 },", //  1. bar(x int, y int) int
            "    { name: baz,       package: 1, receiver: 1, signature:  1, loc: 5, metrics: 2 },", //  2. baz()
            "    { name: cat,       package: 1, receiver: 1, signature: 13, loc: 8, metrics: 3 },", //  3. cat(t ...$Array<T>)
            "    { name: clone,     package: 2, receiver: 2, signature:  4, vis: protected     },", //  4. clone ???
            "    { name: equals,    package: 2, receiver: 2, signature:  7, vis: public        },", //  5. equals ???
            "    { name: finalize,  package: 2, receiver: 2, signature:  1, vis: protected     },", //  6. finalize ???
            "    { name: getClass,  package: 2, receiver: 2, signature:  3, vis: public        },", //  7. getClass ???
            "    { name: hashCode,  package: 2, receiver: 2, signature:  2, vis: public        },", //  8. hashCode ???
            "    { name: notify,    package: 2, receiver: 2, signature:  1, vis: public        },", //  9. notify ???
            "    { name: notifyAll, package: 2, receiver: 2, signature:  1, vis: public        },", // 10. notifyAll ???
            "    { name: toString,  package: 2, receiver: 2, signature:  3, vis: public        },", // 11. toString ???
            "    { name: wait,      package: 2, receiver: 2, signature:  1, vis: public        },", // 12. wait ???
            "    { name: wait,      package: 2, receiver: 2, signature:  5, vis: public        },", // 13. wait ???
            "    { name: wait,      package: 2, receiver: 2, signature:  6, vis: public        }",  // 14. wait ???
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
            "      methods: [ 1, 2, 3 ]",
            "    },",
            "    {", // 2.
            "      name: Object, package: 2, data: 1, interface: 4,",
            "      vis: public,",
            "      methods: [ 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14 ]",
            "    }",
            "  ],",
            "  packages: [",
            "    {", // 1. main package
            "      methods: [ 1, 2, 3 ],",
            "      objects: [ 1 ]",
            "    },",
            "    {", // 2. build-in package
            "      name: java.lang,",
            "      methods: [ 10, 11, 12, 13, 14, 4, 5, 6, 7, 8, 9 ],",
            "      objects: [ 2 ],",
            "    },",
            "  ],",
            "  signatures: [",
            "    {                                    },", //  1. func()
            "    {                     results: [ 2 ] },", //  2. func() int
            "    {                     results: [ 3 ] },", //  3. func( ??? ) ???
            "    {                     results: [ 4 ] },", //  4. func( ??? ) ???
            "    { params: [ 6 ]                      },", //  5. func( ??? ) ???
            "    { params: [ 6, 8 ]                   },", //  6. func( ??? ) ???  
            "    { params: [ 7 ],      results: [ 1 ] },", //  7. func(index int) int
            "    { params: [ 9 ],      results: [ 2 ] },", //  8. func(index int) T
            "    { params: [ 9 ],      results: [ 5 ] },", //  9. func( ??? ) ???
            "    { params: [ 9, 11 ]                  },", // 10. func(index int, value int)
            "    { params: [ 9, 12 ]                  },", // 11. func(index int, value T)
            "    { params: [ 13, 14 ], results: [ 2 ] },", // 12. func(x int, y int) int
            "    { params: [ 10 ], variadic: true     }",  // 13. func(t ...$Array<T>)
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
    public void ClassTestWithTypeParams() {
        final Tester t = Tester.classFromSource(
            "public class Foo<T extends Object> {",
            "  public T value",
            "}");
        t.checkProject(
            "{",
            "  abstracts: [",
            "    { name: clone,     signature: 4 },", //  1. clone() Object
            "    { name: equals,    signature: 7 },", //  2. equals(arg0 Object) boolean
            "    { name: finalize,  signature: 1 },", //  3. finalize()
            "    { name: getClass,  signature: 3 },", //  4. getClass() any // TODO: Is this really "any"?
            "    { name: hashCode,  signature: 2 },", //  5. hashCode() int
            "    { name: notify,    signature: 1 },", //  6. notify()
            "    { name: notifyAll, signature: 1 },", //  7. notifyAll()
            "    { name: toString,  signature: 3 },", //  8. toString() any // TODO: Shouldn't this be "String"?
            "    { name: wait,      signature: 1 },", //  9. wait()
            "    { name: wait,      signature: 5 },", // 10. wait(arg0 long)
            "    { name: wait,      signature: 6 }",  // 11. wait(arg0 long, arg1 int)
            "  ],",
            "  arguments: [",
            "    {             type: basic1         },", // 1. <unnamed> boolean
            "    {             type: basic2         },", // 2. <unnamed> int
            "    {             type: interfaceDesc3 },", // 3. <unnamed> any
            "    {             type: object2        },", // 4. <unnamed> Object
            "    { name: arg0, type: basic3         },", // 5. arg0 long
            "    { name: arg0, type: object2        },", // 6. arg0 Object
            "    { name: arg1, type: basic2         }",  // 7. arg1 int
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
            "    { name: clone,     package: 2, receiver: 2, signature: 4, vis: protected },", //  1. clone() Object
            "    { name: equals,    package: 2, receiver: 2, signature: 7, vis: public    },", //  2. equals(arg0 Object) boolean
            "    { name: finalize,  package: 2, receiver: 2, signature: 1, vis: protected },", //  3. finalize()
            "    { name: getClass,  package: 2, receiver: 2, signature: 3, vis: public    },", //  4. getClass() any // TODO: Is this really "any"?
            "    { name: hashCode,  package: 2, receiver: 2, signature: 2, vis: public    },", //  5. hashCode() int
            "    { name: notify,    package: 2, receiver: 2, signature: 1, vis: public    },", //  6. notify()
            "    { name: notifyAll, package: 2, receiver: 2, signature: 1, vis: public    },", //  7. notifyAll()
            "    { name: toString,  package: 2, receiver: 2, signature: 3, vis: public    },", //  8. toString() any // TODO: Shouldn't this be "String"?
            "    { name: wait,      package: 2, receiver: 2, signature: 1, vis: public    },", //  9. wait()
            "    { name: wait,      package: 2, receiver: 2, signature: 5, vis: public    },", // 10. wait(arg0 long)
            "    { name: wait,      package: 2, receiver: 2, signature: 6, vis: public    }",  // 11. wait(arg0 long, arg1 int)
            "  ],",
            "  objects: [",
            "    {", // 1. Foo<T> { value T }
            "      name: Foo, package: 1, data: 2, interface: 1,",
            "      vis: public, loc: 1,",
            "      typeParams: [ 1 ]",
            "    },",
            "    {", // 2. Object { }
            "      name: Object, package: 2, data: 1, interface: 2,",
            "      methods: [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11 ],",
            "      vis: public",
            "    }",
            "  ],",
            "  packages: [",
            "    { objects: [ 1 ] },", // 1. main package
            "    {",                   // 2. built-in package
            "      methods: [ 1, 10, 11, 2, 3, 4, 5, 6, 7, 8, 9 ],",
            "      name: java.lang,",
            "      objects: [ 2 ]",
            "    }",
            "  ],",
            "  signatures: [",
            "    {                               },", // 1. func()
            "    {                results: [ 2 ] },", // 2. func() int
            "    {                results: [ 3 ] },", // 3. func() any
            "    {                results: [ 4 ] },", // 4. func() Object
            "    { params: [ 5 ]                 },", // 5. func(arg0 long)
            "    { params: [ 5, 7 ]              },", // 6. func(arg0 long, arg1 int)
            "    { params: [ 6 ], results: [ 1 ] }",  // 7. func(arg0 Object) boolean
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
