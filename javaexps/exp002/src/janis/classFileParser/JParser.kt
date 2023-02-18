package janis.classFileParser

import com.sun.tools.javac.jvm.ClassFile.*
import janis.JClass
import java.io.*
import java.util.*
import java.util.jar.*

class JParser(val classes: LinkedList<JClass>) {

    fun parseClass(path: File) {
        var buf: InputStream? = null
        try {
            buf = path.inputStream()
            this.parseClass(buf)
        } catch (ex: Exception) {
            throw IOException("Error reading Class file: ${path.name}", ex)
        } finally {
            buf?.close()
        }
    }

    fun parseJar(path: File) {
        var jar: JarFile? = null
        try {
            jar = JarFile(path)
            for (entry in jar.entries())
                this.parseJarEntry(jar, entry)
        } catch (ex: Exception) {
            throw IOException("Error reading Jar file: ${path.name}", ex)
        } finally {
            jar?.close()
        }
    }

    private fun parseJarEntry(jar: JarFile, entry: JarEntry) {
        if (!isClass(entry)) return

        var buf: InputStream? = null
        try {
            buf = jar.getInputStream(entry)
            this.parseClass(buf)
        } catch (ex: Exception) {
            throw IOException("Error reading Jar Entity: ${entry.name}", ex)
        } finally {
            buf?.close()
        }
    }

    private fun isClass(entry: JarEntry): Boolean {
        val f = File(entry.name)
        return f.isFile && f.extension.lowercase() == "class"
    }

    /**
     *
     * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4>
     */
    private fun parseClass(buf: InputStream) {
        val dat = DataInputStream(buf)

        val magic: Int = dat.readInt()
        if (magic != JAVA_MAGIC)
            throw IOException("Invalid Class file")

        val minorVersion = dat.readUnsignedShort()
        val majorVersion = dat.readUnsignedShort()
        val constantPool = JConstantPool(dat)


        // TODO: Finish
    }


}