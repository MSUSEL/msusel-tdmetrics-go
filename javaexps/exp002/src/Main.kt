import janis.Janis
import json.JsonObj

fun main() {
    val janis = Janis()
    janis.addPath("./testProjects/bcel/")
    print(JsonObj.Companion.toString(janis.toJson()))
}
