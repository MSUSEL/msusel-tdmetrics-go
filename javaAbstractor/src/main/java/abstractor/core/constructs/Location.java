package abstractor.core.constructs;

import java.io.File;

import spoon.reflect.cu.SourcePosition;
import abstractor.core.cmp.Cmp;
import abstractor.core.cmp.CmpGetter;
import abstractor.core.cmp.CmpOptions;
import abstractor.core.json.*;

public class Location implements Jsonable, Comparable<Location>, CmpGetter<Location> {
    public String path;
    public int    line;

    // This is the accumulative offset this location exists relative to all locations.
    public int offset;

    public Location() {}

    public Location(String path, int line) {
       this.path = path;
       this.line = line;
    }
    
    public Location(SourcePosition pos) {
        String path = "unknown";
        int line = 0;

        if (pos.isValidPosition()) {
            final File f = pos.getFile();
            if (f != null) path = f.getPath();
            line = pos.getLine();
        }

       this.path = path;
       this.line = line;
    }
    
    public JsonNode toJson(JsonHelper h) { return JsonValue.of(this.offset); }

    public boolean isValid() { return this.line > 0; }

    public int compareTo(Location o) { return Cmp.compareTo(this, o, null); }

    public Cmp getCmp(Location o, CmpOptions options) {
        return Cmp.or(
            Cmp.defer(this.path, () -> o.path),
            Cmp.defer(this.line, () -> o.line)
        );
    }

    @Override
    public boolean equals(Object o) {
        return o != null && o instanceof Location lo && this.compareTo(lo) == 0;
    }
    
    @Override
    public String toString() { return this.path + ":" + this.line; }
}
