package jAnalyzer

import json.JsonMap
import json.JsonObj
import json.Jsonable
import named.NamedObject
import named.NamedSet

class JPath(override val name: String): NamedObject, Jsonable {
    val children = NamedSet<JPath>()

    override fun toJson(): JsonObj? {
        val m = JsonMap()
        m["children"] = children.toJson()
        return m
    }
}
