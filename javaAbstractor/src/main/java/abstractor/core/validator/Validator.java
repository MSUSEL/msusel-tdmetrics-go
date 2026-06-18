package abstractor.core.validator;

import abstractor.core.log.Logger;

import java.util.ArrayList;

import abstractor.core.constructs.*;

public class Validator {
    final public Logger log;
    final public Project proj;

    public Validator(Logger log, Project proj) {
        this.log = log;
        this.proj = proj;
    }

    public void validate() {
        for (Factory<? extends Construct> f : proj.factories)
            this.validate(f);
    }

    private void validate(Factory<? extends Construct> factory) {
        for (Ref<? extends Construct> ref : factory.getRefSet())
            this.validateRef(factory, ref, "factory ref");
        for (Construct con : factory.getConSet())
            this.validateCon(factory, con);
    }

    private void validateRef(Factory<? extends Construct> factory, Ref<? extends Construct> ref, String usage) {
        if (ref.kind() != factory.kind())
            this.log.error("Expected all references to have the kind of the factory but " + ref +
                " (" + usage + ") was " + ref.kind() + " in " + factory + " with the kind " + factory.kind() + ".");

        if (!this.foundInFactory(ref))
            this.log.error("Expected all references to exist in factory " +
                "but " + ref + " (" + usage + ") was not in factory.");

        if (!ref.isResolved()) {
            this.log.error("Expected all references to be resolved but " + ref + " (" + usage + ") was not resolved.");
            return;
        }

        final Construct con = ref.getResolved();
        if (!this.foundInFactory(con))
            this.log.error("Expected all resolved references to exist in factory " +
                "but " + ref + " (" + usage + ") resolved to " + con + " was not in factory.");
    }

    private boolean foundInFactory(Construct con) {
        final Factory<? extends Construct> factory = this.proj.getFactory(con.kind());
        final Iterable<? extends Construct> set = (con instanceof Ref<?>) ? factory.getRefSet() : factory.getConSet();
        for (Construct other : set) {
            // Use `==` not `equals` to ensure exact reference.
            if (other == con) return true;
        }
        return false;
    }

    private void validateCon(Factory<? extends Construct> factory, Construct con) {
        if (con.kind() != factory.kind())
            this.log.error("Expected all constructs to have the kind of the factory but " + con +
                " was " + con.kind() + " in " + factory + " with the kind " + factory.kind() + ".");

        final int index = con.getIndex();
        final int count = factory.size();
        if (index <= 0 || index > count)
            this.log.error(con + " has invalid index " + index + " should be [1.." + count + "].");

        if (con instanceof Abstract      a) { this.validateAbstract(a);      return; } 
        if (con instanceof Argument      a) { this.validateArgument(a);      return; }
        if (con instanceof Basic         a) { this.validateBasic(a);         return; }
        if (con instanceof Field         a) { this.validateField(a);         return; }
        if (con instanceof InterfaceDecl a) { this.validateInterfaceDecl(a); return; }
        if (con instanceof InterfaceDesc a) { this.validateInterfaceDesc(a); return; }
        if (con instanceof InterfaceInst a) { this.validateInterfaceInst(a); return; }
        if (con instanceof MethodDecl    a) { this.validateMethodDecl(a);    return; }
        if (con instanceof MethodInst    a) { this.validateMethodInst(a);    return; }
        if (con instanceof Metrics       a) { this.validateMetrics(a);       return; }
        if (con instanceof ObjectDecl    a) { this.validateObjectDecl(a);    return; }
        if (con instanceof ObjectInst    a) { this.validateObjectInst(a);    return; }
        if (con instanceof PackageCon    a) { this.validatePackageCon(a);    return; }
        if (con instanceof Selection     a) { this.validateSelection(a);     return; }
        if (con instanceof Signature     a) { this.validateSignature(a);     return; }
        if (con instanceof StructDesc    a) { this.validateStructDesc(a);    return; }
        if (con instanceof TypeParam     a) { this.validateTypeParam(a);     return; }
        if (con instanceof Value         a) { this.validateValue(a);         return; }
        this.log.error(con + " did not have a specific type validation method.");
    }

    private void validateAbstract(Abstract con) {
        this.validateChild(con, con.signature, "signature", false);
    }

    private void validateArgument(Argument con) {
        this.validateChild(con, con.type, "type", false);
    }

    private void validateBasic(Basic con) {
        if (con.name.isBlank())      this.log.error("basic name is black.");
        else if (con.name == "void") this.log.error("basic name is \"void\".");
        else if (con.name == "null") this.log.error("basic name is \"null\".");
    }

    private void validateField(Field con) {
        this.validateChild(con, con.type, "type", false);
    }

    private void validateInterfaceDecl(InterfaceDecl con) {
        this.validateChild(con, con.inter, "inter", false);
        this.validateChildren(con, con.typeParams, "typeParams", true);
        this.validateChildren(con, con.instances,  "instances", true);
    }

    private void validateInterfaceDesc(InterfaceDesc con) {
        this.validateChildren(con, con.abstracts, "abstracts", true);
        this.validateChildren(con, con.inherits,  "inherits", true);
        this.validateChild(con, con.pin, "pin", true);
    }

    private void validateInterfaceInst(InterfaceInst con) {
        this.validateChild(con, con.generic,  "generic", false);
        this.validateChild(con, con.resolved, "resolved", false);
        this.validateChildren(con, con.instanceTypes, "instanceTypes", false);

        final InterfaceDecl gen = con.generic.getResolved();
        if (gen != null) this.validateInstantiation(gen, con, gen.typeParams, con.instanceTypes);
    }

