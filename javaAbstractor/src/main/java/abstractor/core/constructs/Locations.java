package abstractor.core.constructs;

import java.io.File;
import java.util.TreeSet;
import java.util.Map.Entry;
import java.util.TreeMap;

import abstractor.core.json.*;
import spoon.reflect.cu.SourcePosition;

public class Locations implements Jsonable {
    public class Location implements Jsonable, Comparable<Location> {
        public final SourcePosition pos;
        public int offset;

        private Location(SourcePosition pos) {
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
            return this.pos.getLine();
        }

        @Override
        public int compareTo(Location o) {
            return this.toString().compareTo(o.toString());
        }
    }

    private final TreeSet<Location> locs = new TreeSet<Location>();
    private final TreeMap<String, Integer> offsets = new TreeMap<String, Integer>();
    
    public Location create(SourcePosition pos) {
        for (Location loc : this.locs) {
            if (loc.pos == pos) return loc;
        }

        Location loc = new Location(pos);
        this.locs.add(loc);
        return loc;
    }

    public void prepareForOutput() {
        TreeMap<String, Integer> maximums = new TreeMap<String, Integer>();
        for (Location loc : this.locs) {
            if (!loc.isValid()) continue;
            final String path = loc.getPath();
            final int line = loc.getLine();
            final Integer iMax = maximums.get(path);
            final int max = iMax == null ? line : Integer.max((int)iMax, line);
            maximums.put(path, max);
        }

        this.offsets.clear();
        int offset = 0;
        for (String path : maximums.keySet()) {
            this.offsets.put(path, offset);
            offset += maximums.get(path);
        }

        for (Location loc : this.locs) {
            if (!loc.isValid()) continue;
            final String path = loc.getPath();
            final int fileOffset = this.offsets.get(path);
            final int line = loc.getLine();
            loc.offset = line + fileOffset - 1;
        }
    }

    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = new JsonObject();
        for (Entry<String, Integer> file : this.offsets.entrySet())
            obj.put(file.getValue().toString(), file.getKey());
        return obj;
    }
}
