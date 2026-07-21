using Yaml = Commons.Data.Yaml;
using System.Text.RegularExpressions;

namespace GroundTruth;

public class DeclMetrics(Yaml.Object obj) {

    /// <summary>The raw JSON class containing the class information.</summary>
    public readonly Yaml.Object Object = obj;

    /// <summary>Indicates the type of this declaration.</summary>
    public DeclType Type => DeclTypeExt.FromName(this.Object.TryReadString("type"));
    
    /// <summary>The path to the file this class was defined in.</summary>
    public string File => this.Object.TryReadString("file");

    /// <summary>Checks if the file is in src/test/.</summary>
    public bool InTestPath => this.File.StartsWith("src/test/");

    /// <summary>Name of the class with the path dot separated.</summary>
    /// <example>org.apache.bcel.generic.InstructionList</example>
    public string FullName => this.Object.TryReadString("class");

    /// <summary>Name of just the class without any of the path.</summary>
    public string Name => this.FullName.Split('.').Last();

    /// <summary>Raw JSON from Mauricioaniche's CK for this current class/object.</summary>
    /// <see cref="https://github.com/mauricioaniche/ck"/>
    /// <remarks>The documentation on the values from Ck was barrowed from the github repo.</remarks>
    public Yaml.Object Ck => this.Object.TryReadNode("ck")?.AsObject() ?? new Yaml.Object();

    /// <summary>Indicates if the Ck object for this class was defined.</summary>
    public bool HasCk => this.Ck.Count > 0;

    /// <summary>
    /// CBO (Coupling between objects): Counts the number of dependencies a class has.
    /// The tools checks for any type used in the entire class (field declaration, method return types,
    /// variable declarations, etc). It ignores dependencies to Java itself (e.g. java.lang.String).
    /// </summary>
    public int Cbo => this.Ck.TryReadInt("cbo");
 
    /// <summary>
    /// CBO Modified (Coupling between objects): Counts the number of dependencies a class has.
    /// It is very similar to the CKTool's original CBO. However, this metric considers a dependency
    /// from a class as being both the references the type makes to others and the references
    /// that it receives from other types.
    /// </summary>
    public int CboModified => this.Ck.TryReadInt("cboModified");

    /// <summary>
    /// FAN-IN: Counts the number of input dependencies a class has, i.e, the number of classes that
    /// reference a particular class. For instance, given a class X, the fan-in of X would be the number
    /// of classes that call X by referencing it as an attribute, accessing some of its attributes,
    /// invoking some of its methods, etc.
    /// </summary>    
    public int FanIn => this.Ck.TryReadInt("fanin");
    
    /// <summary>
    /// FAN-OUT: Counts the number of output dependencies a class has, i.e, the number of other classes
    /// referenced by a particular class. In other words, given a class X, the fan-out of X is the number
    /// of classes called by X via attributes reference, method invocations, object instances, etc.
    /// </summary>
    public int FanOut => this.Ck.TryReadInt("fanout");

    /// <summary>
    /// DIT (Depth Inheritance Tree): It counts the number of "fathers" a class has. All classes have
    /// DIT at least 1 (everyone inherits java.lang.Object). In order to make it happen, classes must
    /// exist in the project (i.e. if a class depends upon X which relies in a jar/dependency file,
    /// and X depends upon other classes, DIT is counted as 2).
    /// </summary>
    public int Dit => this.Ck.TryReadInt("dit");

    /// <summary>
    /// NOC (Number of Children): It counts the number of immediate subclasses that a particular class has.
    /// </summary>
    public int Noc => this.Ck.TryReadInt("noc");

    /// <summary>
    /// This is the sum of all the WMC for the whole class.
    /// WMC (Weight Method Class) or McCabe's complexity.
    /// It counts the number of branch instructions in a class.
    /// </summary>
    /// <see cref="https://github.com/mauricioaniche/ck/blob/master/src/main/java/com/github/mauricioaniche/ck/metric/WMC.java"/>
    public int Wmc => this.Ck.TryReadInt("wmc");

    /// <summary>
    /// LOC (Lines of code): It counts the lines of count, ignoring empty lines and comments
    /// (i.e., it's Source Lines of Code, or SLOC). The number of lines here might be a bit different
    /// from the original file, as we use JDT's internal representation of the source code to calculate it.
    /// </summary>
    public int Loc => this.Ck.TryReadInt("loc");
    
