package abstractor.core.log;

import java.io.PrintStream;
import java.util.Stack;

public class Logger {
    final static private String defaultIndent = "  ";

    public final Level level;
    private final PrintStream out;
    private final PrintStream err;
    private final Stack<String> indents;
    private int notices;
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

    private void write(PrintStream out, Level level, String color, String text) {
        if (!this.writesLevel(level)) return;

        final String indent = this.indent();
        String head = indent;
        String tail = "";
        if (color != colorReset) {
            head += color;
            tail += colorReset;
        }
        out.println(head + text.replace("\n", "\n"+indent) + tail);
    }
    
    public int noticeCount() { return this.notices; }

    public int warningCount() { return this.warnings; }

    public int errorCount() { return this.errors; }
    
    public void logWithColor(String ansiColor, String text) {
        this.write(this.out, Level.Normal, ansiColor, text);
    }

    public void log(String text) {
        this.logWithColor(colorReset, text);
    }

    public void logIf(boolean condition, String text) {
        if (condition) this.log(text);
    }
    
    public void notice(String text) {
        this.notices++;
        this.write(this.out, Level.Notice, colorBlue, text);
    }

    public void noticeIf(boolean condition, String text) {
        if (condition) this.notice(text);
    }

    public void warning(String text) {
        this.warnings++;
        this.write(this.out, Level.Warning, colorYellow, text);
    }

    public void warningIf(boolean condition, String text) {
        if (condition) this.warning(text);
    }

    public void error(String text) {
        this.errors++;
        this.write(this.err, Level.Error, colorRed, text);
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
}
