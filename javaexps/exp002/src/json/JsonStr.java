package json;

public record JsonStr(String value) implements JsonObj {

    public void write(StringBuilder buf) {
        buf.append(JsonObj.escape(this.value));
    }
}
