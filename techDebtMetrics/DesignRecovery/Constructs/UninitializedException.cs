namespace DesignRecovery.Constructs;

public class UninitializedException(string name) :
   System.Exception("May not access " + name + " until it has been initialized.") { }
