package abstractor.core.json;

import java.io.PrintStream;

public interface JsonNode {
    boolean isEmpty();
    void toString(PrintStream sb, boolean minimize, String indent);
    String toString(boolean minimize);
}
