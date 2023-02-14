package json

import kotlin.collections.HashMap

class JsonMap: HashMap<String, JsonObj?>(), JsonObj {
    override fun write(buf: StringBuilder) {
        buf.append("{")
        var first = true
        for ((key, value) in this) {
            JsonObj.Companion.writeEscapedJsonString(buf, key)
            buf.append(":")
            JsonObj.Companion.write(buf, value)
            if (first) first = false
            else buf.append(",")
        }
        buf.append("}")
    }
}
