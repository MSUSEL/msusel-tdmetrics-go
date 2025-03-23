package abstractor.core;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.SortedSet;
import java.util.TreeMap;
import java.util.TreeSet;

import abstractor.core.constructs.*;

public class Instantiator {
    final private Project proj;
    final private Map<TypeParam, TypeDesc> map;
    final private Set<Construct> inProg;
    
    static public InterfaceInst instantiate(Project proj, InterfaceDecl decl, List<TypeDesc> typeArgs) {
        if (decl == null || typeArgs.size() <= 0 || decl.typeParams.size() <= 0) return null;
        final Instantiator it = new Instantiator(proj, decl.typeParams, typeArgs);
        return it.inst(decl);
    }

    private Instantiator(Project proj, List<TypeParam> tp, List<TypeDesc> ta) {
        this.proj   = proj;
        this.map    = new TreeMap<TypeParam,TypeDesc>();
        this.inProg = new TreeSet<Construct>();
        final int count = Integer.min(tp.size(), ta.size());
        for (int i = 0; i < count; ++i)
            this.map.put(tp.get(i), ta.get(i));
    }

    private TypeDesc lookup(TypeParam tp) {
        final TypeDesc instance = this.map.get(tp);
        return  instance != null ? instance : tp; 
    }

    private InterfaceInst inst(InterfaceDecl d) {
        final List<TypeDesc> its = new ArrayList<TypeDesc>();
        for (TypeParam tp: d.typeParams) its.add(this.lookup(tp));
        final InterfaceDesc resolved = this.inst(d.inter);
        return this.proj.interfaceInsts.addOrGet(new InterfaceInst(d, its, resolved));
    }

    private InterfaceDesc inst(InterfaceDesc d) {
        final SortedSet<Abstract> abstracts = new TreeSet<Abstract>();
        for (Abstract ab : d.abstracts) abstracts.add(this.inst(ab));
        return this.proj.interfaceDescs.addOrGet(new InterfaceDesc(abstracts, d.pin));
    }

    private Abstract inst(Abstract a) {

        // TODO: Implement or determine if needed
        return null;
    }
}
