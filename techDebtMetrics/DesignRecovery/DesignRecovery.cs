using Constructs;

namespace DesignRecovery;

public class DesignRecovery(Project project) {
    public readonly Project Project = project;

    /*
    public Matrix GetParticipation() {
        int methodCount = this.Project.MethodDecls.Count;
        int objectCount = this.Project.ObjectDecls.Count;
        
        // TODO: Need to add synthesised object for basics and projects.

        return new Matrix(methodCount, objectCount);
    }
    */
}
