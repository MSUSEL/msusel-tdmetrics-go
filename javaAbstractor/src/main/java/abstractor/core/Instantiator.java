package abstractor.core;

import java.util.*;

import abstractor.core.constructs.*;
import abstractor.core.require.Require;

public class Instantiator {
    private class Frame {
        final public Frame prior;
        final private TreeMap<Ref<? extends TypeDesc>, Ref<? extends TypeDesc>> subst = new TreeMap<>();
        final private ArrayList<Ref<? extends TypeDesc>> paramOrder = new ArrayList<>();
        private ArrayList<Ref<? extends TypeDesc>> argOrder = null;

        public Frame(Frame prior) {
            this.prior = prior;
            if (this.prior != null) {
                this.subst.putAll(this.prior.subst);
                this.paramOrder.addAll(this.prior.paramOrder);
            }
        }

        public void add(Ref<? extends TypeDesc> param, Ref<? extends TypeDesc> arg) {
            if (this.prior != null) arg = this.prior.replace(arg);
            if (this.subst.put(param, arg) != null) this.paramOrder.remove(param);
            this.paramOrder.add(param);
            this.argOrder = null;
        }

        public Ref<? extends TypeDesc> replace(Ref<? extends TypeDesc> con) {
            final Ref<? extends TypeDesc> other = this.subst.get(con);
            return other != null ? other : con;
        }

        public List<Ref<? extends TypeDesc>> typeArgs() throws Exception {
            if (this.argOrder == null) {
                this.argOrder = new ArrayList<>(this.paramOrder.size());
                for (Ref<? extends TypeDesc> param : this.paramOrder) {
                    final Ref<? extends TypeDesc> arg = this.subst.get(param);
                    Require.notNull(arg, "can not have a null argument for type parameter " + param);
                    this.argOrder.add(arg);
                }
            }
            return this.argOrder;
        }
    }
    
    private Frame topFrame;

    public Instantiator() {
        this.topFrame = null;
    }

    public void pushFrame() { this.topFrame = new Frame(this.topFrame); }

    public void popFrame() throws Exception {
        Require.notNull(this.topFrame, "instantiator has no frame to pop");
        this.topFrame = this.topFrame.prior;
    }

    public void add(Ref<? extends TypeDesc> param, Ref<? extends TypeDesc> value) throws Exception {
        Require.notNull(this.topFrame, "cannot add to an empty instantiator");
        Require.notNull(param, "can not have a null key in an instantiator frame");
        Require.notNull(value, "can not have a null value in an instantiator frame");
        this.topFrame.add(param, value);
    }
    
    public Ref<? extends TypeDesc> replace(Ref<? extends TypeDesc> con) {
        return this.topFrame == null ? con : this.topFrame.replace(con);
    }
    
    public List<Ref<? extends TypeDesc>> typeArgs() throws Exception {
        return this.topFrame == null ? Collections.emptyList() : this.topFrame.typeArgs();
    }
}
