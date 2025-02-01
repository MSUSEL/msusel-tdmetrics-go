public class EntryPoint {
    public static void main(String[] args) {


        final Options options = new Options();
        options.addOption(new Option("d", "debug", false, "Turn on debug."));
        options.addOption(new Option("e", "extract", false, "Turn on extract."));
        options.addOption(new Option("o", "option", true, "Turn on option with argument."));
    }
}
