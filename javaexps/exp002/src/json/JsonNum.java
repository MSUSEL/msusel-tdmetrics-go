package json;

public record JsonNum(double value) implements JsonObj {

    @Override
    public boolean isEmpty() {
        return this.value == 0.0;
    }

    @Override
    public void write(StringBuilder buf) {
        buf.append(this.value);
    }
}
