package abstractor.core.tools;

import java.util.*;

import abstractor.core.cmp.*;
import abstractor.core.constructs.*;
import abstractor.core.json.JsonFormat;
import abstractor.core.json.JsonHelper;
import abstractor.core.log.*;
import abstractor.core.require.Require;

public class Consolidator {
    static public boolean logConsolidations = false;

    final public Logger log;
    final public Project proj;

    public Consolidator(Logger log, Project proj) {
        this.log = log;
        this.proj = proj;
    }

    public void consolidate() throws Exception {
        this.log.log("Consolidating all constructs");
        this.log.push();

        this.setToCompareResolved();
        this.pullPinsOnAnys();
        this.pullPinsOnEmptyStructs();
        this.proj.setAllIndices();

        int collisions;
        do {
            collisions = 0;
            for (Factory<? extends Construct> factory : this.proj.factories)
                collisions += this.consolidate(factory);
            this.log.log("Removed " + collisions + " collisions");
        } while(collisions > 0);
        this.proj.setAllIndices();
        this.log.pop();
    }

    private String conToString(Construct con) {
        final JsonHelper jh = new JsonHelper();
        jh.writeKinds     = true;
        jh.writeIndices   = true;
        jh.writeRefs      = true;
        jh.refSkipResolve = false;
        return JsonFormat.Inline().format(con.toJson(jh));
    }

    /**
     * Removes the pins from all interfaces that have no abstracts
     * (and has no inherits) but have a pin.
     * Without the pin, the consolidator will dedup these interfaces.
     * 
     * This is because all "any" interfaces do not need to be unique.
     */
    private void pullPinsOnAnys() {
        for (InterfaceDesc i : this.proj.interfaceDescs.getConSet()) {
            if (i.pin != null && i.abstracts.size() <= 0 && i.inherits.size() <= 0) {
                if (logConsolidations)
                    this.log.log("unpinning interface description from " + this.conToString(i.pin));
                i.pin = null;
            }
        }
    }

    /**
     * Same as pullPinsOnAnys but for structs.
     */
    private void pullPinsOnEmptyStructs() {
        for (StructDesc s : this.proj.structDescs.getConSet()) {
            if (s.pin != null && s.fields.size() <= 0) {
                if (logConsolidations)
                    this.log.log("unpinning struct description from " + this.conToString(s.pin));
                s.pin = null;
            }
        }
    }

    /**
     * Performs change all the comparison options to use the resolved.
     */
    private void setToCompareResolved() {
        for (Factory<? extends Construct> factory : this.proj.factories)
            this.setToCompareResolved(factory);
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
        for (T con : factory.getConSet()) con.setCmpOptions(options);
        for (Ref<T> ref : factory.getRefSet()) ref.setCmpOptions(options);
        
        // The non-element references is no longer useful but the changed comparisons
        // could cause issues if someone tried to use that set since it is no longer
        // in sorted order, so just clear it out.
        factory.getNonElemRefSet().clear();
    }

    public <T extends Construct> int consolidate(Factory<T> factory) throws Exception {
        final TreeSet<T> conSet = factory.getConSet();
        final int origSize = conSet.size();
        if (origSize <= 1) return 0; // Nothing to consolidate

        // Copy all cons to a list and clear the set so that only
        // the unique cons can be re-added in the new sort order.
        final List<T> conList = new ArrayList<T>(conSet);
        if (isSorted(conList)) return 0;
        conList.sort(Comparator.naturalOrder()); 
        final List<T> squeezedList = this.squeeze(conList);

        // Rebuild the TreeSet using the squeezed sub-list.
        conSet.clear();
        conSet.addAll(squeezedList);
        final int collisions = origSize - conSet.size();

        // If there were collsisions (not just sorting), update all the references
        // so they use the unique values. The non-unique values are already be
        // equal to the unique ones so we can just affirm they are using the unique ones.
        if (collisions > 0) {
            for (Ref<T> ref : factory.getRefSet()) {
                final T oldRes = ref.getResolved();
                final T unique = conSet.floor(oldRes);
                // Since duplicates have been removed then oldRes should be
                // equal to the unique value, otherwise it means that oldRes
                // or a con equal to it is not in the conSet.
                Require.equal(oldRes, unique, "expected all constructs to be found in conSet");
                ref.setResolved(unique);
            }
        }
        return collisions;
    }

    /**
     * Performs an in-place "squeeze" to remove duplicates from given list.
     */
    private <T extends Construct> List<T> squeeze(List<T> conList) {
        final int size = conList.size();
        int uniqueCount = 1;
        T unique = conList.get(0);
        for (int i = 1; i < size; i++) {
            final T con = conList.get(i);
            // Compare against the last confirmed unique node.
            if (unique.equals(con)) {
                // Found another construct that is equal.
                if (logConsolidations)
                    this.log.log("found duplicate: " + this.conToString(con));
                con.setIndex(-100);
                continue;
            }

            // No duplicate found, so shift the construct up in the list.
            if (uniqueCount != i) {
                conList.set(uniqueCount, con);
                con.setIndex(uniqueCount);
            }
            unique = conList.get(uniqueCount);
            uniqueCount++;
        }
        return conList.subList(0, uniqueCount);
    }

    static private <T extends Construct> boolean isSorted(List<T> list) {
        if (list.size() <= 1) return true;
        Iterator<T> iter = list.iterator();
        T cur, prev = iter.next();
        while (iter.hasNext()) {
            cur = iter.next();
            if (prev.compareTo(cur) >= 0) return false;
            prev = cur;
        }
        return true;
    }
}
