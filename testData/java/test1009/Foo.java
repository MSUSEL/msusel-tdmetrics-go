package testData.java.test1009;

public class Foo {
    public interface Bar {
        boolean Run();
    }

    public class Bip<T extends Bar> {
        public class Tak {
            public T x;
        }

        public Tak t = new Tak();
    }

    public class Baz implements Bar {
        public boolean Run() {
            return true;
        }
    }

    public boolean Boop() {
        Bip<Baz> b = new Bip<>();
        b.t.x = new Baz();
        return b.t.x.Run();
    }
    
    public boolean Moog() {
        Bip<Bar> b = new Bip<>();
        b.t.x = new Bar() {
            public boolean Run() {
                return false;
            }
        };
        return b.t.x.Run();
    }
}
