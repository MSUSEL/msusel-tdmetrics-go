package abstractor.core.diff.core;

import abstractor.core.diff.Comparator;

public class WrapComparator implements CostComparator {
    public final Comparator comp;

    public WrapComparator(Comparator comp) { this.comp = comp; }
    
    public int aLength() { return this.comp.aLength(); }
    public int bLength() { return this.comp.bLength(); }

    public boolean equals(int aIndex, int bIndex) {
        return this.comp.equals(aIndex, bIndex);
    }

    public int removeCost(int aIndex) { return 1; }
    public int addCost(int bIndex) { return 1; }
    public int substitutionCost(int aIndex, int bIndex) {
        return this.equals(aIndex, bIndex) ? 0 : 2;
    }
}
