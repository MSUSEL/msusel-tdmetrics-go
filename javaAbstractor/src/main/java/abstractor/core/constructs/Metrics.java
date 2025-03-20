package abstractor.core.constructs;

import java.util.Collections;
import java.util.Set;

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

    public final Set<Construct> invokes;
    public final Set<Construct> reads;
    public final Set<Construct> writes;
    
    public Metrics(Location loc,
        int codeCount, int complexity, int indents, int lineCount,
        boolean getter, boolean setter,
        Set<Construct> invokes, Set<Construct> reads, Set<Construct> writes) {
        this.loc        = loc;
        this.codeCount  = codeCount;
        this.complexity = complexity;
        this.indents    = indents;
        this.lineCount  = lineCount;
        this.getter     = getter;
        this.setter     = setter;
        this.invokes    = Collections.unmodifiableSet(invokes);
        this.reads      = Collections.unmodifiableSet(reads);
        this.writes     = Collections.unmodifiableSet(writes);
    }

    public ConstructKind kind() { return ConstructKind.METRICS; }

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
