package abstractor.core.cmp;

import java.util.HashMap;

public class CmpContext {
    final private CmpOptions options;

    public CmpContext(CmpOptions options) {
        this.options = options != null ? options : new CmpOptions();
    }

    private final HashMap<String, Integer> cache = new HashMap<>();
    
    private String key(Object a, Object b) {
        return System.identityHashCode(a) + ":" + System.identityHashCode(b);
    }
    
    public <T> int compare(T a, T b) {
        if (a == b) return 0;
        if (a == null) return b == null ? 0 : -1;
        if (b == null) return 1;

        // Check cache and prevent loop.
        String cacheKey = key(a, b);
        Integer cachedValue = cache.get(cacheKey);
        if (cachedValue != null) return cachedValue.intValue();
        cache.put(cacheKey, 0);

        // Perform comparison.
        int cmp;
        if (a instanceof CmpGetter<?>) {
            @SuppressWarnings("unchecked")
            CmpGetter<T> ac = (CmpGetter<T>) a;
            cmp = ac.getCmp(b, options).run(this);

        } else if (a instanceof Comparable<?>) {
            @SuppressWarnings("unchecked")
            Comparable<T> ac = (Comparable<T>) a;
            cmp = ac.compareTo(b);

        } else {
            cmp = Integer.compare(a.hashCode(), b.hashCode());
        }

        // Cache results.
        cache.put(cacheKey, cmp);
        return cmp;
    }
}
