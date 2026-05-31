package abstractor.core.cmp;

import java.util.ArrayList;

public class CmpOptions{
    static public boolean shouldUseResolved(CmpOptions op) {
        return op != null? op.useResolved: false;
    }

    public boolean useResolved;
    public boolean debugPrint;

    @Override
    public String toString() {
        ArrayList<String> parts = new ArrayList<>();
        if (this.useResolved) parts.add("useResolved");
        if (this.debugPrint)  parts.add("debugPrint");
        return "cmp:{" + String.join(", ", parts) + "}";
    }
}
