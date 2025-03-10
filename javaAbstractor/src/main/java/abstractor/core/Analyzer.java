package abstractor.core;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.TreeSet;
import java.util.TreeMap;

import spoon.reflect.code.BinaryOperatorKind;
import spoon.reflect.code.CtBinaryOperator;
import spoon.reflect.code.CtBlock;
import spoon.reflect.code.CtCase;
import spoon.reflect.code.CtComment;
import spoon.reflect.code.CtIf;
import spoon.reflect.code.CtLoop;
import spoon.reflect.code.CtStatement;
import spoon.reflect.cu.SourcePosition;
import spoon.reflect.declaration.CtElement;
import spoon.reflect.declaration.CtMethod;
import spoon.reflect.declaration.CtParameter;
import spoon.support.reflect.CtExtendedModifier;
import abstractor.core.constructs.Location;
import abstractor.core.constructs.Method;
import abstractor.core.constructs.Metrics;
import abstractor.core.constructs.TypeDesc;
import abstractor.core.log.Logger;

public class Analyzer {

    private static final boolean logElementTree = true;

    private final Logger log;
    public final Location loc;

    private final Map<Integer, Integer> minColumn;
    private int minLine;
    private int maxLine;
    
    private int complexity;

    private boolean getter;
    private boolean setter;

    private final List<Method> invokes;
    private final List<TypeDesc> reads;
    private final List<TypeDesc> writes;

    public Analyzer(Logger log, Location loc) {
        this.log        = log;
        this.loc        = loc;
        this.minLine    = Integer.MAX_VALUE;
        this.minColumn  = new TreeMap<Integer, Integer>();
        this.invokes    = new ArrayList<Method>();
        this.reads      = new ArrayList<TypeDesc>();
        this.writes     = new ArrayList<TypeDesc>();
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

    public void addMethod(CtMethod<?> m) {
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
            this.detectGetter(onlySt);
            this.detectSetter(onlySt);
        }
    }

    private void addElement(CtElement elem) {
        // Skip over any comments in the code.
        if (elem instanceof CtComment) return;

        this.addPosition(elem.getPosition());
        this.complexity += this.addComplexity(elem);
        
        if (logElementTree) {
            this.log.log("+- (" + elem.getClass().getSimpleName() + ") " + elem);
            this.log.push("|  ");
        }

        for (CtElement child : elem.getDirectChildren())
            this.addElement(child);

        if (logElementTree)
            this.log.pop();
    }

    private void addPosition(SourcePosition pos) {
        if (!pos.isValidPosition()) return;

        // TODO: read about pos.getSourceStart()
        
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

    private void detectGetter(CtStatement st) {
        // TODO: Implement
    }
    
    private void detectSetter(CtStatement st) {
        // TODO: Implement
    }

    /**
     * Gets the McCabe cyclomatic complexity for this element.
     */
    private int addComplexity(CtElement elem) {
        // Add one point for each conditional construct, such as an "if" condition
        if (elem instanceof CtIf ifElem) {
            final CtStatement elseElem = ifElem.getElseStatement();
            if (elseElem == null) return 1;

            // Check for "else if(..) { }" and skip adding complexity since the "if" will add it later.
            if (elseElem instanceof CtIf) return 1;

            // Check for "else { if(..) { }}" and skip adding complexity since the "if" will add it later.
            if (elseElem instanceof CtBlock elseBlock) {
                final List<CtStatement> elseChildren = elseBlock.getStatements();
                if (elseChildren.size() == 1 && elseChildren.get(0) instanceof CtIf) return 1;
            }

            // The "else" part isn't just an "if" so add complexity for the "else".
            return 2;
        }

        // Add one point for each iterative structure (CtDo, CtForEach, CtFor, CtWhile)
        if (elem instanceof CtLoop) return 1;

        // Add one point for each case or default block in a switch statement
        if (elem instanceof CtCase) return 1;
        
        // Add one point for any additional boolean condition, such as the use of && or ||
        if (elem instanceof CtBinaryOperator opElem &&
            (opElem.getKind() == BinaryOperatorKind.AND || opElem.getKind() == BinaryOperatorKind.OR))
            return 1;

        return 0;
    }
}
