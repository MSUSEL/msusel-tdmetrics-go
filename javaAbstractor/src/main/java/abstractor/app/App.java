package abstractor.app;

import abstractor.core.Config;
import abstractor.core.Abstractor;

public class App {
    public static void main(String[] args) {
        final Config cfg = new Config();
        if (!cfg.FromArgs(args, null)) return;
        new Abstractor(cfg).run();
    }
}
