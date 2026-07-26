using Yaml = Commons.Data.Yaml;

namespace GroundTruth;

public class MethodMetrics(DeclMetrics parent, Yaml.Object ckObj, Yaml.Object pmdObj) {
    
    /// <summary>Raw JSON from Mauricioaniche's CK for this current method.</summary>
    /// <see cref="https://github.com/mauricioaniche/ck"/>
    /// <remarks>The documentation on the values from CK was barrowed from the GitHub repo.</remarks>
    public readonly Yaml.Object Ck = ckObj;

    /// <remarks>The parent declaration that this method is declared inside of.</remarks>
    public readonly DeclMetrics Parent = parent;

    /// <summary>Indicates if the Ck object for this method was defined.</summary>
    public bool HasCk => this.Ck.Count > 0;

    /// <summary>Gets the object's full name and the method's name.</summary>
    public string FullName => this.Parent.FullName + "#" + this.Name;

    /// <summary>The name of the method without the path.</summary>
    public string Name {
        get {
            string sig = this.CkSignature;
            if (!string.IsNullOrEmpty(sig)) return sig.Split("/")[0];
            sig = this.PmdSignature;
            if (!string.IsNullOrEmpty(sig)) return sig.Split("(")[0];
            throw new Exception("Name of method in " + Parent.FullName + " was undefined");
        }
    }

    /// <summary>The line number that the method starts on.</summary>
    /// <remarks>
    /// The line numbers for CK and PMD may be slightly different (see CkLine and PmdLine).
    /// This will try to return the PMD line number, then fall back to the CK line number,
    /// then error if no line number could be determined.
    /// </remarks>
    public int Line {
        get {
            int line = this.PmdLine;
            if (line > 0) return line;
            line = this.CkLine;
            if (line > 0) return line;
            throw new Exception("Line of method in " + Parent.FullName + " was undefined");
        }
    }

    /// <summary>Indicates if this method is a constructor or not.</summary>
    public bool Constructor => this.Ck.ReadBool("constructor");

    /// <summary>
    /// The CK line number that the method starts on.
    /// If the line isn't given then -1 is returned.
    /// </summary>
    /// <remarks>
    /// The CK line number seems to always be the line of the open bracket for the method body `{`,
    /// unless the method has no body, then it seems to use the first line of the comment.
    /// </remarks>
    public int CkLine => this.Ck.TryReadInt("line", -1);

    /// <summary>
    /// The string of the CK signature of this method.
    /// If the signature isn't given then an empty string is returned.
    /// </summary>
    /// <example>Variable/2[Mini.ASTIdent,boolean]</example>
    public string CkSignature => this.Ck.TryReadString("method");

    /// <summary>
    /// FAN-IN: Counts the number of input dependencies a class has, i.e, the number of classes that
    /// reference a particular class. For instance, given a class X, the fan-in of X would be the number
    /// of classes that call X by referencing it as an attribute, accessing some of its attributes,
    /// invoking some of its methods, etc.
    /// </summary>  
    public int FanIn => this.Ck.ReadInt("fanin");
    
    /// <summary>
    /// FAN-OUT: Counts the number of output dependencies a class has, i.e, the number of other classes
    /// referenced by a particular class. In other words, given a class X, the fan-out of X is the number
    /// of classes called by X via attributes reference, method invocations, object instances, etc.
    /// </summary>
    public int FanOut => this.Ck.ReadInt("fanout");

    /// <summary>
    /// The number of branch instructions in this method.
    /// WMC (Weight Method Class) or McCabe's complexity.
    /// </summary>
    public int Wmc => this.Ck.ReadInt("wmc");
    
    /// <summary>
    /// LOC (Lines of code): It counts the lines of count, ignoring empty lines and comments
    /// (i.e., it's Source Lines of Code, or SLOC). The number of lines here might be a bit different
    /// from the original file, as we use JDT's internal representation of the source code to calculate it.
    /// This will return 0 if the LOC was not defined in the CK metrics.
    /// </summary>
    public int Loc => this.Ck.TryReadInt("loc");
    
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
    /// RFC (Response for a Class): Counts the number of unique method invocations in a class.
    /// As invocations are resolved via static analysis, this implementation fails when a method
    /// has overloads with same number of parameters, but different types.
    /// </summary>
    public int Rfc => this.Ck.TryReadInt("rfc");
    
    /// <summary>The modifiers for this methods.</summary>
    /// <remarks>Can be decoded using `org.eclipse.jdt.core.dom.Modifier`.</remarks>
    public Modifiers Modifiers => new(this.Ck.TryReadInt("modifiers"));

    /// <summary>Raw JSON from PMD for this current method.</summary>
    /// <see cref="https://docs.pmd-code.org/latest/pmd_rules_java.html"/>
    /// <remarks>The documentation on the values from PMD was barrowed from the PMD website.</remarks>
    public readonly Yaml.Object Pmd = pmdObj;

    /// <summary>Indicates if the PMD object for this method was defined.</summary>
    public bool HasPmd => this.Pmd.Count > 0;

    /// <summary>
    /// The PMD line number that the method starts on.
    /// If the line isn't given then -1 is returned.
    /// </summary>
    /// <remarks>
    /// The PMD metrics' line number seems to be the line of the method's name identifier
    /// (or possibly the first modifier, but not the annotations).
    /// </remarks>
    public int PmdLine => this.Pmd.TryReadInt("line", -1);

    /// <summary>
    /// The string of the PMD signature of this method.
    /// If the signature isn't given then an empty string is returned.
    /// </summary>
    /// <example>Variable(ASTIdent, boolean)</example>
    public string PmdSignature => this.Pmd.TryReadString("signature");

    /// <summary>
    /// NCSS (Non-Commenting Source Statements) metric to determine the number
    /// of lines of code in this method or constructor.
    /// </summary>
    /// <see cref="https://docs.pmd-code.org/latest/pmd_rules_java_design.html#ncsscount"/>
    public int Ncss => this.Pmd.ReadInt("ncss");

    /// <summary>The McCabe’s Cyclomatic Complexity for this method.</summary>
    /// <see cref="https://docs.pmd-code.org/latest/pmd_rules_java_design.html#cyclomaticcomplexity"/>
    /// <see cref="https://docs.pmd-code.org/apidocs/pmd-java/7.26.0/net/sourceforge/pmd/lang/java/metrics/JavaMetrics.html#CYCLO"/>
    public int Cyclo => this.Pmd.ReadInt("cyclo");
    
    /// <summary>
    /// The NPath complexity of a method is the number of acyclic execution paths through that method.
    /// While cyclomatic complexity counts the number of decision points in a method, NPath counts the number of full paths
    /// from the beginning to the end of the block of the method.
    /// That metric grows exponentially, as it multiplies the complexity of statements in the same block. 
    /// </summary>
    /// <see cref="https://docs.pmd-code.org/latest/pmd_rules_java_design.html#npathcomplexity"/>
    /// <see cref="https://docs.pmd-code.org/apidocs/pmd-java/7.26.0/net/sourceforge/pmd/lang/java/metrics/JavaMetrics.html#NPATH"/>
    public int NPath => this.Pmd.ReadInt("npath");

    /// <summary>The number of parameters for this method.</summary>
    /// <see cref="https://docs.pmd-code.org/latest/pmd_rules_java_design.html#excessiveparameterlist"/>
    public int Params => this.Pmd.TryReadInt("params");
}
