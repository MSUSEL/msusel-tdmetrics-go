package janis

import json.*
import java.io.*
import java.util.*
import java.util.jar.*
import java.util.zip.*

class JFiles(): Jsonable {
    private val extensions = arrayOf("java", "jar", "class", "zip", "war")
    private val files = HashMap<String, TreeSet<File>>()

    init { for (ext in this.extensions) this.files[ext] = TreeSet() }

    fun addPath(path: String) = this.addPath(File(path))

    private fun addPath(path: File) {
        if (path.isDirectory) path.listFiles().forEach { addPath(it) }
        else this.files[path.extension.lowercase()]?.add(path)
    }

    fun readFiles(): TreeSet<JClass> {

    }

    override fun toJson(): JsonObj? {
        val m = JsonMap()
        for (ext in this.extensions) {
            val files = this.files[ext]
            if (!files.isNullOrEmpty()) m[ext] = JsonList(files.map { JsonStr(it.path) })
        }
        return m
    }
}
