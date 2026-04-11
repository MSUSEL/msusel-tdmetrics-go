package abstractor.core.cmp;

import java.util.Iterator;
import java.util.List;
import java.util.SortedSet;
import java.util.function.Supplier;

public interface Cmp {
    int run(CmpContext context);

    static public <T> int compareTo(T a, T b, CmpOptions options) {
        return (new CmpContext(options)).compare(a, b);
    }
    
    static public <T> Cmp deferHash(T a, Supplier<T> fetch) {
        return (CmpContext context) -> context.compare(a, fetch.get());
    }

    static public <T> Cmp defer(Comparable<T> a, Supplier<T> fetch) {
        return (CmpContext context) -> context.compare(a, fetch.get());
    }

    static public <T extends Comparable<T>> Cmp deferList(List<? extends T> a, Supplier<List<? extends T>> fetch) {
        return (CmpContext context) -> {
            final List<? extends T> b = fetch.get();
            if (a == null) return b == null ? 0 : -1;
            if (b == null) return 1;
            final int aLen = a.size();
            final int bLen = b.size();
            final int min = Integer.min(aLen, bLen);
            for (int i = 0; i < min; i++) {
                final int cmp = context.compare(a.get(i), b.get(i));
                if (cmp != 0) return cmp;
            }
            return Integer.compare(aLen, bLen);
        };
    }

    static public <T extends Comparable<T>> Cmp deferSet(SortedSet<? extends T> a, Supplier<SortedSet<? extends T>> fetch) {
        return (CmpContext context) -> {
            final SortedSet<? extends T> b = fetch.get();
            if (a == null) return b == null ? 0 : -1;
            if (b == null) return 1; 

            final Iterator<? extends T> aIt = a.iterator();
            final Iterator<? extends T> bIt = b.iterator();
            while (aIt.hasNext() && bIt.hasNext()) {
                final int cmp = context.compare(aIt.next(), bIt.next());
                if (cmp != 0) return cmp;
            }

            if (aIt.hasNext()) return bIt.hasNext() ? 0 : -1;
            return bIt.hasNext() ? 1 : 0;
        };
    }

    static public Cmp or(Cmp ...comparers) {
        return (CmpContext context) -> {
            for (Cmp cmp: comparers) {
                final int result = cmp.run(context);
                if (result != 0) return result;
            }
            return 0;
        };
    }
}
