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

    // Assign one point to account for the start of the method
    // Add one point for each conditional construct, such as an "if" condition
    // Add one point for each iterative structure
    // Add one point for each case or default block in a switch statement
    // Add one point for any additional boolean condition, such as the use of && or ||

    static public JsonObj calculate(CtMethod<?> m) {
        int cc = 1;
        for (CtStatement statement : m.getBody())
            cc += calculate(statement);
        return new JsonInt(cc);
    }

    static private int calculate(CtStatement s) {
        if (s instanceof CtInvocationImpl<?> st) return calculate(st.getArguments());
        else if (s instanceof CtLocalVariableImpl<?> st) return calculate(st.getAssignment());
        else if (s instanceof CtReturnImpl<?> st) return calculate(st.getReturnedExpression());
        else {
            System.out.println("Unsupported statement:");
            System.out.println("   " + s);
            System.out.println("   " + s.getClass().getName());
            return 0;
        }
    }

    static private int calculate(List<CtExpression<?>> le) {
        int cc = 0;
        for (CtExpression<?> e : le) cc += calculate(e);
        return cc;
    }

    static private int calculate(CtExpression<?> e) {
        if (e instanceof CtConstructorCallImpl<?> et) return calculate(et.getArguments());
        else if (e instanceof CtFieldReadImpl<?>) return 0;
        else if (e instanceof CtLiteralImpl<?>) return 0;
        else {
            System.out.println("Unsupported expression:");
            System.out.println("   " + e);
            System.out.println("   " + e.getClass().getName());
            return 0;
        }
    }
}
