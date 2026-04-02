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

    public TreeSet<Ref<? extends Construct>> invokes = new TreeSet<Ref<? extends Construct>>();
    public TreeSet<Ref<? extends Construct>> reads   = new TreeSet<Ref<? extends Construct>>();
    public TreeSet<Ref<? extends Construct>> writes  = new TreeSet<Ref<? extends Construct>>();

    public Metrics() {}

    public Metrics(Location loc) { this.loc = loc; }

    public Metrics(Location loc, int codeCount, int complexity, int indents, int lineCount,
        boolean getter, boolean setter, SortedSet<Ref<? extends Construct>> invokes,
        SortedSet<Ref<? extends Construct>> reads, SortedSet<Ref<? extends Construct>> writes) {
        this.loc        = loc;
        this.codeCount  = codeCount;
        this.complexity = complexity;
        this.indents    = indents;
        this.lineCount  = lineCount;
        this.getter     = getter;
        this.setter     = setter;
        if (invokes != null) this.invokes.addAll(invokes);
        if (reads   != null) this.reads  .addAll(reads);
        if (writes  != null) this.writes .addAll(writes);
    }

    public boolean hasBody() { return this.lineCount > 0; }

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
