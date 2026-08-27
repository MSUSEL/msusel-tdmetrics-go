package abstractor.core.require;

import java.util.Objects;
import java.util.function.Supplier;
import java.util.regex.Pattern;

import abstractor.core.AbstractorException;
import abstractor.core.spoonUtils.SpoonUtils;

import spoon.reflect.declaration.CtMethod;
import spoon.reflect.reference.CtTypeReference;

public class Require {
    private Require() {}
    
    static public boolean enabled = true;

    static private void require(Supplier<Boolean> test, Supplier<String> msg, Supplier<String> defaultMsg) throws Exception {
        if (enabled && !test.get()) {
            String text = msg.get();
            if (text.isBlank()) text = defaultMsg.get();
            throw new AbstractorException(text);
        }
    }

    static public void fail() throws Exception { fail(""); }
    static public void fail(String msg) throws Exception { fail(() -> msg); }
    static public void fail(Supplier<String> msg) throws Exception {
        require(() -> false, msg, () -> "failure: no message");
    }

    static public void require(boolean test) throws Exception { require(test, ""); }
    static public void require(boolean test, String msg) throws Exception { require(() -> test, () -> msg); }
    static public void require(boolean test,  Supplier<String> msg) throws Exception { require(() -> test, msg); }
    static public void require(Supplier<Boolean> test, String msg) throws Exception { require(test, () -> msg); }
    static public void require(Supplier<Boolean> test, Supplier<String> msg) throws Exception {
        require(test, msg, () -> "required test to be true");
    }

    static public <T> void equal(T value, T other) throws Exception { equal(value, other, ""); }
    static public <T> void equal(T value, T other, String msg) throws Exception { equal(value, other, () -> msg); }
    static public <T> void equal(T value, T other, Supplier<String> msg) throws Exception {
        require(() -> Objects.equals(value, other), msg, () -> "required " + value + " to be equal to " + other);
    }

    static public <T> void notEqual(T value, T other) throws Exception { notEqual(value, other, ""); }
    static public <T> void notEqual(T value, T other, String msg) throws Exception { notEqual(value, other, () -> msg); }
    static public <T> void notEqual(T value, T other, Supplier<String> msg) throws Exception {
        require(() -> !Objects.equals(value, other), msg, () -> "required " + value + " to be not equal to " + other);
    }

    static public <T extends Comparable<T>> void lessThan(T value, T other) throws Exception { lessThan(value, other, ""); }
    static public <T extends Comparable<T>> void lessThan(T value, T other, String msg) throws Exception { lessThan(value, other, () -> msg); }
    static public <T extends Comparable<T>> void lessThan(T value, T other, Supplier<String> msg) throws Exception {
        require(() -> value.compareTo(other) < 0, msg, () -> "required " + value + " to be less than " + other);
    }

    static public <T extends Comparable<T>> void lessThanOrEqual(T value, T other) throws Exception { lessThanOrEqual(value, other, ""); }
    static public <T extends Comparable<T>> void lessThanOrEqual(T value, T other, String msg) throws Exception { lessThanOrEqual(value, other, () -> msg); }
    static public <T extends Comparable<T>> void lessThanOrEqual(T value, T other, Supplier<String> msg) throws Exception {
        require(() -> value.compareTo(other) <= 0, msg, () -> "required " + value + " to be less or equal than " + other);
    }

    static public <T extends Comparable<T>> void greaterThan(T value, T other) throws Exception { greaterThan(value, other, ""); }
    static public <T extends Comparable<T>> void greaterThan(T value, T other, String msg) throws Exception { greaterThan(value, other, () -> msg); }
    static public <T extends Comparable<T>> void greaterThan(T value, T other, Supplier<String> msg) throws Exception {
        require(() -> value.compareTo(other) > 0, msg, () -> "required " + value + " to be greater than " + other);
    }

    static public <T extends Comparable<T>> void greaterThanOrEqual(T value, T other) throws Exception { greaterThanOrEqual(value, other, ""); }
    static public <T extends Comparable<T>> void greaterThanOrEqual(T value, T other, String msg) throws Exception { greaterThanOrEqual(value, other, () -> msg); }
    static public <T extends Comparable<T>> void greaterThanOrEqual(T value, T other, Supplier<String> msg) throws Exception {
        require(() -> value.compareTo(other) >= 0, msg, () -> "required " + value + " to be greater or equal than " + other);
    }

    static public <T> void notNull(T value) throws Exception { notNull(value, ""); }
    static public <T> void notNull(T value, String msg) throws Exception { notNull(value, () -> msg); }
    static public <T> void notNull(T value, Supplier<String> msg) throws Exception {
        require(() -> value != null, msg, () -> "required a not null value");
    }

    static public <T> void isNull(T value) throws Exception { isNull(value, ""); }
    static public <T> void isNull(T value, String msg) throws Exception { isNull(value, () -> msg); }
    static public <T> void isNull(T value, Supplier<String> msg) throws Exception {
        require(() -> value == null, msg, () -> "required a null value");
    }

    static public void notBlank(String text) throws Exception { notBlank(text, ""); }
    static public void notBlank(String text, String msg) throws Exception { notBlank(text, () -> msg); }
    static public void notBlank(String text, Supplier<String> msg) throws Exception {
        require(() -> !text.isBlank(), msg, () -> "required a string to not be blank");
    }

    static public void notObjectMethod(CtMethod<?> m) throws Exception { notObjectMethod(m, ""); }
    static public void notObjectMethod(CtMethod<?> m, String msg) throws Exception { notObjectMethod(m, () -> msg); }
    static public void notObjectMethod(CtMethod<?> m, Supplier<String> msg) throws Exception {
        require(() -> !SpoonUtils.isObjectMethod(m), msg, () -> "required a method not to be a default Object method: " + SpoonUtils.describeElem(m));
    }

    static public void notObject(CtTypeReference<?> tr) throws Exception { notObject(tr, ""); }
    static public void notObject(CtTypeReference<?> tr, String msg) throws Exception { notObject(tr, () -> msg); }
    static public void notObject(CtTypeReference<?> tr, Supplier<String> msg) throws Exception {
        require(() -> !SpoonUtils.isObject(tr), msg, () -> "required a type to not be the default Object: " + SpoonUtils.describeElem(tr));
    }

    static private Pattern idPattern;
    static private boolean idMatch(String s) {
        if (idPattern == null) idPattern = Pattern.compile("^[a-zA-Z_$][a-zA-Z0-9_$]*$");
        return idPattern.matcher(s).find();
    }
    
    static public void isIdentifier(String s) throws Exception { isIdentifier(s, ""); }
    static public void isIdentifier(String s, String msg) throws Exception { isIdentifier(s, () -> msg); }
    static public void isIdentifier(String s, Supplier<String> msg) throws Exception {
        require(() -> idMatch(s), msg, () -> "required an identifier: " + s);
    }
}
