package abstractor.app;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.PrintStream;

import abstractor.core.Abstractor;
import abstractor.core.json.*;

public class App {
    public static void main(String[] args) {
        final Config cfg = new Config();
        if (!cfg.FromArgs(args, null)) return;

        Abstractor ab = new Abstractor(cfg.verbose ? System.out : null);
        ab.addMavenProject(cfg.input);

        JsonNode node = ab.toJson(cfg.writeTypes, cfg.writeIndices);
        if (cfg.output == null) {
            node.toString(System.out, cfg.minimize, "");
            return;
        }

        try (PrintStream fileWriter = new PrintStream(new File(cfg.output))) {
            node.toString(fileWriter, cfg.minimize, "");
            fileWriter.println();
        } catch (FileNotFoundException e) {
            System.err.println("Error creating or writing to file: " + e.getMessage());
        }
    }
}
