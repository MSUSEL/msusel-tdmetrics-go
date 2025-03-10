package abstractor.core.constructs;

import java.io.File;

import spoon.reflect.cu.SourcePosition;
import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class Location implements Jsonable, Comparable<Location> {
    public final String path;
    public final int line;
    public int offset;

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

    public Location(String path, int line) {
       this.path = path;
       this.line = line;
    }
    
    public JsonNode toJson(JsonHelper h) { return JsonValue.of(this.offset); }

    public boolean isValid() { return this.line > 0; }

    @Override
    public int compareTo(Location o) {
        if (o == null) return 1;
        return Cmp.or(
            Cmp.defer(this.path, () -> o.path),
            Cmp.defer(this.line, () -> o.line)
        );
    }

    @Override
    public boolean equals(Object o) {
        return o != null &&
            o instanceof Location lo &&
            this.path.equals(lo.path) &&
            this.line == lo.line;
    }
    
    @Override
    public String toString() { return this.path + ":" + this.line; }
}
