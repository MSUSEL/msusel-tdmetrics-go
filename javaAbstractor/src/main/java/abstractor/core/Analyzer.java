package abstractor.core;

import java.util.List;
import java.util.Map;
import java.util.SortedSet;
import java.util.TreeSet;
import java.util.TreeMap;

import spoon.reflect.code.*;
import spoon.reflect.cu.SourcePosition;
import spoon.reflect.declaration.*;
import spoon.reflect.reference.*;
import spoon.support.reflect.CtExtendedModifier;
import spoon.support.reflect.code.CtIfImpl;
import abstractor.core.constructs.*;
import abstractor.core.json.JsonFormat;
import abstractor.core.json.JsonHelper;
import abstractor.core.log.Logger;
import abstractor.core.spoonUtils.SpoonUtils;

public class Analyzer {
    private static final boolean logElementTree = false;
    private static final boolean logUsage       = false;

    private final Abstractor abs;
    private final Logger     log;
    public  final Location   loc;

    private final Map<Integer, Integer> minColumn = new TreeMap<>();
    private int minLine;
    private int maxLine;
    
    private int complexity;

    private boolean getter;
    private boolean setter;

    private final SortedSet<Ref<? extends Construct>> invokes = new TreeSet<>();
    private final SortedSet<Ref<? extends Construct>> reads   = new TreeSet<>();
    private final SortedSet<Ref<? extends Construct>> writes  = new TreeSet<>();

    public Analyzer(Abstractor abs, Location loc) {
        this.abs        = abs;
        this.log        = abs.log;
        this.loc        = loc;
        this.minLine    = Integer.MAX_VALUE;
        this.complexity = 1;
    }

    public Metrics getMetrics() {
        final int lineCount = this.maxLine - this.minLine + 1;
        if (lineCount <= 0) return new Metrics(this.loc);

        final int codeCount = this.minColumn.size();
        final int indents   = this.calcIndent();
        
        if (logElementTree) {
            this.log.log("+- codeCount:  " + codeCount);
            this.log.log("+- complexity: " + this.complexity);
            this.log.log("+- indents:    " + indents);
            this.log.log("+- lineCount:  " + lineCount);
        }

        return new Metrics(this.loc,
            codeCount, this.complexity, indents, lineCount,
            this.getter, this.setter,
            this.invokes, this.reads, this.writes);
    }

    private int calcIndent() {
        final TreeSet<Integer> columns = new TreeSet<>();
        columns.addAll(this.minColumn.values());
        
        final TreeMap<Integer, Integer> indentMap = new TreeMap<>();
        int index = 0;
        for (int col : columns) indentMap.put(col, index++);

        if (logElementTree) {
            this.log.log("+- minCols: " + this.minColumn);
            this.log.log("+- indents: " + indentMap);
        }

	    int indentSum = 0;
        for (int minCol : this.minColumn.values())
            indentSum += indentMap.get(minCol);
        return indentSum;
    }

    public void addMethod(CtExecutable<?> m) throws Exception {
        this.addPosition(m.getPosition());

        if (m instanceof CtMethod<?> method) {
            for (CtExtendedModifier em : method.getExtendedModifiers())
                this.addPosition(em.getPosition());
        }

        for (CtParameter<?> p : m.getParameters())
            this.addPosition(p.getPosition());

        final CtBlock<?> block = m.getBody();
        if (block != null) {
            this.addPosition(block.getPosition());

            final List<CtStatement> stmts = block.getStatements();
            for (CtStatement st : stmts)
                this.addElement(st);

            if (stmts.size() == 1) {
                final CtStatement onlySt = stmts.get(0);
                if (detectGetter(m, onlySt)) this.getter = true;
                else if (detectSetter(m, onlySt)) this.setter = true;
            }
        }
    }

    private void addElement(CtElement elem) throws Exception {
        // Skip over any comments in the code.
        if (elem instanceof CtComment) return;

        if (elem instanceof CtLocalVariable v) {
            for (CtExtendedModifier em : v.getExtendedModifiers())
                this.addPosition(em.getPosition());
        }

        this.addPosition(elem.getPosition());
        this.complexity += addComplexity(elem);
        if (this.addUsage(elem)) {

            if (logElementTree) {
                this.log.log("+- " + SpoonUtils.describeElem(elem));
                this.log.push("|  ");
            }

            for (CtElement child : elem.getDirectChildren()) this.addElement(child);

            if (logElementTree) this.log.pop();
        }
    }

