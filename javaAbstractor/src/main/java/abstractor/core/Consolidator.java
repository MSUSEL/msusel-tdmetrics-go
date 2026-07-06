package abstractor.core;

import java.util.ArrayList;
import java.util.List;

import abstractor.core.cmp.CmpOptions;
import abstractor.core.constructs.*;
import abstractor.core.log.Logger;

public class Consolidator {
    final public Logger log;
    final public Project proj;

    public Consolidator(Logger log, Project proj) {
        this.log = log;
        this.proj = proj;
    }

    public void consolidate() throws Exception {
        this.log.log("Consolidating all constructs");
        this.log.push();

        this.proj.setToCompareResolved();
        this.proj.setAllIndices();
        int collisions = 0;
        do {
            collisions = this.proj.consolidateCons(this.log);
            this.log.log("Removed " + collisions + " collisions");
        } while(collisions > 0);
        this.proj.setAllIndices();
        this.log.pop();
    }


    

    /**
     * Performs change all the comparison options to use the resolved.
     */
    public void setToCompareResolved() {
        for (Factory<? extends Construct> factory : this.factories)
            factory.setToCompareResolved();
    }

    public int consolidateCons(Logger log) throws Exception {
        int collisions = 0;
        for (Factory<? extends Construct> factory : this.factories)
            collisions += factory.consolidateCons(log);
        return collisions;
    }







    /**
     * Change all the comparison options to use the resolved. This should only
     * be called once all constructs have been added to the factory and
     * the code is transitioning from reading in the source to prepare to write
     * the abstraction results. This also clears out sets that may become a
     * problem once the references are resolved.
     */
    public <T extends Construct> void setToCompareResolved(Factory<T> factory) {
        final CmpOptions options = CmpOptions.resolvedCmp();
        for (T con : this.conSet) con.setCmpOptions(options);
        for (Ref<T> ref : this.refSet) ref.setCmpOptions(options);
        
        // The non-element references is no longer useful but the changed comparisons
        // could cause issues if someone tried to use that set since it is no longer
        // in sorted order, so just clear it out.
        this.nonElemRef.clear();
    }

    public int consolidateCons(Logger log) throws Exception {
        final int size = this.conSet.size();
        if (size <= 1) return 0; // Nothing to consolidate

        // Copy all cons to a list and clear the set so that only
        // the unique cons can be re-added in the new sort order.
        final List<T> conList = new ArrayList<T>(this.conSet);
        if (this.isSorted(conList)) return 0;
        conList.sort(Comparator.naturalOrder()); 

        // Perform the in-place "squeeze" to remove duplicates.
        int collisions = 0;
        int uniqueCount = 1;
        T unique = conList.get(0);
        for (int i = 1; i < size; i++) {
            final T con = conList.get(i);
            // Compare against the last confirmed unique node.
            if (!unique.equals(con)) {
                // No conflict found, so shift the construct up in the set.
                if (uniqueCount != i) {
                    conList.set(uniqueCount, con);
                    con.setIndex(uniqueCount);
                }
                unique = conList.get(uniqueCount);
                uniqueCount++;
                continue;
            }

            // Found another construct that is equal so move all references over
            // to the existing construct since the duplicate is about to be removed.
            collisions++;
            for (Ref<T> ref : this.refSet) {
                if (con.equals(ref.getResolved()))
                    ref.setResolved(unique);
            }
            con.setIndex(-100);
        }

        // Rebuild the TreeSet using the squeezed sub-list.
        this.conSet.clear();
        this.conSet.addAll(conList.subList(0, uniqueCount));
        return collisions;
    }

    private boolean isSorted(List<T> list) {
        if (list.size() <= 1) return true;
        Iterator<T> iter = list.iterator();
        T cur, prev = iter.next();
        while (iter.hasNext()) {
            cur = iter.next();
            if (prev.compareTo(cur) > 0) return false;
            prev = cur;
        }
        return true;
    }


}
