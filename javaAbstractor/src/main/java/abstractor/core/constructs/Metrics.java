package abstractor.core.constructs;

import java.util.SortedSet;
import java.util.TreeSet;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class Metrics extends ConstructImp {
    public Location loc;

    public int codeCount;
    public int complexity;
    public int indents;
    public int lineCount;

    public boolean getter;
    public boolean setter;

    public SortedSet<Construct> invokes;
    public SortedSet<Construct> reads;
    public SortedSet<Construct> writes;

    public Metrics() {
        this.invokes = new TreeSet<Construct>();
        this.reads   = new TreeSet<Construct>();
        this.writes  = new TreeSet<Construct>();
    }

    public Metrics(Location loc,
        int codeCount, int complexity, int indents, int lineCount,
        boolean getter, boolean setter,
        SortedSet<Construct> invokes, SortedSet<Construct> reads, SortedSet<Construct> writes) {
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