    private void addPosition(SourcePosition pos) {
        if (!pos.isValidPosition()) return;
        //this.log.log("  adding <"+pos.getLine()+", "+pos.getColumn()+ "> <"+pos.getEndLine()+", "+pos.getEndColumn()+">");
        this.addPosition(pos.getLine(), pos.getColumn());
        this.addPosition(pos.getEndLine(), pos.getEndColumn());
    }

    private void addPosition(int line, int column) {
        if (this.minColumn.containsKey(line)) {
            int col = this.minColumn.get(line);
            column = Integer.min(column, col);
        }
        this.minColumn.put(line, column);
        this.minLine = Integer.min(line, this.minLine);
        this.maxLine = Integer.max(line, this.maxLine);
    }

    static private boolean isSimpleFetch(CtElement elem) {
        if (elem instanceof CtConstructorCall) return false;
        if (elem instanceof CtInvocation)      return false;
        if (elem instanceof CtBinaryOperator)  return false;
        if (elem instanceof CtUnaryOperator)   return false;
        for (CtElement child : elem.getDirectChildren()) {
            if (!isSimpleFetch(child)) return false;
        }
        return true;
    }

    static private boolean isObjectUsed(CtParameterReference<?> ref, CtElement elem) {
        if (elem.equals(ref)) return true;
        for (CtElement child : elem.getDirectChildren()) {
            if (isObjectUsed(ref, child)) return true;
        }
        return false;
    }

    static private boolean detectGetter(CtExecutable<?> m, CtStatement st) {
        if (m.getParameters().size() != 0)  return false;
        if (SpoonUtils.isVoid(m.getType())) return false;
        if (!(st instanceof CtReturn ret))  return false;
        if (!isSimpleFetch(ret))            return false;
        return true;
    }
    
    static private boolean detectSetter(CtExecutable<?> m, CtStatement st) {
        if (m.getParameters().size() > 1)           return false;
        if (!SpoonUtils.isVoid(m.getType()))        return false;
        if (!(st instanceof CtAssignment assign))   return false;
        if (!isSimpleFetch(assign.getAssigned()))   return false;
        if (!isSimpleFetch(assign.getAssignment())) return false;

        // Check for setters may have no parameters for assigning a
        // literal value, e.g. `func(b *Bar) Hide() { b.visible = false }`.
        if (m.getParameters().size() == 0) return true;

        // Make sure the parameter isn't used on the left hand side as in a
        // reversed setter, e.g. `func (b Bar) GetX(x *int) { x* = b.x }`,
        // The parameter may be used on the right hand side or not at all.
        // The parameter may not be used at all if the setter is part of an
        // interface requirement but the value assigned is to a default value.
        final CtParameter<?> param = m.getParameters().get(0);
        final CtParameterReference<?> ref = param.getReference();
        return !isObjectUsed(ref, assign.getAssigned());
    }

    /**
     * Gets the McCabe cyclomatic complexity for this element.
     * @param elem The element to get the complexity metric for.
     * @return The complexity metric for the given element.
     */
    static private int addComplexity(CtElement elem) {
        // Add one point for each conditional construct, such as an "if" condition.
        // The "else" does not add any additional complication.
        if (elem instanceof CtIf) return 1;

        // Add one point for each iterative structure (CtDo, CtForEach, CtFor, CtWhile).
        if (elem instanceof CtLoop) return 1;

        // Add one point for each case block in a switch statement.
        // The "default" does not add any additional complication.
        if (elem instanceof CtCase caseElem) {
            if (caseElem.getCaseExpressions().size() > 0) return 1;
            return 0;
        }
        
        // Add one point for any additional boolean condition, such as the use of && or ||
        if (elem instanceof CtBinaryOperator opElem) {
            if (opElem.getKind() == BinaryOperatorKind.AND) return 1;
            if (opElem.getKind() == BinaryOperatorKind.OR)  return 1;
            return 0;
        }

        return 0;
    }

