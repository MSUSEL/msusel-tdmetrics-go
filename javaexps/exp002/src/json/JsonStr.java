package json;

public record JsonStr(String value) implements JsonObj {

    @Override
    public boolean isEmpty() {
        return this.value.length() == 0;
    }

    @Override
    public void write(StringBuilder buf) {
        buf.append(JsonObj.escape(this.value));
    }
}
