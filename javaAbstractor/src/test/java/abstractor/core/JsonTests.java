package abstractor.core;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

import org.junit.jupiter.api.Test;

import abstractor.core.json.*;

public class JsonTests {
    
    @Test
    public void FormatValueNull() {
        final JsonValue v = JsonValue.ofNull();
        assertTrue(v.isNull());
        assertFalse(v.isString());
        assertFalse(v.isInt());
        assertFalse(v.isDouble());
        assertFalse(v.isBool());

        assertEquals("null", v.toString());
        assertEquals("null", v.asString());
        assertEquals(0, v.asInt());
        assertEquals(0.0, v.asDouble());
        assertFalse(v.asBool());
    }

    @Test
    public void FormatValueBool1() {
        final JsonValue v1 = JsonValue.of(true);
        assertFalse(v1.isNull());
        assertFalse(v1.isString());
        assertFalse(v1.isInt());
        assertFalse(v1.isDouble());
        assertTrue(v1.isBool());

        assertEquals("true", v1.toString());
        assertEquals("true", v1.asString());
        assertEquals(1, v1.asInt());
        assertEquals(1.0, v1.asDouble());
        assertTrue(v1.asBool());
    }
    
    @Test
    public void FormatValueBool2() {
        final JsonValue v2 = JsonValue.of(false);
        assertEquals("false", v2.toString());
        assertEquals("false", v2.asString());
        assertEquals(0, v2.asInt());
        assertEquals(0.0, v2.asDouble());
        assertFalse(v2.asBool());
    }
    
    @Test
    public void FormatValueInt1() {
        final JsonValue v1 = JsonValue.of(321);
        assertFalse(v1.isNull());
        assertFalse(v1.isString());
        assertTrue(v1.isInt());
        assertFalse(v1.isDouble());
        assertFalse(v1.isBool());

        assertEquals("321", v1.toString());
        assertEquals("321", v1.asString());
        assertEquals(321, v1.asInt());
        assertEquals(321.0, v1.asDouble());
        assertTrue(v1.asBool());
    }

    @Test
    public void FormatValueInt2() {
        final JsonValue v2 = JsonValue.of(0);
        assertEquals("0", v2.toString());
        assertEquals("0", v2.asString());
        assertEquals(0, v2.asInt());
        assertEquals(0.0, v2.asDouble());
        assertFalse(v2.asBool());
    }

    @Test
    public void FormatValueInt3() {
        final JsonValue v3 = JsonValue.of(-123);
        assertEquals("-123", v3.toString());
    }

    @Test
    public void FormatValueDouble1() {
        final JsonValue v1 = JsonValue.of(45.321);
        assertFalse(v1.isNull());
        assertFalse(v1.isString());
        assertFalse(v1.isInt());
        assertTrue(v1.isDouble());
        assertFalse(v1.isBool());

        assertEquals("45.321", v1.toString());
        assertEquals("45.321", v1.asString());
        assertEquals(45, v1.asInt());
        assertEquals(45.321, v1.asDouble());
        assertTrue(v1.asBool());
    }

    @Test
    public void FormatValueDouble2() {
        final JsonValue v2 = JsonValue.of(0.0);
        assertEquals("0.0", v2.toString());
        assertEquals("0.0", v2.asString());
        assertEquals(0, v2.asInt());
        assertEquals(0.0, v2.asDouble());
        assertFalse(v2.asBool());
    }

    @Test
    public void FormatValueDouble3() {
        final JsonValue v3 = JsonValue.of(-123.45);
        assertEquals("-123.45", v3.toString());
    }

    @Test
    public void FormatValueString1() {
        final JsonValue v1 = JsonValue.of("cat");
        assertFalse(v1.isNull());
        assertTrue(v1.isString());
        assertFalse(v1.isInt());
        assertFalse(v1.isDouble());
        assertFalse(v1.isBool());

        assertEquals("\"cat\"", v1.toString());
        assertEquals("cat", v1.asString());
        assertEquals(0, v1.asInt());
        assertEquals(0, v1.asDouble());
        assertFalse(v1.asBool());
    }

    @Test
    public void FormatValueString2() {
        final JsonValue v2 = JsonValue.of("0");
        assertEquals("\"0\"", v2.toString());
        assertEquals("0", v2.asString());
        assertEquals(0, v2.asInt());
        assertEquals(0.0, v2.asDouble());
        assertFalse(v2.asBool());
    }

    @Test
    public void FormatValueString3() {
        final JsonValue v3 = JsonValue.of("12.34");
        assertEquals("\"12.34\"", v3.toString());
        assertEquals("12.34", v3.asString());
        assertEquals(12, v3.asInt());
        assertEquals(12.34, v3.asDouble());
        assertTrue(v3.asBool());
    }

    @Test
    public void FormatValueString4() {
        final JsonValue v4 = JsonValue.of("true");
        assertEquals("\"true\"", v4.toString());
        assertEquals("true", v4.asString());
        assertEquals(1, v4.asInt());
        assertEquals(1.0, v4.asDouble());
        assertTrue(v4.asBool());
    }

    @Test
    public void FormatValueString5() {
        final JsonValue v5 = JsonValue.of("");
        assertEquals("\"\"", v5.toString());
        assertEquals("", v5.asString());
        assertEquals(0, v5.asInt());
        assertEquals(0.0, v5.asDouble());
        assertFalse(v5.asBool());
    }

    @Test
    public void FormatValueString6() {
        final JsonValue v6 = JsonValue.of("hell\\o \"world\"");
        assertEquals("\"hell\\\\o \\\"world\\\"\"", v6.toString());
        assertEquals("hell\\o \"world\"", v6.asString());
        assertEquals(0, v6.asInt());
        assertEquals(0.0, v6.asDouble());
        assertFalse(v6.asBool());
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
        // Outputs sorted alphabetically by key.
        assertEquals("{ \"one\": 1, \"three\": 3, \"two\": 2 }", simple.toString(false));
        assertEquals("{\"one\":1,\"three\":3,\"two\":2}", simple.toString(true));
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

    @Test
    public void ParseTrue() throws Exception {
        JsonNode node = JsonNode.parse("true");
        assertTrue(node instanceof JsonNode);
        JsonValue val = (JsonValue)node;
        assertTrue(val.isBool());
        assertTrue(val.asBool());
    }

    @Test
    public void ParseFalse() throws Exception {
        JsonNode node = JsonNode.parse("false");
        assertTrue(node instanceof JsonNode);
        JsonValue val = (JsonValue)node;
        assertTrue(val.isBool());
        assertFalse(val.asBool());
    }
    
    @Test
    public void ParseNull() throws Exception {
        JsonNode node = JsonNode.parse("null");
        assertTrue(node instanceof JsonNode);
        JsonValue val = (JsonValue)node;
        assertTrue(val.isNull());
    }
    
    @Test
    public void ParseSimpleIdent() throws Exception {
        JsonNode node = JsonNode.parse("hello");
        assertTrue(node instanceof JsonNode);
        JsonValue val = (JsonValue)node;
        assertTrue(val.isString());
        assertEquals("hello", val.asString());
    }
}
