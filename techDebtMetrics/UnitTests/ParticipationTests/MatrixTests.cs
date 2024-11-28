using Participation;
using System;

namespace UnitTests.ParticipationTests;

public class MatrixTests {

    [Test]
    public void MatrixTest001() {
        Matrix m = new(4, 4);
        m[0, 0] = 1.2;
        m[1, 2] = 3.4;
        m[2, 1] = 5.6;

        string got = m.ToString();
        string exp = "";
        if (got != exp) {
            Console.WriteLine(got);
            Assert.That(got, Is.EqualTo(exp));
        }
    }
}
