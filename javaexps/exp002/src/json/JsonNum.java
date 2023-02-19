package json;

public record JsonNum(double value) implements JsonObj {

    public void write(StringBuilder buf) {
        buf.append(this.value);
    }
}
