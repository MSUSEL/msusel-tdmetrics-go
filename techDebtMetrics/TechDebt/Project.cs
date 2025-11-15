using Commons.Data.Locations;
using Commons.Data.Yaml;
using Commons.Extensions;
using System;
using System.Collections.Generic;
using System.Linq;
using TechDebt.Exceptions;

namespace TechDebt;

/// <summary>The project stored as only methods and classes.</summary>
public class Project {
    
    /// <summary>The collection of all classes in this project.</summary>
    public readonly SortedSet<Class> Classes = [];
    
    /// <summary>The collection of all methods in this project.</summary>
    public readonly SortedSet<Method> Methods = [];

    /// <summary>Tries to find the class with the given source.</summary>
    /// <param name="source">The source for the class to try to find.</param>
    /// <param name="c">The found class or null.</param>
    /// <returns>True if the class was found, false otherwise.</returns>
    public bool TryFind(Source source, out Class? c) => this.Classes.TryGetValue(new Class(source), out c);
    
    /// <summary>Tries to find the method with the given source.</summary>
    /// <param name="source">The source for the method to try to find.</param>
    /// <param name="m">The found method or null.</param>
    /// <returns>True if the method was found, false otherwise.</returns>
    public bool TryFind(Source source, out Method? m) => this.Methods.TryGetValue(new Method(source), out m);

    /// <summary>Gets or adds a class with the given source.</summary>
    /// <param name="source">The source for the class.</param>
    /// <returns>The class for the given source.</returns>
    public Class GetOrAddClass(Source source) {
        Class c = new(source);
        if (this.Classes.TryGetValue(c, out Class? found)) return found;
        this.Classes.Add(c);
        return c;
    }

    /// <summary>Gets or adds a method with the given source.</summary>
    /// <param name="source">The source for the method.</param>
    /// <returns>The method for the given source.</returns>
    public Method GetOrAddMethod(Source source) {
        Method m = new(source);
        if (this.Methods.TryGetValue(m, out Method? found)) return found;
        this.Methods.Add(m);
        return m;
    }

    /// <summary>Add a new participation between the given method and class.</summary>
    /// <param name="m">The method that the class is participating with.</param>
    /// <param name="value">The value between zero exclusively and one inclusively.</param>
    /// <param name="c">the class that is participating.</param>
    /// <returns>Returns the creates participation.</returns>
    public static Participation Add(Method m, double value, Class c) {
        Participation p = new(m, value, c);
        Add(p);
        return p;
    }

    /// <summary>This adds a participation into the method and class that it connects.</summary>
    /// <remarks>This doesn't add the method nor class to the project.</remarks>
    /// <param name="p">The participation to add.</param>
    public static void Add(Participation p) {
        p.Method.Participation.Add(p);
        p.Class.Participation.Add(p);
    }

    /// <summary>This removes the participation from the class and method in the participation.</summary>
    /// <remarks>This will not remove the class nor method from the project.</remarks>
    /// <param name="p">The participation to remove.</param>
    public static void Remove(Participation p) {
        p.Class.Participation.Remove(p);
        p.Method.Participation.Remove(p);
    }

    /// <summary>Runs participation normalization on all methods.</summary>
    public void Normalize() => this.Methods.ForAll(m => m.Normalize());

    /// <summary>Creates a new empty project.</summary>
    public Project() { }

    /// <summary>Loads a project from a YAML root node.</summary>
    /// <param name="root">The YAML root node to load.</param>
    public Project(Node root) {
        Commons.Data.Yaml.Object obj = root.AsObject();
        Reader locs = Reader.Read(obj.TryReadNode("locs"));
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

    internal class LoaderHelper(Reader locs): IKeyResolver {
        public readonly Reader Locations = locs;
        public List<Class> Classes = [];
        public List<Method> Methods = [];

        public Source ReadSource(Node node) {
            Commons.Data.Yaml.Object obj = node.AsObject();
            Location loc = obj.ReadLocation(this.Locations, "loc");
            string name = obj.ReadString("name");
            return new Source(name, loc);
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
                _ => throw new Exception("Unexpected name for a key: " + name)
            };
        }
    }

    /// <summary>Gets the YAML node for this project that can be saved.</summary>
    /// <returns>The node for this project.</returns>
    public Node ToNode() {
        Factory locFactory = new();
        foreach (Method m in this.Methods) locFactory.Add(m.Source.Location);
        foreach (Class c in this.Classes) locFactory.Add(c.Source.Location);
        Writer locWriter = locFactory.Build();

        ToNodeHelper helper = new(locWriter);
        List<Method> methods = [.. this.Methods];
        List<Class> classes = [.. this.Classes];
        for (int i = 0; i < methods.Count; ++i) helper.Methods.Add(methods[i], i);
        for (int i = 0; i < classes.Count; ++i) helper.Classes.Add(classes[i], i);

        Commons.Data.Yaml.Object obj = new();
        obj.Add("locs", locWriter.Write());
        obj.Add("methods", Commons.Data.Yaml.Array.FromList(helper, methods));
        obj.Add("classes", Commons.Data.Yaml.Array.FromList(helper, classes));
        return obj;
    }

    internal class ToNodeHelper(Writer locs) {
        public readonly Writer Locations = locs;
        public Dictionary<Class, int> Classes = [];
        public Dictionary<Method, int> Methods = [];
    }
}
