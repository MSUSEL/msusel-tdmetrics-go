package abstractor.app;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.PrintStream;

import abstractor.core.*;
import abstractor.core.constructs.Project;
import abstractor.core.json.*;
import abstractor.core.log.*;

public class App {
    
    static public void main(String[] args) throws Exception {
        final Config cfg = new Config();
        if (cfg.FromArgs(args, null)) run(cfg);
    }

    static public boolean run(Config cfg) throws Exception {
        Logger log = new Logger(cfg.verbose, cfg.logOut, cfg.logErr);
        Project proj = new Project();
        Abstractor ab = new Abstractor(log, proj);
        ab.addMavenProject(cfg.input);
        ab.finish();

        JsonHelper h = new JsonHelper();
        h.writeKinds   = cfg.writeKinds;
        h.writeIndices = cfg.writeIndices;
        h.rootPath     = cfg.input + "/";
        JsonNode node = proj.toJson(h);

        JsonFormat fmt = cfg.minimize ? JsonFormat.Minimize() : cfg.defaultFormat;

        if (cfg.output == null) {
            fmt.format(cfg.defaultOut, node, "");
            cfg.defaultOut.println();
        } else {
            try (PrintStream fileWriter = new PrintStream(new File(cfg.output))) {
                fmt.format(fileWriter, node, "");
                fileWriter.println();
            } catch (FileNotFoundException e) {
                log.error("Error creating or writing to file: " + e.getMessage());
            }
        }
 
        if (log.errorCount() > 0) {
            final PrintStream errOut = cfg.logErr != null ? cfg.logErr : System.err;
            errOut.println("Had " + log.errorCount() + " errors");
            return false;
        }

        final PrintStream okOut = cfg.logOut != null ? cfg.logOut : System.out;
        okOut.println("Success");
        return true;
    }
}
