package abstractor.core.constructs;

import java.util.HashMap;
import java.util.TreeSet;

import spoon.reflect.declaration.CtElement;

import abstractor.core.json.*;
import abstractor.core.log.Logger;

public class Factory<T extends Construct> implements Jsonable {
    private final TreeSet<T> set;
    private final HashMap<CtElement, T> byElem;

    public Factory() {
        this.set = new TreeSet<T>();
        this.byElem = new HashMap<CtElement, T>();
    }

    public interface ConstructCreator<T extends Construct> { T create(); }

    public interface FinishConstruct<T extends Construct> { void finish(T con); }

    public T create(Logger log, CtElement elem, String title,
        ConstructCreator<T> c, FinishConstruct<T> f) {
        final T existing = this.byElem.get(elem);
        if (existing != null) return existing;
        
        if (log != null) {
            log.log("Adding " + title);
            log.push();
        }
        try {
            final T newCon = c.create();

            T other = this.set.floor(newCon);
            if (other == newCon) return other;
            this.set.add(newCon);
            this.byElem.put(elem, newCon);

            if (f != null) f.finish(newCon);
            return newCon;
        } finally {
            if (log != null) log.pop();
        }
    }
    
    public T create(Logger log, CtElement elem, String title,
        ConstructCreator<T> c) {
        return this.create(log, elem, title, c, null);
    }

    public void setIndices() {
        int index = 1;
        for (Construct o : this.set) {
            o.setIndex(index);
            index++;
        }
    }

    public JsonNode toJson(JsonHelper h) {
        JsonArray array = new JsonArray();
        for (T t : this.set) array.add(t.toJson(h));
        return array;
    }
}
