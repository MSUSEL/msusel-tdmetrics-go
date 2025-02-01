package abstractor.core;

import java.util.ArrayList;
import java.io.PrintStream;
import java.io.PrintWriter;

import org.apache.commons.cli.*;

public class Config {
    public String input  = "";
    public String output = "";

    public boolean verbose      = false;
    public boolean minimize     = false;
    public boolean writeTypes   = false;
    public boolean writeIndices = false;

    /**
     * This will populate the configs with the given command line arguments
     * @param args The full command line that was passed into the main method.
     * @param out  The output to write errors and help out to.
     *             If null, the default output (System.out) is used.
     * @return True indicates the config was populated correctly.
     *         Otherwise, false if the arguments were invalid or help was requested.
     */
    public boolean FromArgs(String[] args, PrintStream out) {
        if (out == null) out = System.out;
        final Options options = new Options()
            .addOption("i", "input", true,
                "The input path to the directory of the project or package to read. "+
                "The project directory should have a go.mod file.")
            .addOption("o", "output", true,
                "The output file path to write the JSON to. "+
                "If not given, the JSON will be outputted to the console.")
            .addOption("h", "help", false,
                "Prints this message.")
            .addOption("m", "minimize", false,
                "Indicates the JSON output should beb minimized instead of formatted.")
            .addOption("v", "verbose", false,
                "Indicates the abstraction process should output additional status information.");

        try {
            final CommandLine cmd = new DefaultParser().parse(options, args);
            if (cmd.hasOption("help")) {
                final PrintWriter pw = new PrintWriter(out);
                final HelpFormatter hf = new HelpFormatter();
                hf.printHelp(pw, hf.getWidth(), "java-abstractor", null, options, hf.getLeftPadding(), hf.getDescPadding(), null, true);
                pw.flush();
                return false;
            }

            // Read and validate the arguments.
            if (!cmd.hasOption("input")) {
                out.println("Must provide an input folder.");
                out.println("Use -h for help with arguments.");
                return false;
            }
            String input = cmd.getOptionValue("input");
            // TODO: Check that the input exists.

            if (!cmd.hasOption("output")) {
                out.println("Must provide an output file.");
                out.println("Use -h for help with arguments.");
                return false;
            }
            String output = cmd.getOptionValue("output");
            // TODO: Check that the output base directory exists.

            // Write new values to config.
            this.input    = input;
            this.output   = output;
            this.verbose  = cmd.hasOption("verbose");
            this.minimize = cmd.hasOption("minimize");
            return true;

        } catch(ParseException ex) {
            out.println("Invalid arguments: " + ex.getMessage());
            out.println("Use -h for help with arguments.");
            return false;

        } catch(Exception ex) {
            out.println(ex);
            return false;
        }
    }

    public String toString() {
        StringBuilder sb = new StringBuilder();
        sb.append("Input:  ").append(this.input).append("\n");
        sb.append("Output: ").append(this.output);

        ArrayList<String> flags = new ArrayList<String>();
        if (this.verbose) flags.add("Verbose");
        if (this.minimize) flags.add("Minimize");
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
