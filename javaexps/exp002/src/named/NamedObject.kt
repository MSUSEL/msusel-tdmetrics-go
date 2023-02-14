package named

import json.Jsonable

interface NamedObject: Jsonable {
    val name: String
}