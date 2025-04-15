package abstractor.core;

import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

public class ClassTests {

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
