using System.Collections.Generic;
using System.Linq;

namespace Commons.Data.Locations;

/// <summary>This is a factory for creating the location writing.</summary>
/// <remarks>
/// The factory must have all the locations that will be written added to it,
/// or at minimum the location with the highest line number for every file.
/// This is to setup the offsets that are used when writting.
/// </remarks>
public class Factory {
    private readonly Dictionary<string, int> maxLines = [];

    /// <summary>Adds a new location that will be written.</summary>
    /// <param name="loc">The location that will be written.</param>
    public void Add(Location loc) {
        if (loc == Location.Unknown) return; 
        if (this.maxLines.TryGetValue(loc.Path, out int length)) {
            if (loc.LineNo <= length) return;
        }
        this.maxLines[loc.Path] = loc.LineNo;
    }

    /// <summary>Builds the writer that will return offsets for locations.</summary>
    /// <returns>The writer for locations.</returns>
    public Writer Build() {
        Dictionary<string, int> offsets = [];
        int offset = 0;
        foreach (string path in this.maxLines.Keys.Order()) {
            offsets[path] = offset;
            offset += this.maxLines[path];
        }
        return new Writer(offsets);
    }
}
