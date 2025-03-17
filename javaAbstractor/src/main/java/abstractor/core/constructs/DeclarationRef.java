package abstractor.core.constructs;

import java.util.List;

import spoon.reflect.declaration.CtElement;

public class DeclarationRef extends Reference<Declaration> {
    
    public DeclarationRef(CtElement elem, String context, String name, List<TypeDesc> typeParams) {
        super(elem, context, name, typeParams);
    }

    public ConstructKind unresolvedKind() { return ConstructKind.DECLARATION_REF; }
}
