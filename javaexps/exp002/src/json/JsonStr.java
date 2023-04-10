package json;

public final class JsonStr extends JsonValue<String> {
    public JsonStr(String value) {
        super(value);
    }

    @Override
    public void write(StringBuilder buf) {
        buf.append(JsonObj.escape(this.value()));
    }
}
