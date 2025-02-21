package abstractor.core.diff.comparators;

import abstractor.core.diff.Comparator;

public class StringComparator implements Comparator {
    public final String aSource;
    public final String bSource;

    public StringComparator(String aSource, String bSource) {
        this.aSource = aSource;
        this.bSource = bSource;
    }

    public int aLength() { return this.aSource.length(); }
    public int bLength() { return this.bSource.length(); }

    public boolean equals(int aIndex, int bIndex) {
        return this.aSource.charAt(aIndex) == this.bSource.charAt(bIndex);
    }
}
