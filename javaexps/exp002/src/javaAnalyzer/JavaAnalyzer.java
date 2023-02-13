package javaAnalyzer;

import java.io.File;
import java.io.IOException;
import java.nio.file.Path;
import java.util.HashSet;
import java.util.Iterator;

public class JavaAnalyzer {
    private final HashSet<JavaProject> projects;
    private final JavaPath root;

    public JavaAnalyzer() {
        this.projects = new HashSet<>();
        this.root = new JavaPath("");
    }

    public void addPath(String path) throws IOException {
        final File f = new File(path);
        JavaPath head = this.getOrAddPath(f);




    }

    private JavaPath getOrAddPath(final File f) {
        JavaPath parent = this.root;
        Iterator<Path> it = f.toPath().iterator();
        while (it.hasNext()) {
            final String part = it.next().toString();
            JavaPath child = parent.children().get(part);
            if (child == null) {
                child = new JavaPath(part);
                parent.children().add(child);
            }
            parent = child;
        }
        return parent;
    }
}
