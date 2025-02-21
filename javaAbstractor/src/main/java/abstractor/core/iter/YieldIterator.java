package abstractor.core.iter;

import java.util.Iterator;
import java.util.LinkedList;

public class YieldIterator<T> implements Iterator<T> {
    public interface Yield<T> {
        void yield(T value);
        void stop();
    }

    public interface HasNextFn { boolean hasNext(); }
    public interface StepFn<T> { void run(Yield<T> y); }

    private class YieldImp implements Yield<T> {
        public void yield(T value) { addPending(value); }
        public void stop() { callStop(); }
    }

    private final YieldImp yielder;
    private final HasNextFn hasNextFn;
    private final StepFn<T> step;
    private final StepFn<T> finish;
    private LinkedList<T> pending;
    private boolean stopped;

    public YieldIterator(HasNextFn hasNext, StepFn<T> step, StepFn<T> finish) {
        this.yielder   = new YieldImp();
        this.hasNextFn = hasNext;
        this.step      = step;
        this.finish    = finish;
        this.pending   = new LinkedList<T>();
        this.stopped   = false;
    }

    public YieldIterator(HasNextFn hasNext, StepFn<T> step) { this(hasNext, step, null); }    
    
    public YieldIterator(StepFn<T> step) { this(null, step, null); }    
    
    private void addPending(T value) {
        if (!this.stopped) this.pending.addLast(value);
    }
    
    private void callStop() {
        this.stopped = true;
    }
    
    private boolean hasPending() {
        return !this.pending.isEmpty();
    }

    private void seekNext() {
        if (this.hasPending() || this.stopped) return;

        while (this.hasNextFn == null || this.hasNextFn.hasNext()) {
            if (this.step != null) this.step.run(this.yielder);
            if (this.hasPending() || this.stopped) return;
        }

        if (this.finish != null) this.finish.run(this.yielder);
        this.stopped = true;
    }

    @Override
    public boolean hasNext() {
        this.seekNext();
        return this.hasPending();
    }

    @Override
    public T next() {
        this.seekNext();
        return this.pending.pollFirst();
    }
}
