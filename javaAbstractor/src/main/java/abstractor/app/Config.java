package abstractor.app;

import java.util.ArrayList;
import java.io.*;
import java.nio.file.*;

import org.apache.commons.cli.*;

public class Config {
    public String input;
    public String output;

    public boolean verbose;
    public boolean minimize;

    public boolean writeKinds = true;
    public boolean writeIndices = true;

    /**
     * The output to print to when output path is null.
     * This defaults to standard output but may be overwritten
     * to collect the output JSON when testing.
     */
    public PrintStream defaultOut = System.out;

    static private Options getArgsOptions() {
        final Options op = new Options();
        op.addOption("i", "input", true,
            "The input path to the project to read. "+
            "The project directory must have a pom.xml Maven file.");
        
        op.addOption("o", "output", true,
            "The output file path to write the JSON to. "+
            "If not given, the JSON will be outputted to the console.");

        op.addOption("h", "help", false,
            "Prints this message.");
        
        op.addOption("m", "minimize", false,
            "Indicates the JSON output should beb minimized instead of formatted.");
        
        op.addOption("v", "verbose", false,
            "Indicates the abstraction process should output additional status information.");
        return op;
    }

    static private boolean isValidInput(String input) {
        final Path inPath = Path.of(input);
        return Files.exists(inPath) &&
            Files.isDirectory(inPath) &&
            Files.exists(Path.of(input, "pom.xml"));
    }

    static private boolean isValidOutput(String output) {
        if (output == null) return true;
        final boolean ext =
            output.endsWith(".json") ||
            output.endsWith(".yaml") ||
            output.endsWith(".yml");
        if (!ext) return false;

        final Path parPath = Path.of(output).getParent();
        if (parPath == null) return true;
        return Files.exists(parPath) && Files.isDirectory(parPath);
    }
    
    /**
     * This will populate the configs with the given command line arguments.
     * 
     * @see https://commons.apache.org/proper/commons-cli/usage.html
     * @param args The full command line that was passed into the main method.
     * @param out  The output to write errors and help out to.
     *             If null, the default output (System.out) is used.
     * @return True indicates the config was populated correctly.
     *         Otherwise, false if the arguments were invalid or help was requested.
     */
    public boolean FromArgs(String[] args, PrintStream out) {
        if (out == null) out = System.out;
        final Options options = getArgsOptions();

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
                out.println("Must provide an input path to project folder containing a pom.xml Maven file.");
                out.println("Use -h for help with arguments.");
                return false;
            }
            String input = cmd.getOptionValue("input");
            if (!isValidInput(input)) {
                out.println("Input path must be an existing directory containing a pom.xml Maven file.");
                out.println("Use -h for help with arguments.");
                return false;
            }

            String output = cmd.getOptionValue("output");
            if (!isValidOutput(output)) {
                out.println("If an output path is given, it must be in an existing directory and a JSON (or YAML) file.");
                out.println("Use -h for help with arguments.");
                return false;
            }

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
        if (this.writeKinds) flags.add("Write Types");
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
