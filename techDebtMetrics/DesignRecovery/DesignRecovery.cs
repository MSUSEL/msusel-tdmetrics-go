using Constructs;
using Participation;

namespace DesignRecovery;

public class DesignRecovery(Project project) {
    public readonly Project Project = project;

    public Matrix GetMembership() {
        int methodCount = this.Project.MethodDecls.Count;
        int objectCount = this.Project.ObjectDecls.Count;
        Matrix m = new(methodCount, objectCount);
        for (int methodIndex = 0; methodIndex < methodCount; ++methodIndex) {
            MethodDecl method = this.Project.MethodDecls[methodIndex];
            if (method.Receiver is not null)
                m[methodIndex, method.Receiver.Index] = 1.0;
        }
        return m;
    }

    public Matrix GetParticipation() {
        int methodCount = this.Project.MethodDecls.Count;
        int objectCount = this.Project.ObjectDecls.Count;
        // TODO: Need to add synthesised object for basics and projects.
        return new Matrix(methodCount, objectCount);
    }
}
