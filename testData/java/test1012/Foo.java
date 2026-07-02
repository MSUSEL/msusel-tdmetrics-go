public class Foo {
    public class Bar<T> {
        public <U> Bar(T t, U u) {
            this.str = "[" + t + ", " + u + "]";
        }
        public String str;
    }

    public <S, P> String tak(S s, P p) {
        return (new <P>Bar<S>(s, p)).str;
    }

    public String baz() {
        return this.tak(10, "Hello") + ":" + this.tak("Cat", 3.14);
    }
}
