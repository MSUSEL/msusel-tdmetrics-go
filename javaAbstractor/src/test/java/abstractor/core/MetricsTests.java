package abstractor.core;

import static abstractor.core.Testing.*;

import org.junit.jupiter.api.Test;

import abstractor.core.constructs.Project;

public class MetricsTests {
    
    @Test
    public void Empty() {
        final Project proj = classFromSource(
            "public class Foo {",
            "  public void bar() { }",
            "}");

        checkConstruct(proj, "metrics1",
            "{",
            "  codeCount:  1,",
            "  complexity: 1,",
            "  #indents:   0,",
            "  lineCount:  1",
            "}");
    }

    @Test
    public void Simple() {
        final Project proj = classFromSource(
            "public class Foo {",
            "  public int bar() {",
            "    return -1;",
            "  }",
            "}");

        checkConstruct(proj, "metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  3",
            "}");
    }
    
    @Test
    public void SimpleWithExtraIndent() {
        final Project proj = classFromSource(
            "  public class Foo {",
            "      public int bar() {",
            "             return -1;",
            "      }",
            "  }");

        checkConstruct(proj, "metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  3",
            "}");
    }

    @Test
    public void SimpleParams() {
        final Project proj = classFromSource(
            "public class Foo {",
            "   public int bar(int a,",
            "                  int b,",
            "                  int c) {",
            "       return a + b + c;",
            "   }",
            "}");

        checkConstruct(proj, "metrics1",
            "{",
            "  codeCount:  5,",
            "  complexity: 1,",
            "  indents:    5,",
            "  lineCount:  5",
            "}");
    }

    @Test
    public void SimpleWithReturn() {
        final Project proj = classFromSource(
            "public class Foo {",
            "   public int bar(int a) {",
            "      final int x = 4 * a + 1;",
            "      return x;",
            "   }",
            "}");

        checkConstruct(proj, "metrics1",
            "{",
            "  codeCount:  4,",
            "  complexity: 1,",
            "  indents:    2,",
            "  lineCount:  4",
            "}");
    }

    @Test
    public void SimpleWithSpace() {
        final Project proj = classFromSource(
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

        checkConstruct(proj, "metrics1",
            "{",
            "  codeCount:  3,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  8",
            "}");
    }

    @Test
    public void SimpleIf() {
        final Project proj = classFromSource(
            "public class Foo {",
            "  public void bar() {",
            "    int x = 9;",
            "    if (x > 7) {",
            "      x = 4;",
            "    }",
            "    System.out.println(x);",
            "  }",
            "}");

        checkConstruct(proj, "metrics1",
            "{",
            "  codeCount:  7,",
            "  complexity: 2,",
            "  indents:    6,",
            "  lineCount:  7",
            "}");
    }

    @Test
    public void SimpleIfElse() {
        final Project proj = classFromSource(
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
        checkConstruct(proj, "metrics1",
            "{",
            "  codeCount:  10,",
            "  complexity:  2,",
            "  indents:    11,",
            "  lineCount:  10",
            "}");
    }

    @Test
    public void SimpleIfElseIf() {
        final Project proj = classFromSource(
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
        checkConstruct(proj, "metrics1",
            "{",
            "  codeCount:  9,",
            "  complexity: 3,",
            "  indents:    9,",
            "  lineCount:  9",
            "}");
    }

    @Test
    public void SimpleIfElseIfElse() {
        final Project proj = classFromSource(
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
        checkConstruct(proj, "metrics1",
            "{",
            "  codeCount:  11,",
            "  complexity:  3,",
            "  indents:    12,",
            "  lineCount:  11",
            "}");
    }

    @Test
    public void SimpleSwitch() {
        final Project proj = classFromSource(
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
            "    System.out.println(x)",
            "  }",
            "}");
        checkConstruct(proj, "metrics1",
            "{",
            "  codeCount:  12,",
            "  complexity:  3,",
            "  indents:    13,",
            "  lineCount:  12",
            "}");
    }
}
