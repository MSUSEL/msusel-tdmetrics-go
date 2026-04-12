package abstractor.core.validator;

import abstractor.core.log.Logger;
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
        for (Ref<? extends Construct> ref : factory.refSet) this.validate(ref);
        for (Construct con : factory.conSet) this.validate(con);
    }

    private void validate(Ref<? extends Construct> ref) {
        if (!ref.isResolved()) {
            this.log.error("Expected all references to be resolved but " + ref + " was not resolved.");
            return;
        }

        Construct con = ref.getResolved();
        Factory<? extends Construct> factory = this.proj.getFactory(con.kind());
        boolean found = false;
        for (Construct other : factory.conSet) {
            // Use `==` not `equals` to ensure exact reference.
            if (other == con) {
                found = true;
                break;
            }
        }
        if (!found) {
            this.log.error("Expected all resolved references to exist in factory " +
                "but " + ref + " resolved to " + con + " was not in \"" + factory + "\".");
        }
    }

    private void validate(Construct con) {
        // TODO: Add more as needed
    }
}
