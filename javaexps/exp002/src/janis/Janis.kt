package janis

import json.JsonMap
import json.JsonObj
import json.Jsonable
import kotlin.io.path.Path

class Janis: Jsonable {
    private val files = JFiles()

    fun addPath(path: String) = this.files.addPath(path)

    override fun toJson(): JsonObj? {
        val m = JsonMap()
        m["files"] = this.files.toJson()
        return m
    }
}