package abstractor.core.cmp;

public interface Cmp {
    int run();

    public interface Fetch<T> { T run(); }
    
    static public <T> Cmp defer(Comparable<T> a, Fetch<T> fetch) {
        return () -> {
            T other = fetch.run();
            if (a == null) {
                return other == null ? 0 : -1;
            }
            if (other == null) return 1;
            return a.compareTo(other);
        };
    }

    static public int or(Cmp ...comparers) {
        for (Cmp cmp: comparers) {
            final int result = cmp.run();
            if (result != 0) return result;
        }
        return 0;
    }
}
