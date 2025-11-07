namespace TechDebt;

public class TechDebtException(string message) : System.Exception(message) { }

public class NoParticipationException(Method method) :
    TechDebtException("No participation in " + method) {
    public readonly Method Method = method;
}
