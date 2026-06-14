package abstractor.core;

import java.util.*;

import spoon.reflect.declaration.CtElement;

import abstractor.core.cmp.*;
import abstractor.core.constructs.*;
import abstractor.core.json.*;

public class ElementKey implements Comparable<ElementKey>, CmpGetter<ElementKey>, Jsonable {
    static private HashMap<CtElement, String> elemOrder = new HashMap<>();

    static private String getElemOrderKey(CtElement elem) {
        if (elem == null) return "null";
        return elemOrder.computeIfAbsent(elem, k -> {
            // Use the hash and position to try to order the elements as consistently as possible.
            // There shouldn't be randomness in how the AST is processed but I don't want
            // to fighting bugs because of elements being added in random order.
            final String pos = k.getPosition().toString();
            return k.hashCode() + "-" + pos + "-" + elemOrder.size();
        });
    }

    final private CtElement elem;
    final private List<Ref<? extends TypeDesc>> typeArgs = new ArrayList<>();

    public ElementKey(CtElement elem) {
        this.elem = elem;
    }

    public ElementKey(CtElement elem, List<Ref<? extends TypeDesc>> typeArgs) {
        this.elem = elem;
        if (typeArgs != null) this.typeArgs.addAll(typeArgs);
    }

    @Override
    public int hashCode() {
        return this.elem != null ? this.elem.hashCode() : 0;
    }

    @Override
    public Cmp getCmp(ElementKey other, CmpOptions options) {
        return Cmp.or("ElementKey",
            Cmp.defer(getElemOrderKey(this.elem), () -> getElemOrderKey(other.elem), "elem"),
            Cmp.deferList(this.typeArgs, () -> other.typeArgs, "typeArgs")
        );
    }

    public int compareTo(ElementKey c) {
        return Cmp.compareTo(this, c, new CmpOptions());
    }

    @Override
    public boolean equals(Object obj) {
        return obj instanceof ElementKey c && this.compareTo(c) == 0;
    }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = new JsonObject();
        obj.put("elem", this.elem.toStringDebug());
        if (this.typeArgs != null) {
            JsonArray arr = new JsonArray();
            for (Ref<? extends TypeDesc> ta : this.typeArgs)
                arr.add(ta.toJson(h));
            obj.putNotEmpty("typeArgs", arr);
        }
        return obj;
    }

    @Override
    public String toString() {
        JsonHelper jh = new JsonHelper();
        jh.writeKinds     = true;
        jh.writeIndices   = true;
        jh.writeRefs      = true;
        jh.refSkipResolve = true;
        return JsonFormat.Relaxed().format(this.toJson(jh));
    }
}
