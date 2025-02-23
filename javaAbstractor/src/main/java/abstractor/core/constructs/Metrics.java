package abstractor.core.constructs;

import java.util.List;

import spoon.reflect.declaration.CtField;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class Metrics extends ConstructImp {
    public final Location loc;

    public final int codeCount;
    public final int complexity;
    public final int indents;
    public final int lineCount;

    public final boolean getter;
    public final boolean setter;

    public final List<Method> invokes;
    public final List<TypeDesc> reads;
    public final List<TypeDesc> writes;
    
    public Metrics(CtField<?> src, Location loc,
        int codeCount, int complexity, int indents, int lineCount,
        boolean getter, boolean setter,
        List<Method> invokes, List<TypeDesc> reads, List<TypeDesc> writes) {
        super(src);



        
    }

    public String kind() { return "metrics"; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        // TODO: | `loc`        | ◯ | ◯ | The [location](#locations) offset. |
        
        // TODO: | `codeCount`  | ⬤ | ◯ | The number of lines in the method that are not comments or empty. |
        // TODO: | `complexity` | ⬤ | ◯ | The cyclomatic complexity of the method. |
        // TODO: | `indents`    | ⬤ | ◯ | The indent complexity of the method. |
        // TODO: | `lineCount`  | ⬤ | ◯ | The number of lines in the method. |
        
        // TODO: | `getter`     | ⬤ | ◯ | True indicates the method is a getter pattern. |
        // TODO: | `setter`     | ⬤ | ◯ | True indicates the method is a setter pattern. |
        
        // TODO: | `invokes`    | ⬤ | ◯ | List of [keys](#keys) to methods that were invoked in the method. |
        // TODO: | `reads`      | ⬤ | ◯ | List of [keys](#keys) to types that were read from in the method. |
        // TODO: | `writes`     | ⬤ | ◯ | List of [keys](#keys) to types that were written to in the method. |
        
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c)
            // TODO: Fill out
        );
    }   
}
