using Commons.Data.Yaml;
using System;
using System.Collections.Generic;
using System.Linq;

namespace Commons.Data.Locations;

public class Writer(Dictionary<string, int> offsets) {
    private readonly Dictionary<string, int> offsets = offsets;

    public int Offset(Location loc) {
        if (this.offsets.TryGetValue(loc.Path, out int offset))
            return loc.LineNo - 1 + offset;
        throw new Exception("The path was not found: " + loc.Path);
    }

    public Node Write() {
        Yaml.Object obj = new();
        foreach (string path in this.offsets.Keys.Order()) {

        }




        return obj;
    }
}
