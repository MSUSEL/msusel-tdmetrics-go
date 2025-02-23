package abstractor.core.cmp;

import java.util.Iterator;
import java.util.List;
import java.util.SortedSet;

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

    static public <T extends Comparable<T>> Cmp deferList(List<? extends T> a, Fetch<List<? extends T>> fetch) {
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

    static public <T extends Comparable<T>> Cmp deferSet(SortedSet<? extends T> a, Fetch<SortedSet<? extends T>> fetch) {
        return () -> {
            final SortedSet<? extends T> b = fetch.run();
            if (a == null) return b == null ? 0 : -1;
            if (b == null) return 1; 

            final Iterator<? extends T> aIt = a.iterator();
            final Iterator<? extends T> bIt = b.iterator();
            while (aIt.hasNext() && bIt.hasNext()) {
                final int cmp = compare(aIt.next(), bIt.next());
                if (cmp != 0) return cmp;
            }

            if (aIt.hasNext()) return bIt.hasNext() ? 0 : -1;
            return bIt.hasNext() ? 1 : 0;
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
