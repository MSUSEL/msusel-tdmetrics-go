package abstractor.core;

import spoon.Launcher;
import spoon.reflect.declaration.CtClass;

public class Abstractor {
    private final Config cfg;

    public Abstractor(Config cfg) {
        this.cfg = cfg;
    }

    public void run() {
        try {
            System.out.println("Config:"+this.cfg);
            System.out.println("===========================");


            CtClass l = Launcher.parseClass("class A { void m() { System.out.println(\"yeah\");} }");
            System.out.println(l);


        } catch(Exception ex) {
            System.out.println(ex);
        }
    }
}
