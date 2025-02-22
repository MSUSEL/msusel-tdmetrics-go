package abstractor.core.constructs;

import java.util.TreeSet;

import spoon.reflect.declaration.CtElement;

import abstractor.core.json.*;

public class Factory<T extends Construct> implements Jsonable {
    private final TreeSet<T> set = new TreeSet<T>();

    public T findWithSource(CtElement source) {
        for (T t : this.set) {
            if (t.source() == source) return t;
        }
        return null;
    }

    public TryAddResult<T> tryAdd(T t) {
        T other = this.set.floor(t);
        if (other == t) return new TryAddResult<T>(other, true);
        this.set.add(t);
        return new TryAddResult<T>(t, false);
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
