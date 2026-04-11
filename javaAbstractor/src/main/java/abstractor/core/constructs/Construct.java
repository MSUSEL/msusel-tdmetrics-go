package abstractor.core.constructs;

import abstractor.core.cmp.CmpGetter;
import abstractor.core.cmp.CmpOptions;
import abstractor.core.json.Jsonable;

public interface Construct extends Comparable<Construct>, CmpGetter<Construct>, Jsonable {
    public void setCmpOptions(CmpOptions options);
    public CmpOptions getCmpOptions();
    public void setIndex(int index);
    public int getIndex();
    public ConstructKind kind();
}
