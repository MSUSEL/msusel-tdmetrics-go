package abstractor.core.constructs;

import java.util.List;

import spoon.reflect.declaration.CtElement;

public class TypeDescRef extends Reference<TypeDesc> implements TypeDeclaration {
    
    public TypeDescRef(CtElement elem, String context, String name, List<TypeDesc> typeParams) {
        super(elem, context, name, typeParams);
    }

    public ConstructKind kind() { return ConstructKind.TYPE_DESC_REF; }
}
