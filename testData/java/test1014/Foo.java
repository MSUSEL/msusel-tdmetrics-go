package testData.java.test1014;

public class Foo {
  public enum Bar {
    LEFT,
    RIGHT,
    UP,
    DOWN
  }

  public Bar barFrom(String name) {
    return Enum.valueOf(Bar.class, name);
  }

  public int indexOf(Bar value) {
    int index = 0;
    for (Enum e : Bar.values()) {
      if (e.compareTo(value) == 0) return index;
      index++;
    }
    return -1;
  }
}
