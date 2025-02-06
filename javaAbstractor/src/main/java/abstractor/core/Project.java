package abstractor.core;

import java.util.*;
import spoon.reflect.declaration.*;

public class Project {

    public final TreeSet<CtPackage> packages = new TreeSet<CtPackage>(
        Comparator.comparing((CtPackage p) -> p.getQualifiedName()));
    
    public final TreeSet<CtClass<?>> objects = new TreeSet<CtClass<?>>(
        Comparator.comparing((CtClass<?> c) -> c.getQualifiedName()));
        
    public final TreeSet<CtInterface<?>> interfaceDecls = new TreeSet<CtInterface<?>>(
        Comparator.comparing((CtInterface<?> i) -> i.getQualifiedName()));
}
