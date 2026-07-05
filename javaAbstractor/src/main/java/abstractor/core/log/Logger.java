package abstractor.core.log;

import java.io.PrintStream;
import java.time.Duration;
import java.util.Stack;

public class Logger {
    static public final LogColor noticeColor  = LogColor.blue;
    static public final LogColor warningColor = LogColor.yellow;
    static public final LogColor errorColor   = LogColor.red;
    static private final String defaultIndent = "  ";

    public final Level level;
    private final PrintStream out;
    private final PrintStream err;
    private final Stack<String> indents;
    private int notices;
    private int warnings;
    private int errors;
    
    public Logger(Level level) { this(level, null, null); }

    public Logger(Level level, PrintStream out, PrintStream err) {
        this.level = level;
        this.out = out != null ? out : System.out;
        this.err = err != null ? err : System.err;
        this.indents = new Stack<String>();
        this.warnings = 0;
        this.errors = 0;
    }
    
    private String indent() {
        return this.indents.empty() ? "" : this.indents.peek();
    }

    public boolean writesLevel(Level level) {
        return this.level.Contains(level);
    }

    private void write(PrintStream out, Level level, LogColor color, String text) {
        if (!this.writesLevel(level)) return;

        final String indent = this.indent();
        String head = indent;
        String tail = "";
        if (color != LogColor.reset) {
            head += color;
            tail += LogColor.reset;
        }
        out.println(head + text.replace("\n", "\n"+indent) + tail);
    }
    
    public int noticeCount() { return this.notices; }

    public int warningCount() { return this.warnings; }

    public int errorCount() { return this.errors; }
    
    public void logWithColor(LogColor color, String text) {
        this.write(this.out, Level.Normal, color, text);
    }

    public void log(String text) {
        this.logWithColor(LogColor.reset, text);
    }

    public void logIf(boolean condition, String text) {
        if (condition) this.log(text);
    }
    
    public void notice(String text) {
        this.notices++;
        this.write(this.out, Level.Notice, noticeColor, text);
    }

    public void noticeIf(boolean condition, String text) {
        if (condition) this.notice(text);
    }

    public void warning(String text) {
        this.warnings++;
        this.write(this.out, Level.Warning, warningColor, text);
    }

    public void warningIf(boolean condition, String text) {
        if (condition) this.warning(text);
    }

    public void error(String text) {
        this.errors++;
        this.write(this.err, Level.Error, errorColor, text);
    }

    public void errorIf(boolean condition, String text) {
        if (condition) this.error(text);
    }
    
    public void push() { this.push(Level.Normal, defaultIndent); }
    
    public void push(Level level) { this.push(level, defaultIndent); }

    public void push(String indent) { this.push(Level.Normal, indent); }
    
    public void push(Level level, String indent) {
        if (!this.writesLevel(level)) return;
        this.indents.push(this.indent() + indent);
    }

    public void pop() { this.pop(Level.Normal); }
    
    public void pop(Level level) {
        if (!this.writesLevel(level)) return;
        if (!this.indents.empty()) this.indents.pop();
    }

    @FunctionalInterface public interface ThrowingRunnable { void run() throws Exception; }

    public void measure(String label, ThrowingRunnable func) throws Exception {
        final long start = System.nanoTime();
        this.log("Starting " + label + "...");
        func.run();
        final long stop = System.nanoTime();
        Duration elapsed = Duration.ofNanos(stop - start);
        this.log("Finished " + label + " (" + elapsed + ")");
    }
}
