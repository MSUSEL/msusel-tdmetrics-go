package abstractor.core.cmp;

public interface Cmp {
    int run();

    static public <T> Cmp cmp(Comparable<T> a, T b) {
        return () -> a.compareTo(b);
    }

    public interface Fetch<T> { T run(); }
    
    static public <T> Cmp defer(Comparable<T> a, Fetch<T> fetch) {
        return () -> a.compareTo(fetch.run());
    }

    static public int or(Cmp ...comparers) {
        for (Cmp cmp: comparers) {
            final int result = cmp.run();
            if (result != 0) return result;
        }
        return 0;
    }
}
