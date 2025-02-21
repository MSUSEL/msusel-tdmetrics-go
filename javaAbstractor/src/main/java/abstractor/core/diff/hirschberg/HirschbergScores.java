package abstractor.core.diff.hirschberg;

import abstractor.core.diff.core.CostComparator;
import abstractor.core.diff.core.SubComparator;

public class HirschbergScores {
    private int[] front;
    private int[] back;
    private int[] other;

    public HirschbergScores(int length) {
        this.allocateVectors(length);
    }

    private void allocateVectors(int length) {
        if (length < 0) length = 0;
        this.front = new int[length];
        this.back  = new int[length];
        this.other = new int[length];
    }

    public int length() {
        return (this.back == null) ? 0 : this.back.length;
    }

    private void swap() {
        final int[] temp = this.front;
        this.front = this.back;
        this.back = temp;
    }

    private void store() {
        final int[] temp = this.other;
        this.other = this.back;
        this.back = temp;
    }

    private void calculate(CostComparator cont) {
        final int bLen = cont.bLength();
        if (this.length() < bLen+1)
            this.allocateVectors(bLen + 1);

        final int aLen = cont.aLength();
        this.back[0] = 0;
        for (int j = 1; j <= bLen; j++)
            this.back[j] = this.back[j-1] + cont.addCost(j-1);

        for (int i = 1; i <= aLen; i++) {
            final int removeCost = cont.removeCost(i-1);
            this.front[0] = this.back[0] + removeCost;
            for (int j = 1; j <= bLen; j++) {
                this.front[j] = Math.min(
                    this.back[j-1] + cont.substitutionCost(i-1, j-1),
                    Math.min(
                        this.back[j] + removeCost,
                        this.front[j-1] + cont.addCost(j-1)));
            }

            this.swap();
        }
    }

    private int findPivot(int bLength) {
        int index = 0;
        int min = this.other[0] + this.back[bLength];
        for (int j = 1; j <= bLength; j++) {
            final int value = this.other[j] + this.back[bLength-j];
            if (value < min) {
                min = value;
                index = j;
            }
        }
        return index;
    }

    public record SplitResult(int aMid, int bMid) { }

    public SplitResult split(SubComparator cont) {
        final int aLen = cont.aLength();
        final int bLen = cont.bLength();

        final int aMid = aLen / 2;
        this.calculate(cont.sub(0, aMid, 0, bLen));
        this.store();
        this.calculate(cont.sub(aMid, aLen, 0, bLen).reversed());
        final int bMid = this.findPivot(bLen);

        return new SplitResult(aMid, bMid);
    }
}
