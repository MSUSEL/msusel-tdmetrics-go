namespace DesignRecovery.Constructs;

public class MissingDataException(string name) :
   System.Exception("Missing JSON value for " + name + ".") { }
