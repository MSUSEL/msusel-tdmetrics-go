package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.Collection;

public class Group {
    public final Construct[] cons;

    public Group(Construct first, Collection<? extends Construct> rest) {
        ArrayList<Construct> conList = new ArrayList<>();
        conList.add(first);
        conList.addAll(conList);
        Construct[] cons = new Construct[conList.size()];
        this.cons = conList.toArray(cons);
    }

    // TODO: Finish Implementing

}
