package json

class JsonNum(private val value: Double): JsonObj {
    override fun write(buf: StringBuilder) { buf.append(value) }
}
