package abstractor.core.diff.wagner;

import java.util.Iterator;

import abstractor.core.diff.Comparator;
import abstractor.core.diff.DiffStep;
import abstractor.core.diff.StepType;
import abstractor.core.diff.core.Algorithm;
import abstractor.core.diff.core.CostComparator;
import abstractor.core.iter.IterableWrapper;
import abstractor.core.iter.YieldIterator;

public class Wagner implements Algorithm {
    private int[] costs;

    public Wagner(int size) {
        this.costs = new int[0];
        this.allocateMatrix(size);
    }

    private void allocateMatrix(int size) {
        if (size > 0) this.costs = new int[size];
    }

    @Override
    public boolean noResizeNeeded(CostComparator comp) {
        return this.costs.length >= comp.aLength()*comp.bLength();
    }

    @Override
    public Iterable<DiffStep> diff(CostComparator comp) {
        int size = comp.aLength()*comp.bLength();
        if (this.costs.length < size) this.allocateMatrix(size);
        this.setCosts(comp);
        return new IterableWrapper<DiffStep>(() -> this.walkPath(comp));
    }

    private void setCosts(CostComparator comp) {
        final int aLen = comp.aLength();
        final int bLen = comp.bLength();

        final int start = comp.substitutionCost(0, 0);
        this.costs[0] = start;

        for (int i = 1, value = start; i < aLen; i++) {
            value = Math.min(value + 1,
                i+comp.substitutionCost(i, 0));
            this.costs[i] = value;
        }

        for (int j = 1, k = aLen, value = start; j < bLen; j++, k+=aLen) {
            value = Math.min(value + 1,
                j+comp.substitutionCost(0, j));
            this.costs[k] = value;
        }

        for (int j = 1, k = aLen+1, k2 = 1, k3 = 0; j < bLen; j++, k++, k2++, k3++) {
            for (int i = 1, value = this.costs[k-1]; i < aLen; i++, k++, k2++, k3++) {
                value = Math.min(value+1,
                    Math.min(this.costs[k2] + 1,
                    this.costs[k3] + comp.substitutionCost(i, j)));
                this.costs[k] = value;
            }
        }
    }

    private int getCost(int i, int j, int aLen) {
        if (i < 0) return j + 1;
        if (j < 0) return i + 1;
        return this.costs[i + j*aLen];
    }

    private Iterator<DiffStep> walkPath(Comparator comp) {
        enum StepType {
            none,
            removed,
            added,
            equal,
            substitute
        }

        class Indices {
            int i = comp.aLength() - 1;
            int j = comp.bLength() - 1;
        }

        final Indices n = new Indices();
        final int aLen = comp.aLength();
        return new YieldIterator<DiffStep>(
            () -> n.i >= 0 && n.j >= 0,
            (YieldIterator.Yield<DiffStep> y) -> {
                final int aCost = this.getCost(n.i-1, n.j,   aLen);
                final int bCost = this.getCost(n.i,   n.j-1, aLen);
                final int cCost = this.getCost(n.i-1, n.j-1, aLen);
                final int minCost = Math.min(aCost, Math.min(bCost, cCost));

                StepType stepType = StepType.none;
                if (aCost == minCost)
                    stepType = StepType.removed;
                
                if (bCost == minCost)
                    stepType = StepType.added;
                
                if (cCost == minCost) {
                    if (comp.equals(n.i, n.j)) {
                        stepType = StepType.equal;
                    } else if (stepType == StepType.none)
                        stepType = StepType.substitute;
                }

                switch (stepType) {
                    case none: break;
                    case removed:
                        n.i--;
                        y.yield(DiffStep.Removed(1));
                        break;
                
                    case added:
                        n.j--;
                        y.yield(DiffStep.Added(1));
                        break;
                    
                    case equal:
                        n.i--;
                        n.j--;
                        y.yield(DiffStep.Equal(1));
                        break;
                    
                    case substitute:
                        n.i--;
                        n.j--;
                        y.yield(DiffStep.Added(1));
                        y.yield(DiffStep.Removed(1));
                        break;
                }
            },
            (YieldIterator.Yield<DiffStep> y) -> {
                y.yield(DiffStep.Removed(n.i + 1));
                y.yield(DiffStep.Added(n.j + 1));
            });
    }
}
