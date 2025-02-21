package abstractor.core.diff.core;

import abstractor.core.diff.DiffStep;

public interface Algorithm {
    boolean noResizeNeeded(CostComparator comp);
    Iterable<DiffStep> diff(CostComparator comp);
}
