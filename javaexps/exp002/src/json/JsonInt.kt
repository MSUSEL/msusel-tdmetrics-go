package json

class JsonInt(private val value: Int) : JsonObj {
    override fun write(buf: StringBuilder) {
        buf.append(value)
    }
}
