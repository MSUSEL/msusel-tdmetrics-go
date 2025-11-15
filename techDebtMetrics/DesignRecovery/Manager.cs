using Constructs;
using System.Collections.Generic;
using TechDebt;

namespace DesignRecovery;

/// <summary>Assists algorithms by containing containing the input and output project.</summary>
public class Manager {
    public readonly Constructs.Project Project;
    public readonly TechDebt.Project Tech;
    
    public readonly Dictionary<ObjectDecl, Class> ToObject;
    public readonly Dictionary<MethodDecl, Method> ToMethod;

    public Manager(Constructs.Project project) {
        this.Project = project;
        this.Tech = new();
        this.ToObject = [];
        this.ToMethod = [];

        foreach (ObjectDecl cIn in this.Project.ObjectDecls) {
            Class cOut = new(new Source(cIn.Name, cIn.Location));
            this.Tech.Classes.Add(cOut);
            this.ToObject.Add(cIn, cOut);
        }

        foreach(MethodDecl mIn in this.Project.MethodDecls) {
            Method mOut = new(new Source(mIn.Name, mIn.Location));
            this.Tech.Methods.Add(mOut);
            this.ToMethod.Add(mIn, mOut);
            if (mIn.Metrics is not null) {
                Metrics met = mIn.Metrics;
                mOut.CyclomaticComplexity = met.Complexity;
                mOut.CodeCount = met.CodeCount;

                foreach (IConstruct c in met.Reads) {
                    //mOut.Reads.Add(c);


                }

                // TODO: Copy over more information
            }
        }
    }
}
