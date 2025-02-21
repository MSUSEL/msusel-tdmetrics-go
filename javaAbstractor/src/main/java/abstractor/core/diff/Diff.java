package abstractor.core.diff;

import java.util.Iterator;
import java.util.LinkedList;
import java.util.List;

import abstractor.core.diff.comparators.*;
import abstractor.core.diff.core.*;
import abstractor.core.diff.core.SubComparator.ReduceResult;
import abstractor.core.diff.hirschberg.Hirschberg;
import abstractor.core.diff.wagner.Wagner;
import abstractor.core.iter.ExpandIterator;
import abstractor.core.iter.Iter;
import abstractor.core.iter.IterableWrapper;
import abstractor.core.iter.YieldIterator;

public class Diff {
    static public final int defaultWagnerThreshold = 500;

    private final Algorithm algorithm;

    public Diff(Algorithm algorithm) { this.algorithm = algorithm; }

    public Diff() { this(-1, defaultWagnerThreshold); }

    public Diff(int length, int size) {
        this(new Hirschberg(length, new Wagner(size)));
    }

    public Iterable<DiffStep> Path(Comparator comp) {
        if (comp == null) Iter.<DiffStep>Empty();
        return new IterableWrapper<DiffStep>(() -> {
            final LinkedList<Iterator<DiffStep>> parts = new LinkedList<Iterator<DiffStep>>();
            SubComparator cont = new SubComparator(new ReverseComparator(new WrapComparator(comp)));

            ReduceResult red = cont.reduce();
            cont = red.sub(); 
            if (red.back() > 0)
                parts.add(Iter.<DiffStep>SingleIterator(
                    DiffStep.Equal(red.back())
                ));

            if (cont.isEndCase()) parts.add(cont.endCase().iterator());
            else parts.add(this.algorithm.diff(cont).iterator());

            if (red.front() > 0)
                parts.add(Iter.<DiffStep>SingleIterator(
                    DiffStep.Equal(red.front())
                ));

            return new Simplifier(new ExpandIterator<DiffStep>(parts.iterator()));
        });
    }
    
    public Iterable<DiffStep> Path(String aSource, String bSource) {
        return Path(new StringComparator(aSource, bSource));
    }
    
    public <T> Iterable<DiffStep> Path(List<T> aSource, List<T> bSource) {
        return Path(new ListComparator<T>(aSource, bSource));
    }
    
    public <T> Iterable<DiffStep> Path(T[] aSource, T[] bSource) {
        return Path(new ArrayComparator<T>(aSource, bSource));
    }

    private interface Getter { Object get(int index); }

    private Iterable<String> PlusMinus(Iterable<DiffStep> path,
        Getter aGetter, Getter bGetter,
        String equalPrefix, String addedPrefix, String removedPrefix) {
        return new IterableWrapper<String>(
            () -> {
                final Iterator<DiffStep> pathIt = path.iterator();
                class Indices {
                    public int aIndex = 0;
                    public int bIndex = 0;
                }
                final Indices n = new Indices();
                return new YieldIterator<String>(
                    () -> pathIt.hasNext(),
                    (YieldIterator.Yield<String> y) -> {
                        final DiffStep step = pathIt.next();
                        switch (step.type()) {
                            case equal:
                                for (int i = step.count()-1; i >= 0; i--) {
                                    y.yield(equalPrefix + aGetter.get(n.aIndex));
                                    n.aIndex++;
                                    n.bIndex++;
                                }
                                break;

                            case added:
                                for (int i = step.count()-1; i >= 0; i--) {
                                    y.yield(addedPrefix + bGetter.get(n.bIndex));
                                    n.bIndex++;
                                }
                                break;

                            case removed:
                                for (int i = step.count()-1; i >= 0; i--) {
                                    y.yield(removedPrefix + aGetter.get(n.aIndex));
                                    n.aIndex++;
                                }
                                break;
                        }
                    });
            });
    }

    private static final String defaultPlusMinusEqualPrefix   = " ";
    private static final String defaultPlusMinusAddedPrefix   = "+";
    private static final String defaultPlusMinusRemovedPrefix = "-";

    public <T> Iterable<String> PlusMinus(List<T> aSource, List<T> bSource,
        String equalPrefix, String addedPrefix, String removedPrefix) {
        return PlusMinus(this.Path(aSource, bSource),
            (int index) -> aSource.get(index),
            (int index) -> bSource.get(index),
            equalPrefix, addedPrefix, removedPrefix);
    }

    public <T> Iterable<String> PlusMinus(List<T> aSource, List<T> bSource) {
        return PlusMinus(aSource, bSource, defaultPlusMinusEqualPrefix,
            defaultPlusMinusAddedPrefix, defaultPlusMinusRemovedPrefix);
    }
    
    public <T> Iterable<String> PlusMinus(T[] aSource, T[] bSource,
        String equalPrefix, String addedPrefix, String removedPrefix) {
        return PlusMinus(this.Path(aSource, bSource),
            (int index) -> aSource[index],
            (int index) -> bSource[index],
            equalPrefix, addedPrefix, removedPrefix);
    }

    public <T> Iterable<String> PlusMinus(T[] aSource, T[] bSource) {
        return PlusMinus(aSource, bSource, defaultPlusMinusEqualPrefix,
            defaultPlusMinusAddedPrefix, defaultPlusMinusRemovedPrefix);
    }
    
    public <T> Iterable<String> PlusMinusByChar(String aSource, String bSource,
        String equalPrefix, String addedPrefix, String removedPrefix) {
        return PlusMinus(this.Path(aSource, bSource),
            (int index) -> aSource.charAt(index),
            (int index) -> bSource.charAt(index),
            equalPrefix, addedPrefix, removedPrefix);
    }

