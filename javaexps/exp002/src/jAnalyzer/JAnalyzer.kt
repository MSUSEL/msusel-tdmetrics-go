package jAnalyzer

import json.EObject
import json.Exportable
import kotlin.io.path.Path

class JAnalyzer: Exportable {
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

    override fun export(): EObject? {
        TODO("Not yet implemented")
    }
}