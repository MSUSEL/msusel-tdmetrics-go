package json;

import javax.lang.model.type.UnknownTypeException;

public interface JsonObj {
    void write(StringBuilder buf);

    static String escape(String value) {
        return "\"" + value.
            replace("\\", "\\\\").
            replace("\t", "\\t").
            replace("\n", "\\n").
            replace("\r", "\\r").
            replace("\b", "\\b").
            replace("\"", "\\\"") + "\"";
    }

    static JsonObj convert(Object data) {
        if (data == null) return null;
        if (data instanceof JsonObj)  return (JsonObj) data;
        if (data instanceof Jsonable) return ((Jsonable) data).toJson();
        if (data instanceof Integer)  return new JsonInt((int) data);
        if (data instanceof Boolean)  return new JsonBool((boolean) data);
        if (data instanceof Double)   return new JsonNum((double) data);
        if (data instanceof String)   return new JsonStr((String) data);
        if (data instanceof Iterable) return new JsonList((Iterable) data);
        return new JsonStr(data.toString());
    }

    static void write(StringBuilder buf, Object obj) {
        if (obj == null) buf.append("null");
        else convert(obj).write(buf);
    }

    static String toString(Object obj) {
        StringBuilder buf = new StringBuilder();
        write(buf, convert(obj));
        return buf.toString();
    }
}
