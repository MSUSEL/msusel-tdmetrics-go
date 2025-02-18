package abstractor.core.cmp;

import java.util.List;

public interface Cmp {
    int run();

    public interface Fetch<T> { T run(); }
    
    static private <T> int compare(Comparable<T> a, T b) {
        if (a == null) return b == null ? 0 : -1;
        if (b == null) return 1;
        return a.compareTo(b);
    }

    static public <T> Cmp defer(Comparable<T> a, Fetch<T> fetch) {
        return () -> compare(a, fetch.run());
    }

    static public <T> Cmp deferList(List<? extends Comparable<T>> a, Fetch<List<? extends T>> fetch) {
        return () -> {
            final List<? extends T> b = fetch.run();
            if (a == null) return b == null ? 0 : -1;
            if (b == null) return 1;
            final int aLen = a.size();
            final int bLen = b.size();
            final int min = Integer.min(aLen, bLen);
            for (int i = 0; i < min; i++) {
                final int cmp = compare(a.get(i), b.get(i));
                if (cmp != 0) return cmp;
            }
            return Integer.compare(aLen, bLen);
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
