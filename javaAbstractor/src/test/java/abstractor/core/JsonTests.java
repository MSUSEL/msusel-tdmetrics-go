package abstractor.core;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static abstractor.core.Testing.*;

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
        checkJson(JsonFormat.Expand(),   simple, "[ ]");
        checkJson(JsonFormat.Normal(),   simple, "[ ]");
        checkJson(JsonFormat.Relaxed(),  simple, "[ ]");
        checkJson(JsonFormat.Minimize(), simple, "[]");
    }

    @Test
    public void FormatArraySimple() {
        JsonArray simple = new JsonArray(
            JsonValue.of(1),
            JsonValue.of(2),
            JsonValue.of(3));
        checkJson(JsonFormat.Expand(), simple,
            "[",
            "  1,",
            "  2,",
            "  3",
            "]");
        checkJson(JsonFormat.Normal(), simple, "[ 1, 2, 3 ]");
        checkJson(JsonFormat.Relaxed(), simple,
            "[",
            "  1,",
            "  2,",
            "  3",
            "]");
        checkJson(JsonFormat.Minimize(), simple, "[1,2,3]");
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
        checkJson(JsonFormat.Expand(), complex,
            "[",
            "  [",
            "    1,",
            "    2,",
            "    3",
            "  ],",
            "  [",
            "    4,",
            "    5,",
            "    6",
            "  ]",
            "]");
        checkJson(JsonFormat.Normal(), complex,
            "[",
            "  [ 1, 2, 3 ],",
            "  [ 4, 5, 6 ]",
            "]");
        checkJson(JsonFormat.Relaxed(), complex,
            "[",
            "  [ 1, 2, 3 ],",
            "  [ 4, 5, 6 ]",
            "]");
        checkJson(JsonFormat.Minimize(), complex, "[[1,2,3],[4,5,6]]");
    }

    @Test
    public void FormatObjectEmpty() {
        JsonObject simple = new JsonObject();
        assertEquals("{ }", JsonFormat.Expand().format(simple));
        assertEquals("{ }", JsonFormat.Normal().format(simple));
        assertEquals("{ }", JsonFormat.Relaxed().format(simple));
        assertEquals("{}",  JsonFormat.Minimize().format(simple));
    }

    @Test
    public void FormatObjectSimple() {
        JsonObject simple = new JsonObject();
        simple.put("one", JsonValue.of(1));
        simple.put("two", JsonValue.of(2));
        simple.put("three", JsonValue.of(3));
        // Outputs sorted alphabetically by key.
        checkJson(JsonFormat.Expand(), simple,
            "{",
            "  \"one\": 1,",
            "  \"three\": 3,",
            "  \"two\": 2",
            "}");
        checkJson(JsonFormat.Normal(),   simple, "{ \"one\": 1, \"three\": 3, \"two\": 2 }");
        checkJson(JsonFormat.Relaxed(),  simple, 
            "{",
            "  one: 1,",
            "  three: 3,",
            "  two: 2",
            "}");
        checkJson(JsonFormat.Minimize(), simple, "{\"one\":1,\"three\":3,\"two\":2}");
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
        checkJson(JsonFormat.Expand(), complex,
            "{",
            "  \"one\": [",
            "    1,",
            "    2,",
            "    3",
            "  ],",
            "  \"two\": [",
            "    4,",
            "    5,",
            "    6",
            "  ]",
            "}");
        checkJson(JsonFormat.Normal(), complex,
            "{",
            "  \"one\": [ 1, 2, 3 ],",
            "  \"two\": [ 4, 5, 6 ]",
            "}");
        checkJson(JsonFormat.Relaxed(), complex,
            "{",
            "  one: [ 1, 2, 3 ],",
            "  two: [ 4, 5, 6 ]",
            "}");
        checkJson(JsonFormat.Minimize(), complex, "{\"one\":[1,2,3],\"two\":[4,5,6]}");
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
    
    @Test
    public void ParseComplexIdent() throws Exception {
        JsonNode node = JsonNode.parse("$_hello_42");
        assertTrue(node instanceof JsonNode);
        JsonValue val = (JsonValue)node;
        assertTrue(val.isString());
        assertEquals("$_hello_42", val.asString());
    }

    @Test
    public void ParseEmptyQuote() throws Exception {
        JsonNode node = JsonNode.parse("\"\"");
        assertTrue(node instanceof JsonNode);
        JsonValue val = (JsonValue)node;
        assertTrue(val.isString());
        assertEquals("", val.asString());
    }

    @Test
    public void ParseEscapedQuote() throws Exception {
        JsonNode node = JsonNode.parse("\"hello \\\"'\\nworld\\t!\"");
        assertTrue(node instanceof JsonNode);
        JsonValue val = (JsonValue)node;
        assertTrue(val.isString());
        assertEquals("hello \"'\nworld\t!", val.asString());
    }

    @Test
    public void ParseHexQuote() throws Exception {
        JsonNode node = JsonNode.parse("\"\\uF12A\"");
        assertTrue(node instanceof JsonNode);
        JsonValue val = (JsonValue)node;
        assertTrue(val.isString());
        assertEquals("\uF12A", val.asString());
    }

    @Test
    public void ParseZeroInteger() throws Exception {
        JsonNode node = JsonNode.parse("0");
        assertTrue(node instanceof JsonNode);
        JsonValue val = (JsonValue)node;
        assertTrue(val.isInt());
        assertEquals(0, val.asInt());
    }

    @Test
    public void ParsePosInteger() throws Exception {
        JsonNode node = JsonNode.parse("1234");
        assertTrue(node instanceof JsonNode);
        JsonValue val = (JsonValue)node;
        assertTrue(val.isInt());
        assertEquals(1234, val.asInt());
    }

    @Test
    public void ParseNegInteger() throws Exception {
        JsonNode node = JsonNode.parse("-246");
        assertTrue(node instanceof JsonNode);
        JsonValue val = (JsonValue)node;
        assertTrue(val.isInt());
        assertEquals(-246, val.asInt());
    }

    @Test
    public void ParseDecimalReal() throws Exception {
        JsonNode node = JsonNode.parse("3.14");
        assertTrue(node instanceof JsonNode);
        JsonValue val = (JsonValue)node;
        assertTrue(val.isDouble());
        assertEquals(3.14, val.asDouble());
    }

    @Test
    public void ParseNegDecimalReal() throws Exception {
        JsonNode node = JsonNode.parse("-0.2");
        assertTrue(node instanceof JsonNode);
        JsonValue val = (JsonValue)node;
        assertTrue(val.isDouble());
        assertEquals(-0.2, val.asDouble());
    }
    
    @Test
    public void ParseExpReal() throws Exception {
        JsonNode node = JsonNode.parse("1e3");
        assertTrue(node instanceof JsonNode);
        JsonValue val = (JsonValue)node;
        assertTrue(val.isDouble());
        assertEquals(1000.0, val.asDouble());
    }
    
    @Test
    public void ParseDecAndExpReal() throws Exception {
        JsonNode node = JsonNode.parse("1.02e03");
        assertTrue(node instanceof JsonNode);
        JsonValue val = (JsonValue)node;
        assertTrue(val.isDouble());
        assertEquals(1020.0, val.asDouble());
    }

    @Test
    public void ParseNegDecAndExpReal() throws Exception {
        JsonNode node = JsonNode.parse("-124.10e-2");
        assertTrue(node instanceof JsonNode);
        JsonValue val = (JsonValue)node;
        assertTrue(val.isDouble());
        assertEquals(-1.241, val.asDouble());
    }
    
    @Test
    public void ParseEmptyArray() throws Exception {
        JsonNode node = JsonNode.parse(" [ ] ");
        assertTrue(node instanceof JsonArray);
        JsonArray arr = (JsonArray)node;
        assertTrue(arr.isEmpty());
        assertEquals("[ ]", arr.toString());
    }
    
    @Test
    public void ParseComplexNumberArray() throws Exception {
        JsonNode node = JsonNode.parse("[ 1,    2,\t3, 4, ]");
        assertTrue(node instanceof JsonArray);
        JsonArray arr = (JsonArray)node;
        assertFalse(arr.isEmpty());
        assertEquals("[ 1, 2, 3, 4 ]", arr.toString());
    }
    
    @Test
    public void ParseArrayOfArray() throws Exception {
        JsonNode node = JsonNode.parse("[ [ ], [ hello ], [ 1, 2 ] ]");
        assertTrue(node instanceof JsonArray);
        JsonArray arr = (JsonArray)node;
        assertFalse(arr.isEmpty());
        assertEquals(
            "[\n" +
            "  [ ],\n" +
            "  [ \"hello\" ],\n" +
            "  [ 1, 2 ]\n" +            
            "]", arr.toString());
    }

    @Test
    public void ParseEmptyObject() throws Exception {
        JsonNode node = JsonNode.parse("{ }");
        assertTrue(node instanceof JsonObject);
        JsonObject obj = (JsonObject)node;
        assertTrue(obj.isEmpty());
        assertEquals("{ }", obj.toString());
    }
    
    @Test
    public void ParseComplexObject() throws Exception {
        JsonNode node = JsonNode.parse("{ hello: world, 12: 34, xyz: [true, false], }");
        assertTrue(node instanceof JsonObject);
        JsonObject obj = (JsonObject)node;
        assertFalse(obj.isEmpty());
        assertEquals(
            "{\n" +
            "  \"12\": 34,\n" +
            "  \"hello\": \"world\",\n" +
            "  \"xyz\": [ true, false ]\n" +
            "}", obj.toString());
    }

    @Test
    public void ParseComment() throws Exception {
        JsonNode node = JsonNode.parse(
            "[ 1, # hello world\n" +
            " 2, #\n" +
            "] #byeee");
        assertTrue(node instanceof JsonArray);
        JsonArray arr = (JsonArray)node;
        assertFalse(arr.isEmpty());
        assertEquals("[ 1, 2 ]", arr.toString());
    }
}
