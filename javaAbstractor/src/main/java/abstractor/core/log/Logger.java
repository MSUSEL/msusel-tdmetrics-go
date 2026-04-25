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

    // See https://j8ahmed.com/2021/09/13/day-37-learning-ansi-escape-codes/
    static public final String colorBlack   = "\u001b[30m";
    static public final String colorRed     = "\u001b[31m";
    static public final String colorGreen   = "\u001b[32m";
    static public final String colorYellow  = "\u001b[33m";
    static public final String colorBlue    = "\u001b[34m";
    static public final String colorMagenta = "\u001b[35m";
    static public final String colorCyan    = "\u001b[36m";
    static public final String colorWhite   = "\u001b[37m";
    static public final String colorReset   = "\u001b[0m";
    
    public Logger(boolean verbose) { this(verbose, null, null); }

    public Logger(boolean verbose, PrintStream out, PrintStream err) {
        this.verbose = verbose;
        this.out = out != null ? out : System.out;
        this.err = err != null ? err : System.err;
        this.indents = new Stack<String>();
        this.warnings = 0;
        this.errors = 0;
    }
    
    private String indent() {
        return this.indents.empty() ? "" : this.indents.peek();
    }

    private void write(PrintStream out, String color, String text) {
        final String indent = this.indent();
        String head = indent;
        String tail = "";
        if (color != colorReset) {
            head += color;
            tail += colorReset;
        }
        out.println(head + text.replace("\n", "\n"+indent) + tail);
    }
    
    public int warningCount() { return this.warnings; }

    public int errorCount() { return this.errors; }
    
    public void logWithColor(String ansiColor, String text) {
        if (!this.verbose) return;
        this.write(this.out, ansiColor, text);
    }

    public void log(String text) {
        this.logWithColor(colorReset, text);
    }

    public void logIf(boolean condition, String text) {
        if (condition) this.log(text);
    }
    
    public void warning(String text) {
        if (!this.verbose) return;
        this.warnings++;
        this.write(this.out, colorYellow, text);
    }

    public void warningIf(boolean condition, String text) {
        if (condition) this.warning(text);
    }

    public void error(String text) {
        this.errors++;
        this.write(this.err, colorRed, text);
    }

    public void errorIf(boolean condition, String text) {
        if (condition) this.error(text);
    }
    
    public void push() { this.push("  "); }
    
    public void push(String indent) {
        this.indents.push(this.indent() + indent);
    }

    public void pop() {
        if (!this.indents.empty()) this.indents.pop();
    }
}
