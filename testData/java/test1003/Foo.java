package testData.java.test1003;

public class Foo<T extends Object> {
    private T value;

    public T getValue() { return this.value; }
    public void setValue(T value) { this.value = value; }
}
