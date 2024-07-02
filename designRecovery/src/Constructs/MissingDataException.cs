namespace designRecovery.src.Constructs;

internal class MissingDataException(string name) :
    Exception("Missing JSON value for "+name+".") { }
