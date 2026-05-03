package abstractor.core;

import org.junit.jupiter.api.Test;

/**
 * Tests that the abstractor handles various Java constructs without crashing.
 * These tests verify robustness — they don't check exact output,
 * just that no exceptions are thrown and no errors are logged.
 */
public class RobustnessTests {

    @Test
    public void AnnotationUsage() {
        Tester.classesFromSource(
            "import java.lang.annotation.*;",
            "@Retention(RetentionPolicy.RUNTIME)",
            "public class Foo {",
            "  @Deprecated",
            "  public void bar() { }",
            "}");
    }

    @Test
    public void WildcardParameter() {
        Tester.classesFromSource(
            "import java.util.*;",
            "public class Foo {",
            "  public void bar(List<?> items) { }",
            "}");
    }

    @Test
    public void BoundedWildcard() {
        Tester.classesFromSource(
            "public class Foo {",
            "  public void bar(Comparable<? extends Number> comp) { }",
            "}");
    }

    @Test
    public void LowerBoundedWildcard() {
        Tester.classesFromSource(
            "import java.util.*;",
            "public class Foo {",
            "  public void bar(List<? super Integer> items) { }",
            "}");
    }

    @Test
    public void AnonymousClassInMethod() {
        Tester.classesFromSource(
            "public class Foo {",
            "  public void bar() {",
            "    Runnable r = new Runnable() {",
            "      @Override",
            "      public void run() { }",
            "    };",
            "  }",
            "}");
    }

    @Test
    public void LambdaInMethod() {
        Tester.classesFromSource(
            "public class Foo {",
            "  public void bar() {",
            "    Runnable r = () -> { };",
            "  }",
            "}");
    }

    @Test
    public void ExternalTypeUsage() {
        Tester.classesFromSource(
            "import java.util.*;",
            "public class Foo {",
            "  public String name;",
            "  public List<String> items;",
            "  public Map<String, Integer> lookup;",
            "}");
    }

    @Test
    public void EnumDeclaration() {
        Tester.classesFromSource(
            "public class Foo {",
            "  public enum Color { RED, GREEN, BLUE }",
            "  public Color color;",
            "}");
    }

    @Test
    public void GenericClass() {
        Tester.classesFromSource(
            "public class Foo<T> {",
            "  private T value;",
            "  public T getValue() { return value; }",
            "  public void setValue(T v) { this.value = v; }",
            "}");
    }

    @Test
    public void NestedNamedClass() {
        Tester.classesFromSource(
            "public class Foo {",
            "  public static class Inner {",
            "    public int x;",
            "  }",
            "  public Inner inner;",
            "}");
    }
}
