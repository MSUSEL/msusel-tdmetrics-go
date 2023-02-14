package json

import java.util.TreeMap

class JsonMap: TreeMap<String, JsonObject>(), JsonObject
