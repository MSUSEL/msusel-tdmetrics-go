package abstractor.core.diff.hirschberg;

import java.util.Stack;

import abstractor.core.diff.DiffStep;
import abstractor.core.diff.core.Algorithm;
import abstractor.core.diff.core.CostComparator;
import abstractor.core.diff.core.SubComparator;
import abstractor.core.diff.core.SubComparator.ReduceResult;
import abstractor.core.diff.hirschberg.HirschbergScores.SplitResult;
import abstractor.core.iter.*;

public class Hirschberg implements Algorithm {
    private final HirschbergScores scores;
    private final Algorithm hybrid;
    
    public Hirschberg(int length) {
        this.scores = new HirschbergScores(length);
        this.hybrid = null;
    }
    
    public Hirschberg(int length, Algorithm hybrid) {
        this.scores = new HirschbergScores(length);
        this.hybrid = hybrid;
    }

    public boolean noResizeNeeded(CostComparator comp) {
        return this.scores.length() >= comp.bLength() + 1;
    }

    private record SubRem(SubComparator sub, int remainder) { }

    public Iterable<DiffStep> diff(CostComparator comp) {
        return new IterableWrapper<DiffStep>(() -> {
            Stack<SubRem> stack = new Stack<SubRem>();
            stack.push(new SubRem(new SubComparator(comp), 0));
            
            return new YieldIterator<DiffStep>(
                () -> !stack.isEmpty(),
                (YieldIterator.Yield<DiffStep> y)-> {
                    final SubRem pair = stack.pop();
                    SubComparator cur = pair.sub;
                    final int remainder = pair.remainder;
        
                    if (remainder > 0) y.yield(DiffStep.Equal(remainder));
                    if (cur == null) return;
        
                    ReduceResult red = cur.reduce();
                    cur = red.sub();
                    if (red.back() > 0) y.yield(DiffStep.Equal(red.back()));
                    stack.push(new SubRem(null, red.front()));
        
                    if (cur.isEndCase()) {
                        for (DiffStep step : cur.endCase())
                            y.yield(step);
                            return;
                    }
        
                    if (this.hybrid != null && this.hybrid.noResizeNeeded(cur)) {
                        for (DiffStep step : this.hybrid.diff(cur))
                            y.yield(step);
                            return;
                    }
        
                    final int aLen = cur.aLength();
                    final int bLen = cur.bLength();
                    SplitResult mid = this.scores.split(cur);
                    stack.push(new SubRem(cur.sub(0, mid.aMid(), 0, mid.bMid()), 0));
                    stack.push(new SubRem(cur.sub(mid.aMid(), aLen, mid.bMid(), bLen), 0));
                });
        });
    }
}
