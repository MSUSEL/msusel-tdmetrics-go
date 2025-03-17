package abstractor.core.constructs;

import java.util.List;

import spoon.reflect.declaration.CtElement;

public class TypeParamRef extends Reference<TypeParam> implements TypeDesc {
    
    public TypeParamRef(CtElement elem, String context, String name, List<TypeDesc> typeParams) {
        super(elem, context, name, typeParams);
    }

    public ConstructKind unresolvedKind() { return ConstructKind.TYPE_PARAM_REF; }
}
