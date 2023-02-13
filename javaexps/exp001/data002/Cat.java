package javaexps.exp001.data002;

public class Cat {
    private final String name;
    private Location loc;

    public Cat(String name) {
        this.name = name;
        this.loc = new Location(0, 0);
    }

    public String getName() {
        return this.name;
    }

    public Location getLocation() {
        return this.loc;
    }

    public void move(Location loc) {
        this.loc = loc;
    }

    public String toString() {
        return this.name + " @ " + this.loc;
    }
}
