package abstractor.core.constructs;

import abstractor.core.json.*;

public interface Construct extends Comparable<Construct>, Jsonable {
    public void setIndex(int index);
    public int getIndex();
    public Object source();
    public String kind();
}
