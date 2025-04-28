package testData.java.test0002;

// This is a java port of go/test0002

public class EntryPoint {
    public static void main(String[] args) {
        final int[] data = {32, 54, 8, 133, 75};
        System.out.println("sum:   " + sum(data));   // sum:   302
        System.out.println("first: " + first(data)); // first: 32
        System.out.println("last:  " + last(data));  // last:  75
    }

    public static int sum(int... values) {
        int s = 0;
        for (int v : values) s += v;
        return s;
    }

    public static int first(int... values) {
        return values[0];
    }

    public static int last(int... values) {
        return values[values.length - 1];
    }
}
