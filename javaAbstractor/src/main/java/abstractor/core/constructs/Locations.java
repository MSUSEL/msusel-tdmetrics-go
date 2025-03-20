package abstractor.core.constructs;

import java.util.TreeSet;
import java.io.File;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Map.Entry;
import java.util.TreeMap;

import spoon.reflect.cu.SourcePosition;

import abstractor.core.json.*;

public class Locations implements Jsonable {
    private final TreeSet<Location> locs = new TreeSet<Location>();
    private final TreeMap<String, Integer> offsets = new TreeMap<String, Integer>();
    
    public Location create(SourcePosition pos) {
        return this.create(new Location(pos));
    }
    
    public Location create(String path, int line) {
        return this.create(new Location(path, line));
    }

    private Location create(Location loc) {
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

    static private String cleanPath(String path) {
        if (path == null) path = "";
        try { path = new File(path).getCanonicalPath(); }
        catch (Exception ex) { }
        return path.replace("\\", "/");
    }

    static private String subPath(String root, String path) {
        if (path.startsWith(root)) path = path.substring(root.length());
        if (path.startsWith("/"))  path = path.substring(1);
        return path;
    }

    public JsonNode toJson(JsonHelper h) {
        final String root = cleanPath(h.rootPath);
        JsonObject obj = new JsonObject();
        for (Entry<String, Integer> file : this.offsets.entrySet()) {
            String path = subPath(root, cleanPath(file.getKey()));
            obj.put(file.getValue().toString(), path);
        }
        return obj;
    }
}
