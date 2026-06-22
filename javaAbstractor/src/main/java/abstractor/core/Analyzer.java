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

import abstractor.core.constructs.*;
import abstractor.core.log.Logger;
import abstractor.core.spoonUtils.SpoonUtils;

public class Analyzer {
    private static final boolean logElementTree = true; // TODO: Restore to false;
    private static final boolean logUsage       = true; // TODO: Restore to false;

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
        this.addUsage(elem);

        if (logElementTree) {
            this.log.log("+- " + SpoonUtils.describeElem(elem));
            this.log.push("|  ");
        }

        for (CtElement child : elem.getDirectChildren()) this.addElement(child);

        if (logElementTree) this.log.pop();
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

    private void addUsage(CtElement elem) throws Exception {
        if (elem instanceof CtInvocation          in) { this.addInvocationUsage(in);          return; }
        if (elem instanceof CtFieldRead           fr) { this.addFieldReadUsage(fr);           return; }
        if (elem instanceof CtTypeAccess          ta) { this.addTypeAccessUsage(ta);          return; }
        if (elem instanceof CtAssignment          as) { this.addAssignmentUsage(as);          return; }
        if (elem instanceof CtTypeReference       tr) { this.addTypeReferenceUsage(tr);       return; }
        if (elem instanceof CtPackageReference    pr) { this.addPackageReferenceUsage(pr);    return; }
        if (elem instanceof CtFieldReference      fr) { this.addFieldReferenceUsage(fr);      return; }
        if (elem instanceof CtExecutableReference er) { this.addExecutableReferenceUsage(er); return; }
        if (elem instanceof CtLiteral             lt) { this.addLiteralUsage(lt);             return; }
        if (elem instanceof CtThisAccess          ta) { this.addThisAccess(ta);               return; }
        if (elem instanceof CtVariableRead        vr) { this.addVariableRead(vr);             return; }
        if (elem instanceof CtParameterReference  pr) { this.addParameterReference(pr);       return; }
        if (elem instanceof CtConstructorCall     cc) { this.addConstructorCall(cc);          return; }

        // Use to see elements (may produce a lot of output).
        this.log.notice("unimplemented addUsage: " + SpoonUtils.describeElem(elem));
    }

    private void addInvocationUsage(CtInvocation<?> in) throws Exception {
        if (logUsage) this.log.log("addUsage.CtInvocation: " + SpoonUtils.describeElem(in));
        this.addExecutableReferenceUsage(in.getExecutable());
    }

    private void addFieldReadUsage(CtFieldRead<?> fr) throws Exception {
        if (logUsage) this.log.log("addUsage.CtFieldRead: " + SpoonUtils.describeElem(fr));
        final CtFieldReference<?> ref = fr.getVariable();
        final CtField<?> field = ref.getFieldDeclaration();
        if (field == null) {
            this.log.warning("addUsage.CtFieldRead: no field declaration for " + SpoonUtils.describeElem(ref));
            return;
        }
        final Ref<Selection> sel = this.abs.addSelection(field);
        this.reads.add(sel);
    }

    private void addTypeAccessUsage(CtTypeAccess<?> ta) throws Exception {
        if (logUsage) this.log.log("addUsage.CtTypeAccess: " + SpoonUtils.describeElem(ta));
        final Ref<? extends TypeDesc> acc = this.abs.addTypeDesc(ta.getAccessedType());
        this.reads.add(acc);
    }

    private void addAssignmentUsage(CtAssignment<?,?> as) throws Exception {
        if (logUsage) this.log.log("addUsage.CtAssignment: " + SpoonUtils.describeElem(as));

        // TODO: Implement
        this.log.warning("|  unimplemented CtAssignment: " + SpoonUtils.describeElem(as));
        this.log.notice("|  getType:     " + SpoonUtils.describeElem(as.getType()));
        this.log.notice("|  getAssigned: " + SpoonUtils.describeElem(as.getAssigned()));
    }

