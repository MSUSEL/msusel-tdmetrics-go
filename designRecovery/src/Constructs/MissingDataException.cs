namespace designRecovery.src.Constructs;

internal class MissingDataException : System.Exception {
    public MissingDataException(string nodeName, string valueName) :
        base("Missing JSON value for "+nodeName+"."+valueName+".") { }
}
