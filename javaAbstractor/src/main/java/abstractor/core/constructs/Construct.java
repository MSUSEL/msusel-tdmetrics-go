package abstractor.core.constructs;

import abstractor.core.json.Jsonable;

public interface Construct extends Comparable<Construct>, Jsonable {
    public void setIndex(int index);
    public int getIndex();
    public ConstructKind kind();
}
