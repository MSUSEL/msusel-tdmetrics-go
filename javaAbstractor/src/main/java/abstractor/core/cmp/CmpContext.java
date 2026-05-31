package abstractor.core.cmp;

import java.util.HashMap;

public class CmpContext {
    final private CmpOptions options;
    final private String prefix;

    public CmpContext(CmpOptions options) {
        this(options, "");
    }

    public CmpContext(CmpOptions options, String prefix) {
        this.options = options != null ? options : new CmpOptions();
        this.prefix = prefix;
    }

    static private String join(String a, String b) {
        if (a.isBlank()) return b;
        if (b.isBlank()) return a;
        return a + "." + b;
    }

    public CmpContext subContext(String prefix) {
        return new CmpContext(this.options, join(this.prefix, prefix));
    }

    private final HashMap<String, Integer> cache = new HashMap<>();
    
    private String key(Object a, Object b) {
        return System.identityHashCode(a) + ":" + System.identityHashCode(b);
    }
    
    public <T> int compare(T a, T b) {
        return this.compare(a, b, "");
    }

    public <T> int compare(T a, T b, String name) {
        final int cmp = this.innerCompare(a, b, name);
        if (this.options.debugPrint) {
            name = join(this.prefix, name);
            final String header = name.isBlank() ? "" : name+": ";
            String type = "unknown";
            if (a != null) type = a.getClass().getSimpleName();
            else if (b != null) type = b.getClass().getSimpleName();
            System.out.println(header + "compare<"+type+">(" + a + ", " + b + ") => " + cmp);
        }
        return cmp;
    }
    
    private <T> int innerCompare(T a, T b, String name) {
        if (a == b) return 0;
        if (a == null) return b == null ? 0 : -1;
        if (b == null) return 1;

        // Check cache and prevent loop.
        final String cacheKey = key(a, b);
        final Integer cachedValue = cache.get(cacheKey);
        if (cachedValue != null) return cachedValue.intValue();
        cache.put(cacheKey, 0);

        // Perform comparison.
        int cmp;
        if (a instanceof CmpGetter<?>) {
            @SuppressWarnings("unchecked")
            final CmpGetter<T> ac = (CmpGetter<T>) a;
            final CmpContext subContext = this.subContext(name);
            cmp = ac.getCmp(b, options).run(subContext);

        } else if (a instanceof Comparable<?>) {
            @SuppressWarnings("unchecked")
            final Comparable<T> ac = (Comparable<T>) a;
            cmp = ac.compareTo(b);

        } else {
            cmp = Integer.compare(a.hashCode(), b.hashCode());
        }

        // Cache results.
        cache.put(cacheKey, cmp);
        return cmp;
    }
}
