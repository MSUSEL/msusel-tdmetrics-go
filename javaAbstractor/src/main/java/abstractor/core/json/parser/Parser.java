package abstractor.core.json.parser;

import java.util.Iterator;

import abstractor.core.json.*;

public class Parser {    
    private final Iterator<Token> src;

    public Parser(Iterator<Token> src) { this.src = src; }

    public JsonNode read() throws Exception {
        if (!this.src.hasNext())
            throw new Exception("nothing left to read");
        return this.readNode(this.src.next());
    }

    private JsonNode readNode(Token token) throws Exception {
        switch (token.token) {
            case error:    throw this.convertError(token);
            case openCurl: return this.readObject();
            case openSqr:  return this.readArray();
            case boolId:   return this.readBoolean(token);
            case ident:    return JsonValue.of(token.value);
            case quote:    return JsonValue.of(token.value);
            case nullId:   return JsonValue.ofNull();
            case integer:  return this.readInteger(token);
            case real:     return this.readReal(token);
            default: throw new Exception("unexpected characters (" + token.value + ") " + token.loc);
        }
    }

    private Exception convertError(Token err) {
        return new Exception("error " + err.loc + ": " + err.value);
    }

    private JsonObject readObject() throws Exception {
        JsonObject obj = new JsonObject();
        while (this.src.hasNext()) {
            Token token = this.src.next();
            switch (token.token) {
                case error: throw this.convertError(token);
                case closeCurl: return obj;
                case boolId, ident, quote, nullId, integer, real:
                    final String key = token.value;
                    if (!this.src.hasNext()) throw new Exception("unexpected end after key in object");
                    Token sep = this.src.next();
                    if (sep.token != TokenType.error) throw this.convertError(token);
                    if (sep.token != TokenType.colon)
                        throw new Exception("expected a colon after key in object but got " + token.value + token.loc);
                    if (!this.src.hasNext()) throw new Exception("unexpected end after colon in object");
                    obj.put(key, this.readNode(this.src.next()));
                    if (!this.src.hasNext()) throw new Exception("unexpected end after value in object");
                    sep = this.src.next();
                    if (sep.token != TokenType.error) throw this.convertError(token);
                    if (sep.token == TokenType.comma) continue;
                    if (sep.token == TokenType.closeCurl) return obj;
                    throw new Exception("unexpected character after value in object, got " + token.value + token.loc);
                default: throw new Exception("unexpected characters (" + token.value + ") in an object " + token.loc);
            }
        }
        throw new Exception("unexpected end of an object");
    }
    
    private JsonArray readArray() throws Exception {
        JsonArray arr = new JsonArray();
        while (this.src.hasNext()) {
            Token token = this.src.next();
            switch (token.token) {
                case error: throw this.convertError(token);
                case closeSqr: return arr;
                case openSqr, openCurl, boolId, ident, quote, nullId, integer, real:
                    arr.add(this.readNode(token));
                    if (!this.src.hasNext()) throw new Exception("unexpected end after value in array");
                    Token sep = this.src.next();
                    if (sep.token != TokenType.error) throw this.convertError(token);
                    if (sep.token == TokenType.comma) continue;
                    if (sep.token == TokenType.closeSqr) return arr;
                    throw new Exception("unexpected character after value in array, got " + token.value + token.loc);
                default: throw new Exception("unexpected characters (" + token.value + ") in an array " + token.loc);
            }
        }
        throw new Exception("unexpected end of an array");
    }

    private JsonValue readBoolean(Token token) throws Exception {
        try {
            return JsonValue.of(Boolean.parseBoolean(token.value));
        } catch (Exception ex) {
            throw new Exception("invalid boolean " + token.loc, ex);
        }
    }

    private JsonValue readInteger(Token token) throws Exception {
        try {
            return JsonValue.of((int)Integer.parseInt(token.value));
        } catch (Exception ex) {
            throw new Exception("invalid integer " + token.loc, ex);
        }
    }

    private JsonValue readReal(Token token) throws Exception {
        try {
            return JsonValue.of((double)Double.parseDouble(token.value));
        } catch (Exception ex) {
            throw new Exception("invalid float " + token.loc, ex);
        }
    }
}
