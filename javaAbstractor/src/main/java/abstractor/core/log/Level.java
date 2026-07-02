package abstractor.core.log;

public enum Level {
    All(0),     // Outputs all logs
    Normal(1),  // Output Normal, Notice, Warning, and Error logs
    Notice(2),  // Output Notice, Warning, and Error logs
    Warning(3), // Output Warning and Error logs 
    Error(4),   // Output only Error logs
    None(5);    // Disable logging

    final public int value;
    private Level(int value) { this.value = value; }

    public boolean Contains(Level level) {
        return this.value <= level.value;
    }
}
