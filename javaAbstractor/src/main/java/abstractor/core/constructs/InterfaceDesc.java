package abstractor.core.constructs;

import java.util.SortedSet;
import java.util.TreeSet;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class InterfaceDesc extends ConstructImp implements TypeDesc {
    public final SortedSet<Abstract> abstracts;
    public final SortedSet<InterfaceDesc> inherits;
    public final Construct pin;  

    public InterfaceDesc(SortedSet<Abstract> abstracts) {
        this(abstracts, null);
    }

    public InterfaceDesc(SortedSet<Abstract> abstracts, Construct pin) {
        this.abstracts = unmodifiableSortedSet(abstracts);
        this.inherits = new TreeSet<InterfaceDesc>();
        this.pin = pin;
    }

    public ConstructKind kind() { return ConstructKind.INTERFACE_DESC; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        // Not needed: `approx`, `exact`, `hint`
        obj.put("abstracts",        indexSet(this.abstracts));
        obj.putNotEmpty("inherits", indexSet(this.inherits));
        obj.putNotEmpty("pin",      key(this.pin));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.deferSet(this.abstracts, () -> ((InterfaceDesc)c).abstracts),
            Cmp.defer(this.pin,          () -> ((InterfaceDesc)c).pin)
        );
    }   
}
