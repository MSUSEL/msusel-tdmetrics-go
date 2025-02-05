package abstractor.core.log;

public class Null implements Logger {
    public void log(String text) {}
    public void push() {}
    public void push(String indent) {}
    public void pop() {}
}