    /// <summary>
    /// TCC (Tight Class Cohesion): Measures the cohesion of a class with a value range from 0 to 1.
    /// TCC measures the cohesion of a class via direct connections between visible methods,
    /// two methods or their invocation trees access the same class variable.
    /// </summary>
    public double Tcc => this.Ck.TryReadDouble("tcc");
    
    /// <summary>
    /// LCC (Loose Class Cohesion): Similar to TCC but it further includes the number of indirect
    /// connections between visible classes for the cohesion calculation.
    /// Thus, the constraint LCC >= TCC holds always.
    /// </summary>
    public double Lcc => this.Ck.TryReadDouble("lcc");
    
    /// <summary>
    /// NOSI (Number of static invocations): Counts the number of invocations to static methods.
    /// It can only count the ones that can be resolved by the JDT.
    /// </summary>
    public int Nosi => this.Ck.TryReadInt("nosi");

    /// <summary>
    /// LCOM (Lack of Cohesion of Methods): Calculates LCOM metric. This is the very first version of metric,
    /// which is not reliable. LCOM-HS can be better (hopefully, you will send us a pull request).
    /// </summary>
    public int Lcom => this.Ck.TryReadInt("lcom");

    /// <summary>
    /// LCOM* (Lack of Cohesion of Methods): This metric is a modified version of the current version
    /// of LCOM implemented in CK Tool. LCOM* is a normalized metric that computes the lack of cohesion
    /// of class within a range of 0 to 1. Then, the closer to 1 the value of LCOM* in a class, the less
    /// the cohesion degree of this respective class. The closer to 0 the value of LCOM* in a class,
    /// the most the cohesion of this respective class. This implementation follows the third version
    /// of LCOM* defined in [1].
    /// </summary>
    /// <remarks>
    /// Reference: [1] Henderson-Sellers, Brian, Larry L. Constantine and Ian M. Graham.
    /// “Coupling and cohesion (towards a valid metrics suite for object-oriented analysis and design).”
    /// Object Oriented Systems 3 (1996): 143-158.
    /// </remarks>
    public int Lcom2 => this.Ck.TryReadInt("lcom*");

    /// <summary>
    /// RFC (Response for a Class): Counts the number of unique method invocations in a class.
    /// As invocations are resolved via static analysis, this implementation fails when a method
    /// has overloads with same number of parameters, but different types.
    /// </summary>
    public int Rfc => this.Ck.TryReadInt("rfc");

    /// <summary>Raw JSON from PMD for this current class/object.</summary>
    /// <see cref="https://docs.pmd-code.org/latest/pmd_rules_java.html"/>
    /// <remarks>The documentation on the values from PMD was barrowed from the PMD website.</remarks>
    public Yaml.Object Pmd  => this.Object.TryReadNode("pmd")?.AsObject() ?? new Yaml.Object();
    
    /// <summary>Indicates if the PMD object for this class was defined.</summary>
    public bool HasPmd => this.Pmd.Count > 0;

    /// <summary>
    /// NCSS (Non-Commenting Source Statements) metric to determine the sum of
    /// the number of lines of code in a class.
    /// </summary>
    /// <see cref="https://docs.pmd-code.org/latest/pmd_rules_java_design.html#ncsscount"/>
    public int Ncss => this.Pmd.TryReadInt("ncss");

    /// <summary>The highest NCSS (Non-Commenting Source Statements) of any method in this class.</summary>
    public int NcssHighest => this.Pmd.TryReadInt("ncss_highest");

    /// <summary>The sum of McCabe’s Cyclomatic Complexity for all methods in this class.</summary>
    /// <see cref="https://docs.pmd-code.org/latest/pmd_rules_java_design.html#cyclomaticcomplexity"/>
    /// <see cref="https://docs.pmd-code.org/apidocs/pmd-java/7.26.0/net/sourceforge/pmd/lang/java/metrics/JavaMetrics.html#CYCLO"/>
    public int CycloTotal => this.Pmd.TryReadInt("cyclo_total");

    /// <summary>The highest Cyclomatic Complexity of any method in this class.</summary>
    public int CycloHighest => this.Pmd.TryReadInt("cyclo_highest");

