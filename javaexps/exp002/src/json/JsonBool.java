package json;

public record JsonBool(boolean value) implements JsonObj {

    @Override
    public boolean isEmpty() {
        return !this.value;
    }

    @Override
    public void write(StringBuilder buf) {
        buf.append(this.value);
    }
}
