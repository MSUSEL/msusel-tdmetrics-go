package abstractor.core.constructs;

import java.util.List;

import spoon.reflect.declaration.CtTypeParameter;

public class TypeParamRef extends Reference<TypeParam> implements TypeDesc {
    
    public TypeParamRef(CtTypeParameter elem, String context, String name, List<TypeDesc> typeParams) {
        super(elem, context, name, typeParams);
    }

    public ConstructKind kind() { return ConstructKind.TYPE_PARAM_REF; }
}
