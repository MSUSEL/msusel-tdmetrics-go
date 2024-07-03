using System.Collections.Generic;
using System.IO;
using System.Text.Json.Nodes;

namespace DesignRecovery.Constructs;

public class Project {
    public readonly Types Types;

    private readonly List<Package> inPackages = [];
    public IReadOnlyList<Package> Packages => this.inPackages.AsReadOnly();

    static public Project FromJsonFile(string path) {
        using StreamReader r = new(path);
        return FromJson(r.ReadToEnd());
    }

    static public Project FromJson(string text) {
        JsonNode node = JsonNode.Parse(text) ??
            throw new IOException("Failed to parse Json text.");
        return new Project(node);
    }

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
