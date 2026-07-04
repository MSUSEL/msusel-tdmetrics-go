package abstractor.core.require;

import java.util.Objects;
import java.util.regex.Pattern;

import abstractor.core.AbstractorException;
import abstractor.core.spoonUtils.SpoonUtils;

import spoon.reflect.declaration.CtMethod;
import spoon.reflect.reference.CtTypeReference;

public class Require {
    private Require() {}
    
    static public boolean enabled = true;

    static public void fail() throws Exception { fail(""); }
    static public void fail(String msg) throws Exception {
        if (enabled) throw new AbstractorException(!msg.isBlank()? msg: "failure: no message");
    }

    static public void require(boolean test) throws Exception { require(test, ""); }
    static public void require(boolean test, String msg) throws Exception {
        if (!test) fail(!msg.isBlank()? msg: "required test to be true");
    }

    static public <T> void equal(T value, T other) throws Exception { equal(value, other, ""); }
    static public <T> void equal(T value, T other, String msg) throws Exception {
        require(Objects.equals(value, other), !msg.isBlank()? msg: "required " + value + " to be equal to " + other);
    }

    static public <T> void notEqual(T value, T other) throws Exception { notEqual(value, other, ""); }
    static public <T> void notEqual(T value, T other, String msg) throws Exception {
        require(!Objects.equals(value, other), !msg.isBlank()? msg: "required " + value + " to be not equal to " + other);
    }

    static public <T extends Comparable<T>> void lessThan(T value, T other) throws Exception { lessThan(value, other, ""); }
    static public <T extends Comparable<T>> void lessThan(T value, T other, String msg) throws Exception {
        require(value.compareTo(other) < 0, !msg.isBlank()? msg: "required " + value + " to be less than " + other);
    }

    static public <T extends Comparable<T>> void lessThanOrEqual(T value, T other) throws Exception { lessThanOrEqual(value, other, ""); }
    static public <T extends Comparable<T>> void lessThanOrEqual(T value, T other, String msg) throws Exception {
        require(value.compareTo(other) <= 0, !msg.isBlank()? msg: "required " + value + " to be less or equal than " + other);
    }

    static public <T extends Comparable<T>> void greaterThan(T value, T other) throws Exception { greaterThan(value, other, ""); }
    static public <T extends Comparable<T>> void greaterThan(T value, T other, String msg) throws Exception {
        require(value.compareTo(other) > 0, !msg.isBlank()? msg: "required " + value + " to be greater than " + other);
    }

    static public <T extends Comparable<T>> void greaterThanOrEqual(T value, T other) throws Exception { greaterThanOrEqual(value, other, ""); }
    static public <T extends Comparable<T>> void greaterThanOrEqual(T value, T other, String msg) throws Exception {
        require(value.compareTo(other) >= 0, !msg.isBlank()? msg: "required " + value + " to be greater or equal than " + other);
    }

    static public <T> void notNull(T value) throws Exception { notNull(value, ""); }
    static public <T> void notNull(T value, String msg) throws Exception {
        require(value != null, !msg.isBlank()? msg: "required a not null value");
    }

    static public <T> void isNull(T value) throws Exception { isNull(value, ""); }
    static public <T> void isNull(T value, String msg) throws Exception {
        require(value == null, !msg.isBlank()? msg: "required a null value");
    }

    static public void notBlank(String text) throws Exception { notBlank(text, ""); }
    static public void notBlank(String text, String msg) throws Exception {
        require(!text.isBlank(), !msg.isBlank()? msg: "required a string to not be blank");
    }

    static public void notObjectMethod(CtMethod<?> m) throws Exception { notObjectMethod(m, ""); }
    static public void notObjectMethod(CtMethod<?> m, String msg) throws Exception {
        require(!SpoonUtils.isObjectMethod(m), !msg.isBlank()? msg: "required a method not to be a default Object method: " + SpoonUtils.describeElem(m));
    }

    static public void notObject(CtTypeReference<?> tr) throws Exception { notObject(tr, ""); }
    static public void notObject(CtTypeReference<?> tr, String msg) throws Exception {
        require(!SpoonUtils.isObject(tr), !msg.isBlank()? msg: "required a type to not be the default Object: " + SpoonUtils.describeElem(tr));
    }

    static private Pattern idPattern;
    static private boolean idMatch(String s) {
        if (idPattern == null) idPattern = Pattern.compile("^[a-zA-Z_$][a-zA-Z0-9_$]*$");
        return idPattern.matcher(s).find();
    }
    
    static public void isIdentifier(String s) throws Exception { isIdentifier(s, ""); }
    static public void isIdentifier(String s, String msg) throws Exception {
        require(idMatch(s), !msg.isBlank()? msg: "required an identifier: " + s);
    }
}
