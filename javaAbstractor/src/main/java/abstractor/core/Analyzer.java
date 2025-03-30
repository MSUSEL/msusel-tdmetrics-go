package abstractor.core;

import java.util.List;
import java.util.Map;
import java.util.SortedSet;
import java.util.TreeSet;
import java.util.TreeMap;

import spoon.reflect.code.BinaryOperatorKind;
import spoon.reflect.code.CtAssignment;
import spoon.reflect.code.CtBinaryOperator;
import spoon.reflect.code.CtBlock;
import spoon.reflect.code.CtCase;
import spoon.reflect.code.CtComment;
import spoon.reflect.code.CtConstructorCall;
import spoon.reflect.code.CtFieldRead;
import spoon.reflect.code.CtIf;
import spoon.reflect.code.CtInvocation;
import spoon.reflect.code.CtLocalVariable;
import spoon.reflect.code.CtLoop;
import spoon.reflect.code.CtReturn;
import spoon.reflect.code.CtStatement;
import spoon.reflect.code.CtTypeAccess;
import spoon.reflect.code.CtUnaryOperator;
import spoon.reflect.cu.SourcePosition;
import spoon.reflect.declaration.CtElement;
import spoon.reflect.declaration.CtExecutable;
import spoon.reflect.declaration.CtMethod;
import spoon.reflect.declaration.CtParameter;
import spoon.reflect.reference.CtParameterReference;
import spoon.reflect.reference.CtTypeReference;
import spoon.support.reflect.CtExtendedModifier;

import abstractor.core.constructs.*;
import abstractor.core.log.Logger;

public class Analyzer {

    private static final boolean logElementTree = false;

    private final Abstractor abs;
    private final Logger log;
    public final Location loc;

    private final Map<Integer, Integer> minColumn;
    private int minLine;
    private int maxLine;
    
    private int complexity;

    private boolean getter;
    private boolean setter;

    private final SortedSet<Construct> invokes;
    private final SortedSet<Construct> reads;
    private final SortedSet<Construct> writes;

    public Analyzer(Abstractor abs, Location loc) {
        this.abs        = abs;
        this.log        = abs.log;
        this.loc        = loc;
        this.minLine    = Integer.MAX_VALUE;
        this.minColumn  = new TreeMap<Integer, Integer>();
        this.invokes    = new TreeSet<Construct>();
        this.reads      = new TreeSet<Construct>();
        this.writes     = new TreeSet<Construct>();
        this.complexity = 1;
    }

    public Metrics getMetrics() {
        final int lineCount = this.maxLine - this.minLine + 1;
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
        final TreeSet<Integer> columns = new TreeSet<Integer>();
        columns.addAll(this.minColumn.values());
        
        final TreeMap<Integer, Integer> indentMap = new TreeMap<Integer, Integer>();
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

    public void addMethod(CtMethod<?> m) throws Exception {
        this.addPosition(m.getPosition());

        for (CtExtendedModifier em : m.getExtendedModifiers())
            this.addPosition(em.getPosition());

        for (CtParameter<?> p : m.getParameters())
            this.addPosition(p.getPosition());

        final CtBlock<?> block = m.getBody();
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

    static private String formatElem(CtElement elem) {
        String head = elem.toStringDebug().replaceAll("\\s+", " ");
        if (head.length() > 40) head = head.substring(0, 36) + "...";
        final SourcePosition pos = elem.getPosition();
        final String tail = pos.isValidPosition() ?
            " @ "+pos.getLine() + ":" + pos.getColumn() : "";
        final String type = elem.getClass().getSimpleName();
        return "(" + type + ") " + head + tail;
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
            this.log.log("+- " + formatElem(elem));
            this.log.push("|  ");
        }

        for (CtElement child : elem.getDirectChildren())
            this.addElement(child);

        if (logElementTree) this.log.pop();
    }

    private void addPosition(SourcePosition pos) {
        if (!pos.isValidPosition()) return;
        //this.log.log("  adding <"+pos.getLine()+", "+pos.getColumn()+
        //    "> <"+pos.getEndLine()+", "+pos.getEndColumn()+">");
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

    static private boolean isVoid(CtTypeReference<?> tr) {
        return tr.isPrimitive() && tr.getSimpleName().equals("void");
    }

    static private boolean isSimpleFetch(CtElement elem) {
        if (elem instanceof CtConstructorCall) return false;
        if (elem instanceof CtInvocation) return false;
        if (elem instanceof CtBinaryOperator) return false;
        if (elem instanceof CtUnaryOperator) return false;
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

    static private boolean detectGetter(CtMethod<?> m, CtStatement st) {
        if (m.getParameters().size() != 0) return false;
        if (isVoid(m.getType())) return false;
        if (!(st instanceof CtReturn ret)) return false;
        if (!isSimpleFetch(ret)) return false;
        return true;
    }
    
    static private boolean detectSetter(CtMethod<?> m, CtStatement st) {
        if (m.getParameters().size() > 1) return false;
        if (!isVoid(m.getType())) return false;
        if (!(st instanceof CtAssignment assign)) return false;
        if (!isSimpleFetch(assign.getAssigned())) return false;
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
            if (opElem.getKind() == BinaryOperatorKind.OR) return 1;
            return 0;
        }

        return 0;
    }

    private void addUsage(CtElement elem) throws Exception {
        // List<Method> invokes;
        // List<TypeDesc> reads;
        // List<TypeDesc> writes;

        if (elem instanceof CtInvocation inv) {
            final CtExecutable<?> ex = inv.getExecutable().getDeclaration();
            if (ex instanceof CtMethod<?> method) {
                DeclarationRef ref = this.abs.create(this.abs.proj.declRefs, method,
                    "add invocation " + method.getSimpleName(),
                    ()-> {
                        final String pkgPath = method.getClass().getName();
                        final String name    = method.getSimpleName();
                        final List<TypeDesc> tp = null; // TODO: Implement
                        return new DeclarationRef(method, pkgPath, name, tp);
                    });
                this.invokes.add(ref);
            }
            return;
        }
        if (elem instanceof CtFieldRead fr) {

            // TODO: Implement

            return;
        }
        if (elem instanceof CtTypeAccess ta) {
            
            // TODO: Implement

            return;
        }
        if (elem instanceof CtAssignment as) {

            // TODO: Implement

            return;
        }

        // TODO: Use to see elements.
        //this.log.log("addUsage: "+formatElem(elem));
    }
}
