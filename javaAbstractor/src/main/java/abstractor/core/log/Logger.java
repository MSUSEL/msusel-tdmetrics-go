package abstractor.core.log;

public interface Logger {
    void log(String text);
    void push();
    void push(String indent);
    void pop();
}
