public class Foo {
    public void bar(int x) {
        ((funky) () -> {
            baz(x + 42);
        }).execute();
    }
            
    public void baz(int value) { }
            
    interface funky {
        void execute();
    }
}
