package abstractor.core.cmp;

import java.util.*;

import abstractor.core.log.Logger;

public class CmpContext {
    final private CmpOptions options;
    final private String prefix;
    final private Logger log;
    final private HashMap<CacheKey, Integer> cache;

    public CmpContext(CmpOptions options) {
        this(options, "", null);
    }

    public CmpContext(CmpOptions options, String prefix) {
        this(options, prefix, null);
    }

    public CmpContext(CmpOptions options, Logger log) {
        this(options, "", log);
    }

    public CmpContext(CmpOptions options, String prefix, Logger log) {
        this(options, prefix, log, new HashMap<>());
    }

    private CmpContext(CmpOptions options, String prefix, Logger log, HashMap<CacheKey, Integer> cache) {
        this.options = options != null ? options : new CmpOptions();
        this.prefix = prefix;
        this.log = log;
        this.cache = cache;
    }

    static private String join(String a, String b) {
        if (a.isBlank()) return b;
        if (b.isBlank()) return a;
        return a + "." + b;
    }

    private void debugPrint(String msg) {
        if (this.log != null) this.log.log(msg);
        else System.out.println(msg);
    }

    public CmpContext subContext(String prefix) {
        return new CmpContext(this.options, join(this.prefix, prefix), this.log, this.cache);
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

            final String msg = header + "compare<" + type + ">(" + a + ", " + b + ") => " + cmp;
            this.debugPrint(msg);
        }
        return cmp;
    }
    
    private <T> int innerCompare(T a, T b, String name) {
        if (a == b) return 0;
        if (a == null) return b == null ? 0 : -1;
        if (b == null) return 1;

        // Check cache and prevent loop.
        final CacheKey cacheKey = new CacheKey(a, b);
        final Integer cachedValue = this.cache.get(cacheKey);
        if (cachedValue != null) return cachedValue.intValue();
        this.cache.put(cacheKey, 0);

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
        this.cache.put(cacheKey, cmp);
        return cmp;
    }
}
