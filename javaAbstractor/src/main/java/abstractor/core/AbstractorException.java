package abstractor.core;

/**
 * Marker exception type thrown by the {@link Abstractor} and {@link Analyzer}
 * (and other classes in {@code abstractor.core}) to differentiate failures
 * originating inside the abstractor from exceptions raised by external
 * libraries such as Spoon.
 *
 * <p>Catch this type at frame boundaries when you want to handle only
 * abstractor-internal failures and let library exceptions propagate (or
 * vice versa).
 */
public class AbstractorException extends RuntimeException {
    private static final long serialVersionUID = 1L;

    public AbstractorException() {
        super();
    }

    public AbstractorException(String message) {
        super(message);
    }

    public AbstractorException(String message, Throwable cause) {
        super(message, cause);
    }

    public AbstractorException(Throwable cause) {
        super(cause);
    }
}
