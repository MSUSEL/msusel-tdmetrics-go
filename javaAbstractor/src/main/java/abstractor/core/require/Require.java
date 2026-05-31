package abstractor.core.require;

import java.util.Objects;

import abstractor.core.AbstractorException;
import abstractor.core.spoonUtils.SpoonUtils;
import spoon.reflect.declaration.CtMethod;
import spoon.reflect.reference.CtTypeReference;

public class Require {
    private Require() {}
    
    static public void fail() throws Exception {
        throw new AbstractorException("failure: no message");
    }

    static public void failure(String msg) throws Exception {
        throw new AbstractorException(!msg.isBlank()? msg: "failure: no message");
    }
    
    static public void failure(Exception ex) throws Exception {
        if (ex == null) return;
        if (ex instanceof AbstractorException ae) throw ae;
        throw new AbstractorException(ex);
    }

    static public void require(boolean test) throws Exception {
        if (!test) throw new AbstractorException("required test to be true");
    }

    static public void require(boolean test, String msg) throws Exception {
        if (!test) throw new AbstractorException(!msg.isBlank()? msg: "required test to be true");
    }

    static public <T> void equal(T value, T other) throws Exception {
        if (!Objects.equals(value, other)) {
            throw new AbstractorException("required " + value + " to be equal to " + other);
        }
    }

    static public <T> void equal(T value, T other , String msg) throws Exception {
        if (!Objects.equals(value, other)) {
            throw new AbstractorException(!msg.isBlank()? msg: "required " + value + " to be equal to " + other);
        }
    }

    static public <T> void notEqual(T value, T other) throws Exception {
        if (Objects.equals(value, other)) {
            throw new AbstractorException("required " + value + " to be not equal to " + other);
        }
    }

    static public <T> void notEqual(T value, T other , String msg) throws Exception {
        if (Objects.equals(value, other)) {
            throw new AbstractorException(!msg.isBlank()? msg: "required " + value + " to be not equal to " + other);
        }
    }

    static public <T> void notNull(T value) throws Exception {
        if (value == null) {
            throw new AbstractorException("required a not null value");
        }
    }

    static public <T> void notNull(T value, String msg) throws Exception {
        if (value == null) {
            throw new AbstractorException(!msg.isBlank()? msg: "required a not null value");
        }
    }

    static public <T> void isNull(T value) throws Exception {
        if (value != null) {
            throw new AbstractorException("required a null value");
        }
    }

    static public <T> void isNull(T value, String msg) throws Exception {
        if (value != null) {
            throw new AbstractorException(!msg.isBlank()? msg: "required a null value");
        }
    }

    static public void notBlank(String text) throws Exception {
        if (text.isBlank()) {
            throw new AbstractorException("required a string to not be blank");
        }
    }

    static public void notBlank(String text, String msg) throws Exception {
        if (text.isBlank()) {
            throw new AbstractorException(!msg.isBlank()? msg: "required a string to not be blank");
        }
    }

    static public void notObjectMethod(CtMethod<?> m) throws Exception {
        if (SpoonUtils.isObjectMethod(m)) {
            throw new AbstractorException("required a method not to be a default Object method: " + SpoonUtils.describeElem(m));
        }
    }

    static public void notObjectMethod(CtMethod<?> m, String msg) throws Exception {
        if (SpoonUtils.isObjectMethod(m)) {
            throw new AbstractorException(!msg.isBlank()? msg: "required a method not to be a default Object method: " + SpoonUtils.describeElem(m));
        }
    }

    static public void notObject(CtTypeReference<?> tr) throws Exception {
        if (SpoonUtils.isObject(tr)) {
            throw new AbstractorException("required a type to not be the default Object: " + SpoonUtils.describeElem(tr));
        }
    }

    static public void notObject(CtTypeReference<?> tr, String msg) throws Exception {
        if (SpoonUtils.isObject(tr)) {
            throw new AbstractorException(!msg.isBlank()? msg: "required a type to not be the default Object: " + SpoonUtils.describeElem(tr));
        }
    }
}
