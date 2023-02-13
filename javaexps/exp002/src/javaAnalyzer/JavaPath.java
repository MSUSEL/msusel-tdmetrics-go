package javaAnalyzer;

import util.NamedObject;
import util.NamedSet;

public class JavaPath implements NamedObject {
    private final String name;
    private final NamedSet<JavaPath> children;

    JavaPath(final String name) {
        this.name = name;
        this.children = new NamedSet<>();
    }

    public NamedSet<JavaPath> children() {
        return this.children;
    }

    public String name() {
        return this.name;
    }
}
