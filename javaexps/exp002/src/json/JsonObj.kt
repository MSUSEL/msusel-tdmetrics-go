package json

interface JsonObj {
    fun write(buf: StringBuilder)

    companion object {
        fun writeEscapedJsonString(buf: StringBuilder, value: String) {
            buf.append("\"")
            for (c in value) {
                buf.append(
                    when (c) {
                        '\\' -> "\\\\"
                        '\t' -> "\\t"
                        '\n' -> "\\n"
                        '\r' -> "\\r"
                        '\b' -> "\\b"
                        '"' -> "\\\""
                        else -> c
                    }
                )
            }
            buf.append("\"")
        }

        fun write(buf: StringBuilder, obj: JsonObj?) {
            if (obj == null) buf.append("null")
            else obj.write(buf)
        }

        fun toString(obj: JsonObj?): String {
            val buf = StringBuilder()
            write(buf, obj)
            return buf.toString()
        }
    }
}
