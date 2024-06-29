namespace designRecovery.src.Constructs;

internal class MissingDataException : System.Exception {
    public MissingDataException(string name) :
        base("Missing JSON value for "+name+".") { }
}