    private void validateMethodDecl(MethodDecl con) {
        // Receiver maybe nil in Go, but Java should always have a receiver that is the containing class.
        this.validateChild(con, con.receiver,  "receiver", false);
        this.validateChild(con, con.signature, "signature", false);
        this.validateChildren(con, con.typeParams, "typeParams", true);
        this.validateChildren(con, con.instances,  "instances", true);
        this.validateChild(con, con.metrics, "metrics", true);
    }

    private void validateMethodInst(MethodInst con) {
        // Receiver maybe nil in Go, but Java should always have a receiver that is the containing class.
        this.validateChild(con, con.receiver, "receiver", false);
        this.validateChild(con, con.generic,  "generic", false);
        this.validateChild(con, con.resolved, "resolved", false);
        this.validateChildren(con, con.instanceTypes, "instanceTypes", false);

        final MethodDecl gen = con.generic.getResolved();
        if (gen != null) this.validateInstantiation(gen, con, gen.typeParams, con.instanceTypes);
    }

    private void validateMetrics(Metrics con) {
        this.validateChildren(con, con.invokes, "invokes", true);
        this.validateChildren(con, con.reads,   "reads", true);
        this.validateChildren(con, con.writes,  "writes", true);
    }

    private void validateObjectDecl(ObjectDecl con) {
        this.validateChild(con, con.struct, "struct", false);
        this.validateChild(con, con.inter,  "inter", false);
        this.validateChildren(con, con.methodDecls, "methodDecls", true);
        this.validateChildren(con, con.typeParams,  "typeParams", true);
        this.validateChildren(con, con.instances,   "instances", true);
    }

    private void validateObjectInst(ObjectInst con) {
        this.validateChild(con, con.generic,      "generic", false);
        this.validateChild(con, con.resData,      "resData", false);
        this.validateChild(con, con.resInterface, "resInterface", false);
        this.validateChildren(con, con.instanceTypes, "instanceTypes", false);
        this.validateChildren(con, con.methods,       "methods", true);

        final ObjectDecl gen = con.generic.getResolved();
        if (gen != null) this.validateInstantiation(gen, con, gen.typeParams, con.instanceTypes);
    }

    private void validatePackageCon(PackageCon con) {
        this.validateChildren(con, con.imports,        "imports", true);
        this.validateChildren(con, con.interfaceDecls, "interfaceDecls", true);
        this.validateChildren(con, con.methodDecls,    "methodDecls", true);
        this.validateChildren(con, con.objectDecls,    "objectDecls", true);
        this.validateChildren(con, con.values,         "values", true);
    }

    private void validateSelection(Selection con) {
        if (con.name.isBlank()) this.log.error("selection name is black in " + con + ".");
        this.validateChild(con, con.origin, "origin", false);
    }

    private void validateSignature(Signature con) {
        this.validateChildren(con, con.params,  "params", true);
        this.validateChildren(con, con.results, "results", true);
    }

    private void validateStructDesc(StructDesc con) {
        this.validateChildren(con, con.fields, "fields", true);
    }

    private void validateTypeParam(TypeParam con) {
        if (con.name.isBlank()) this.log.error("type param name is black in " + con + ".");
        this.validateChild(con, con.type, "type", false);
    }

    private void validateValue(Value con) {
        this.validateChild(con, con.type,    "type",    false);
        this.validateChild(con, con.metrics, "metrics", true);
    }

    private void validateInstantiation(Construct decl, Construct inst, ArrayList<Ref<TypeParam>> typeParams, ArrayList<Ref<? extends TypeDesc>> instanceTypes) {
        final int typeParamSize = typeParams.size();
        if (typeParamSize <= 0)
            this.log.error("the declaration " + decl + " had instances but has " + typeParamSize + " type parameters.");

        final int instanceTypeSize = instanceTypes.size();
        if (typeParamSize != instanceTypeSize)
            this.log.error("the declaration " + decl + " had " + typeParamSize + " (" + typeParams + ") type parameters which did not match " +
                "the instance " + inst + " that had " + instanceTypeSize + " (" + instanceTypes + ") instance types.");

        final int count = Integer.min(typeParamSize, instanceTypeSize);
        boolean match = true;
        for (int i = 0; i < count; i++) {
            if (!typeParams.get(i).equals(instanceTypes.get(i))) {
                match = false;
                break;
            }
        }
        if (match)
            this.log.error("the declaration " + decl + " and the instance " + inst + " had the same (" + count + ") type parameters.");
    }
    
    private void validateChildren(Construct parent, Iterable<? extends Ref<? extends Construct>> children, String usage, boolean maybeEmpty) {
        if (children == null) {
            this.log.error("the collection of children (" + usage + ") in " + parent + " is null.");
            return;
        }
        int index = 0;
        for (Ref<? extends Construct> child : children) {
            this.validateChild(parent, child, usage + "[" + index + "]", false);
            index++;
        }
        if (!maybeEmpty && index == 0) this.log.error("the collection of children (" + usage + ") in " + parent + " is empty.");
    }
    
    private void validateChild(Construct parent, Ref<? extends Construct> child, String usage, boolean maybeNull) {
        if (child == null) {
            if (!maybeNull) this.log.error("the child (" + usage + ") in a " + parent.kind() + " is null: " + parent);
            return;
        }

        final Factory<? extends Construct> factory = this.proj.getFactory(child.kind());
        this.validateRef(factory, child, parent + "." +usage);
    }
}
