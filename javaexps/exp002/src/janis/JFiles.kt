package janis

import json.*
import java.io.File
import java.util.TreeSet

class JFiles(): Jsonable {
    private val javaFiles  = TreeSet<File>()
    private val jarFiles   = TreeSet<File>()
    private val classFiles = TreeSet<File>()
    private val zipFiles   = TreeSet<File>()
    private val warFiles   = TreeSet<File>()

    fun addPath(path: String) = this.addPath(File(path))

    private fun addPath(path: File) {
        if (path.isDirectory) path.listFiles().forEach { addPath(it) }
        else when(path.extension.lowercase()) {
            "java"  -> this.javaFiles.add(path)
            "jar"   -> this.jarFiles.add(path)
            "class" -> this.classFiles.add(path)
            "zip"   -> this.zipFiles.add(path)
            "war"   -> this.warFiles.add(path)
        }
    }

    private fun isFile(f: File, ext: String) = f.isFile && f.extension.lowercase() == ext;

    override fun toJson(): JsonObj? {
        val m = JsonMap()
        if (this.javaFiles.isNotEmpty())  m["java"]  = JsonList(this.javaFiles.map   { JsonStr(it.name) })
        if (this.jarFiles.isNotEmpty())   m["jar"]   = JsonList(this.jarFiles.map   { JsonStr(it.name) })
        if (this.classFiles.isNotEmpty()) m["class"] = JsonList(this.classFiles.map { JsonStr(it.name) })
        if (this.zipFiles.isNotEmpty())   m["zip"]   = JsonList(this.zipFiles.map   { JsonStr(it.name) })
        if (this.warFiles.isNotEmpty())   m["war"]   = JsonList(this.warFiles.map   { JsonStr(it.name) })
        return m
    }
}
