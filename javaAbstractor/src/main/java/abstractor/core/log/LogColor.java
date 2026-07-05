package abstractor.core.log;

/**
 * Colors that can be used for logs.
 *
 * See https://j8ahmed.com/2021/09/13/day-37-learning-ansi-escape-codes/
 * See https://gist.github.com/JBlond/2fea43a3049b38287e5e9cefc87b2124
 */
public enum LogColor {
    black   ("\u001b[30m"),
    red     ("\u001b[31m"),
    green   ("\u001b[32m"),
    yellow  ("\u001b[33m"),
    blue    ("\u001b[34m"),
    magenta ("\u001b[35m"),
    cyan    ("\u001b[36m"),
    white   ("\u001b[37m"),
    reset   ("\u001b[0m"),

    brightBlack   ("\u001b[30;1m"),
    brightRed     ("\u001b[31;1m"),
    brightGreen   ("\u001b[32;1m"),
    brightYellow  ("\u001b[33;1m"),
    brightBlue    ("\u001b[34;1m"),
    brightMagenta ("\u001b[35;1m"),
    brightCyan    ("\u001b[36;1m"),
    brightWhite   ("\u001b[37;1m");

    private final String value;
    LogColor(String value) { this.value = value; }

    public String toString() { return this.value; }
}
