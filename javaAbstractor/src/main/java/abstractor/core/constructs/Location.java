package abstractor.core.constructs;

import java.io.File;

import spoon.reflect.cu.SourcePosition;

import abstractor.core.json.*;

public class Location implements Jsonable, Comparable<Location> {
    public final SourcePosition pos;
    public int offset;

    public Location(SourcePosition pos) {
        this.pos = pos;
    }

    public JsonNode toJson(JsonHelper h) {
        return JsonValue.of(this.offset);
    }

    public boolean isValid() {
        return this.pos.isValidPosition();
    }

    public String getPath() {
        final File f = this.pos.getFile();
        return f != null ? f.getPath() : "unknown";
    }

    public int getLine() {
        final int line = this.pos.getLine();
        return line < 1 ? 1 : line;
    }

    @Override
    public int compareTo(Location o) {
        if (o == null) return 1;
        if (this.pos.isValidPosition()) return o.pos.isValidPosition() ? 0 : -1;
        if (o.pos.isValidPosition()) return 1;

        int cmp = this.pos.toString().compareTo(o.pos.toString());
        if (cmp != 0) return cmp;
        return Integer.compare(this.pos.getSourceStart(), o.pos.getSourceStart());
    }
}
