package json

import java.util.TreeMap

class JsonMap: TreeMap<String, JsonObj?>(), JsonObj {
    override fun write(buf: StringBuilder) {
        buf.append("{")
        var first = true
        for ((key, value) in this.toSortedMap()) {
            JsonObj.Companion.writeEscapedJsonString(buf, key)
            buf.append(":")
            JsonObj.Companion.write(buf, value)
            if (first) first = false
            else buf.append(",")
        }
        buf.append("}")
    }
}
