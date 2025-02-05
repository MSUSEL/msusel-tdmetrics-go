package abstractor.core.log;

import java.io.PrintStream;
import java.util.Stack;

public class Stream implements Logger {
    private final PrintStream out;
    private final Stack<String> indents;
    
    public Stream() { this(null); }

    public Stream(PrintStream out) {
        if (out == null) out = System.out;
        this.out = out;
        this.indents = new Stack<String>();
    }

    public void log(String text) {
        final String indent = this.indent();
        this.out.println(indent + text.replace("\n", "\n"+indent));
    }
    
    private String indent() {
        return this.indents.empty() ? "" : this.indents.peek();
    }
    
    public void push() {
        this.push("  ");
    }
    
    public void push(String indent) {
        this.indents.push(this.indent() + indent);
    }

    public void pop() {
        if (!this.indents.empty()) this.indents.pop();
    }
}
