using Constructs.Data;
using System.Collections.Generic;
using System.Linq;

namespace Constructs.Tooling;

public class Locations {
    static public readonly Location Unknown = new(0, "<unknown>");

    private readonly List<Location> offsets;

    internal Locations(Node? node) {
        if (node is null) {
            this.offsets = [];
            return;
        }

        this.offsets = new(node.AsObject().Children.
            Select(child => new Location(int.Parse(child.Key), child.Value.AsString())).
            Order());
    }

    private int findRange(int offset) {
        int index = this.offsets.BinarySearch(new Location(offset, ""));
        return index switch {
            >= 0 => index,   // exact match
            ~0 => -1,        // not found
            _ => ~index - 1, // nearest match
        };
    }

    public Location this[int offset] {
        get {
            int index = this.findRange(offset);
            if (index < 0) return Unknown;
            Location pos = this.offsets[index];
            return new(offset - pos.LineNo, pos.Path);
        }
    }
}