    /// <summary>The number of unique attributes, local variables, and return types within an object.</summary>
    /// <see cref="https://docs.pmd-code.org/latest/pmd_rules_java_design.html#couplingbetweenobjects"/>
    public int Coupling => this.Pmd.TryReadInt("coupling");

    /// <summary>
    /// Sum of the statistical complexity of the operations in the class
    /// from the PMD God class rule.
    /// </summary>
    /// <see cref="https://docs.pmd-code.org/latest/pmd_rules_java_design.html#godclass"/>
    /// <see cref="https://docs.pmd-code.org/apidocs/pmd-java/7.26.0/net/sourceforge/pmd/lang/java/metrics/JavaMetrics.html#WEIGHED_METHOD_COUNT"/>
    public int GodWmc => this.Pmd.TryReadInt("god_wmc");

    /// <summary>
    /// Number of usages of foreign attributes, both directly and through accessors.
    /// "Foreign" hier means "not belonging to this", although field accesses to fields declared
    /// in the enclosing class are not considered foreign. 
    /// </summary>
    /// <see cref="https://docs.pmd-code.org/latest/pmd_rules_java_design.html#godclass"/>
    /// <see cref="https://docs.pmd-code.org/apidocs/pmd-java/7.26.0/net/sourceforge/pmd/lang/java/metrics/JavaMetrics.html#ACCESS_TO_FOREIGN_DATA"/>
    public int GodAtfd => this.Pmd.TryReadInt("god_atfd");

    /// <summary>
    /// The relative number of method pairs of a class that access in common at least one attribute of the measured class.
    /// TCC only counts direct attribute accesses, that is, only those attributes that are accessed in the body of the method.
    /// The value is a double between 0 and 1.
    /// </summary>
    /// <see cref="https://docs.pmd-code.org/latest/pmd_rules_java_design.html#godclass"/>
    /// <see cref="https://docs.pmd-code.org/apidocs/pmd-java/7.26.0/net/sourceforge/pmd/lang/java/metrics/JavaMetrics.html#TIGHT_CLASS_COHESION"/>
    public double GodTcc => this.Pmd.TryReadDouble("god_tcc_pct");

    /// <summary>The number of public methods and public attributes in this class.</summary>
    /// <see cref="https://docs.pmd-code.org/latest/pmd_rules_java_design.html#excessiveparameterlist"/>
    public int PublicCount => this.Pmd.TryReadInt("public_count");

    /// <summary>The methods that are members of this class.</summary>
    public IReadOnlyList<MethodMetrics> Methods {
        get {
            if (field is not null) return field;

            Dictionary<int, Yaml.Object> ck  = this.methodsByLine("methods_ck");
            Dictionary<int, Yaml.Object> pmd = this.methodsByLine("methods_pmd");

            List<int> lines = [.. ck.Keys.Union(pmd.Keys)];
            lines.Sort();

            List<MethodMetrics> methods = [];
            foreach (int line in lines) {
                ck.TryGetValue(line, out Yaml.Object? ckMethod);
                pmd.TryGetValue(line, out Yaml.Object? pmdMethod);
                methods.Add(new(this, ckMethod ?? new Yaml.Object(), pmdMethod ?? new Yaml.Object()));
            }

            field = methods.AsReadOnly();
            return field;
        }
    }

    /// <summary>Gets the methods from the named JSON object keyed by the methods' line.</summary>
    /// <param name="name">The name of the JSON obejct to get.</param>
    /// <returns>The set of method objects keyed by the methods' line.</returns>
    private Dictionary<int, Yaml.Object> methodsByLine(string name) {
        Dictionary<int, Yaml.Object> ck = [];
        Yaml.Array ckArray = this.Object.TryReadNode(name)?.AsArray() ?? new Yaml.Array();
        foreach (Yaml.Node ckNode in ckArray.Items) {
            Yaml.Object ckMethod = ckNode.AsObject();
            int line = ckMethod.TryReadInt("line");
            if (ck.ContainsKey(line))
                throw new Exception("Multiple methods on the same line " + line + " in " + name + " for class " + this.FullName);
            ck[line] = ckMethod;
        }
        return ck;
    }
}
