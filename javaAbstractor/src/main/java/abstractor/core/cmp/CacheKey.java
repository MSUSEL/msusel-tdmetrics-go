package abstractor.core.cmp;

public class CacheKey {
    final public Object a;
    final public Object b;
    final public int hash;

    public CacheKey(Object a, Object b) {
        this.a = a;
        this.b = b;
        final String s = a.hashCode() + ":" + b.hashCode();
        this.hash = s.hashCode(); 
    }

    @Override
    public int hashCode() { return this.hash; }

    @Override
    public boolean equals(Object obj) {
        return obj != null &&
            obj instanceof CacheKey ck &&
            this.hash == ck.hash &&
            this.a.equals(ck.a) &&
            this.b.equals(ck.b);
    }
}
