import janis.Janis
import json.JsonObj

fun main() {
    val janis = Janis()
    janis.addPath("./testProjects/bcel/")
    janis.build()
    print(JsonObj.Companion.toString(janis.toJson()))
}
