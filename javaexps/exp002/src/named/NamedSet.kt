package named

import json.JsonMap
import json.JsonObj
import json.Jsonable
import kotlin.collections.Collection
import kotlin.collections.HashMap

class NamedSet<T: NamedObject>: Collection<T>, Jsonable {
    private val data = HashMap<String, T>()
    
    override val size: Int get() = this.data.size

    override fun isEmpty() = this.data.isEmpty()

    override fun containsAll(elements: Collection<T>) = elements.all { this.contains(it) }

    override fun contains(element: T) = this.data[element.name] == element

    fun contains(name: String) = this.data.containsKey(name)

    operator fun get(name: String) = this.data[name]

    override fun iterator() = this.data.values.iterator()

    fun add(element: T): Boolean {
        if (this.contains(element.name)) return false
        this.data[element.name] = element
        return true
    }

    fun remove(element: T): Boolean {
        if (!contains(element)) return false
        data.remove(element.name)
        return true
    }

    fun remove(name: String): Boolean {
        if (!contains(name)) return false
        data.remove(name)
        return true
    }

    override fun toJson(): JsonObj? {
        val m = JsonMap()
        for ((name, elem) in this.data)
            m[name] = elem.toJson()
        return m
    }
}