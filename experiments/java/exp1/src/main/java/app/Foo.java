package app;

public class Foo {
    public class FooBar<T extends FooBar.Bar> {
        public interface Bar {
            boolean Run1();
        }

        public T x;
    }

    public class FooBaz<T extends FooBaz.Baz> {
        public interface Baz {
            int Run2();
        }

        public T x;
    }
    
    public class Pizza<T extends FooBar.Bar & FooBaz.Baz> {
        public T x;
    }
}
