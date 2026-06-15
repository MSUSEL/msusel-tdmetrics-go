package testData.java.test1010;

public class Foo<T extends Foo<>.Bar> {
    public interface Bar {
        boolean Run();
    }

    public T x;
}
