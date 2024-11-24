namespace Constructs.Exceptions;

public class InvalidDataException(string name) :
   System.Exception("Invalid data in JSON value for " + name + ".") { }
