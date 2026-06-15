package testData.java.test1008;

public class Foo<T extends Object> {
    public T x;

    public Foo<Integer> Foo1() { return new Foo<>(); };
    public Foo<String>  Foo2() { return new Foo<>(); };
    public <U> Foo<U>   Foo3() { return new Foo<>(); };
}
