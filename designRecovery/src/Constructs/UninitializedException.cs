namespace designRecovery.src.Constructs;

internal class UninitializedException : Exception {
    public UninitializedException(string nodeName, string valueName) :
        base("May not access "+nodeName+"."+valueName+" until it has been initialized.") { }
}
