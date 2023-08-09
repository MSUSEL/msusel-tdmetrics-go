import janis.Janis;
import json.JsonMap;

class Main {
    public static void main(String[] args) {
        JsonMap data = Janis.read("./expProjects/bcel/src/main/java/org/apache/bcel");
        data.put("version", "1.0.0");
        data.put("language", "Java");
        System.out.println(data);
    }
}
