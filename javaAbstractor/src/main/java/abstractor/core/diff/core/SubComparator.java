package abstractor.core.diff.core;

import java.util.ArrayList;

import abstractor.core.diff.DiffStep;
import abstractor.core.iter.Iter;

public class SubComparator implements CostComparator {
    private final CostComparator comp;
    private final int aOffset;
    private final int aLen;
    private final int bOffset;
    private final int bLen;
    
    private SubComparator(CostComparator comp, int aOffset, int aLength, int bOffset, int bLength) {
        this.comp    = comp;
        this.aOffset = aOffset;
        this.aLen    = aLength;
        this.bOffset = bOffset;
        this.bLen    = bLength;
    }

    public SubComparator(CostComparator comp) {
        this(comp, 0, comp.aLength(), 0, comp.bLength());
    }

    public SubComparator sub(int aLow, int aHigh, int bLow, int bHigh) {
        return new SubComparator(this.comp, this.aOffset+aLow, aHigh-aLow, this.bOffset+bLow, bHigh-bLow);
    }

    public ReverseComparator reversed() {
        return new ReverseComparator(this);
    }
    
    public int aLength() { return this.aLen; }
    public int bLength() { return this.bLen; }

    public boolean equals(int aIndex, int bIndex) {
        return this.comp.equals(this.aOffset + aIndex, this.bOffset + bIndex);
    }

    public int removeCost(int aIndex) {
        return this.comp.removeCost(this.aOffset + aIndex);
    }

    public int addCost(int bIndex) {
        return this.comp.addCost(this.bOffset + bIndex);
    }
    
    public int substitutionCost(int aIndex, int bIndex) {
        return this.comp.substitutionCost(this.aOffset + aIndex, this.bOffset+ bIndex);
    }

    public boolean isEndCase() {
        return this.aLen <= 1 || this.bLen <= 1;
    }
    
    private int matchFront(int width) {
        int i = this.aOffset;
        int j = this.bOffset;
        for (int front = 0; front < width; front++) {
            if (!this.comp.equals(i, j)) return front;
            i++;
            j++;
        }
        return width;
    }

    private int matchBack(int width) {
        int i = this.aLen - 1 + this.aOffset;
        int j = this.bLen - 1 + this.bOffset;
        for (int back = 0; back < width; back++) {
            if (!this.comp.equals(i, j)) return back;
            i--;
            j--;
        }
        return width;
    }

    public record ReduceResult(SubComparator sub, int front, int back) { }

    public ReduceResult reduce() {
        int max   = Math.min(this.aLen, this.bLen);
        int front = this.matchFront(max);
        int back  = this.matchBack(max-front);

        SubComparator sub = new SubComparator(this.comp,
            this.aOffset+front, this.aLen-front-back,
            this.bOffset+front, this.bLen-front-back);

        return new ReduceResult(sub, front, back);
    }

    private Iterable<DiffStep> aEdge() {
        ArrayList<DiffStep> result = new ArrayList<DiffStep>();
        if (this.aLen <= 0) {
            result.add(DiffStep.Added(this.bLen));
            return result;
        }

        int split = -1;
        for (int i = 0, j = this.bOffset; i < this.bLen; i++, j++) {
            if (this.comp.equals(0, j)) {
                split = j;
                break;
            }
        }

        if (split < 0) {
            result.add(DiffStep.Added(this.bLen));
            result.add(DiffStep.Removed(1));
        } else {
            result.add(DiffStep.Added(this.bLen - split - 1));
            result.add(DiffStep.Equal(1));
            result.add(DiffStep.Added(split));
        }
        return result;
    }

    private Iterable<DiffStep> bEdge() {
        ArrayList<DiffStep> result = new ArrayList<DiffStep>();
        if (this.bLen <= 0) {
            result.add(DiffStep.Removed(this.aLen));
            return result;
        }

        int split = -1;
        for (int i = 0, j = this.aOffset; i < this.aLen; i++, j++) {
            if (this.comp.equals(j, 0)) {
                split = j;
                break;
            }
        }

        if (split < 0) {
            result.add(DiffStep.Added(1));
            result.add(DiffStep.Removed(this.aLen));
        } else {
            result.add(DiffStep.Removed(this.aLen - split - 1));
            result.add(DiffStep.Equal(1));
            result.add(DiffStep.Removed(split));
        }
        return result;
    }

    public Iterable<DiffStep> endCase() {
        if (this.aLen <= 1) return this.aEdge();
        if (this.bLen <= 1) return this.bEdge();
        return Iter.<DiffStep>Empty();
    }
}
