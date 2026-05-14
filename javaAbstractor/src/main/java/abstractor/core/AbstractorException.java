package abstractor.core;

/**
 * Marker exception type thrown by the {@link Abstractor} and {@link Analyzer}
 * to differentiate failures originating inside the abstractor from exceptions
 * raised by external libraries such as Spoon.
 */
public class AbstractorException extends Exception {
    public AbstractorException()                                { super(); }
    public AbstractorException(String message)                  { super(message); }
    public AbstractorException(String message, Throwable cause) { super(message, cause); }
    public AbstractorException(Throwable cause)                 { super(cause); }
}
