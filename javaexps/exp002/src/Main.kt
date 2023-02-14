import json.Yaml
import jAnalyzer.JAnalyzer

fun main() {
    var jana = JAnalyzer()
    jana.addPath("/Users/grantnelson/personal/PHD-Work/Checkstyle")

    val y = Yaml()
    y.write(jana.export())
    print(y)
}
