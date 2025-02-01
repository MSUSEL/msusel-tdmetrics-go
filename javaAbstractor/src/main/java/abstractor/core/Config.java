package abstractor.core;

import java.util.ArrayList;

public class Config {
    public String input  = "";
    public String output = "";

    public boolean verbose;
    public boolean writeTypes;
    public boolean writeIndices;

    public String toString() {
        StringBuilder sb = new StringBuilder();
        sb.append("Input:  ").append(this.input).append("\n");
        sb.append("Output: ").append(this.output);

        ArrayList<String> flags = new ArrayList<String>();
        if (this.verbose) flags.add("Verbose");
        if (this.writeTypes) flags.add("Write Types");
        if (this.writeIndices) flags.add("Write Indices");
    
        if (!flags.isEmpty()) {
            sb.append("\nFlags:  ");
            boolean first = true;
            for (String flag : flags) {
                if (first) first = false;
                else sb.append(", ");
                sb.append(flag);
            }
        }
        return sb.toString();
    }
}
