package json

import java.util.*

class JsonMap : TreeMap<String, JsonObj?>, JsonObj {
    constructor(data: Map<String, JsonObj?> = emptyMap()) : super() {
        this.putAll(data)
    }

    constructor(data: Map<String, Jsonable?>) : this(convertValue(data) { it?.toJson() })
    constructor(data: Map<String, Boolean>) : this(convertValue(data) { JsonBool(it) })
    constructor(data: Map<String, String>) : this(convertValue(data) { JsonStr(it) })
    constructor(data: Map<String, Int>) : this(convertValue(data) { JsonInt(it) })
    constructor(data: Map<String, Double>) : this(convertValue(data) { JsonNum(it) })

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

    companion object {
        private fun <T, U, V> convertValue(data: Map<T, U>, valueTransform: (U) -> V): Map<T, V> {
            val result = TreeMap<T, V>()
            for ((key, value) in data) result[key] = valueTransform(value)
            return result
        }
    }

}
