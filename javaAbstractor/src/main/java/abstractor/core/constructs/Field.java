package abstractor.core.constructs;

import spoon.reflect.declaration.CtField;

public class Field extends Construct {
    private CtField<?> field;

    static public Field Create(Project proj, CtField<?> src) {
        Field existing = proj.fields.findWithSource(src);
        if (existing != null) return existing;

        Field f = new Field(src);
        existing = proj.fields.tryAdd(f);
        if (existing != null) return existing;
        
        return f;
    }

    private Field(CtField<?> field) {
        this.field = field;
    }

    public Object source() { return this.field; }

    public String kind() { return "field"; }

    // TODO: | `name`     | ◯ | ◯ | The string name for the field. |
    // TODO: | `type`     | ◯ | ◯ | [Key](#keys) for any [type description](#type-descriptions). |
    
}
