package abstractor.core;

import java.util.*;

import abstractor.core.constructs.*;
import abstractor.core.require.Require;
import abstractor.core.log.*;

public class Instantiator {
    private class Frame {
        final Frame prior;
        final TreeMap<Ref<? extends TypeDesc>, Ref<? extends TypeDesc>> subst = new TreeMap<>();

        public Frame(Frame prior) {
            this.prior = prior;
            if (this.prior != null)
                this.subst.putAll(this.prior.subst);
        }

        public void add(Ref<TypeParam> key, Ref<? extends TypeDesc> value) {
            if (this.prior != null) value = this.prior.replace(value);
            this.subst.put(key, value);
        }

        public Ref<? extends TypeDesc> replace(Ref<? extends TypeDesc> con) {
            final Ref<? extends TypeDesc> other = this.subst.get(con);
            return other != null ? other : con;
        }
    }
    
    public final Logger  log;
    public final Project proj;

    private Frame topFrame;

    public Instantiator(Logger log, Project proj) {
        this.log      = log;
        this.proj     = proj;
        this.topFrame = null;
    }

    private void pushFrame() { this.topFrame = new Frame(this.topFrame); }

    private void popFrame() throws Exception {
        Require.notNull(this.topFrame, "instantiator has no frame to pop");
        this.topFrame = this.topFrame.prior;
    }

    private void addReplacement(Ref<TypeParam> key, Ref<? extends TypeDesc> value) throws Exception {
        Require.notNull(this.topFrame, "cannot add to an empty instantiator");
        Require.notNull(key, "can not have a null key in an instantiator frame");
        Require.notNull(value, "can not have a null value in an instantiator frame");
        this.topFrame.add(key, value);
    }
    
    private Ref<? extends TypeDesc> replace(Ref<? extends TypeDesc> con) {
        return this.topFrame == null ? con : this.topFrame.replace(con);
    }

    public ObjectInst instantiate(ObjectDecl decl, List<Ref<? extends TypeDesc>> typeArgs) {

        // TODO: Implement

        return null;
    }

    public InterfaceInst instantiate(InterfaceDecl decl, List<Ref<? extends TypeDesc>> typeArgs) {

        // TODO: Implement

        return null;
    }

    public MethodInst instantiate(MethodDecl decl, List<Ref<? extends TypeDesc>> typeArgs) {

        // TODO: Implement

        return null;
    }
}
