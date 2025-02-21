package abstractor.core.diff.core;

public class ReverseComparator implements CostComparator {
    private final CostComparator comp;
    
    public ReverseComparator(CostComparator comp) { this.comp = comp; }
    
    public int aLength() { return this.comp.aLength(); }
    public int bLength() { return this.comp.bLength(); }

    public boolean equals(int aIndex, int bIndex) {
        return this.comp.equals(this.aLength() - 1 - aIndex, this.bLength() - 1 - bIndex);
    }
    
    public int removeCost(int aIndex) {
        return this.comp.removeCost(this.aLength() - 1 - aIndex);
    }

    public int addCost(int bIndex) {
        return this.comp.addCost(this.bLength() - 1 - bIndex);
    }
    
    public int substitutionCost(int aIndex, int bIndex) {
        return this.comp.substitutionCost(this.aLength() - 1 - aIndex, this.bLength() - 1 - bIndex);
    }
}
