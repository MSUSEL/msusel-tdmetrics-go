package abstractor.core;

/**
 * An {@link AbstractorException} variant indicating that the underlying
 * failure has already been reported via the {@code Logger}. It is thrown
 * purely as a control-flow signal to escape the current construct being
 * processed and unwind to a higher frame in the abstractor pipeline.
 *
 * <p>Catchers of this exception should not log it again;
 * they should simply abandon the in-progress construct and continue with
 * the next one.
 */
public class LoggedException extends AbstractorException {
    private static final long serialVersionUID = 1L;

    public LoggedException() {
        super();
    }

    public LoggedException(String message) {
        super(message);
    }

    public LoggedException(String message, Throwable cause) {
        super(message, cause);
    }

    public LoggedException(Throwable cause) {
        super(cause);
    }
}
