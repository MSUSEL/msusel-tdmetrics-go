using Constructs;
using TechDebt;

namespace DesignRecovery;

/// <summary>
/// This gets the membership from the project as defined in the by the receivers.
/// Any method without a receiver will not be given a membership.
/// </summary>
/// <remarks>This does not perform any design recovery.</remarks>
public class Membership : IAlgorithm {
    public void Run(Manager man) {
        foreach(MethodDecl mIn in man.Project.MethodDecls) {
            if (mIn.Receiver is not null) {
                Method mOut = man.ToMethod[mIn];
                Class cOut = man.ToObject[mIn.Receiver];
                mOut.Add(1.0, cOut);
            }
        }
    }
}
