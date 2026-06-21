package app;

import app.Nook;
import app.Foo;

public class Main {
    public static void main(String[] args) {
        Nook n = new Nook();
        
        Foo f0 = new Foo();

        Foo.FooBar<Nook> f1 = f0.new FooBar<Nook>();
        f1.x = n;
        
        Foo.FooBaz<Nook> f2 = f0.new FooBaz<Nook>();
        f2.x = n;

        Foo.Pizza<Nook> f3 = f0.new Pizza<Nook>();
        f3.x = n;

        f1.x.Run1();
        f2.x.Run2();
        f3.x.Run1();
        f2.x.Run2();
    }
}
