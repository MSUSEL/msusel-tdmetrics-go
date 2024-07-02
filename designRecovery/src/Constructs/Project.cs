using System.Text.Json.Nodes;

namespace designRecovery.src.Constructs;

internal class Project {
    public readonly Types Types;

    private readonly List<Package> inPackages = [];
    public IReadOnlyList<Package> Packages => this.inPackages.AsReadOnly();

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
                throw new MissingDataException("packages["+i+"]");
            this.inPackages.Add(new Package(packageNode));
        }

        TypeGetter getter = new(this);
        this.Types.Initialize(getter, types);
        for (int i = 0; i < packageCount; i++){
            JsonNode packageNode = packagesArr[i] ??
                throw new MissingDataException("packages["+i+"]");
            this.inPackages[i].Initialize(getter, packageNode);
        }
    }
}
