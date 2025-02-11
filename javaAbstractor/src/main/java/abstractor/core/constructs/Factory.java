package abstractor.core.constructs;

import java.util.TreeSet;

import abstractor.core.json.*;

public class Factory<T extends Construct> implements Jsonable {
    private final TreeSet<T> set = new TreeSet<T>();

    public T findWithSource(Object source) {
        for (T t : this.set) {
            if (t.source() == source) return t;
        }
        return null;
    }

    public boolean containsSource(Object source) {
        for (T t : this.set) {
            if (t.source() == source) return true;
        }
        return false;
    }

    public T tryAdd(T t) {
        T other = this.set.floor(t);
        if (other == t) return other;
        this.set.add(t);
        return t;
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