    public String conString(Ref<? extends Construct> c) {
        JsonHelper jh = new JsonHelper();
        jh.writeKinds     = true;
        jh.writeIndices   = true;
        jh.writeRefs      = true;
        jh.refSkipResolve = false;
        return JsonFormat.Inline().format(c.refToJson(jh));
    }

    private void addInvoke(Ref<? extends Construct> c) {
        if (c == null)  return;
        if (logUsage) this.log.log("Adding invocation: " + this.conString(c));
        this.invokes.add(c);
    }

    private void addRead(Ref<? extends Construct> c) {
        if (c == null) return;
        if (logUsage) this.log.log("Adding read: " + this.conString(c));
        this.reads.add(c);
    }

    private void addWrite(Ref<? extends Construct> c) {
        if (c == null) return;
        if (logUsage) this.log.log("Adding write: " + this.conString(c));
        this.writes.add(c);
    }

    private boolean addUsage(CtElement elem) throws Exception {
        // Skip these and their children.
        if (elem instanceof CtAnnotation) return false;

        // Skip these but their children will still be checked.
        if (elem instanceof CtBinaryOperator) return true;
        if (elem instanceof CtBlock         ) return true;
        if (elem instanceof CtCase          ) return true;
        if (elem instanceof CtFor           ) return true;
        if (elem instanceof CtForEach       ) return true;
        if (elem instanceof CtIfImpl        ) return true;
        if (elem instanceof CtReturn        ) return true;
        if (elem instanceof CtSwitch        ) return true;
        if (elem instanceof CtThisAccess    ) return true;
        if (elem instanceof CtUnaryOperator ) return true;

        // Process these and their children.
        if (elem instanceof CtAssignment             as) { this.addAssignmentUsage(as);             return true; }
        if (elem instanceof CtClass                  cl) { this.addClassUsage(cl);                  return true; }
        if (elem instanceof CtConstructor            cc) { this.addConstructorUsage(cc);            return true; }
        if (elem instanceof CtConstructorCall        cc) { this.addConstructorCallUsage(cc);        return true; }
        if (elem instanceof CtExecutableReference    er) { this.addExecutableReferenceUsage(er);    return true; }
        if (elem instanceof CtFieldRead              fr) { this.addFieldReadUsage(fr);              return true; }
        if (elem instanceof CtFieldReference         fr) { this.addFieldReferenceUsage(fr);         return true; }
        if (elem instanceof CtFieldWrite             fw) { this.addFieldWriteUsage(fw);             return true; }
        if (elem instanceof CtInvocation             in) { this.addInvocationUsage(in);             return true; }
        if (elem instanceof CtLiteral                lt) { this.addLiteralUsage(lt);                return true; }
        if (elem instanceof CtLocalVariable          lv) { this.addLocalVariableUsage(lv);          return true; }
        if (elem instanceof CtLocalVariableReference lr) { this.addLocalVariableReferenceUsage(lr); return true; }
        if (elem instanceof CtMethod                 mt) { this.addMethodUsage(mt);                 return true; }
        if (elem instanceof CtPackageReference       pr) { this.addPackageReferenceUsage(pr);       return true; }
        if (elem instanceof CtParameterReference     pr) { this.addParameterReferenceUsage(pr);     return true; }
        if (elem instanceof CtTypeAccess             ta) { this.addTypeAccessUsage(ta);             return true; }
        if (elem instanceof CtTypeReference          tr) { this.addTypeReferenceUsage(tr);          return true; }
        if (elem instanceof CtVariableRead           vr) { this.addVariableReadUsage(vr);           return true; }
        if (elem instanceof CtVariableWrite          vw) { this.addVariableWrite(vw);               return true; }

        // Use to see elements (may produce a lot of output).
        this.log.notice("unimplemented addUsage: " + SpoonUtils.describeElem(elem));
        return true;
    }

    private void addInvocationUsage(CtInvocation<?> in) throws Exception {
        if (logUsage) this.log.log("addUsage.CtInvocation: " + SpoonUtils.describeElem(in));
        // Try to narrow the invocation edge to a per-call MethodInst using the
        // actual type arguments visible at this call site. Falls back to the
        // generic MethodDecl when no narrowing is possible.
        final Ref<? extends Construct> narrowed = this.abs.addMethodInstForCall(in);
        if (narrowed != null) {
            this.addInvoke(narrowed);
            return;
        }
        this.addExecutableReferenceUsage(in.getExecutable());
    }

