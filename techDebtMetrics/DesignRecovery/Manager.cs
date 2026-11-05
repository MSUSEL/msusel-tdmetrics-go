using Constructs;
using System.Collections.Generic;
using TechDebt;

namespace DesignRecovery;

/// <summary>Assists algorithms by containing the input and output project.</summary>
public class Manager {
    public readonly Constructs.Project Project;
    public readonly TechDebt.Project Tech;
    
    public readonly Dictionary<ObjectDecl, Class> ToObject;
    public readonly Dictionary<MethodDecl, Method> ToMethod;

    /// <summary>Creates a new manager for the given project.</summary>
    /// <param name="project">The project to read the classes and methods from.</param>
    public Manager(Constructs.Project project) {
        this.Project = project;
        this.Tech = new();
        this.ToObject = [];
        this.ToMethod = [];

        foreach (ObjectDecl cIn in this.Project.ObjectDecls) {
            Class cOut = new(new Source(cIn.Name, cIn.Location));
            if (!this.Tech.Classes.Add(cOut))
                throw new System.Exception("Class already exists in project: " + cOut);
            this.ToObject.Add(cIn, cOut);
        }

        foreach (MethodDecl mIn in this.Project.MethodDecls) {
            Method mOut = new(new Source(mIn.Name, mIn.Location));
            if (!this.Tech.Methods.Add(mOut))
                throw new System.Exception("Method already exists in project: " + mOut);
            this.ToMethod.Add(mIn, mOut);
        }
        
        foreach (MethodDecl mIn in this.Project.MethodDecls) {
            Method mOut = this.ToMethod[mIn];
            Metrics? met = mIn.Metrics;
            if (met is null) {
                // Methods without metrics, e.g. abstract methods, are counted as one by PMD.
                mOut.CyclomaticComplexity = 1;
                continue;
            }
            
            // Using PmdCyclo since that matches the ground truth values found in the TD dataset.
            mOut.CyclomaticComplexity = met.PmdCyclo;
            mOut.CodeCount = met.CodeCount;

            foreach (IConstruct c in met.Reads) {
                if (this.tryGetClass(c, out Class? cOut)) mOut.Reads.Add(cOut);
            }

            foreach (IConstruct c in met.Writes) {
                if (this.tryGetClass(c, out Class? cOut)) mOut.Writes.Add(cOut);
            }

            foreach (IConstruct c in met.Invokes) {
                if (this.tryGetMethod(c, out Method? m)) mOut.Invokes.Add(m);
            }
        }
    }

    // TODO: Need to add synthesised object for basics, interfaces, and other non-object types.
    private bool tryGetClass(IConstruct c, [System.Diagnostics.CodeAnalysis.NotNullWhen(true)] out Class? cOut) {
        ObjectDecl? decl = c switch {
            ObjectDecl d => d,
            ObjectInst i => i.Generic,
            _ => null
        };
        cOut = null;
        return decl is not null && this.ToObject.TryGetValue(decl, out cOut);
    }

    private bool tryGetMethod(IConstruct c, [System.Diagnostics.CodeAnalysis.NotNullWhen(true)] out Method? mOut) {
        MethodDecl? decl = c switch {
            MethodDecl d => d,
            MethodInst i => i.Generic,
            _ => null
        };
        mOut = null;
        return decl is not null && this.ToMethod.TryGetValue(decl, out mOut);
    }

    /// <summary>Gets the class for the given object declaration.</summary>
    public Class GetClass(ObjectDecl obj) => this.ToObject[obj];

    /// <summary>Gets the method for the given method declaration.</summary>
    public Method GetMethod(MethodDecl method) => this.ToMethod[method];

    /// <summary>Gets the Weighted Method Count (WMC) for the given object declaration.</summary>
    public double Wmc(ObjectDecl obj) => this.GetClass(obj).WeigthedMethodCount;

    /// <summary>Validates the technical debt project.</summary>
    /// <returns>The set of problems that were found, or empty if valid.</returns>
    public Validator.Problem[] Validate() => Validator.Validate(this.Tech);
}
