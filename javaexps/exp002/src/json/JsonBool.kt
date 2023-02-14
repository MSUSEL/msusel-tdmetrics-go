package json

class JsonBool(private val value: Boolean): JsonObj {
    override fun write(buf: StringBuilder) { buf.append(value) }
}
