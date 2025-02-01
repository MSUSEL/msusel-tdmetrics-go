package abstractor.app;

import org.apache.commons.cli.Options;
import org.apache.commons.cli.Option;

import abstractor.core.Config;
import abstractor.core.Abstractor;

public class App {
    public static void main(String[] args) {
        final Options options = new Options();
        options.addOption(new Option("i", "input", false, "The input path"));
        options.addOption(new Option("o", "output", false, "The output path"));


        Config cfg = new Config();




        Abstractor ab = new Abstractor(cfg);
        ab.run();
    }
}
