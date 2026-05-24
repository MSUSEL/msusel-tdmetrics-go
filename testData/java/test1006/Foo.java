package testData.java.test1006;

public class Foo<T1 extends Object> {
    public class Bar<T2 extends Object> {
        public Bar(T1 x, T2 y) {
            this.x = x;
            this.y = y;
        }

        final public T1 x;
        final public T2 y;
    }

    public Bar<Integer> value1;
    public Bar<String>  value2;
}
