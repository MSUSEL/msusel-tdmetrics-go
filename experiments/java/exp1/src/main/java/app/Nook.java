package app;

import app.Foo;

public class Nook implements Foo.FooBar.Bar, Foo.FooBaz.Baz {
    public boolean Run1() {
        System.out.println("Run1");
        return false;
    }

    public int Run2() {
        System.out.println("Run2");
        return 42;
    }
}
