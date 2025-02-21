package abstractor.core.diff.comparators;

import abstractor.core.diff.Comparator;

public class ArrayComparator<T> implements Comparator {
    public final T[] aSource;
    public final T[] bSource;
  
    public ArrayComparator(T[] aSource, T[] bSource) {
        this.aSource = aSource;
        this.bSource = bSource;
    }

    public int aLength() { return this.aSource.length; }
    public int bLength() { return this.bSource.length; }

    public boolean equals(int aIndex, int bIndex) {
        final T aValue = this.aSource[aIndex];
        final T bValue = this.bSource[bIndex];
        if (aValue == null) return bValue == null;
        if (bValue == null) return false;
        return aValue.equals(bValue);
    }
}
