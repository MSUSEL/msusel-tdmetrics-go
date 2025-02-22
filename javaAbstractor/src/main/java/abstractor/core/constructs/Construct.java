package abstractor.core.constructs;

import spoon.reflect.declaration.CtElement;

import abstractor.core.json.Jsonable;

public interface Construct extends Comparable<Construct>, Jsonable {
    public CtElement source();
    public void setIndex(int index);
    public int getIndex();
    public String kind();
}
