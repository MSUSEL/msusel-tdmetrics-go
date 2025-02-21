package abstractor.core.diff.comparators;

import java.util.List;

import abstractor.core.diff.Comparator;

public class ListComparator<T> implements Comparator {
    public final List<T> aSource;
    public final List<T> bSource;

    public ListComparator(List<T> aSource, List<T> bSource) {
        this.aSource = aSource;
        this.bSource = bSource;
    }

    public int aLength() { return this.aSource.size(); }
    public int bLength() { return this.bSource.size(); }

    public boolean equals(int aIndex, int bIndex) {
        final T aValue = this.aSource.get(aIndex);
        final T bValue = this.bSource.get(bIndex);
        if (aValue == null) return bValue == null;
        if (bValue == null) return false;
        return aValue.equals(bValue);
    }
}
