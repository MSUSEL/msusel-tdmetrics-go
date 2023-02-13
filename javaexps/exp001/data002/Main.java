package javaexps.exp001.data002;

class Main {

    public static void main(final String[] args) {
        Cat c1 = new Cat("patches");
        c1.move(new Location(10, 3));
        System.out.println(c1);
    }    
}
