using Participation;
using System;

namespace UnitTests.ParticipationTests;

public class MatrixTests {

    [Test]
    public void MatrixTest001() {
        Matrix m = new(4, 4);
        this.checkString(m.Serialize(),
            "0 4x4", "", "", "", "");
        this.checkString(m.ToString(),
            "[[0.0000, 0.0000, 0.0000, 0.0000],",
            " [0.0000, 0.0000, 0.0000, 0.0000],",
            " [0.0000, 0.0000, 0.0000, 0.0000],",
            " [0.0000, 0.0000, 0.0000, 0.0000]]");

        m[0, 0] = 1.2;
        this.checkString(m.Serialize(),
            "0 4x4", "0:1.2000", "", "", "");
        this.checkString(m.ToString(),
            "[[1.2000, 0.0000, 0.0000, 0.0000],",
            " [0.0000, 0.0000, 0.0000, 0.0000],",
            " [0.0000, 0.0000, 0.0000, 0.0000],",
            " [0.0000, 0.0000, 0.0000, 0.0000]]");

        m[1, 2] = 3.4;
        m[2, 1] = 5.6;
        this.checkString(m.Serialize(),
            "0 4x4", "0:1.2000", "2:3.4000", "1:5.6000", "");
        this.checkString(m.ToString(),
            "[[1.2000, 0.0000, 0.0000, 0.0000],",
            " [0.0000, 0.0000, 3.4000, 0.0000],",
            " [0.0000, 5.6000, 0.0000, 0.0000],",
            " [0.0000, 0.0000, 0.0000, 0.0000]]");

        m[1, 1] = 7.8;
        m[2, 2] = 9.1;
        this.checkString(m.Serialize(),
            "0 4x4", "0:1.2000", "1:7.8000 2:3.4000", "1:5.6000 2:9.1000", "");
        this.checkString(m.ToString(),
            "[[1.2000, 0.0000, 0.0000, 0.0000],",
            " [0.0000, 7.8000, 3.4000, 0.0000],",
            " [0.0000, 5.6000, 9.1000, 0.0000],",
            " [0.0000, 0.0000, 0.0000, 0.0000]]");
    }

    private void checkString(string got, params string[] expLines) {
        string exp = string.Join('\n', expLines);
        if (got != exp) {
            Console.WriteLine("Expected:");
            Console.WriteLine(exp);
            Console.WriteLine();
            Console.WriteLine("Gotten:");
            Console.WriteLine(got);
            Assert.That(got, Is.EqualTo(exp));
        }
    }
}
