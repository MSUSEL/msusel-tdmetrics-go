package javaexps.exp001.data002;

public class Location {
    public final int X;
    public final int Y;

    public Location(int x, int y) {
        this.X = x;
        this.Y = y;
    }

    public String toString() {
        return "x: " + this.X + ", y:" + this.Y;
    }
}
