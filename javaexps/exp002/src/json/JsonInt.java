package json;

public record JsonInt(int value) implements JsonObj {

    public void write(StringBuilder buf) {
        buf.append(this.value);
    }
}
