package janis

import java.io.*
import java.util.*
import java.util.jar.*

class JParser(val classes: LinkedList<JClass>) {

    fun parseJava(path: File) { /* Currently java files aren't parsed. */ }

    fun parseClass(path: File) {
        var buf: BufferedReader? = null
        try {
            buf = path.bufferedReader()
            parseClass(buf)
        } finally {
            buf?.close()
        }
    }

    fun parseJar(path: File) {
        var jar: JarFile? = null
        try {
            jar = JarFile(path)
            for (entry in jar.entries()) {
                if (File(entry.name).extension.lowercase() == "class") {
                    var buf: BufferedReader? = null
                    try {
                        buf = BufferedReader(InputStreamReader(jar.getInputStream(entry)))
                        parseClass(buf)
                    } finally {
                        buf?.close()
                    }
                }
            }

        } finally {
            jar?.close()
        }
    }

    private fun parseClass(buf: BufferedReader) {

    }
}