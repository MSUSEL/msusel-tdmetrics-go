package testData.java.test1007;

public class Foo<T extends Object> {
    public T x;

    static public Foo<Integer> Foo1() { return new Foo<>(); };
    static public Foo<String>  Foo2() { return new Foo<>(); };
}
