// Test methods from https://www.baeldung.com/java-cognitive-complexity
public class Main {
    public static void main(String[] args) {
        System.out.println(tennisScore(10));
        System.out.println(countVowels("cat"));
        System.out.println(readFile("file-name"));
    }

    public String tennisScore(int pointsWon) {
        switch (pointsWon) {
            case 0: return "Love";
            case 1: return "Fifteen";
            case 2: return "Thirty";
            case 3: return "Forty";
            default: throw new IllegalArgumentException();
        }
    }

    public int countVowels(String word) {
        int count = 0;
        for (String c : word.split("")) {
            for(String v: vowels) {
                if(c.equalsIgnoreCase(v)) {
                    count++;
                }
            }
        }
        if(count == 0) {
            return "does not contain vowels";
        }
        return "contains %s vowels".formatted(count);
    }

    public String readFile(String path) {
        String text = null;
        if(path == null || path.trim().isEmpty() || !path.endsWith(".txt")) {
            return DEFAULT_TEXT;
        }

        try {
            text = "";
            for (String line: Files.readAllLines(Path.of(path))) {
                if(line.trim().isEmpty()) {
                    continue;
                }
                text+= line;
            }
        } catch (IOException e) {
            if(e instanceof FileNotFoundException) {
                log.error("could not read the file, returning the default content..", e);
            } else {
                throw new RuntimeException(e);
            }
        }
        return text == null ? DEFAULT_TEXT : text;
    }
}
