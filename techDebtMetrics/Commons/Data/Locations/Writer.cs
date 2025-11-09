using Commons.Data.Yaml;
using System;
using System.Collections.Generic;
using System.Linq;

namespace Commons.Data.Locations;

/// <summary>A writer used for writing location information to Yaml files.</summary>
/// <param name="offsets">The offsets that will be used.</param>
public class Writer(Dictionary<string, int> offsets) {
    private readonly Dictionary<string, int> offsets = offsets;

    /// <summary>Gets an offset for the given location.</summary>
    /// <param name="loc">The location to get the offset for.</param>
    /// <returns>The offset for the given location.</returns>
    public int Offset(Location loc) {
        if (this.offsets.TryGetValue(loc.Path, out int offset))
            return loc.LineNo - 1 + offset;
        throw new Exception("The path was not found: " + loc.Path);
    }

    /// <summary>Writes a node with the location information for decoding offsets.</summary>
    /// <returns>The node containing location information.</returns>
    public Node Write() {
        Yaml.Object obj = new();
        foreach (string path in this.offsets.Keys.Order()) {
            int offset = this.offsets[path];
            obj.Add(offset.ToString(), offset);
        }
        return obj;
    }
}
