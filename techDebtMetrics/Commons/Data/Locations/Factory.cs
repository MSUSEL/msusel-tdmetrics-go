using System.Collections.Generic;
using System.IO;
using System.Linq;

namespace Commons.Data.Locations;

public class Factory {
    private readonly Dictionary<string, int> maxLines = [];

    public void Add(Location loc) {
        if (loc == Location.Unknown) return; 
        if (this.maxLines.TryGetValue(loc.Path, out int length)) {
            if (loc.LineNo <= length) return;
        }
        this.maxLines[loc.Path] = loc.LineNo;
    }

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
