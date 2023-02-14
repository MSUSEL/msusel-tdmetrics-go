package json

class JsonStr(private val value: String): JsonObj {
    override fun write(buf: StringBuilder) { JsonObj.Companion.writeEscapedJsonString(buf, value) }
}
