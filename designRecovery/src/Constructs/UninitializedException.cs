namespace designRecovery.src.Constructs;

internal class UninitializedException : Exception {
    public UninitializedException(string name) :
        base("May not access "+name+" until it has been initialized.") { }
}