    private void addFieldReadUsage(CtFieldRead<?> fr) throws Exception {
        if (logUsage) this.log.log("addUsage.CtFieldRead: " + SpoonUtils.describeElem(fr));
        final CtFieldReference<?> ref = fr.getVariable();

        final CtTypeReference<?> tr = ref.getDeclaringType();
        if (tr != null) {
            if (tr.isArray()) return;
        }

        final CtField<?> field = ref.getFieldDeclaration();
        if (field == null) {
            this.log.warning("addUsage.CtFieldRead: no field declaration for " + SpoonUtils.describeElem(ref));
            return;
        }
        this.addRead(this.abs.addSelection(field));
    }

    private void addTypeAccessUsage(CtTypeAccess<?> ta) throws Exception {
        if (logUsage) this.log.log("addUsage.CtTypeAccess: " + SpoonUtils.describeElem(ta));
        this.addRead(this.abs.addTypeDesc(ta.getAccessedType()));
    }

    private void addAssignmentUsage(CtAssignment<?,?> as) throws Exception {
        if (logUsage) this.log.log("addUsage.CtAssignment: " + SpoonUtils.describeElem(as));
        final CtTypeReference<?> tr = as.getType();
        if (!SpoonUtils.isVoid(tr)) this.addWrite(this.abs.addTypeDesc(tr));
    }

    private void addTypeReferenceUsage(CtTypeReference<?> tr) throws Exception {
        if (logUsage) this.log.log("addUsage.CtTypeReference: " + SpoonUtils.describeElem(tr));
        this.addListOfTypeArgsUsages(tr.getActualTypeArguments());
        if (!SpoonUtils.isVoid(tr)) this.addRead(this.abs.addTypeDesc(tr));
    }

    private void addPackageReferenceUsage(CtPackageReference pr) throws Exception {
        if (logUsage) this.log.log("addUsage.CtPackageReference: " + SpoonUtils.describeElem(pr));
        this.abs.addPackage(pr.getDeclaration());
    }

    private void addFieldWriteUsage(CtFieldWrite<?> fw) throws Exception {
        if (logUsage) this.log.log("addUsage.CtFieldWrite: " + SpoonUtils.describeElem(fw));
        final CtField<?> field = fw.getVariable().getFieldDeclaration();
        if (field == null) {
            this.log.warning("addUsage.CtFieldWrite: no field declaration for " + SpoonUtils.describeElem(fw));
            return;
        }
        this.addWrite(this.abs.addSelection(field));
    }
    
    private void addFieldReferenceUsage(CtFieldReference<?> fr) throws Exception {
        if (logUsage) this.log.log("addUsage.CtFieldReference: " + SpoonUtils.describeElem(fr));

        final CtTypeReference<?> tr = fr.getDeclaringType();
        if (tr != null) {
            if (tr.isArray()) return;
        }

        final CtField<?> field = fr.getFieldDeclaration();
        if (field == null) {
            this.log.warning("addUsage.CtFieldReference: no field declaration for " + SpoonUtils.describeElem(fr));
            return;
        }
        this.addRead(this.abs.addSelection(field));
    }

    private void addListOfTypeArgsUsages(List<CtTypeReference<?>> typeArgs) throws Exception {
        // TODO: Should we create a reference?
        if (typeArgs.size() > 0) this.log.notice("Unused type args: [" + SpoonUtils.describeElems(typeArgs) + "]");
    }
    
    private void addExecutableReferenceUsage(CtExecutableReference<?> er) throws Exception {
        if (logUsage) this.log.log("addUsage.CtExecutableReference: " + SpoonUtils.describeElem(er));
        this.addListOfTypeArgsUsages(er.getActualTypeArguments());
        if (er.isImplicit()) {
            this.log.notice("Ignoring implicit executable reference: " + SpoonUtils.describeElem(er));
            return;
        }

        // When this reference is the executable of a call, the parent handler
        // (addInvocationUsage / addConstructorCallUsage) has already added the
        // right invocation edge (a narrowed MethodInst when possible). Skip
        // here to avoid also adding the generic decl as a duplicate edge.
        final CtElement parent = er.getParent();
        if (parent instanceof CtAbstractInvocation) return;

        final CtExecutable<?> ex = er.getDeclaration();

        // TODO: Finished null investigation.
        if (ex == null) return;
        //Require.notNull(ex, "executable target was null for " + SpoonUtils.describeElem(er));

        if (ex instanceof CtMethod      mt) { this.addMethodUsage(mt);      return; }
        if (ex instanceof CtConstructor ct) { this.addConstructorUsage(ct); return; }

        this.log.warning("addUsage.CtExecutableReference: expected method for " + SpoonUtils.describeElem(er) + " and "+ SpoonUtils.describeElem(ex));
    }

