package abstractor.core.diff;

public record DiffStep(StepType type, int count) {

    static public DiffStep Equal(int count) {
        return new DiffStep(StepType.equal, count);
    }

    static public DiffStep Added(int count) {
        return new DiffStep(StepType.added, count);
    }
    
    static public DiffStep Removed(int count) {
        return new DiffStep(StepType.removed, count);
    }
    
    public final boolean isEqual()   { return this.type == StepType.equal; }
    public final boolean isAdded()   { return this.type == StepType.added; }
    public final boolean isRemoved() { return this.type == StepType.removed; }
}
