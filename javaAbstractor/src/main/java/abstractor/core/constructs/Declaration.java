package abstractor.core.constructs;

import java.util.TreeSet;

public interface Declaration extends Construct {
    public Ref<PackageCon> pkgRef();
    public Ref<? extends Construct> getNest();
    public void setNest(Ref<? extends Construct> nest) throws Exception;
    public TreeSet<Ref<? extends TypeDesc>> getNestedTypes();
}
