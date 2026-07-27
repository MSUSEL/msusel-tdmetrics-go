using Commons.Extensions;
using Yaml = Commons.Data.Yaml;

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
    /// This is the WMC (Weight Method Class), that is the sum of all the McCabe's Cyclomatic Complexity (CC)
    /// for the whole class. It counts the number of branch instructions in the class's methods.
    /// </summary>
    /// <remarks>Since CK gets this wrong, this will use PMD's CycloTotal.</remarks>
    public int Wmc => this.CycloTotal;

    /// <summary>
    /// This is the CK version of WMC that is incorrect.
    /// </summary>
    /// <remarks>
    /// WARNING: CK's algorithm for finding WMC is Incorrect!!!
    ///
    /// 1. Counting inline binary comparison operators (`==`, `>`, `<`, etc.) is incorrect:
    /// 
    ///    If CK is incrementing the CC for standard relational/comparison operators, it is deviating
    ///    from McCabe's definition. CC is a measure of the *Control Flow Graph* (CFG) — specifically,
    ///    decision points that branch the execution path. A simple boolean comparison like `a == b`
    ///    does not branch the code unless it is evaluated inside a control flow statement
    ///    (like an `if` or `while`), and in those cases, the `if` or `while` itself is what generates
    ///    the +1 complexity, not the operator.
    ///
    /// 2. Undercounting nested ternary operators (`?:`) is incorrect (and buggy):
    ///
    ///    A ternary operator (`? :`) is functionally identical to an `if-else` block. Therefore,
    ///    **every single instance** of a ternary operator should add +1 to the complexity.
    ///    If `int a = c1 ? (c2 ? 12 : 34) : (c3 ? 56 : 78)` is only yielding +1, CK is undercounting.
    ///    Because CK is built on top of the Eclipse JDT (Java Development Tools) AST parser, this specific
    ///    bug usually happens when an AST Visitor's `visit(ConditionalExpression node)` method calculates
    ///    the count but fails to `return true;` (which tells the parser to recursively traverse the child
    ///    nodes of that expression).
    ///
    /// 3. Counting bitwise-AND (`&`) and bitwise-OR (`|`) is incorrect:
    ///
    ///    - **Logical operators** (`&&`, `||`) are short-circuiting. `if (A && B)` essentially behaves as
    ///      `if (A) { if (B) { ... } }`. Because the second condition might not execute based on the first,
    ///      there is a hidden branch in the control flow graph. They *should* add +1.
    ///    - **Bitwise operators** (`&`, `|`) do not short-circuit. Both sides of the expression are always
    ///      fully evaluated linearly before the operator is applied. They create no branches and *should not*
    ///      add to the CC.
    /// 
    /// 4. Other possible issues that I haven't notices yet.
    /// </remarks>
    /// <see cref="https://github.com/mauricioaniche/ck/blob/master/src/main/java/com/github/mauricioaniche/ck/metric/WMC.java"/>
    public int CkWmc => this.Ck.TryReadInt("wmc");

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

    /// <summary>The sum of McCabe’s Cyclomatic Complexity for all methods in this class, i.e. WMC (Weighted Method Count)</summary>
    /// <see cref="https://docs.pmd-code.org/latest/pmd_rules_java_design.html#cyclomaticcomplexity"/>
    /// <see cref="https://docs.pmd-code.org/apidocs/pmd-java/7.26.0/net/sourceforge/pmd/lang/java/metrics/JavaMetrics.html#CYCLO"/>
    /// <see cref="https://github.com/pmd/pmd/blob/main/pmd-java/src/main/java/net/sourceforge/pmd/lang/java/metrics/internal/CycloVisitor.java"/>
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
            List<Yaml.Object> ck  = this.methodsByLine("methods_ck");
            List<Yaml.Object> pmd = this.methodsByLine("methods_pmd");
            List<MethodMetrics> methods = [.. ck.Merge(pmd, compareByLines).
                Squish(squishByNamePmdFirst).
                Squish(squishByNameCkFirst).
                Select((t) => tupleToMethod(this, t))];
            field = methods.AsReadOnly();
            return field;
        }
    }

    static private int compareByLines(Yaml.Object ck, Yaml.Object pmd) =>
        ck.TryReadInt("line", -1) - pmd.TryReadInt("line", -1);

    static private Tuple<Yaml.Object?, Yaml.Object?>? squishByNamePmdFirst(Tuple<Yaml.Object?, Yaml.Object?> prev, Tuple<Yaml.Object?, Yaml.Object?> cur) {
        // Look for this shape since PMD's line will normally be less than CK's line if they can match:
        // prev: [  -  | PMD ]
        // cur:  [ CK  |  -  ]
        if (prev.Item1 is not null || prev.Item2 is null) return null;
        if (cur.Item1 is null || cur.Item2 is not null) return null;

        string ckSig = cur.Item1.TryReadString("method");
        if (string.IsNullOrEmpty(ckSig)) return null;
        ckSig = ckSig.Split("/")[0];

        string pmdSig = prev.Item2.TryReadString("signature");
        if (string.IsNullOrEmpty(pmdSig)) return null;
        pmdSig = pmdSig.Split("(")[0];

        if (ckSig != pmdSig) return null;
        return new(cur.Item1, prev.Item2);
    }
    
    static private Tuple<Yaml.Object?, Yaml.Object?>? squishByNameCkFirst(Tuple<Yaml.Object?, Yaml.Object?> prev, Tuple<Yaml.Object?, Yaml.Object?> cur) {
        // Look for this shape so that this will handle the rare case where PMD's line is greater than CK's line if they can match:
        // cur:  [ CK  |  -  ]
        // prev: [  -  | PMD ]
        if (prev.Item1 is null || prev.Item2 is not null) return null;
        if (cur.Item1 is not null || cur.Item2 is null) return null;

        string ckSig = prev.Item1.TryReadString("method");
        if (string.IsNullOrEmpty(ckSig)) return null;
        ckSig = ckSig.Split("/")[0];

        string pmdSig = cur.Item2.TryReadString("signature");
        if (string.IsNullOrEmpty(pmdSig)) return null;
        pmdSig = pmdSig.Split("(")[0];

        if (ckSig != pmdSig) return null;
        return new(prev.Item1, cur.Item2);
    }

    static private MethodMetrics tupleToMethod(DeclMetrics parent, Tuple<Yaml.Object?, Yaml.Object?> t) =>
        new(parent, t.Item1 ?? new Yaml.Object(), t.Item2 ?? new Yaml.Object());

    /// <summary>Gets the methods from the named JSON object sorted by the methods' line.</summary>
    /// <param name="name">The name of the JSON object to get.</param>
    /// <returns>The set of method objects sorted by the methods' line.</returns>
    private List<Yaml.Object> methodsByLine(string name) {
        SortedDictionary<int, Yaml.Object> obj = [];
        Yaml.Array arr = this.Object.TryReadNode(name)?.AsArray() ?? new Yaml.Array();
        foreach (Yaml.Node n in arr.Items) {
            Yaml.Object method = n.AsObject();
            int line = method.TryReadInt("line", -1);
            if (line <= 0)
                throw new Exception("Line number was not defined or invalid in " + name + " for class " + this.FullName);
            if (obj.ContainsKey(line))
                throw new Exception("Multiple methods on the same line " + line + " in " + name + " for class " + this.FullName);
            obj[line] = method;
        }
        return [.. obj.Values];
    }
}