    private void addMethodUsage(CtMethod<?> mt) throws Exception {
        if (logUsage) this.log.log("addUsage.CtMethod: " + SpoonUtils.describeElem(mt));
        this.addInvoke(this.abs.addDeclaration(mt));
    }

    private void addConstructorUsage(CtConstructor<?> ct) throws Exception {
        if (logUsage) this.log.log("addUsage.CtConstructor: " + SpoonUtils.describeElem(ct));
        this.addInvoke(this.abs.addDeclaration(ct));
    }

    private void addLiteralUsage(CtLiteral<?> lt) throws Exception {
        if (logUsage) this.log.log("addUsage.CtLiteral: " + SpoonUtils.describeElem(lt));
        this.addRead(this.abs.addTypeDesc(lt.getType()));
    }

    private void addVariableReadUsage(CtVariableRead<?> vr) throws Exception {
        if (logUsage) this.log.log("addUsage.CtVariableRead: " + SpoonUtils.describeElem(vr));
        this.addRead(this.abs.addTypeDesc(vr.getType()));
    }

    private void addParameterReferenceUsage(CtParameterReference<?> pr) throws Exception {
        if (logUsage) this.log.log("addUsage.CtParameterReference: " + SpoonUtils.describeElem(pr));
        this.addRead(this.abs.addTypeDesc(pr.getType()));
    }

    private void addConstructorCallUsage(CtConstructorCall<?> cc) throws Exception {
        if (logUsage) this.log.log("addUsage.CtConstructorCall: " + SpoonUtils.describeElem(cc));
        final CtExecutableReference<?> execRef = cc.getExecutable();

        final CtExecutable<?> exec = execRef.getDeclaration();
        if (exec instanceof CtConstructor<?> ctor) {
            if (ctor.isImplicit()) { // default constructor called.
                final CtTypeReference<?> dt = execRef.getDeclaringType();
                final Ref<? extends TypeDesc> td = this.abs.addTypeDesc(dt);
                if (td != null) this.writes.add(td);
            } else {
                // Narrow the ctor call to a MethodInst using the call's actual
                // type arguments; falls back to the generic MethodDecl.
                final Ref<? extends Construct> narrowed = this.abs.addMethodInstForCall(cc);
                this.addInvoke(narrowed != null ? narrowed : this.abs.addMethodDeclForConstructor(ctor));
            }
        } else {
            this.log.warning("addUsage.CtConstructorCall: expected constructor for " +
                SpoonUtils.describeElem(exec) + " from "+ SpoonUtils.describeElem(cc));
        }
    }

    private void addLocalVariableReferenceUsage(CtLocalVariableReference<?> lr) throws Exception {
        if (logUsage) this.log.log("addUsage.CtLocalVariableReference: " + SpoonUtils.describeElem(lr));
        this.addRead(this.abs.addTypeDesc(lr.getType()));
    }

    private void addLocalVariableUsage(CtLocalVariable<?> lv) throws Exception {
        if (logUsage) this.log.log("addUsage.CtLocalVariable: " + SpoonUtils.describeElem(lv));
        this.addWrite(this.abs.addTypeDesc(lv.getType()));
    }

    private void addVariableWrite(CtVariableWrite<?> vw) throws Exception {
        if (logUsage) this.log.log("addUsage.CtVariableWrite: " + SpoonUtils.describeElem(vw));
        this.addWrite(this.abs.addTypeDesc(vw.getType()));
    }

    private void addClassUsage(CtClass<?> c) throws Exception {
        if (logUsage) this.log.log("addUsage.CtClass: " + SpoonUtils.describeElem(c));
        this.addRead(this.abs.addObjectDecl(c));
    }
}
