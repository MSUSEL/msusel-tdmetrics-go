package janis

import json.*
import java.util.*

class JClass(
    val name: String,
    val path: String,
    var packageName: String = "",
    var isAbstract: Boolean = false): Jsonable {
    val imports = TreeSet<String>()

    override fun toJson(): JsonObj? {
        val m = JsonMap()
        m["name"] = JsonStr(this.name)
        m["path"] = JsonStr(this.path)
        m["isAbstract"] = JsonBool(this.isAbstract)
        m["package"] = JsonStr(this.packageName)
        m["imports"] = JsonList(this.imports.map { JsonStr(it) })
        return m
    }
}
