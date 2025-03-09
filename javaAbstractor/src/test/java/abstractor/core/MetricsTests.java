package abstractor.core;

import static abstractor.core.Testing.*;

import org.junit.jupiter.api.Test;

import abstractor.core.constructs.Project;

public class MetricsTests {
    
    @Test
    public void MetricsTestWithEmpty() {
        final Project proj = classFromSource(
            "public class Foo {",
            "  public void bar() { }",
            "}");

        checkConstruct(proj, "metrics1",
            "{",
            "  complexity: 1,",
            "  lineCount:  1",
            "}");
    }

    @Test
    public void MetricsTestWithLiteralGetters() {
        final Project proj = classFromSource(
            "public class Foo {",
            "  public int bar() {",
            "    return -1;",
            "  }",
            "}");

        checkConstruct(proj, "metrics1",
            "{",
            "  codeCount:  1,",
            "  complexity: 1,",
            "  indents:    1,",
            "  lineCount:  4",
            "}");
    }
}
