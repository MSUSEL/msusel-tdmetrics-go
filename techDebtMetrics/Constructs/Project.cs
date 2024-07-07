using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;
using System.IO;
using System.Text;

namespace Constructs;

/// <summary>
/// A project represents all the packages to completely describe
/// a program or a library with all of it's import types also described.
/// </summary>
public class Project : IConstruct {

    /// <summary>The collection of type descriptions used in the project.</summary>
    public readonly Types Types;

    /// <summary>The collection of packages in the project.</summary>
    public IReadOnlyList<Package> Packages => this.inPackages.AsReadOnly();
    private readonly List<Package> inPackages = [];

    /// <summary>Loads a project from a YAML file.</summary>
    /// <param name="path">The file path to read from.</param>
    /// <returns>The project that was read from the YAML.</returns>
    static public Project FromFile(string path) {
        using StreamReader r = new(path);
        return FromText(r.ReadToEnd());
    }

    /// <summary>Loads a project from a string containing YAML data.</summary>
    /// <param name="text">The string containing YAML data to read from.</param>
    /// <returns>The project that was read from the YAML.</returns>
    static public Project FromText(string text) =>
        new(Data.Node.Parse(text));

    /// <summary>Creates a new project.</summary>
    /// <param name="root">The YAML root node to load.</param>
    internal Project(Data.Node root) {
        Data.Object obj = root.AsObject();

        Data.Object types = obj["types"]?.AsObject() ??
            throw new MissingDataException("types");
        this.Types = new Types(types);

        Data.Array packagesArr = obj["packages"].AsArray();
        int packageCount = packagesArr.Count;
        for (int i = 0; i < packageCount; i++)
            this.inPackages.Add(new Package(packagesArr[i]));

        TypeGetter getter = new(this);
        (this.Types as IInitializable).Initialize(getter, types);
        for (int i = 0; i < packageCount; i++)
            (this.inPackages[i] as IInitializable).Initialize(getter, packagesArr[i]);
    }

    public void ToStub(Journal j) =>
        j.Write(this.Packages, separator: "\n\n");

    public string ToStub() =>
        new Journal().Write(this).ToString();
}
