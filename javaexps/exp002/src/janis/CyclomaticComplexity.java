package janis;

import java.util.List;

import json.*;

import spoon.reflect.code.*;
import spoon.reflect.declaration.CtMethod;
import spoon.support.reflect.code.*;

/**
 * This calculates the Cognitive Complexity of a single method.
 * <a href="https://www.baeldung.com/java-cognitive-complexity">Java CC</a>
 */
public class CyclomaticComplexity {
    private CyclomaticComplexity() {}

    // These steps need to be handled.
    // [x] Assign one point to account for the start of the method
    // [x] Add one point for each conditional construct, such as an "if" condition
    // [ ] Add one point for each iterative structure
    // [x] Add one point for each case or default block in a switch statement
    // [x] Add one point for any additional boolean condition, such as the use of && or ||

    static public JsonObj calculate(CtMethod<?> m) {
        int cc = 1;
        for (CtStatement s : m.getBody()) cc += calcState(s);
        return new JsonInt(cc);
    }

    static private int calcListState(List<CtStatement> ls) {
        int cc = 0;
        for (CtStatement s : ls) cc += calcState(s);
        return cc;
    }

    static private int calcState(CtStatement s) {
        if (s == null) return 0;
        if (s instanceof CtAssignmentImpl    st) return calcExp(st.getAssigned()) + calcExp(st.getAssignment());
        if (s instanceof CtBlockImpl         st) return calcListState(st.getStatements());
        if (s instanceof CtCommentImpl         ) return 0;
        if (s instanceof CtConditionalImpl   st) return 1 + calcExp(st.getCondition()) + calcExp(st.getThenExpression()) + calcExp(st.getElseExpression());
        if (s instanceof CtContinueImpl        ) return 0;
        if (s instanceof CtForImpl           st) return 1 + calcListState(st.getForInit()) + calcExp(st.getExpression()) + calcListState(st.getForUpdate()) + calcState(st.getBody());
        if (s instanceof CtForEachImpl       st) return 1 + calcExp(st.getExpression()) + calcState(st.getBody());
        if (s instanceof CtIfImpl            st) return 1 + calcExp(st.getCondition()) + calcState(st.getThenStatement()) + calcState(st.getElseStatement());
        if (s instanceof CtInvocationImpl    st) return calcListExp(st.getArguments());
        if (s instanceof CtLocalVariableImpl st) return calcExp(st.getAssignment());
        if (s instanceof CtReturnImpl        st) return calcExp(st.getReturnedExpression());
        if (s instanceof CtSwitchImpl        st) return calcExp(st.getSelector()) + calcListCase(st.getCases());
        if (s instanceof CtThrowImpl         st) return calcExp(st.getThrownExpression());
        if (s instanceof CtTryImpl           st) return calcListCatcher(st.getCatchers()) + calcState(st.getBody());
        if (s instanceof CtUnaryOperatorImpl st) return calcExp(st.getOperand());
        return unsupported("statement", s);
    }

    static private int calcListExp(List<CtExpression> le) {
        int cc = 0;
        for (CtExpression e : le) cc += calcExp(e);
        return cc;
    }

    static private int calcExp(CtExpression e) {
        if (e == null) return 0;
        if (e instanceof CtBinaryOperatorImpl  et) return calcBinaryOp(et);
        if (e instanceof CtConditionalImpl     et) return 1 + calcExp(et.getCondition()) + calcExp(et.getThenExpression()) + calcExp(et.getElseExpression());
        if (e instanceof CtConstructorCallImpl et) return calcListExp(et.getArguments());
        if (e instanceof CtFieldReadImpl         ) return 0;
        if (e instanceof CtInvocationImpl      et) return calcListExp(et.getArguments());
        if (e instanceof CtLiteralImpl           ) return 0;
        if (e instanceof CtVariableReadImpl      ) return 0;
        if (e instanceof CtVariableWriteImpl     ) return 0;
        return unsupported("expression", e);
    }

    static private int calcBinaryOp(CtBinaryOperatorImpl op) {
        int cc = ((op.getKind() == BinaryOperatorKind.AND) ||
                  (op.getKind() == BinaryOperatorKind.OR)) ? 1 : 0;
        return cc +  calcExp(op.getLeftHandOperand()) + calcExp(op.getLeftHandOperand());
    }

    static private  int calcListCatcher(List<CtCatch> lc) {
        int cc = 0;
        for (CtCatch c : lc) cc += calcCatcher(c);
        return cc;
    }

    static private  int calcCatcher(CtCatch c) {
        return 1 + calcState(c.getBody());
    }

    static private int calcListCase(List<CtCase> lc) {
        int cc = 0;
        for (CtCase c : lc) cc += 1 + calcListExp(c.getCaseExpressions());
        return cc;
    }

    static private int unsupported(String groupName, Object o) {
        System.out.println("CC: Unsupported " + groupName + ": " + o.getClass().getName());
        final String edge = "   +----------------";
        final String indent = "\n   | ";
        System.out.println(edge + indent + o.toString().replace("\n", indent) + "\n" + edge);
        return 0;
    }
}
