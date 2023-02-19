import janis.Janis;

class Main {
    public static void main(String[] args) {
        Janis j = new Janis();
        j.read("./testProjects/bcel/src/main/java/org/apache/bcel");
        j.write();
    }
}
