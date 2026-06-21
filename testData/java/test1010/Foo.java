package testData.java.test1010;

public class Foo<T> {
    public void Bar(T value) { Baz(value); }
    public void Bar(int value) { Baz(value); }
    public <U> void Baz(U value) {}

    static public void Do() {
        new Foo<String>().Bar("Hello");
    }
}