    private void addTypeReferenceUsage(CtTypeReference<?> tr) throws Exception {
        if (logUsage) this.log.log("addUsage.CtTypeReference: " + SpoonUtils.describeElem(tr));

        // TODO: Should we create a reference?
        List<CtTypeReference<?>> typeArgs = tr.getActualTypeArguments();
        if (typeArgs.size() > 0) this.log.notice("addUsage.CtTypeReference: unused type args: " + SpoonUtils.describeElems(typeArgs));

        if (!SpoonUtils.isVoid(tr)) this.abs.addTypeDesc(tr);
    }

    private void addPackageReferenceUsage(CtPackageReference pr) throws Exception {
        if (logUsage) this.log.log("addUsage.CtPackageReference: " + SpoonUtils.describeElem(pr));
        this.abs.addPackage(pr.getDeclaration());
    }
    
    private void addFieldReferenceUsage(CtFieldReference<?> fr) throws Exception {
        if (logUsage) this.log.log("addUsage.CtFieldReference: " + SpoonUtils.describeElem(fr));
        final CtField<?> field = fr.getFieldDeclaration();
        if (field == null) {
            this.log.warning("addUsage.CtFieldReference: no field declaration for " + SpoonUtils.describeElem(fr));
            return;
        }
        final Ref<Selection> sel = this.abs.addSelection(field);
        this.reads.add(sel);
    }
    
    private void addExecutableReferenceUsage(CtExecutableReference<?> er) throws Exception {
        if (logUsage) this.log.log("addUsage.CtExecutableReference: " + SpoonUtils.describeElem(er));
        
        // TODO: Should we create a reference?
        List<CtTypeReference<?>> typeArgs = er.getActualTypeArguments();
        if (typeArgs.size() > 0) this.log.notice("addUsage.CtExecutableReference: unused type args: " + SpoonUtils.describeElems(typeArgs));

        final CtExecutable<?> ex = er.getDeclaration();
        if (ex instanceof CtMethod      mt) { this.addMethodUsage(mt);      return; }
        if (ex instanceof CtConstructor ct) { this.addConstructorUsage(ct); return; }
        
        this.log.warning("addUsage.CtExecutableReference: expected method for " + SpoonUtils.describeElem(er) + " and "+ SpoonUtils.describeElem(ex));
    }

    private void addMethodUsage(CtMethod<?> mt) throws Exception {
        if (logUsage) this.log.log("addUsage.CtMethod: " + SpoonUtils.describeElem(mt));
        final Ref<? extends Construct> decl = this.abs.addDeclaration(mt);
        if (decl != null) this.invokes.add(decl);
    }

    private void addConstructorUsage(CtConstructor<?> ct) throws Exception {
        if (logUsage) this.log.log("addUsage.CtConstructor: " + SpoonUtils.describeElem(ct));
        final Ref<? extends Construct> decl = this.abs.addDeclaration(ct);
        if (decl != null) this.invokes.add(decl);
    }

    private void addLiteralUsage(CtLiteral<?> lt) throws Exception {
        if (logUsage) this.log.log("addUsage.CtLiteral: " + SpoonUtils.describeElem(lt));
        this.reads.add(this.abs.addTypeDesc(lt.getType()));
    }

    private void addThisAccess(CtThisAccess<?> ta) throws Exception {
        if (logUsage) this.log.log("addUsage.CtThisAccess: " + SpoonUtils.describeElem(ta));
        // Ignore since this will continue onto the type for "this".
    }

    private void addVariableRead(CtVariableRead<?> vr) throws Exception {
        if (logUsage) this.log.log("addUsage.CtVariableRead: " + SpoonUtils.describeElem(vr));

        // TODO: Implement

        this.log.warning("unimplemented addUsage for CtVariableRead: " + SpoonUtils.describeElem(vr));
    }

    private void addParameterReference(CtParameterReference<?> pr) throws Exception {
        if (logUsage) this.log.log("addUsage.CtParameterReference: " + SpoonUtils.describeElem(pr));

        // TODO: Implement

        this.log.warning("unimplemented addUsage for CtParameterReference: " + SpoonUtils.describeElem(pr));
    }

    private void addConstructorCall(CtConstructorCall<?> cc) {
        if (logUsage) this.log.log("addUsage.CtConstructorCall: " + SpoonUtils.describeElem(cc));

        // TODO: Implement

        this.log.warning("unimplemented addUsage for CtConstructorCall: " + SpoonUtils.describeElem(cc));
    }
}
