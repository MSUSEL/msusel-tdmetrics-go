package json

import java.rmi.UnexpectedException

class Yaml {
    private val buf = StringBuilder()

    fun clear() = this.buf.clear()

    fun write(obj: EObject?) {
        this.buf.appendLine("---")
        write(obj ?: ENull(), "")
    }

    private fun write(obj: EObject, indent: String) {
        when(obj) {
            is ENull -> this.buf.appendLine("~")
            is EBool -> this.buf.appendLine(obj.value)
            is EInt -> this.buf.appendLine(obj.value)
            is EFloat -> this.buf.appendLine(obj.value)
            is EString -> this.writeString(obj.value, indent)
            is EComment -> this.writeComment(obj.comment, obj.inner, indent)
            is EMap -> this.writeMap(obj, indent)
            is EList -> this.writeList(obj, indent)
            else -> throw UnexpectedException("unexpected EObject implementation: ${obj::class.simpleName}")
        }
    }

    private fun writeString(value: String, indent: String) =
        value.splitToSequence('\n').joinTo(this.buf, "\\n$indent", "\"", "\"")

    private fun writeComment(comment: String, inner: EObject, indent: String) {
        comment.splitToSequence('\n').joinTo(this.buf, "\\n$indent#", "#")
        write(inner, indent)
    }

    private fun writeMap(value: EMap, indent: String) {

    }

    private fun writeList(value: EList, indent: String) {

    }

    override fun toString() = this.buf.toString()
}