package json;

public final class JsonNum extends JsonValue<Double> {
    public JsonNum(double value) {
        super(value);
    }

    @Override
    public String toString() {
        String str = this.value().toString();
        if (str.contains(".") || str.contains("e") || str.contains("E"))
            return str;
        return str+".0";
    }
}