    public <T> Iterable<String> PlusMinusByChar(String aSource, String bSource) {
        return PlusMinusByChar(aSource, bSource, defaultPlusMinusEqualPrefix,
            defaultPlusMinusAddedPrefix, defaultPlusMinusRemovedPrefix);
    }
    
    public <T> Iterable<String> PlusMinusByLine(String aSource, String bSource,
        String equalPrefix, String addedPrefix, String removedPrefix) {
        String[] aLines = aSource.split("\n");
        String[] bLines = bSource.split("\n");
        return PlusMinus(this.Path(aLines, bLines),
            (int index) -> aSource.charAt(index),
            (int index) -> bSource.charAt(index),
            equalPrefix, addedPrefix, removedPrefix);
    }

    public <T> Iterable<String> PlusMinusByLine(String aSource, String bSource) {
        return PlusMinusByLine(aSource, bSource, defaultPlusMinusEqualPrefix,
            defaultPlusMinusAddedPrefix, defaultPlusMinusRemovedPrefix);
    }

    private Iterable<String> Merge(Iterable<DiffStep> path,
        Getter aGetter, Getter bGetter,
        String startChange, String middleChange, String endChange) {
        return new IterableWrapper<String>(
            () -> {
                final Iterator<DiffStep> pathIt = path.iterator();
                class Indices {
                    public int aIndex = 0;
                    public int bIndex = 0;
                    public StepType prevState = StepType.equal;
                }
                final Indices n = new Indices();
                return new YieldIterator<String>(
                    () -> pathIt.hasNext(),
                    (YieldIterator.Yield<String> y) -> {
                        final DiffStep step = pathIt.next();
                        switch (step.type()) {
                            case equal:
                                switch (n.prevState) {
                                    case added:
                                        y.yield(endChange);
                                        break;
                                    case removed:
                                        y.yield(middleChange);
                                        y.yield(endChange);
                                        break;
                                    default: break;
                                }
                                for (int i = step.count() - 1; i >= 0; i--) {
                                    y.yield("" + aGetter.get(n.aIndex));
                                    n.aIndex++;
                                    n.bIndex++;
                                }
                                break;
            
                            case added:
                                switch (n.prevState) {
                                    case equal:
                                        y.yield(startChange);
                                        y.yield(middleChange);
                                        break;
                                    case removed:
                                        y.yield(middleChange);
                                        break;
                                    default: break;
                                }
                                for (int i = step.count() - 1; i >= 0; i--) {
                                    y.yield("" + bGetter.get(n.bIndex));
                                    n.bIndex++;
                                }
                                break;
            
                            case removed:
                                switch (n.prevState) {
                                    case equal:
                                        y.yield(startChange);
                                        break;
                                    case added:
                                        y.yield(endChange);
                                        y.yield(startChange);
                                        break;
                                    default: break;
                                }
                                for (int i = step.count() - 1; i >= 0; i--) {
                                    y.yield("" + aGetter.get(n.aIndex));
                                    n.aIndex++;
                                }
                                break;
                        }
                        n.prevState = step.type();
                    },
                    (YieldIterator.Yield<String> y) -> {
                        switch (n.prevState) {
                            case added:
                                y.yield(endChange);
                                break;
                            case removed:
                                y.yield(middleChange);
                                y.yield(endChange);
                                break;
                            default: break;
                        }
                    });
            });
    }

    private static final String defaultMergeStartChange  = "<<<<<<<<";
    private static final String defaultMergeMiddleChange = "========";
    private static final String defaultMergeEndChange    = ">>>>>>>>";

    public <T> Iterable<String> Merge(List<T> aSource, List<T> bSource,
        String startChange, String middleChange, String endChange) {
        return Merge(this.Path(aSource, bSource),
            (int index) -> aSource.get(index),
            (int index) -> bSource.get(index),
            startChange, middleChange, endChange);
    }

    public <T> Iterable<String> Merge(List<T> aSource, List<T> bSource) {
        return Merge(aSource, bSource, defaultMergeStartChange,
            defaultMergeMiddleChange, defaultMergeEndChange);
    }
    
    public <T> Iterable<String> Merge(T[] aSource, T[] bSource,
        String startChange, String middleChange, String endChange) {
        return Merge(this.Path(aSource, bSource),
            (int index) -> aSource[index],
            (int index) -> bSource[index],
            startChange, middleChange, endChange);
    }

    public <T> Iterable<String> Merge(T[] aSource, T[] bSource) {
        return Merge(aSource, bSource, defaultMergeStartChange,
            defaultMergeMiddleChange, defaultMergeEndChange);
    }
    
    public <T> Iterable<String> MergeByChar(String aSource, String bSource,
        String startChange, String middleChange, String endChange) {
        return Merge(this.Path(aSource, bSource),
            (int index) -> aSource.charAt(index),
            (int index) -> bSource.charAt(index),
            startChange, middleChange, endChange);
    }

    public <T> Iterable<String> MergeByChar(String aSource, String bSource) {
        return MergeByChar(aSource, bSource, defaultMergeStartChange,
            defaultMergeMiddleChange, defaultMergeEndChange);
    }
    
    public <T> Iterable<String> MergeByLine(String aSource, String bSource,
        String startChange, String middleChange, String endChange) {
        String[] aLines = aSource.split("\n");
        String[] bLines = bSource.split("\n");
        return Merge(this.Path(aLines, bLines),
            (int index) -> aSource.charAt(index),
            (int index) -> bSource.charAt(index),
            startChange, middleChange, endChange);
    }

    public <T> Iterable<String> MergeByLine(String aSource, String bSource) {
        return MergeByLine(aSource, bSource, defaultMergeStartChange,
            defaultMergeMiddleChange, defaultMergeEndChange);
    }
}
