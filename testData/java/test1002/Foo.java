package testData.java.test1002;

public class Foo {
  int bar(int x, int y) {
    return x + y*2;
  }

  void baz() {
    System.out.println("Baz");
  }

  void cat(int ...t) {
    System.out.println(t);
  }
}
