package abstractor.core.diff.core;

import java.util.Iterator;
import java.util.LinkedList;

import abstractor.core.diff.DiffStep;

public class Simplifier implements Iterator<DiffStep> {
    private final Iterator<DiffStep> steps;
    private int addedRun, removedRun, equalRun;
    private LinkedList<DiffStep> pending;

    public Simplifier(Iterator<DiffStep> steps) {
        this.steps      = steps;
        this.addedRun   = 0;
        this.removedRun = 0;
        this.equalRun   = 0;
        this.pending    = new LinkedList<DiffStep>();
    }

    private void addEqual() {
        if (this.equalRun > 0) {
            this.pending.addLast(DiffStep.Equal(this.equalRun));
            this.equalRun = 0;
        }
    }
   
    private void addRemoved() {
        if (this.removedRun > 0) {
            this.pending.addLast(DiffStep.Removed(this.removedRun));
            this.removedRun = 0;
        }
    }

    private void addAdded() {
        if (this.addedRun > 0) {
            this.pending.addLast(DiffStep.Added(this.addedRun));
            this.addedRun = 0;
        }
    }

    @Override
    public boolean hasNext() {
        return !this.pending.isEmpty() || this.steps.hasNext() ||
            this.addedRun > 0 || this.removedRun > 0 || this.equalRun > 0;
    }

    @Override
    public DiffStep next() {
        if (!this.pending.isEmpty()) return this.pending.pollFirst();

        while (this.steps.hasNext()) {
            final DiffStep step = this.steps.next();
            if (step.count() <= 0) continue;

            switch (step.type()) {

                case added:
                    this.addEqual();
                    this.addedRun += step.count();
                    break;

                case removed:
                    this.addEqual();
                    this.removedRun += step.count();
                    break;

                case equal:
                    this.addRemoved();
                    this.addAdded();
                    this.equalRun += step.count();
                    break;
            }
            if (!this.pending.isEmpty()) return this.pending.pollFirst();
        }

        this.addRemoved();
        this.addAdded();
        this.addEqual();
        return this.pending.pollFirst();
    }
}
