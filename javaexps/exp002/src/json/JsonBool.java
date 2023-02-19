package json;

public record JsonBool(boolean value) implements JsonObj {

    public void write(StringBuilder buf) {
        buf.append(this.value);
    }
}
