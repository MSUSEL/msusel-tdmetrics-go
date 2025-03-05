package abstractor.core;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.TreeSet;
import java.util.TreeMap;

import spoon.reflect.code.CtBlock;
import spoon.reflect.code.CtStatement;
import spoon.reflect.cu.SourcePosition;
import spoon.reflect.declaration.CtElement;

import abstractor.core.constructs.Location;
import abstractor.core.constructs.Method;
import abstractor.core.constructs.Metrics;
import abstractor.core.constructs.TypeDesc;
import abstractor.core.log.Logger;

public class Analyzer {
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
        this.log       = log;
        this.loc       = loc;
        this.minColumn = new TreeMap<Integer, Integer>();
        this.invokes   = new ArrayList<Method>();
        this.reads     = new ArrayList<TypeDesc>();
        this.writes    = new ArrayList<TypeDesc>();
    }

    public Metrics getMetrics() {
        final int lineCount = this.maxLine - this.minLine + 1;
        final int codeCount = this.minColumn.size();
        final int indents   = this.calcIndent();

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

	    int indentSum = 0;
        for (int minCol : this.minColumn.values())
            indentSum += indentMap.get(minCol);
        return indentSum;
    }
    
    public void addBlock(CtBlock<?> block) {
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
        this.addPosition(elem.getPosition());
        //if (elem instanceof )

        // CtAssertImpl
        // CtAssignmentImpl
        // CtBlockImpl
        // CtBreakImpl
        // CtCaseImpl
        // CtClassImpl
        // CtCodeSnippetStatementImpl
        // CtCommentImpl
        // CtConstructorCallImpl
        // CtContinueImpl
        // CtDoImpl
        // CtEnumImpl
        // CtForEachImpl
        // CtForImpl
        // CtIfImpl
        // CtInterfaceImpl
        // CtInvocationImpl
        // CtJavaDocImpl
        // CtLocalVariableImpl
        // CtLoopImpl
        // CtNewClassImpl
        // CtOperatorAssignmentImpl
        // CtRecordImpl
        // CtReturnImpl
        // CtStatementImpl
        // CtSwitchImpl
        // CtSynchronizedImpl
        // CtThrowImpl
        // CtTryImpl
        // CtTryWithResourceImpl
        // CtUnaryOperatorImpl
        // CtWhileImpl
        // CtYieldStatementImpl
        
        this.log.log("+- (" + elem.getClass().getSimpleName() + ") " + elem);
        this.log.push("|  ");
        for (CtElement child : elem.getDirectChildren())
            this.addElement(child);
        this.log.pop();
    }

    private void addPosition(SourcePosition pos) {
        if (!pos.isValidPosition()) return;
        this.addPosition(pos.getLine(), pos.getColumn());   
        // TODO: check if the end needs to subtract one.
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

    }
    
    private void detectSetter(CtStatement st) {
        
    }
}
