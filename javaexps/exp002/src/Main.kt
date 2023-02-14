import jAnalyzer.JAnalyzer
import json.JsonObj

fun main() {
    var jana = JAnalyzer()
    jana.addPath("/Users/grantnelson/personal/PHD-Work/Checkstyle")
    print(JsonObj.Companion.toString(jana.toJson()))
}
