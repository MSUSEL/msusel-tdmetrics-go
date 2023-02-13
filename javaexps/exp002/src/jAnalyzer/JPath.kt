package jAnalyzer

import named.NamedObject
import named.NamedSet

class JPath(override val name: String): NamedObject {
    val children = NamedSet<JPath>()
}
