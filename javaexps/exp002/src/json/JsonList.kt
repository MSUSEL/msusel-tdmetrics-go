package json

class JsonList : ArrayList<JsonObj?>, JsonObj {
    constructor() : super()
    constructor(elements: Iterable<JsonObj?>) : super(elements.map { it })
    constructor(elements: Iterable<Jsonable?>) : super(elements.map { it?.toJson() })
    constructor(elements: Iterable<Boolean>) : super(elements.map { JsonBool(it) })
    constructor(elements: Iterable<String>) : super(elements.map { JsonStr(it) })
    constructor(elements: Iterable<Int>) : super(elements.map { JsonInt(it) })
    constructor(elements: Iterable<Double>) : super(elements.map { JsonNum(it) })

    override fun write(buf: StringBuilder) {
        buf.append("[")
        var first = true
        for (elem in this) {
            if (first) first = false
            else buf.append(",")
            JsonObj.Companion.write(buf, elem)
        }
        buf.append("]")
    }
}
