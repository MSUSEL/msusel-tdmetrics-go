package testData.java.test1014;

public class Foo {
  final public int Bar;
  public Foo(Boolean value) {
    if (value) {
      this.Bar = 42;
      return;
    }
    this.Bar = 38;
  }

  public Foo(boolean value) {
    this.Bar = value ? 123 : 456;
  }
}
