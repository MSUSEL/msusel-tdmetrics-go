package json;

public record JsonInt(int value) implements JsonObj {

    @Override
    public boolean isEmpty() {
        return this.value == 0;
    }

    @Override
    public void write(StringBuilder buf) {
        buf.append(this.value);
    }
}
