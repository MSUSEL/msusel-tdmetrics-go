package abstractor.core;

import java.util.HashMap;
import java.util.TreeSet;

import abstractor.core.log.Logger;

public class Indexer {
    private final Logger log;
    private final HashMap<Object, Integer> indices = new HashMap<Object, Integer>();

    public Indexer(Logger log, Project proj) {
        this.log = log;
        this.setIndices(proj.packages);
        this.setIndices(proj.objects);
        this.setIndices(proj.interfaceDecls);
    }

    private void setIndices(TreeSet<?> constructs) {
        int index = 1;
        for (Object o : constructs) {
            if (this.indices.containsKey(o))
                this.log.error("Index already set for " + o + ".");
            else {
                this.indices.put(o, index);
                index++;
            }
        }
    }

    public int indexOf(Object o) {
        Integer index = this.indices.get(o);
        if (index == null) {
            this.log.error("Failed to find index for " + o + ".");
            return 0;
        }
        return index;
    }
}
