namespace designRecovery.src.Constructs;

internal class UninitializedException(string name) :
    Exception("May not access "+name+" until it has been initialized.") { }
