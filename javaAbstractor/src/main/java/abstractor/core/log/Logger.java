package abstractor.core.log;

import java.io.PrintStream;
import java.util.Stack;

public class Logger {
    public final boolean verbose;
    private final PrintStream out;
    private final PrintStream err;
    private final Stack<String> indents;
    private int errors;
    
    public Logger(boolean verbose) { this(verbose, null, null); }

    public Logger(boolean verbose, PrintStream out, PrintStream err) {
        this.verbose = verbose;
        this.out = out != null ? out : System.out;
        this.err = err != null ? err : System.err;
        this.indents = new Stack<String>();
        this.errors = 0;
    }

    public int errorCount() { return this.errors; }

    public void log(String text) {
        if (!this.verbose) return;
        final String indent = this.indent();
        this.out.println(indent + text.replace("\n", "\n"+indent));
    }

    public void error(String text) {
        this.errors++;
        final String indent = this.indent();
        this.err.println(indent + text.replace("\n", "\n"+indent));
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
