package abstractor.core.constructs;

import java.util.SortedSet;
import java.util.TreeSet;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class InterfaceDesc extends ConstructImp implements TypeDesc {
    public final TreeSet<Ref<Abstract>>      abstracts = new TreeSet<Ref<Abstract>>();
    public final TreeSet<Ref<InterfaceDesc>> inherits  = new TreeSet<Ref<InterfaceDesc>>();
    public       Ref<? extends Construct>    pin;  
    
    public InterfaceDesc() {}

    public InterfaceDesc(SortedSet<Ref<Abstract>> abstracts) {
        this(abstracts, null);
    }

    public InterfaceDesc(SortedSet<Ref<Abstract>> abstracts, Ref<? extends Construct> pin) {
        if (abstracts != null) this.abstracts.addAll(abstracts);
        this.pin = pin;
    }

    public ConstructKind kind() { return ConstructKind.INTERFACE_DESC; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        // Not needed: `approx`, `exact`, `hint`
        obj.put(        "abstracts", indexSet(this.abstracts));
        obj.putNotEmpty("inherits",  indexSet(this.inherits));
        obj.putNotEmpty("pin",       key(this.pin));
        return obj;
    }

    @Override
    public Cmp getCmp(Construct c) {
        return Cmp.or(super.getCmp(c),
            Cmp.deferSet(this.abstracts, () -> ((InterfaceDesc)c).abstracts),
            Cmp.defer(   this.pin,       () -> ((InterfaceDesc)c).pin)
        );
    }   
}
