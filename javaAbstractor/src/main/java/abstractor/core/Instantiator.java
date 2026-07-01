package abstractor.core;

import java.util.*;

import abstractor.core.constructs.*;
import abstractor.core.json.*;
import abstractor.core.require.Require;

public class Instantiator {
    private class Frame {
        final public Frame prior;
        final private TreeMap<Ref<? extends TypeDesc>, Ref<? extends TypeDesc>> subst = new TreeMap<>();
        final private ArrayList<Ref<? extends TypeDesc>> paramOrder = new ArrayList<>();
        private ArrayList<Ref<? extends TypeDesc>> argOrder = null;
        private int nestCount;

        public Frame(Frame prior) {
            this.prior = prior;
        }

        public void copyFromPrior() {
            if (prior == null) return;
            this.subst.putAll(prior.subst);
            this.paramOrder.addAll(prior.paramOrder);
            this.nestCount = this.paramOrder.size();
        }

        public void add(Ref<? extends TypeDesc> param, Ref<? extends TypeDesc> arg) throws Exception {
            if (this.prior != null) arg = this.prior.replace(arg);
            if (this.subst.put(param, arg) != null) {
                final int index = this.paramOrder.indexOf(param);
                if (index >= 0) {
                    this.paramOrder.remove(index);
                    if (index < this.nestCount) this.nestCount--;
                }
            }
            this.paramOrder.add(param);
            this.argOrder = null;
        }

        public Ref<? extends TypeDesc> replace(Ref<? extends TypeDesc> con) {
            final Ref<? extends TypeDesc> other = this.subst.get(con);
            return other != null ? other : con;
        }

        public List<Ref<? extends TypeDesc>> typeArgsWithNest() throws Exception {
            if (this.argOrder != null) return this.argOrder;

            this.argOrder = new ArrayList<>(this.paramOrder.size());
            for (Ref<? extends TypeDesc> param : this.paramOrder) {
                final Ref<? extends TypeDesc> arg = this.subst.get(param);
                Require.notNull(arg, "can not have a null argument for type parameter " + param);
                this.argOrder.add(arg);
            }
            return this.argOrder;
        }
        
        public List<Ref<? extends TypeDesc>> immediateTypeArgs() throws Exception {
            final List<Ref<? extends TypeDesc>> full = this.typeArgsWithNest();
            return full.subList(this.nestCount, full.size());
        }

        @Override
        public String toString() {
            final JsonHelper jh = new JsonHelper();
            jh.writeKinds     = true;
            jh.writeIndices   = true;
            jh.writeRefs      = true;
            jh.refSkipResolve = false;
            List<String> parts = new ArrayList<>(this.paramOrder.size());
            for (int i = 0; i < this.paramOrder.size(); i++) {
                final Ref<? extends TypeDesc> param = this.paramOrder.get(i);
                final Ref<? extends TypeDesc> arg   = this.subst.get(param);
                final String paramStr = JsonFormat.Inline().format(param.toJson(jh));
                final String argStr   = JsonFormat.Inline().format(arg.toJson(jh));
                final String header   = i < this.nestCount
                    ? "nest." + i + ". "
                    : (i - this.nestCount) + ". ";
                parts.add(header + paramStr + " => " + argStr);    
            }
            return "[\n\t" + String.join("\n\t", parts) + "\n]";
        }
    }
    
    private Frame topFrame;

    public Instantiator() {
        this.topFrame = null;
    }

    public void pushFrame() {
        this.topFrame = new Frame(this.topFrame);
        this.topFrame.copyFromPrior();
    }

    public void pushCleanFrame() {
        this.topFrame = new Frame(this.topFrame);
    }

    public void popFrame() throws Exception {
        Require.notNull(this.topFrame, "instantiator has no frame to pop");
        this.topFrame = this.topFrame.prior;
    }

    /**
     * Adds a parameter to argument map to the current frame. The argument
     * is automatically updates with prior matching parameters. If there are
     * matching parameters in the current frame, the older one is removed and
     * the new parameter is added to the end.
     * 
     * This is designed so that if <T> is nested inside of a <T, U>, the new T
     * overrides the older T, so we then are left with <U> from the nest and <T>
     * for the current type params, shown like <U; T>, where `;` separates
     * the nest from the current.
     */
    public void add(Ref<? extends TypeDesc> param, Ref<? extends TypeDesc> arg) throws Exception {
        Require.notNull(this.topFrame, "cannot add to an empty instantiator");
        Require.notNull(param, "can not have a null type parameter in an instantiator frame");
        Require.notNull(arg, "can not have a null the argument in an instantiator frame");
        this.topFrame.add(param, arg);
    }

    public Ref<? extends TypeDesc> replace(Ref<? extends TypeDesc> con) {
        return this.topFrame == null ? con : this.topFrame.replace(con);
    }

    public List<Ref<? extends TypeDesc>> typeArgs(boolean withNest) throws Exception {
        if (this.topFrame == null) return Collections.emptyList();
        if (withNest) return this.topFrame.typeArgsWithNest();
        return this.topFrame.immediateTypeArgs();
    }

    @Override
    public String toString() {
        if (this.topFrame == null) return "<null>";
        return this.topFrame.toString();
    }
}
