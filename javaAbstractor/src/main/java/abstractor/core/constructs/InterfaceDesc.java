package abstractor.core.constructs;

import java.util.SortedSet;
import java.util.TreeSet;

import spoon.reflect.declaration.CtField;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class InterfaceDesc extends ConstructImp implements TypeDesc {
    public final SortedSet<Abstract> abstracts;
    public final TreeSet<InterfaceDesc> inherits;
    public final Construct pin;  

    public InterfaceDesc(CtField<?> src, SortedSet<Abstract> abstracts) {
        this(src, abstracts, null);
    }

    public InterfaceDesc(CtField<?> src, SortedSet<Abstract> abstracts, Construct pin) {
        super(src);
        this.abstracts = unmodifiableSortedSet(abstracts);
        this.inherits = new TreeSet<InterfaceDesc>();
        this.pin = pin;
    }

    public String kind() { return "interfaceDesc"; }

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
