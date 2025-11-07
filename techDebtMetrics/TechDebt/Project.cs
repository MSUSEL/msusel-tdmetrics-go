using Commons.Data.Locations;
using Commons.Data.Reader;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using TechDebt.Exceptions;

namespace TechDebt;

public class Project {
    
    /// <summary>The collection of all classes in this project.</summary>
    public readonly SortedSet<Class> Classes = [];
    
    /// <summary>The collection of all methods in this project.</summary>
    public readonly SortedSet<Method> Methods = [];

    public bool TryFind(Source source, out Class? c) => this.Classes.TryGetValue(new Class(source), out c);

    public bool TryFind(Source source, out Method? m) => this.Methods.TryGetValue(new Method(source), out m);

    public static void Add(Method m, double value, Class c) => Add(new Participation(m, value, c));

    public static void Add(Participation p) {
        p.Method.Participation.Add(p);
        p.Class.Participation.Add(p);
    }

    public static void Remove(Participation p) {
        p.Class.Participation.Remove(p);
        p.Method.Participation.Remove(p);
    }

    public void Normalize() => this.Methods.ForAll(Normalize);

    public static void Normalize(Method m) {
        if (m.Participation.Count <= 0) throw new NoParticipationException(m);

        List<(Class Class, double Value)> entries = [.. m.Participation.Select(p => (p.Class, p.Value))];
        entries.Sort((p1, p2) => p1.Value.CompareTo(p2.Value));

        void apply(Func<double, double> h) {
            for (int i = 0; i < entries.Count; i++)
                entries[i] = (entries[i].Class, h(entries[i].Value));
        }

        double min = entries.Min(p => p.Value);
        if (min < 0.0) apply(v => v + min);

        double max = entries.Max(p => p.Value);
        if (Math.IsZero(max)) apply(v => v + 1.0);

        double sum = entries.Sum(p => p.Value);
        apply(v => v / sum);

        entries.RemoveAll(e => Math.IsZero(e.Value));

        // Do a second normalization pass after removing nearly zero entries.
        sum = entries.Sum(p => p.Value);
        apply(v => v / sum);

        m.Participation.ForAll(p => p.Class.Participation.Remove(p));
        m.Participation.Clear();
        entries.ForAll(e => Add(m, e.Value, e.Class));
    }

    /// <summary>Loads a project from a YAML root node.</summary>
    /// <param name="root">The YAML root node to load.</param>
    public Project(Node root) {
        Commons.Data.Reader.Object obj = root.AsObject();
        Locations locs = Locations.Read(obj.TryReadNode("locs"));
        LoaderHelper lh = new(locs);

        obj.PreallocateList("classes", lh.Classes, n => new Class(lh.ReadSource(n)));
        obj.PreallocateList("methods", lh.Methods, n => new Method(lh.ReadSource(n)));
        
        obj.InitializeList(lh, "classes", lh.Classes);
        obj.InitializeList(lh, "methods", lh.Methods);

        foreach (Method m in lh.Methods) {
            if (this.Methods.Contains(m))
                throw new Exception("Method already exists in project: " + m);
            this.Methods.Add(m);
        }

        foreach (Class c in lh.Classes) {
            if (this.Classes.Contains(c))
                throw new Exception("Class already exists in project: " + c);
            this.Classes.Add(c);
        }
    }

    internal class LoaderHelper(Locations locs): IKeyResolver {
        public readonly Locations Locations = locs;
        public List<Class> Classes = [];
        public List<Method> Methods = [];

        public Source ReadSource(Node node) {
            Commons.Data.Reader.Object obj = node.AsObject();
            Location loc = obj.ReadLocation(this.Locations, "loc");
            string name = obj.ReadString("name");
            return new Source(name, loc.Path, loc.LineNo);
        }

        /// <summary>Reads the given index from the given source as part of reading the given key.</summary>
        /// <typeparam name="T">The type of the list to read from.</typeparam>
        /// <param name="key">The key that is being processed.</param>
        /// <param name="index">The index from the key used to read a value from the given list.</param>
        /// <param name="source">The list get an item at the given index from.</param>
        /// <returns>The item from the given list at the given index.</returns>
        static private T readKeyIndex<T>(string key, int index, IReadOnlyList<T> source) {
            if (index < 0 || index >= source.Count)
                throw new Exception("Key " + key + " out of range [0.." + source.Count + "): " + index);
            return source[index];
        }

        /// <summary>Reads a single key from the given project.</summary>
        /// <see cref="docs/genFeatureDef.md#keys"/>
        /// <param name="key">The key of the value to read.</param>
        /// <param name="project">The project to read a key from.</param>
        /// <returns>The read key from the project.</returns>
        public object FindData(string name, int index) {
            return name switch {
                "class" => readKeyIndex(name, index, this.Classes),
                "method" => readKeyIndex(name, index, this.Methods),
                _ => throw new InvalidDataException(name)
            };
        }
    }

    public string Serialize() {

        return "";
    }
}
