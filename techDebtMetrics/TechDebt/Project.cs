using System;
using System.IO;
using System.Linq;
using SCG = System.Collections.Generic;

namespace TechDebt;

public class Project {

    public readonly SCG.SortedSet<Class> Classes = [];

    public readonly SCG.SortedSet<Method> Methods = [];

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

        SCG.List<(Class Class, double Value)> entries = [.. m.Participation.Select(p => (p.Class, p.Value))];
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

    public void Save(string path) => File.WriteAllText(path, this.Serialize());

    public void Load(string path) => this.Deserialize(File.ReadAllText(path));

    public string Serialize() {

        return "";
    }

    public void Deserialize(string text) {
    }
}
