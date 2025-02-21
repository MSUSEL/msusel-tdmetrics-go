package abstractor.core.diff;

public interface Comparator {
    int aLength();
    int bLength();
    boolean equals(int aIndex, int bIndex);
}
