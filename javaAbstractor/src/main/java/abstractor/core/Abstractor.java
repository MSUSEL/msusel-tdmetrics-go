package abstractor.core;

public class Abstractor {
    private final Config cfg;

    public Abstractor(Config cfg) {
        this.cfg = cfg;
    }

    public void run() {
        System.out.println(this.cfg);
    }
}
