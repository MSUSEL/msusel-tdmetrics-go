package jAnalyzer

import json.JsonMap
import json.JsonObj
import json.Jsonable
import kotlin.io.path.Path

class JAnalyzer: Jsonable {
    private val root = JPath(".")

    fun addPath(path: String): JPath {
        var parent = this.root
        for (sub in Path(path)) {
            val part = sub.toString()
            var child = parent.children[part]
            if (child == null) {
                child = JPath(part)
                parent.children.add(child)
            }
            parent = child
        }
        return parent
    }

    override fun toJson(): JsonObj? {
        val m = JsonMap()
        m["root"] = root.toJson()
        return m
    }
}