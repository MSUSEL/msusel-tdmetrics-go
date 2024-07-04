using System.Collections.Generic;
using System.IO;
using System.Text.Json.Nodes;
using System.Text.RegularExpressions;

namespace DesignRecovery.Constructs;

/// <summary>
/// A project represents all the packages to completely describe
/// a program or a library with all of it's import types also described.
/// </summary>
public partial class Project {
    [GeneratedRegex(@"#.*$", RegexOptions.Compiled)]
    private static partial Regex JsonCommentMatch();

    /// <summary>
    /// The collection of type descriptions used in the project.
    /// </summary>
    public readonly Types Types;

    /// <summary>
    /// The collection of packages in the project.
    /// </summary>
    public IReadOnlyList<Package> Packages => this.inPackages.AsReadOnly();
    private readonly List<Package> inPackages = [];

    /// <summary>
    /// Loads a project from a JSON file with optional YAML style comments.
    /// </summary>
    /// <param name="path">The file path to read from.</param>
    /// <returns>The project that was read from the JSON.</returns>
    static public Project FromJsonFile(string path) {
        using StreamReader r = new(path);
        return FromJson(r.ReadToEnd());
    }

    /// <summary>
    /// Loads a project from a string containing JSON data with optional YAML style comments.
    /// </summary>
    /// <param name="text">The string containing JSON data to read from.</param>
    /// <returns>The project that was read from the JSON.</returns>
    static public Project FromJson(string text) {
        text = JsonCommentMatch().Replace(text, "");
        JsonNode node = JsonNode.Parse(text) ??
            throw new IOException("Failed to parse Json text.");
        return new Project(node);
    }

    /// <summary>Creates a new project.</summary>
    /// <param name="node">The JSON node to load.</param>
    public Project(JsonNode node) {
        JsonObject obj = node.AsObject();

        JsonObject types = obj["types"]?.AsObject() ??
            throw new MissingDataException("types");
        this.Types = new Types(types);

        JsonArray packagesArr = obj["packages"]?.AsArray() ??
            throw new MissingDataException("packages");
        int packageCount = packagesArr.Count;
        for (int i = 0; i < packageCount; i++) {
            JsonNode packageNode = packagesArr[i] ??
                throw new MissingDataException("packages[" + i + "]");
            this.inPackages.Add(new Package(packageNode));
        }

        TypeGetter getter = new(this);
        this.Types.Initialize(getter, types);
        for (int i = 0; i < packageCount; i++) {
            JsonNode packageNode = packagesArr[i] ??
                throw new MissingDataException("packages[" + i + "]");
            this.inPackages[i].Initialize(getter, packageNode);
        }
    }
}
