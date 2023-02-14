package json

import java.util.TreeMap

class JsonMap: TreeMap<String, JsonObj?>(), JsonObj {
    override fun write(buf: StringBuilder) {
        buf.append("{")
        var first = true
        for ((key, value) in this.toSortedMap()) {
            if (first) first = false
            else buf.append(",")
            JsonObj.Companion.writeEscapedJsonString(buf, key)
            buf.append(":")
            JsonObj.Companion.write(buf, value)
        }
        buf.append("}")
    }
}
