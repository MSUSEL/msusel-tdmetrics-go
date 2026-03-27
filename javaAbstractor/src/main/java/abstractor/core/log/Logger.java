package abstractor.core.log;

import java.io.PrintStream;
import java.util.Stack;

public class Logger {
    public final boolean verbose;
    private final PrintStream out;
    private final PrintStream err;
    private final Stack<String> indents;
    private int warnings;
    private int errors;
    
    public Logger(boolean verbose) { this(verbose, null, null); }

    public Logger(boolean verbose, PrintStream out, PrintStream err) {
        this.verbose = verbose;
        this.out = out != null ? out : System.out;
        this.err = err != null ? err : System.err;
        this.indents = new Stack<String>();
        this.warnings = 0;
        this.errors = 0;
    }
    
    public int warningCount() { return this.warnings; }

    public int errorCount() { return this.errors; }

    public void log(String text) {
        if (!this.verbose) return;
        final String indent = this.indent();
        this.out.println(indent + text.replace("\n", "\n"+indent));
    }

    public void logIf(boolean condition, String text) {
        if (condition) this.log(text);
    }
    
    public void warning(String text) {
        if (!this.verbose) return;
        this.warnings++;
        final String indent = this.indent();
        this.out.println(indent + text.replace("\n", "\n"+indent));
    }

    public void warningIf(boolean condition, String text) {
        if (condition) this.warning(text);
    }

    public void error(String text) {
        this.errors++;
        final String indent = this.indent();
        this.err.println(indent + text.replace("\n", "\n"+indent));
    }

    public void errorIf(boolean condition, String text) {
        if (condition) this.error(text);
    }
    
    private String indent() {
        return this.indents.empty() ? "" : this.indents.peek();
    }
    
    public void push() { this.push("  "); }
    
    public void push(String indent) {
        this.indents.push(this.indent() + indent);
    }

    public void pop() {
        if (!this.indents.empty()) this.indents.pop();
    }
}
