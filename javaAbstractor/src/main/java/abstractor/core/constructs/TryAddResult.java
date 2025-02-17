package abstractor.core.constructs;

public class TryAddResult<T extends Construct> {
    public final T value;
    public final boolean existed;

    public TryAddResult(T value, boolean existed) {
        this.value = value;
        this.existed = existed;
    }    
}
