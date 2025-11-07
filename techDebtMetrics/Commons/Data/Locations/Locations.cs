using Commons.Data.Reader;
using System.Collections.Generic;
using System.Linq;

namespace Commons.Data.Locations;

/// <summary>A tool for looking up a file path and line number for source code.</summary>
public class Locations {

    /// <summary>The location returned for an unknown location.</summary>
    static public readonly Location Unknown = new(0, "<unknown>");

    private readonly List<Location> offsets = [];

    /// <summary>Creates a new location lookup for the data in the given node.</summary>
    /// <see cref="docs/genFeatureDef.md#locations"/>
    /// <param name="node">The node to load the location data from.</param>
    static public Locations Read(Node? node) {
        Locations loc = new();
        if (node is null) return loc;

        loc.offsets.AddRange(node.AsObject().Children.
            Select(child => new Location(int.Parse(child.Key), child.Value.AsString())).
            Order());
        return loc;
    }

    /// <summary>
    /// Finds the index of the location that has the closest
    /// less than or equal to offset relative to the given offset.
    /// </summary>
    /// <param name="offset">The offset for the location to find.</param>
    /// <returns>The index of the location with the offset or -1 if none found.</returns>
    private int findRange(int offset) {
        int index = this.offsets.BinarySearch(new Location(offset, ""));
        return index switch {
            >= 0 => index,   // exact match
            ~0 => -1,        // not found
            _ => ~index - 1, // nearest match
        };
    }

    /// <summary>Gets the location for the given offset.</summary>
    /// <see cref="../../../docs/genFeatureDef.md#locations"/>
    /// <param name="offset">The offset is a specific line number and file path in the source code.</param>
    /// <returns>The location for the given offset, or the unknown offset if the offset is unknown.</returns>
    public Location this[int offset] {
        get {
            int index = this.findRange(offset);
            if (index < 0) return Unknown;
            Location pos = this.offsets[index];
            return new(offset - pos.LineNo, pos.Path);
        }
    }
}
