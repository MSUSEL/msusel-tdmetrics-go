package abstractor.core.spoonUtils;

import java.io.File;

import spoon.reflect.cu.SourcePosition;
import spoon.reflect.declaration.CtElement;
import spoon.reflect.declaration.CtExecutable;
import spoon.reflect.declaration.CtMethod;
import spoon.reflect.declaration.CtNamedElement;
import spoon.reflect.declaration.CtPackage;
import spoon.reflect.declaration.CtType;
import spoon.reflect.declaration.CtTypeInformation;
import spoon.reflect.reference.CtReference;
import spoon.reflect.reference.CtTypeReference;

final public class SpoonUtils {
    private SpoonUtils() { }

    static public String normalizePath(String path) {
        return path.replaceAll("\\\\", "/");
    }

    static public String packageName(CtPackage pkg) {
        if (pkg == null) return "<java.lang>";
        final String name = pkg.getQualifiedName();
        return name.isBlank() ? "<unnamed>" : name;
    }

    static public String packagePath(CtPackage pkg) {
        if (pkg == null) return "";
        final SourcePosition pos = pkg.getPosition();
        if (!pos.isValidPosition()) return "";
        
        final File file = pos.getFile();
        if (file == null) return "";

        final String path = normalizePath(file.getPath());
        final String tail = "/package-info.java";
        if (!path.endsWith(tail)) return path;
        return path.substring(0, path.length()-tail.length());
    }

    /** Short description for logs (does not throw). */
    static public String describeElem(CtElement elem) {
        if (elem == null) return "(null)";
        if (elem instanceof CtPackage pkg)
           return packageName(pkg);
        if (elem instanceof CtNamedElement ne) {
            try { return ne.getSimpleName(); }
            catch (Exception ignore) { }
        }
        if (elem instanceof CtTypeInformation ti) {
            try { return ti.getQualifiedName(); }
            catch (Exception ignore) { }
        }
        if (elem instanceof CtExecutable<?> ex) {
            try { return ex.getSignature(); }
            catch (Exception ignore) { }
        }
        if (elem instanceof CtReference ref) {
            try { return ref.getSimpleName(); }
            catch (Exception ignore) { }
        }
        return elem.getClass().getName();
    }

    static public boolean isVoid(CtTypeReference<?> tr) {
        return tr.isPrimitive() && tr.getSimpleName().equals("void");
    }

    static public boolean isNull(CtTypeReference<?> tr) {
        return CtTypeReference.NULL_TYPE_NAME.equals(tr.getSimpleName());
    }

    static public boolean isObject(CtTypeReference<?> tr) {
        final String objName = "java.lang.Object";
        return objName.equals(tr.getQualifiedName());
    }

    /**
     * This determines if the given method is a method on the base Object.
     * Since all Objects inherits the base Object, adding those methods are
     * just additional unneeded noise in the abstraction.
     */
    static public boolean isObjectMethod(CtMethod<?> m) {
        if (m == null) return false;

        // TODO: CHECK this
        final CtTypeReference<?> objectRef  = m.getFactory().Type().objectType();
        final CtType<?>          objectDecl = objectRef.getTypeDeclaration();
        if (objectDecl == null) return false;

        final String sig = m.getSignature();
        for (CtMethod<?> objectMethod : objectDecl.getMethods()) {
            if (sig.equals(objectMethod.getSignature())) return true;
        }
        return false;
    }
}
