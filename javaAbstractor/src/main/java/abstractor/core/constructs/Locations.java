package abstractor.core.constructs;

import java.util.TreeSet;
import java.util.Map.Entry;
import java.util.TreeMap;

import spoon.reflect.cu.SourcePosition;

import abstractor.core.json.*;

public class Locations implements Jsonable {
    private final TreeSet<Location> locs = new TreeSet<Location>();
    private final TreeMap<String, Integer> offsets = new TreeMap<String, Integer>();
    
    public Location create(SourcePosition pos) {
        final Location loc = new Location(pos);

        final Location existing = this.locs.floor(loc);
        if (loc.equals(existing)) return existing;

        this.locs.add(loc);
        return loc;
    }

    private TreeMap<String, Integer> getMaximums() {
        TreeMap<String, Integer> maximums = new TreeMap<String, Integer>();
        for (Location loc : this.locs) {
            if (!loc.isValid()) continue;
            final String path = loc.path;
            final int line = loc.line;
            final Integer iMax = maximums.get(path);
            final int max = iMax == null ? line : Integer.max((int)iMax, line);
            maximums.put(path, max);
        }
        return maximums;
    }

    private void setOffsets(TreeMap<String, Integer> maximums) {
        this.offsets.clear();
        int offset = 1;
        for (String path : maximums.keySet()) {
            this.offsets.put(path, offset);
            offset += maximums.get(path);
        }
    }

    private void updateLocations() {
        for (Location loc : this.locs) {
            if (!loc.isValid()) continue;
            final String path = loc.path;
            final int fileOffset = this.offsets.get(path);
            final int line = loc.line;
            loc.offset = line + fileOffset - 1;
        }
    }

    public void prepareForOutput() {
        this.setOffsets(this.getMaximums());
        this.updateLocations();
    }

    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = new JsonObject();
        for (Entry<String, Integer> file : this.offsets.entrySet())
            obj.put(file.getValue().toString(), file.getKey());
        return obj;
    }
}
