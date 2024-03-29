package json;

public abstract class JsonValue<T> implements JsonObj {
    private T val;
    private T defaultVal;
    private boolean shouldOmit;

    protected JsonValue(T value) {
        this.val = value;
    }

    public void setValue(T value) {
        this.val = value;
    }

    public T value() {
        return this.val;
    }

    public void setDefaultValue(T defaultValue) {
        this.defaultVal = defaultValue;
    }

    public T defaultValue() {
        return this.defaultVal;
    }

    @Override
    public void setOmitOnDefault(boolean omit) {
        this.shouldOmit = omit;
    }

    @Override
    public boolean omit() {
        return this.shouldOmit && this.val == this.defaultVal;
    }

    @Override
    public void write(StringBuilder buf) {
        buf.append(this.val);
    }

    @Override
    public void removeOmitted() {
        // No effect on values
    }

    @Override
    public boolean equals(Object o) {
        return o instanceof JsonValue<?> other &&
            this.val.equals(other.val);
    }

    @Override
    public void assertCompare(JsonObj obj) throws Exception {
        if (!(obj instanceof JsonValue<?> other))
            throw new Exception("Expected type to be " + this.getClass().getName() + " but got " + obj.getClass().getName());
        if (!this.val.equals(other.val))
            throw new Exception("Expected value to be " + this + " but got "+ other);
    }

    @Override
    public String toString() {
        return JsonObj.toString(this);
    }
}
