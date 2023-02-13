package export

import java.util.LinkedList
import java.util.TreeMap

interface EObject

class ENull: EObject

class EBool(val value: Boolean): EObject

class EInt(val value: Int): EObject

class EFloat(val value: Double): EObject

class EString(val value: String): EObject

class EComment(val comment: String, val inner: EObject): EObject

class EMap: TreeMap<String, EObject>(), EObject

class EList: LinkedList<EObject>(), EObject
