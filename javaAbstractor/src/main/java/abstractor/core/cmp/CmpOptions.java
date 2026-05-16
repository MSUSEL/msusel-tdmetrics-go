package abstractor.core.cmp;

public class CmpOptions implements Comparable<CmpOptions>{
    public boolean useResolved;

    @Override
    public String toString() { return "cmp: useResolved=" + this.useResolved; }
    
    public int compareTo(CmpOptions c) {
        return Boolean.compare(this.useResolved, c.useResolved);
    }

    public boolean equals(Object obj) {
        return obj instanceof CmpOptions c && this.compareTo(c) == 0;
    }
}
