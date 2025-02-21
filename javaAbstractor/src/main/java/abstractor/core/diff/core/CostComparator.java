package abstractor.core.diff.core;

import abstractor.core.diff.Comparator;

public interface CostComparator extends Comparator {
    int removeCost(int aIndex);
    int addCost(int bIndex);
    int substitutionCost(int aIndex, int bIndex);
}
