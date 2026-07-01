package abstractor.core;

import org.junit.jupiter.api.Test;

public class MetricsTests {
    
    @Test
    public void Empty() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  public void bar() { }",
            "}");
        t.checkConstructCount("metrics", 0);
    }

    @Test
    public void Simple() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  public int bar() {",
            "    return -1;",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  3,",
            "  reads: [ basic1 ],",
            "}");
    }
    
    @Test
    public void SimpleWithExtraIndent() {
        final Tester t = Tester.classesFromSource(
            "  public class Foo {",
            "      public int bar() {",
            "             return -1;",
            "      }",
            "  }");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  3,",
            "  reads: [ basic1 ],",
            "}");
    }

    @Test
    public void SimpleParams() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "   public int bar(int a,",
            "                  int b,",
            "                  int c) {",
            "       return a + b + c;",
            "   }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  5,",
            "  complexity: 1,",
            "  indents:    5,",
            "  lineCount:  5,",
            "  reads: [ basic1 ],",
            "}");
    }

    @Test
    public void SimpleWithReturn() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "   public int bar(int a) {",
            "      final int x = 4 * a + 1;",
            "      return x;",
            "   }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  4,",
            "  complexity: 1,",
            "  indents:    2,",
            "  lineCount:  4,",
            "  reads:  [ basic1 ],",
            "  writes: [ basic1 ],",
            "}");
    }

    @Test
    public void SimpleWithSpace() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  public int bar(int a) {",
            "    // Bacon is tasty",
            "    ",
            "    return a + 3;",
            "    /* This is not a comment",
            "       it is a sandwich */",
            "    ",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  8,",
            "  reads: [ basic1 ],",
            "}");
    }
    
    @Test
    public void SimpleSiblingInvoke() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  public void bar(int x) {",
            "    this.tak(x);",
            "  }",
            "  public void tak(int x) { }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  3,",
            "  invokes: [ method2 ],",
            "  reads:   [ basic1, object1 ],",
            "}");
    }
    
    @Test
    public void SimpleNeighborInvoke() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  public void bar(int x) {",
            "    Boop.tak(x);",
            "  }",
            "  public class Boop {",
            "    static public void tak(int x) { }",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  3,",
            "  invokes: [ method2 ],",
            "  reads:   [ basic1, object1, object2 ],",
            "}");
    }
    
    @Test
    public void SimpleDoubleSelectInvoke() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  public void bar(int x) {",
            "    Boop.baz.tak(x);",
            "  }",
            "  public class Boop {",
            "    public class Baz {",
            "      static public void tak(int x) { }",
            "    }",
            "    static public Baz baz;",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  3,",
            "  invokes: [ method2 ],",
            "  reads:   [ basic1, object1, object2, object3, selection1 ],",
            "}");
    }
    
    @Test
    public void SimpleSystemPrint() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  public void bar(int x) {",
            "    System.out.println(x);",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  3,",
            "  reads: [ basic6, object9, object10, selection1 ],",
            // TODO: Add invokes for println (for this and other tests with println)
            "}");
    }

    @Test
    public void SimpleIf() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  public void bar() {",
            "    int x = 9;",
            "    if (x > 7) {",
            "      x = 4;",
            "    }",
            "    System.out.println(x);",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  7,",
            "  complexity: 2,",
            "  indents:    6,",
            "  lineCount:  7,",
            "  reads:  [ basic1, basic6, object9, object10, selection1 ],",
            "  writes: [ basic6 ],",
            // TODO: Add invokes for println (for this and other tests with println)
            "}");
    }

    @Test
    public void SimpleIfElse() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  public void bar() {",
            "    int x = 9;",
            "    if (x > 7) {",
            "      x = 4;",
            "    } else {",
            "      x = 2;",
            "      System.out.println(\"cat\");",
            "    }",
            "    System.out.println(x);",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  10,",
            "  complexity:  2,",
            "  indents:    11,",
            "  lineCount:  10,",
            "  reads:  [ basic1, basic6, basic8, object9, object10, selection1 ],",
            "  writes: [ basic6 ],",
            // TODO: Add invokes for println (for this and other tests with println)
            "}");
    }

    @Test
    public void SimpleIfElseIf() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  public void bar() {",
            "    int x = 9;",
            "    if (x > 7) {",
            "      x = 4;",
            "    } else if (x > 4) {",
            "      x = 3;",
            "    }",
            "    System.out.println(x)",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  9,",
            "  complexity: 3,",
            "  indents:    9,",
            "  lineCount:  9,",
            "  reads:  [ basic1, basic6, object9, object10, selection1 ],",
            "  writes: [ basic6 ],",
            // TODO: Add invokes for println (for this and other tests with println)
            "}");
    }

    @Test
    public void SimpleIfElseIfElse() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  public void bar() {",
            "    int x = 9;",
            "    if (x > 7) {",
            "      x = 4;",
            "    } else if (x > 4) {",
            "      x = 3;",
            "    } else {",
            "      x = 2;",
            "    }",
            "    System.out.println(x);",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  11,",
            "  complexity:  3,",
            "  indents:    12,",
            "  lineCount:  11,",
            "  reads:  [ basic1, basic6, object9, object10, selection1 ],",
            "  writes: [ basic6 ],",
            // TODO: Add invokes for println (for this and other tests with println)
            "}");
    }

    @Test
    public void SimpleSwitch() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  public void bar() {",
            "    int x = 9;",
            "    switch (x) {",
            "    case 7:",
            "      x = 4;",
            "    case 4:",
            "      x = 3;",
            "    default:",
            "      x = 2;",
            "    }",
            "    System.out.println(x);",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  12,",
            "  complexity:  3,",
            "  indents:    13,",
            "  lineCount:  12,",
            "  reads: [ basic6, object9, object10, selection1 ],",
            "  writes: [ basic6 ],",
            // TODO: Add invokes for println (for this and other tests with println)
            "}");
    }

    @Test
    public void SimpleForLoop() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  public void bar() {",
            "    for (int i = 0; i < 10; i++) {",
            "      System.out.println(i);",
            "    }",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  5,",
            "  complexity: 2,",
            "  indents:    4,",
            "  lineCount:  5,",
            "  reads:  [ basic1, basic6, object9, object10, selection1 ],",
            "  writes: [ basic6 ],",
            // TODO: Add invokes for println (for this and other tests with println)
            "}");
    }

    @Test
    public void SimpleLogicalOr() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  public int bar(int x) {",
            "    if (x < 0 || x > 10) {",
            "      x = 4;",
            "    }",
            "    return x;",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  6,",
            "  complexity: 3,",
            "  indents:    5,",
            "  lineCount:  6,",
            "  reads:  [ basic1, basic2 ],",
            "  writes: [ basic2 ],",
            "}");
    }

    @Test
    public void OneLineLogicalOr() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  public boolean bar(int x) {",
            "    return x < 0 || x > 10;",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 2,",
            "  indents:    1,",
            "  lineCount:  3,",
            "  reads: [ basic1, basic2 ],",
            "}");
    }

    @Test
    public void SimpleLogicalAnd() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  public boolean bar(int x) {",
            "    if (x >= 0 && x < 10) {",
            "      x = 4;",
            "    }",
            "    return x;",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  6,",
            "  complexity: 3,",
            "  indents:    5,",
            "  lineCount:  6,",
            "  reads:  [ basic1, basic2 ],",
            "  writes: [ basic2 ],",
            "}");
    }

    @Test
    public void OneLineLogicalAnd() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  public boolean bar(int x) {",
            "    return x >= 0 && x < 10;",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 2,",
            "  indents:    1,",
            "  lineCount:  3,",
            "  reads: [ basic1, basic2 ],",
            "}");
    }

    @Test
    public void OneLineBinaryOp() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  public boolean bar(int x) {",
            "    return x >= 0 & x < 10;",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  3,",
            "  reads: [ basic1, basic2 ],",
            "}");
    }

    @Test
    public void GetterWithSelect() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  private int x;",
            "  public int bar() {",
            "    return this.x;",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  3,",
            "  getter:     true,",
            "  reads: [ basic1, object1, selection1 ],",
            "}");
    }
    
    @Test
    public void GetterWithParenthesesAndCast() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  private double x;",
            "  public int bar() {",
            "    return ((int)((this.x)));",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  3,",
            "  getter:     true,",
            "  reads: [ basic1, basic2, object1, selection1 ],",
            "}");
    }

    @Test
    public void SetterWithSelect() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  private int x;",
            "  public void bar(int x) {",
            "    this.x = x;",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  3,",
            "  setter:     true,",
            "  reads:  [ basic1, object1, selection1 ],",
            "  writes: [ basic1, selection1 ],",
            "}");
    }

    @Test
    public void SetterWithParenthesesAndCast() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  private int x;",
            "  public void bar(double x) {",
            "    this.x = ((int)((x)));",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  3,",
            "  setter:     true,",
            "  reads:  [ basic1, basic2, object1, selection1 ],",
            "  writes: [ basic2, selection1 ],",
            "}");
    }

    @Test
    public void SetterWithFieldValue() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  private int x;",
            "  public void bar(Foo other) {",
            "    this.x = other.x;",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  3,",
            "  setter:     true,",
            "  reads:  [ basic1, object1, selection1 ],",
            "  writes: [ basic1, selection1 ],",
            "}");
    }

    @Test
    public void SetterNotReversed() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  private int x;",
            "  public void bar(Foo other) {",
            "    other.x = this.x;",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  3,",
            "  reads:  [ basic1, object1, selection1 ],",
            "  writes: [ basic1, selection1 ],",
            "}");
    }

    @Test
    public void LiteralSetter() {
        final Tester t = Tester.classesFromSource(
            "public class Foo {",
            "  private boolean visible;",
            "  public void hide() {",
            "    this.visible = false;",
            "  }",
            "}");
        t.checkConstruct("metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  3,",
            "  setter:     true,",
            "  reads:  [ basic1, object1, selection1 ],",
            "  writes: [ basic1, selection1 ],",
            "}");
    }
}
