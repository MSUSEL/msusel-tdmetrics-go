package abstractor.core.constructs;

import java.util.List;

import spoon.reflect.declaration.CtField;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class Metrics extends ConstructImp {
    public final Location loc;

    public final int codeCount;
    public final int complexity;
    public final int indents;
    public final int lineCount;

    public final boolean getter;
    public final boolean setter;

    public final List<Method> invokes;
    public final List<TypeDesc> reads;
    public final List<TypeDesc> writes;
    
    public Metrics(CtField<?> src, Location loc,
        int codeCount, int complexity, int indents, int lineCount,
        boolean getter, boolean setter,
        List<Method> invokes, List<TypeDesc> reads, List<TypeDesc> writes) {
        super(src);
        this.loc        = loc;
        this.codeCount  = codeCount;
        this.complexity = complexity;
        this.indents    = indents;
        this.lineCount  = lineCount;
        this.getter     = getter;
        this.setter     = setter;
        this.invokes    = invokes;
        this.reads      = reads;
        this.writes     = writes;
    }

    public String kind() { return "metrics"; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.putNotEmpty("loc",        loc.toJson(h));
        obj.putNotEmpty("codeCount",  this.codeCount);
        obj.putNotEmpty("complexity", this.complexity);
        obj.putNotEmpty("indents",    this.indents);
        obj.putNotEmpty("lineCount",  this.lineCount);
        obj.putNotEmpty("getter",     this.getter);
        obj.putNotEmpty("setter",     this.setter);
        obj.putNotEmpty("invokes",    keySet(this.invokes));
        obj.putNotEmpty("reads",      keySet(this.reads));
        obj.putNotEmpty("writes",     keySet(this.writes));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.loc, () -> ((Metrics)c).loc)
        );
    }   
}
