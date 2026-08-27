package abstractor.core.tools;

import java.util.*;

import abstractor.core.constructs.*;
import abstractor.core.json.*;
import abstractor.core.log.*;
import abstractor.core.require.Require;

public class Instantiator {
    public class Frame {
        final private TreeMap<Ref<? extends TypeDesc>, Ref<? extends TypeDesc>> subst = new TreeMap<>();
        final private ArrayList<Ref<? extends TypeDesc>> paramOrder = new ArrayList<>();
        private ArrayList<Ref<? extends TypeDesc>> argOrder = null;
        private int nestCount;

        private Frame() {}

        private void copyFrom(Frame other) {
            this.subst.putAll(other.subst);
            this.paramOrder.addAll(other.paramOrder);
            this.nestCount = this.paramOrder.size();
        }

        private void copyImmediateFrom(Frame other) {
            final int fullCount = other.paramOrder.size();
            for (int i = other.nestCount; i < fullCount; i++) {
                final Ref<? extends TypeDesc> param = other.paramOrder.get(i);
                this.paramOrder.add(param);
                final Ref<? extends TypeDesc> sub = other.subst.get(param);
                this.subst.put(param, sub);
            }
            this.nestCount = this.paramOrder.size();
        }

        private void add(Ref<? extends TypeDesc> param, Ref<? extends TypeDesc> arg, Logger log) throws Exception {
            if (log != null) log.log("add(" + param + ", " + arg + ")");

            if (this.subst.put(param, arg) != null) {
                final int index = this.paramOrder.indexOf(param);
                if (log != null) log.log("  prior found at " + index);
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
            final int size = this.paramOrder.size();
            if (size <= 0) return "[ ]";

            final JsonHelper jh = new JsonHelper();
            jh.writeKinds     = true;
            jh.writeIndices   = true;
            jh.writeRefs      = true;
            jh.refSkipResolve = false;
            List<String> parts = new ArrayList<>(size);
            for (int i = 0; i < size; i++) {
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

    private class StackNode {
        final public Frame frame;
        final public StackNode prior;

        public StackNode(StackNode prior) {
            this.frame = new Frame();
            this.prior = prior;
        }
    }
    
    private StackNode topFrame = null;

    public void pushFrame() {
        final StackNode node = new StackNode(this.topFrame);
        if (this.topFrame != null)
            node.frame.copyFrom(this.topFrame.frame);
        this.topFrame = node;
    }
    
    public void pushFrameCopy(Frame copy) throws Exception {
        final StackNode node = new StackNode(this.topFrame);
        if (copy != null)
            node.frame.copyFrom(copy);
        this.topFrame = node;
    }

    public void pushCleanFrame() {
        this.topFrame = new StackNode(this.topFrame);
    }

    public void pushImmediateFrame() {
        final StackNode node = new StackNode(this.topFrame);
        if (this.topFrame != null)
            node.frame.copyImmediateFrom(this.topFrame.frame);
        this.topFrame = node;
    }

    public void popFrame() throws Exception {
        Require.notNull(this.topFrame, "instantiator has no frame to pop");
        this.topFrame = this.topFrame.prior;
    }

    public Frame copyFrame() {
        Frame frame = new Frame();
        if (this.topFrame != null)
            frame.copyFrom(this.topFrame.frame);
        return frame;
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
    public void add(Ref<? extends TypeDesc> param, Ref<? extends TypeDesc> arg, Logger log) throws Exception {
        Require.notNull(this.topFrame, "cannot add to an empty instantiator");
        Require.notNull(param, "can not have a null type parameter in an instantiator frame");
        Require.notNull(arg, "can not have a null the argument in an instantiator frame");

        // TODO: Is this actually needed? How should we handle pushing a clean frame or from another frame if it is?
        //if (this.topFrame.prior != null)
        //    arg = this.topFrame.prior.frame.replace(arg);
        this.topFrame.frame.add(param, arg, log);
    }

    public void add(Ref<? extends TypeDesc> param, Ref<? extends TypeDesc> arg) throws Exception {
        this.add(param, arg, null);
    }

    public Ref<? extends TypeDesc> replace(Ref<? extends TypeDesc> con) {
        return this.topFrame == null ? con : this.topFrame.frame.replace(con);
    }

    public List<Ref<? extends TypeDesc>> typeArgs() throws Exception {
        return this.typeArgs(true);
    }

    public List<Ref<? extends TypeDesc>> typeArgs(boolean withNest) throws Exception {
        if (this.topFrame == null) return Collections.emptyList();
        if (withNest) return this.topFrame.frame.typeArgsWithNest();
        return this.topFrame.frame.immediateTypeArgs();
    }

    @Override
    public String toString() {
        if (this.topFrame == null) return "<null>";
        return this.topFrame.frame.toString();
    }
}
