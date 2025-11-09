using Commons.Extensions;
using System.Collections.Generic;
using System.Linq;

namespace TechDebt;

public class Validator {

    public record Problem(string Message);

    public record MethodParticipationSum(Method Method, double Sum) :
        Problem("Mathod's participation sum of " + Sum + " is invalid for " + Method);

    public record WrongParticipationMethod(Method Method, Participation Participation) :
        Problem("Participation in " + Method + " was for " + Participation.Method);

    public record WrongParticipationClass(Class Class, Participation Participation) :
        Problem("Participation in " + Class + " was for " + Participation.Class);

    public record ParticipationValue(Participation Participation) :
        Problem("Participation between " + Participation.Method + " and " + Participation.Class + " has the invalid value of " + Participation.Value);

    public record MethodMissingFromProject(Participation Participation) :
        Problem("Participation's " + Participation.Method + " is not in the project");

    public record ClassMissingFromProject(Participation Participation) :
        Problem("Participation's " + Participation.Class + " is not in the project");

    public record MethodMissingParticipation(Participation Participation) :
        Problem("Participation's " + Participation.Method + " does not contain this participation");

    public record ClassMissingParticipation(Participation Participation) :
        Problem("Participation's " + Participation.Class + " does not contain this participation");

    public record DuplicateParticipation(Participation Participation1, Participation Participation2) :
        Problem("Multiple participations between " + Participation1.Method + " and " + Participation1.Class);

    public record MethodReadNotInProject(Method Method, Class Read) :
        Problem("The class read by " + Method + ", " + Read + ", is not in the project");

    public record MethodWriteNotInProject(Method Method, Class Write) :
        Problem("The class written by " + Method + ", " + Write + ", is not in the project");

    public record MethodInvokeNotInProject(Method Method, Method Invoke) :
        Problem("The method invoked by " + Method + ", " + Invoke + ", is not in the project");

    /// <summary>This will validate the project to ensure all the data is setup correctly.</summary>
    /// <param name="project">The project to validate.</param>
    /// <returns>The set of problems that were found, or empty if valid.</returns>
    public static Problem[] Validate(Project project) => new Validator(project).validate();

    private readonly Project project;
    private readonly List<Problem> problems = [];

    private Validator(Project project) => this.project = project;

    private Problem[] validate() {
        this.project.Methods.ForAll(this.validate);
        this.project.Classes.ForAll(this.validate);
        return [.. this.problems];
    }

    private void log(Problem p) => this.problems.Add(p);

    private void validate(Method m) {
        double sum = 0.0;
        Dictionary<Class, Participation> seen = [];
        foreach (Participation p in m.Participation) {
            if (p.Method != m)
                this.log(new WrongParticipationMethod(m, p));

            if (Math.LessThanEqual(p.Value, 0.0) || Math.GreaterThan(p.Value, 1.0))
                this.log(new ParticipationValue(p));
            sum += p.Value;

            if (!this.project.Classes.Contains(p.Class))
                this.log(new ClassMissingFromProject(p));

            if (!p.Class.Participation.Contains(p))
                this.log(new ClassMissingParticipation(p));

            if (seen.TryGetValue(p.Class, out Participation other))
                this.log(new DuplicateParticipation(p, other));
            else seen.Add(p.Class, p);
        }
        if (!Math.IsOne(sum))
            this.log(new MethodParticipationSum(m, sum));

        m.Reads.WhereNot(this.project.Classes.Contains).ForAll(r => this.log(new MethodReadNotInProject(m, r)));
        m.Writes.WhereNot(this.project.Classes.Contains).ForAll(w => this.log(new MethodWriteNotInProject(m, w)));
        m.Invokes.WhereNot(this.project.Methods.Contains).ForAll(i => this.log(new MethodInvokeNotInProject(m, i)));
    }

    private void validate(Class c) {
        Dictionary<Method, Participation> seen = [];
        foreach (Participation p in c.Participation) {
            if (p.Class != c)
                this.log(new WrongParticipationClass(c, p));

            if (Math.LessThanEqual(p.Value, 0.0))
                this.log(new ParticipationValue(p));

            if (!this.project.Methods.Contains(p.Method))
                this.log(new MethodMissingFromProject(p));

            if (!p.Method.Participation.Contains(p))
                this.log(new MethodMissingParticipation(p));

            if (seen.TryGetValue(p.Method, out Participation other))
                this.log(new DuplicateParticipation(p, other));
            else seen.Add(p.Method, p);
        }
    }
}
