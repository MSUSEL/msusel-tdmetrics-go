package abstractor.core;

import static org.junit.jupiter.api.Assertions.assertEquals;

import org.junit.jupiter.api.Test;

import abstractor.core.json.*;

public class JsonTests {
    
    @Test
    public void FormatValueNull() {
        assertEquals("null", JsonValue.ofNull().toString());
    }

    @Test
    public void FormatValueBool() {
        assertEquals("true", JsonValue.of(true).toString());
        assertEquals("false", JsonValue.of(false).toString());
    }
    
    @Test
    public void FormatValueInt() {
        assertEquals("0", JsonValue.of(0).toString());
        assertEquals("-123", JsonValue.of(-123).toString());
        assertEquals("321", JsonValue.of(321).toString());
    }

    @Test
    public void FormatValueFloat() {
        assertEquals("0.0", JsonValue.of(0.0).toString());
        assertEquals("-123.45", JsonValue.of(-123.45).toString());
        assertEquals("45.321", JsonValue.of(45.321).toString());
    }

    @Test
    public void FormatValueString() {
        assertEquals("\"\"", JsonValue.of("").toString());
        assertEquals("\"cat\"", JsonValue.of("cat").toString());
        assertEquals("\"hell\\\\o \\\"world\\\"\"", JsonValue.of("hell\\o \"world\"").toString());
    }

    @Test
    public void FormatArrayEmpty() {
        JsonArray simple = new JsonArray();
        assertEquals("[ ]", simple.toString(false));
        assertEquals("[]", simple.toString(true));
    }

    @Test
    public void FormatArraySimple() {
        JsonArray simple = new JsonArray(
            JsonValue.of(1),
            JsonValue.of(2),
            JsonValue.of(3));
        assertEquals("[ 1, 2, 3 ]", simple.toString(false));
        assertEquals("[1,2,3]", simple.toString(true));
    }
    
    @Test
    public void FormatArrayComplex() {
        JsonArray complex = new JsonArray(
            new JsonArray(
                JsonValue.of(1),
                JsonValue.of(2),
                JsonValue.of(3)),
            new JsonArray(
                JsonValue.of(4),
                JsonValue.of(5),
                JsonValue.of(6)));
        assertEquals("[\n" +
                     "  [ 1, 2, 3 ],\n" +
                     "  [ 4, 5, 6 ]\n" +
                     "]", complex.toString(false));
        assertEquals("[[1,2,3],[4,5,6]]", complex.toString(true));
    }

    @Test
    public void FormatObjectEmpty() {
        JsonObject simple = new JsonObject();
        assertEquals("{ }", simple.toString(false));
        assertEquals("{}", simple.toString(true));
    }

    @Test
    public void FormatObjectSimple() {
        JsonObject simple = new JsonObject();
        simple.put("one", JsonValue.of(1));
        simple.put("two", JsonValue.of(2));
        simple.put("three", JsonValue.of(3));
        assertEquals("{ \"one\": 1, \"two\": 2, \"three\": 3 }", simple.toString(false));
        assertEquals("{\"one\":1,\"two\":2,\"three\":3}", simple.toString(true));
    }

    @Test
    public void FormatObjectComplex() {
        JsonObject complex = new JsonObject();
        complex.put("one", new JsonArray(
            JsonValue.of(1),
            JsonValue.of(2),
            JsonValue.of(3)));
        complex.put("two", new JsonArray(
            JsonValue.of(4),
            JsonValue.of(5),
            JsonValue.of(6)));
        assertEquals("{\n" +
                     "  \"one\": [ 1, 2, 3 ],\n" +
                     "  \"two\": [ 4, 5, 6 ]\n" +
                     "}", complex.toString(false));
        assertEquals("{\"one\":[1,2,3],\"two\":[4,5,6]}", complex.toString(true));
    }
}
