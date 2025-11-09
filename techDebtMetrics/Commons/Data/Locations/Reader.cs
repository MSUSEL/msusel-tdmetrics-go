using Commons.Data.Yaml;
using Commons.Extensions;
using System;
using System.Collections.Generic;
using System.Linq;

namespace Commons.Data.Locations;

/// <summary>A tool for looking up a file path and line number for source code.</summary>
public class Reader {

    /// <summary>The file locations</summary>
    private readonly SortedList<int, string> offsets = [];

    /// <summary>Creates a new location lookup for the data in the given node.</summary>
    /// <see cref="docs/genFeatureDef.md#locations"/>
    /// <param name="node">The node to load the location data from.</param>
    static public Reader Read(Node? node) {
        Reader loc = new();
        if (node is null) return loc;
        node.AsObject().Children.
            ForAll(child => loc.offsets.Add(child.Key.AsInt(), child.Value.AsString()));
        return loc;
    }

    /// <summary>Gets the location for the given offset.</summary>
    /// <see cref="../../../docs/genFeatureDef.md#locations"/>
    /// <param name="offset">The offset is a specific line number and file path in the source code.</param>
    /// <returns>The location for the given offset, or the unknown offset if the offset is unknown.</returns>
    public Location this[int offset] {
        get {
            (int index, bool _) = this.offsets.Keys.BinarySearch(offset.CompareTo);
            if (index < 0) return Location.Unknown;
            int start = this.offsets.GetKeyAtIndex(index);
            string path = this.offsets[index];
            return new(offset - start + 1, path);
        }
    }
}
