package janis

import json.JsonMap
import json.JsonObj
import json.Jsonable

/**
 *
 * @see <https://spoon.gforge.inria.fr/>
 */
class Janis: Jsonable {

    fun addPath(path: String) {
    }

    fun build() {

    }

    override fun toJson(): JsonObj? {
        return JsonMap()
    }
}