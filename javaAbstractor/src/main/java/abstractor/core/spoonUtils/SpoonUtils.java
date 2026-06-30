package abstractor.core.spoonUtils;

import java.io.File;
import java.util.*;

import spoon.reflect.cu.SourcePosition;
import spoon.reflect.declaration.CtElement;
import spoon.reflect.declaration.CtExecutable;
import spoon.reflect.declaration.CtMethod;
import spoon.reflect.declaration.CtNamedElement;
import spoon.reflect.declaration.CtPackage;
import spoon.reflect.declaration.CtType;
import spoon.reflect.declaration.CtTypeInformation;
import spoon.reflect.declaration.CtTypeParameter;
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

    /** Short description of an element that can be used for logs. */
    static public String describeElem(CtElement elem) {
        return describeElem(elem, true, true);
    }

    /** Short description of an element that can be used for logs. */
    static public String describeElem(CtElement elem, boolean showType) {
        return describeElem(elem, showType, true);
    }

    /** Short description of an element that can be used for logs. */
    static public String describeElem(CtElement elem, boolean showType, boolean showPos) {
        if (elem == null) return "(null)";
        if (elem instanceof CtPackage pkg) return packageName(pkg);
        
        String header = "";
        if (showType) {
            header = "(" + elem.getClass().getSimpleName() + ") ";
        }

        String tail = "";
        if (showPos) {
            final SourcePosition pos = elem.getPosition();
            if (pos.isValidPosition()) tail = " @ "+pos.getLine() + ":" + pos.getColumn();
        }

        if (elem instanceof CtNamedElement ne) {
            try { return header + ne.getSimpleName() + tail; }
            catch (Exception ignore) { }
        }
        if (elem instanceof CtTypeInformation ti) {
            try { return header + ti.getQualifiedName() + tail; }
            catch (Exception ignore) { }
        }
        if (elem instanceof CtExecutable<?> ex) {
            try { return header + ex.getSignature() + tail; }
            catch (Exception ignore) { }
        }
        if (elem instanceof CtReference ref) {
            try { return header + ref.getSimpleName() + tail; }
            catch (Exception ignore) { }
        }
        return elem.getClass().getName() + tail;
    }

    static public String describeElems(Iterable<? extends CtElement> elems) {
        ArrayList<String> descs = new ArrayList<>();
        for (CtElement elem : elems)
            descs.add(describeElem(elem, false));
        return String.join(", ",  descs);
    }

    /**
     * Returns a CtTypeReference for the given type populated with the formal
     * type parameters of the type and its declaring-type chain. Spoon's
     * {@code CtType.getReference()} returns a raw reference (empty
     * actualTypeArguments), which is unsuitable for callers that need to walk
     * the type-parameter chain (e.g. constructor result types).
     */
    static public CtTypeReference<?> parameterizedRef(CtType<?> type) {
        if (type == null) return null;
        final CtTypeReference<?> ref = type.getReference();
        for (CtTypeParameter tp : type.getFormalCtTypeParameters())
            ref.addActualTypeArgument(tp.getReference());
        final CtType<?> declaring = type.getDeclaringType();
        if (declaring != null) ref.setDeclaringType(parameterizedRef(declaring));
        return ref;
    }

    static public String describeGeneric(CtTypeReference<?> tr) {
        final List<CtTypeReference<?>> typeArgs = tr.getActualTypeArguments();
        String tail = "";
        if (typeArgs.size() > 0) {
            tail = "<"+SpoonUtils.describeElems(typeArgs)+">";
        }
        return SpoonUtils.describeElem(tr, false) + tail;
    }

    static public boolean isVoid(CtTypeReference<?> tr) {
        return tr.isPrimitive() && tr.getSimpleName().equals("void");
    }

    static public boolean isNull(CtTypeReference<?> tr) {
        return CtTypeReference.NULL_TYPE_NAME.equals(tr.getSimpleName());
    }

    static public boolean isObject(CtTypeInformation ti) {
        final String objName = "java.lang.Object";
        return objName.equals(ti.getQualifiedName());
    }

    /**
     * This determines if the given method is a method on the base Object.
     * Since all Objects inherits the base Object, adding those methods are
     * just additional unneeded noise in the abstraction.
     */
    static public boolean isObjectMethod(CtMethod<?> m) {
        if (m == null) return false;

        final CtTypeReference<?> objectRef = m.getFactory().Type().objectType();
        assert(isObject(objectRef));

        final CtType<?> objectDecl = objectRef.getTypeDeclaration();
        assert(objectDecl != null);

        final String sig = m.getSignature();
        assert(sig != null);
        for (CtMethod<?> objectMethod : objectDecl.getMethods()) {
            if (sig.equals(objectMethod.getSignature())) return true;
        }
        return false;
    }
}
