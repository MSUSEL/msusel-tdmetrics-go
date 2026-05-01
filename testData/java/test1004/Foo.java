package testData.java.test1004;

import java.lang.annotation.*;
import java.util.*;

@Retention(RetentionPolicy.RUNTIME)
@interface MyAnnotation {
    String value() default "";
}

@MyAnnotation("test")
public class Foo {
    public int x;

    public void bar(List<?> items) {
        for (Object item : items) {
            Runnable r = new Runnable() {
                @Override
                public void run() {}
            };
        }
    }

    public void baz(Comparable<? extends Number> comp) {
    }
}
