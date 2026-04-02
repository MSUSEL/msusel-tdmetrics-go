package abstractor.core.cmp;

import java.util.Iterator;
import java.util.List;
import java.util.SortedSet;
import java.util.function.Supplier;

import spoon.reflect.declaration.CtElement;

public interface Cmp {
    int run();
    
    static private <T> int compare(Comparable<T> a, T b) {
        if (a == null) return b == null ? 0 : -1;
        if (b == null) return 1;
        return a.compareTo(b);
    }
    
    static private <T> int compareHash(T a, T b) {
        if (a == null) return b == null ? 0 : -1;
        if (b == null) return 1;
        return Integer.compare(a.hashCode(), b.hashCode());
    }

    static public Cmp defer(CtElement a, Supplier<CtElement> fetch) {
        return () -> compareHash(a, fetch.get());
    }

    static public <T> Cmp defer(Comparable<T> a, Supplier<T> fetch) {
        return () -> compare(a, fetch.get());
    }

    static public <T extends Comparable<T>> Cmp deferList(List<? extends T> a, Supplier<List<? extends T>> fetch) {
        return () -> {
            final List<? extends T> b = fetch.get();
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

    static public <T extends Comparable<T>> Cmp deferSet(SortedSet<? extends T> a, Supplier<SortedSet<? extends T>> fetch) {
        return () -> {
            final SortedSet<? extends T> b = fetch.get();
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
