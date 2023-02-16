package janis

import json.*
import java.io.*
import java.util.*

/**
 * This is a tool for collecting a specific set of files from a path to a java project.
 */
class JFiles(): Jsonable {
    private val extensions = arrayOf("java", "jar", "class", "zip", "war")
    private val files = HashMap<String, TreeSet<File>>()

    init { for (ext in this.extensions) this.files[ext] = TreeSet() }

    /**
     * If the given path is a file, and it is one of the expected extension types
     * then the file will be added otherwise it will be ignored.
     * If the path is a directory then all the files within the directory and
     * child directories recursively will be added if they are the expected extension types.
     */
    fun addPath(path: String) = this.addPath(File(path))

    private fun addPath(path: File) {
        if (path.isDirectory) path.listFiles().forEach { addPath(it) }
        else this.files[path.extension.lowercase()]?.add(path)
    }

    /**
     * Gets all the file of the given extension type. An empty collection is returned
     * if the given extension type isn't one of the expected extension types.
     */
    fun getFiles(ext: String) = this.files[ext] ?: emptyList()

    /** Gets a JSON map with all the files grouped by extension. */
    override fun toJson(): JsonObj? {
        val m = JsonMap()
        for (ext in this.extensions) {
            val files = this.files[ext]
            if (!files.isNullOrEmpty())
                m[ext] = JsonList(files.map { JsonStr(it.path) })
        }
        return m
    }
}
