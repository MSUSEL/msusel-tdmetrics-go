package abstractor.core;

import org.junit.jupiter.api.Test;

import abstractor.core.constructs.Project;

public class MetricsTests {
    
    @Test
    public void MetricsTestWithLiteralGetters() {
        final Project proj = Testing.classFromSource(
            "public class Foo {",
            "  public int bar() {",
            "    return -1;",
            "  }",
            "}");

        Testing.checkConstruct(proj, "metrics1",
            "{ codeCount: 1, complexity: 1, lineCount: 4 }");
    }
}
