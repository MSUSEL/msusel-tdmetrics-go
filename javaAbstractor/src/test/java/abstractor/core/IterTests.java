package abstractor.core;

import static org.junit.jupiter.api.Assertions.assertEquals;

import java.util.Iterator;
import java.util.LinkedList;

import org.junit.jupiter.api.Test;

import abstractor.core.iter.ExpandIterator;
import abstractor.core.iter.Iter;
import abstractor.core.iter.PushBackIterator;
import abstractor.core.iter.YieldIterator;

public class IterTests {
    
    @Test
    public void IterEmptyTest() {
        final Iterable<String> it = Iter.<String>Empty();
        assertEquals("", String.join("|", it));
        assertEquals("", String.join("|", it));
    }

    @Test
    public void IterSingleTest() {
        final Iterable<String> it = Iter.<String>Single("Hello");
        assertEquals("Hello", String.join("|", it));
        assertEquals("Hello", String.join("|", it));
    }

    @Test
    public void IterArrayTest() {
        final Iterable<String> it = Iter.<String>Array(new String[]{"Hello", "Mad", "World"});
        assertEquals("Hello|Mad|World", String.join("|", it));
        assertEquals("Hello|Mad|World", String.join("|", it));
    }

    @Test
    public void PushBackIteratorTest() {
        final Iterable<String> it = Iter.<String>Array(new String[]{"Hello", "Mad", "World"});
        String result = "";
        final PushBackIterator<String> pb = new PushBackIterator<String>(it.iterator());
        while (pb.hasNext()) {
            String value = pb.next();
            if (value == "Mad") {
                pb.pushBack("Blue");
                pb.pushBack("Small");
            } else result += "[" + value + "]";
        }
        assertEquals("[Hello][Small][Blue][World]", result);
    }

    @Test
    public void ExpandIteratorTest() {
        final LinkedList<Iterator<String>> input = new LinkedList<Iterator<String>>();
        input.add(Iter.<String>Empty().iterator());
        input.add(Iter.<String>Array(new String[]{"Hello", "World"}).iterator());
        input.add(Iter.<String>Empty().iterator());
        input.add(Iter.<String>Array(new String[]{"Goodbye", "Moon"}).iterator());
        input.add(Iter.<String>Empty().iterator());

        final Iterator<String> it = new ExpandIterator<String>(input.iterator());
        
        String result = "";
        while (it.hasNext()) result += "[" + it.next() + "]";
        assertEquals("[Hello][World][Goodbye][Moon]", result);
    }

    @Test
    public void YieldIteratorTestWithOnlyStep() {
        class Dat { public int countDown = 10; }
        final Dat n = new Dat();
        final Iterator<String> it = new YieldIterator<String>(
            (YieldIterator.Yield<String> y) -> {
                n.countDown--;
                y.yield("before:"+n.countDown);
                if (n.countDown == 6) y.stop();
                y.yield("after:"+n.countDown);
            });

        String result = "";
        while (it.hasNext()) result += "[" + it.next() + "]";
        assertEquals("[before:9][after:9][before:8][after:8][before:7][after:7][before:6]", result);
    }
    
    @Test
    public void YieldIteratorTestWithAll() {
        class Dat { public int countDown = 10; }
        final Dat n = new Dat();
        final Iterator<String> it = new YieldIterator<String>(
            () -> n.countDown > 6,
            (YieldIterator.Yield<String> y) -> {
                n.countDown--;
                y.yield("at:"+n.countDown);
            },
            (YieldIterator.Yield<String> y) -> {
                y.yield("done");
            });

        String result = "";
        while (it.hasNext()) result += "[" + it.next() + "]";
        assertEquals("[at:9][at:8][at:7][at:6][done]", result);
    }
}
