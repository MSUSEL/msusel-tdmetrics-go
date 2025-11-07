namespace Constructs.Exceptions;

public class ConstructsException(string message) :
   System.Exception(message) { }

public class UninitializedException(string name) :
   ConstructsException("May not access " + name + " until it has been initialized.") { }
