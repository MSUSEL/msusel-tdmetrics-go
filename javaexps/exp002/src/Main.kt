import janis.Janis
import json.JsonObj

fun main() {
    val janis = Janis()
    janis.addPath("C:\\Data\\Code\\Apache\\bcel-6.5.0-src")
    print(JsonObj.Companion.toString(janis.toJson()))
}
