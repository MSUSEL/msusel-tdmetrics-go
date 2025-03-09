package abstractor.app;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.PrintStream;

import abstractor.core.*;
import abstractor.core.constructs.Project;
import abstractor.core.json.*;
import abstractor.core.log.*;

public class App {
    
    public static void main(String[] args) {
        final Config cfg = new Config();
        if (!cfg.FromArgs(args, null)) return;

        Logger log = new Logger(cfg.verbose);
        Project proj = new Project();
        Abstractor ab = new Abstractor(log, proj);
        ab.addMavenProject(cfg.input);

        JsonHelper h = new JsonHelper();
        h.writeKinds = cfg.writeKinds;
        h.writeIndices = cfg.writeIndices;
        JsonNode node = proj.toJson(h);

        JsonFormat fmt = cfg.minimize ? JsonFormat.Minimize() : JsonFormat.Normal();

        if (cfg.output == null) {
            fmt.format(System.out, node, "");
            System.out.println();
        } else {
            try (PrintStream fileWriter = new PrintStream(new File(cfg.output))) {
                fmt.format(fileWriter, node, "");
                fileWriter.println();
            } catch (FileNotFoundException e) {
                log.error("Error creating or writing to file: " + e.getMessage());
            }
        }

        if (log.errorCount() > 0)
            System.err.println("Had " + log.errorCount() + " errors");
        else System.out.println("Success");
    }
}
