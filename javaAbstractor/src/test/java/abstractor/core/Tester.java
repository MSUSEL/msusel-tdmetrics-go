package abstractor.core;

import java.io.*;
import java.nio.charset.StandardCharsets;
import java.nio.file.*;
import java.util.List;

import org.junit.jupiter.api.*;

import abstractor.core.constructs.*;
import abstractor.core.diff.Diff;
import abstractor.core.iter.Iter;
import abstractor.core.json.*;
import abstractor.core.log.Logger;

public class Tester {
    static final int logOutputLinesMax = 2000;
    static final int diffOutputLinesMax = 500;
    
    static public Tester classesFromSource(String ...lines) {
        final Tester t = new Tester(4);
        t.addClassesFromSource(lines);
        return t;
    }

    static private String getTestMethodName(int depth) {
        final StackTraceElement[] ste = Thread.currentThread().getStackTrace();
        return  ste[depth].getClassName() + "." + ste[depth].getMethodName();
    }

    private final ByteArrayOutputStream buffer;
    private final Logger log;
    private final Project proj;
    private final Abstractor ab;

    public Tester() { this(3); }

    public Tester(int depth) {
        this.buffer = new ByteArrayOutputStream();
        PrintStream ps = new PrintStream(buffer);
        this.log = new Logger(true, ps, ps);

        this.log.log("");
        this.log.log("Testing " + getTestMethodName(depth) + "...");

        this.proj = new Project();
        this.ab = new Abstractor(this.log, this.proj);
    }

    public void printLogs() {
        this.printLogs(null);
    }

    public void printLogs(String logFile) {
        printLogs(this.buffer, logFile);
    }

    static public boolean hasFilePath(String file) {
        return file != null && !file.isBlank();
    }

    static public void deleteFile(String file) {
         if (hasFilePath(file)) new File(file).delete();
    }

    static public void printLogs(ByteArrayOutputStream logBuf, String logFile) {
        String logStr = logBuf.toString();
        logBuf.reset();

        // clear out old log file
        deleteFile(logFile);
        
        if (logStr.isBlank()) {
            System.out.println("No logs");
            System.out.flush();
            return;
        }

        if (hasFilePath(logFile)) {
            try {
                BufferedWriter writer = new BufferedWriter(new FileWriter(logFile));
                writer.append(logStr);
                writer.flush();
                writer.close();
                System.out.println("Logs written to "+logFile);
            } catch (Exception ex) {
                Assertions.fail("Failed to write log file", ex);
            }
            return;
        }

        if (logStr.lines().count() < logOutputLinesMax) {
            System.out.println("===[ Logs ]=======================");
            System.out.println(logStr);
            System.out.println("==================================");
        } else {
            System.out.println("Logs too large, check files");
        }
        System.out.flush();
    }

    /**
     * Adds a file containing one or more classes.
     * @param path The path of to the file containing one or more classes.
     */
    public void addClassFromFile(String path) {
        String source = "";
        try {
            final byte[] data = Files.readAllBytes(Paths.get(path));
            source = new String(data, StandardCharsets.UTF_8);
        } catch(Exception ex) {
            System.out.println("failed to read path: " + path + ": " + ex);
        }
        this.addClassesFromSource(source);
    }

    /**
     * Adds one or more classes from the given lines of source code.
     * @param lines The lines of code that contains one or more classes.
     */
    public void addClassesFromSource(String ...lines) {
        try {
            this.ab.prepareClassesFromSource(lines);
            this.ab.performAbstraction();
            this.ab.validate();
        } catch (Exception ex) {
            this.printLogs();
            Assertions.fail(ex);
        }

        if (this.log.errorCount() > 0) {
            this.printLogs();
            Assertions.fail("expected zero errors");
        }
    }

    static private String genJson(Jsonable j) {
        if (j == null) return "null";
        final JsonHelper jh = new JsonHelper();
        return JsonFormat.Relaxed().format(j.toJson(jh));
    }

