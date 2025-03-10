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
}
