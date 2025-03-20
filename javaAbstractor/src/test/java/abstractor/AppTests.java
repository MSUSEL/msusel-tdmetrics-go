package abstractor;

import static org.junit.jupiter.api.Assertions.fail;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Collections;
import java.util.List;

import org.junit.jupiter.api.Test;

public class AppTests {
    
    @Test
    public void test0001() { run("test0001"); }

    static private void run(String testName) {
        final String testPath = "../testData/java/" + testName;
        runCommand("mvn clean compile assembly:single");
        runCommand("java -jar .\\target\\abstractor-0.1-jar-with-dependencies.jar -i " + testPath);
    }

    static private List<String> runCommand(String command) {
        BufferedReader stdInput = null;
        BufferedReader stdError = null;
        try {
            System.out.println(">>|<< "+Paths.get("").toAbsolutePath());


            final Process proc = Runtime.getRuntime().exec(command);
            stdInput = new BufferedReader(new InputStreamReader(proc.getInputStream()));
            stdError = new BufferedReader(new InputStreamReader(proc.getErrorStream()));

            final int result = proc.waitFor();
            if (result == 0) return stdInput.lines().toList();

            dumpLog(stdInput);
            dumpLog(stdError);
            fail("Non-zero exit value: " + result);

        } catch (Exception ex) {
            dumpLog(stdInput);
            dumpLog(stdError);
            fail(ex);
        }
        return Collections.emptyList();
    }

    static private void dumpLog(BufferedReader log) {
        if (log == null) return;
        try {
            String s = log.readLine();
            while (s != null) {
                System.out.println(s);
                s = log.readLine();
            }
        } catch(Exception ex) { } // ignore error
    }
}
