package json

import kotlin.collections.ArrayList

class JsonList: ArrayList<JsonObj?>(), JsonObj {
    override fun write(buf: StringBuilder) {
        buf.append("[")
        var first = true
        for (elem in this) {
            JsonObj.Companion.write(buf, elem)
            if (first) first = false
            else buf.append(",")
        }
        buf.append("]")
    }
}
