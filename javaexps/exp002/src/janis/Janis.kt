package janis

import json.*
import java.util.*

class Janis: Jsonable {
    private val files = JFiles()
    private val classes = LinkedList<JClass>()

    fun addPath(path: String) = this.files.addPath(path)

    fun build() {
        this.classes.clear()
        val parser = JParser(this.classes)

        // Read all java files.
        this.files.getFiles("java").forEach { parser.parseJava(it) }

        // Read all class files.
        this.files.getFiles("class").forEach { parser.parseClass(it) }

        // Read all jar, zip, and war files as jar files.
        this.files.getFiles("jar").forEach { parser.parseJar(it) }
        this.files.getFiles("zip").forEach { parser.parseJar(it) }
        this.files.getFiles("war").forEach { parser.parseJar(it) }
    }

    override fun toJson(): JsonObj? {
        val m = JsonMap()
        m["files"] = this.files.toJson()
        m["classes"] = JsonList(this.classes.map { it.toJson() })
        return m
    }
}