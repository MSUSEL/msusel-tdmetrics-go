package testData.java.test1002;

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
}
