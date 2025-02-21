package abstractor.core;

import static org.junit.jupiter.api.Assertions.assertEquals;

import org.junit.jupiter.api.Test;

import abstractor.core.diff.Diff;

public class DiffTests {
    
    private void checkPlusMinus(String[] aSource, String[] bSource, String[] exp) {
        final String result = String.join("|", new Diff().PlusMinus(aSource, bSource));
        final String expStr = String.join("|", exp);
        assertEquals(result, expStr);
    }
    
    @Test
    public void PlusMinus01() {
        checkPlusMinus(
            new String[] { "cat" },
            new String[] { "cat" },
            new String[] { " cat" });
    }

    @Test
    public void PlusMinus02() {
        checkPlusMinus(
            new String[] { "cat" },
            new String[] { "dog" },
            new String[] { "-cat", "+dog" });
    }

    @Test
    public void PlusMinus03() {
        checkPlusMinus(
            new String[] { "A", "G", "T", "A", "C", "G", "C", "A" },
            new String[] { "T", "A", "T", "G", "C" },
            new String[] { "+T", " A", "-G", " T", "-A", "-C", " G", " C", "-A" });
    }

    @Test
    public void PlusMinus04() {
        this.checkPlusMinus(
            new String[] { "cat", "dog" },
            new String[] { "cat", "horse" },
            new String[] { " cat", "-dog", "+horse" });
    }

    @Test
    public void PlusMinus05() {
        this.checkPlusMinus(
            new String[] { "cat", "dog" },
            new String[] { "cat", "horse", "dog" },
            new String[] { " cat", "+horse", " dog" });
    }

    @Test
    public void PlusMinus06() {
        this.checkPlusMinus(
            new String[] { "cat", "dog", "pig" },
            new String[] { "cat", "horse", "dog" },
            new String[] { " cat", "+horse", " dog", "-pig" });
    }

    @Test
    public void PlusMinus07() {
        this.checkPlusMinus(
            new String[] { "Mike", "Ted", "Mark", "Jim" },
            new String[] { "Ted", "Mark", "Bob", "Bill" },
            new String[] { "-Mike", " Ted", " Mark", "-Jim", "+Bob", "+Bill" });
    }

    @Test
    public void PlusMinus08() {
        this.checkPlusMinus(
            new String[] { "k", "i", "t", "t", "e", "n" },
            new String[] { "s", "i", "t", "t", "i", "n", "g" },
            new String[] { "-k", "+s", " i", " t", " t", "-e", "+i", " n", "+g" });
    }

    @Test
    public void PlusMinus09() {
        this.checkPlusMinus(
            new String[] { "s", "a", "t", "u", "r", "d", "a", "y" },
            new String[] { "s", "u", "n", "d", "a", "y" },
            new String[] { " s", "-a", "-t", " u", "-r", "+n", " d", " a", " y" });
    }

    @Test
    public void PlusMinus10() {
        this.checkPlusMinus(
            new String[] { "s", "a", "t", "x", "r", "d", "a", "y" },
            new String[] { "s", "u", "n", "d", "a", "y" },
            new String[] { " s", "-a", "-t", "-x", "-r", "+u", "+n", " d", " a", " y" });
    }
    
    private void checkMerge(String[] aSource, String[] bSource, String[] exp) {
        final String result = String.join("\n", new Diff().Merge(aSource, bSource));
        final String expStr = String.join("\n", exp);
        assertEquals(result, expStr);
    }

    @Test
    public void Merge01() {
        checkMerge(
            new String[] {
                "function A() int {",
                "  return 10",
                "}",
                "",
                "function C() int {",
                "  a := 12",
                "  return a",
                "}" },
            new String[] {
                "function A() int {",
                "  return 10",
                "}",
                "",
                "function B() int {",
                "  return 11",
                "}",
                "",
                "function C() int {",
                "  return 12",
                "}" },
            new String[] {
                "function A() int {",
                "  return 10",
                "}",
                "",
                "<<<<<<<<",
                "========",
                "function B() int {",
                "  return 11",
                "}",
                "",
                ">>>>>>>>",
                "function C() int {",
                "<<<<<<<<",
                "  a := 12",
                "  return a",
                "========",
                "  return 12",
                ">>>>>>>>",
                "}"});
    }
}
