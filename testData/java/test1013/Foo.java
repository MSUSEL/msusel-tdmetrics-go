package testData.java.test1013;

public class Foo {
  int bar(int x, int y) {
    return x + y*2;
  }

  void baz() {
    int a = 0;
  }

  void cat(int ...t) {
    int b = t.length;
  }

  @Override
  public String toString() {
      return "Foo(" + super.toString() + ")";
  }

  @Override
  public int hashCode() {
      return super.hashCode() + 42;
  }
}