    public void checkJsonWithFile(Jsonable j, String absFile, String diffFile, String gotFile, String logFile) {
        final String result = genJson(j);
        final String exp = this.formatJsonFromFile(absFile);
        this.assertLines(exp, result, diffFile, gotFile, logFile);
    }
    
    public void checkProjectWithFile(String absFile, String diffFile, String gotFile, String logFile) {
        this.checkJsonWithFile(this.proj, absFile, diffFile, gotFile, logFile);
    }

    public void checkJson(Jsonable j, String ...lines) {
        final String result = genJson(j);
        final String exp = this.formatJson(lines);
        this.assertLines(exp, result, null, null, null);
    }
    
    public void checkProject(String ...lines) {
        this.checkJson(this.proj, lines);
    }

    public void checkConstructCount(String kind, int expCount) {
        final Factory<?> factory = this.proj.getFactory(kind);
        if (factory == null) {
            this.printLogs();
            Assertions.fail("Unable to find kind, " + kind + ", in given project.");
            return;
        }
        if (expCount != factory.size()) {
            Assertions.fail("Factory expected to have " + expCount + " construct(s), but it has " + factory.size() + ".");
            return;
        }
    }

    public void checkConstruct(String key, String ...lines) {
        final Construct con = this.proj.getConstructWithKey(key);
        if (con == null) {
            this.printLogs();
            Assertions.fail("unable to find key, " + key + ", in given project");
            return;
        }
        this.checkJson(con, lines);
    }

    public String formatJson(String ...lines) {
        try {
            return JsonFormat.Relaxed().format(JsonNode.parse(lines));
        } catch(Exception ex) {
            this.printLogs();
            Assertions.fail(ex);
            return "This should be unreachable.";
        }
    }

    public String formatJsonFromFile(String path) {
        try {
            return JsonFormat.Relaxed().format(JsonNode.tryParseFile(path));
        } catch(Exception ex) {
            this.printLogs();
            Assertions.fail(ex);
            return "This should be unreachable.";
        }
    }

    public void assertLines(String exp, String result, String diffFile, String gotFile, String logFile) {
        if (!exp.equals(result)) {
            this.printLogs(logFile);
            printDiff(exp, result, diffFile, gotFile);
            Assertions.fail("unexpected lines (see diff)");
        } else {
            // clear out any old file
            deleteFile(logFile);
            deleteFile(gotFile);
            deleteFile(diffFile);
        }
    }

    static public void printDiff(String exp, String result, String diffFile, String gotFile) {
        final List<String> lines = Iter.ToList(new Diff().PlusMinusByLine(exp, result));

        // clear out any old file
        deleteFile(gotFile);
        deleteFile(diffFile);

        if (hasFilePath(gotFile)) {
            try {
                BufferedWriter writer = new BufferedWriter(new FileWriter(gotFile));
                writer.append(result);
                writer.flush();
                writer.close();
                System.out.println("Unexpected \"got\" lines written to "+gotFile);
            } catch (Exception ex) {
                Assertions.fail("Failed to write \"got\" file", ex);
            }
        }

        if (hasFilePath(diffFile)) {
            try {
                BufferedWriter writer = new BufferedWriter(new FileWriter(diffFile));
                for (String line : lines){
                    writer.append(line);
                    writer.append("\n");
                }
                writer.flush();
                writer.close();
                System.out.println("Unexpected lines diff written to "+diffFile);
            } catch (Exception ex) {
                Assertions.fail("Failed to write diff file", ex);
            }
        }

        if (lines.size() < diffOutputLinesMax) {
            System.out.print("Error: unexpected lines (-exp, +got):");
            for (String line : lines) {
                System.out.print("\n\t");
                System.out.print(line);
            }
            System.out.println();
        } else {
            System.out.println("Error: unexpected lines are too large, check files");
        }
        System.out.flush();
    }
}
