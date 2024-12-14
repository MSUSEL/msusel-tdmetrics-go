using Participation;
using System;

namespace UnitTests.ParticipationTests;

public class MatrixTests {

    [Test]
    public void SerializeAndString() {
        Matrix m = new(4, 4);
        checkString(m.Serialize(),
            "0 4x4\n\n\n\n");
        checkString(m.ToString(),
            "[[ -, -, -, - ],",
            " [ -, -, -, - ],",
            " [ -, -, -, - ],",
            " [ -, -, -, - ]]");

        m[0, 0] = 1.2;
        checkString(m.Serialize(),
            "0 4x4\n0:1.2\n\n\n");
        checkString(m.ToString(),
            "[[ 1.2, -, -, - ],",
            " [  -,  -, -, - ],",
            " [  -,  -, -, - ],",
            " [  -,  -, -, - ]]");

        m[1, 2] = 3.4;
        m[2, 1] = 5.6;
        checkString(m.Serialize(),
            "0 4x4\n0:1.2\n2:3.4\n1:5.6\n");
        checkString(m.ToString(),
            "[[ 1.2,  -,   -,  - ],",
            " [  -,   -,  3.4, - ],",
            " [  -,  5.6,  -,  - ],",
            " [  -,   -,   -,  - ]]");

        m[1, 1] = 7.8;
        m[2, 2] = 9.1;
        checkString(m.Serialize(),
            "0 4x4\n0:1.2\n1:7.8 2:3.4\n1:5.6 2:9.1\n");
        checkString(m.ToString(),
            "[[ 1.2,  -,   -,  - ],",
            " [  -,  7.8, 3.4, - ],",
            " [  -,  5.6, 9.1, - ],",
            " [  -,   -,   -,  - ]]");

        m[1, 1] = 11.1;
        m[2, 2] = 22.2;
        m[0, 3] = -0.33;
        checkString(m.Serialize(),
            "0 4x4\n0:1.2 3:-0.33\n1:11.1 2:3.4\n1:5.6 2:22.2\n");
        checkString(m.ToString(),
            "[[ 1.2,   -,    -,  -0.33 ],",
            " [  -,  11.1,  3.4,   -   ],",
            " [  -,   5.6, 22.2,   -   ],",
            " [  -,    -,    -,    -   ]]");

        m[1, 1] = 0.0;
        m[2, 2] = 0.0;
        m[1, 2] = 0.0;
        checkString(m.Serialize(),
            "0 4x4\n0:1.2 3:-0.33\n\n1:5.6\n");
        checkString(m.ToString(),
            "[[ 1.2,  -,  -, -0.33 ],",
            " [  -,   -,  -,   -   ],",
            " [  -,  5.6, -,   -   ],",
            " [  -,   -,  -,   -   ]]");

        m = new(new double[,] {
            { 1.0, 0.0, 0.0, 0.0 },
            { 0.0, 0.0, 1.0, 0.0 },
            { 0.0, 1.0, 0.5, 0.5 },
            { 0.0, 0.0, 0.0, 0.0 } });
        checkString(m.ToString(),
            "[[ 1.0,  -,   -,   -  ],",
            " [  -,   -,  1.0,  -  ],",
            " [  -,  1.0, 0.5, 0.5 ],",
            " [  -,   -,   -,   -  ]]");

        m = new(new double[,] {
            { 1.0, 0.1 },
            { 2.0, 0.2 },
            { 3.0, 0.4 },
            { 4.0, 0.8 } });
        checkString(m.ToString(),
            "[[ 1.0, 0.1 ],",
            " [ 2.0, 0.2 ],",
            " [ 3.0, 0.4 ],",
            " [ 4.0, 0.8 ]]");

        m = new(new double[,] {
            { 123.0, 0.0 },
            { 0.456, 0.0 },
            { 0.0, 123.456 } });
        checkString(m.ToString(),
            "[[ 123.0,      -    ],",
            " [   0.456,    -    ],",
            " [    -,    123.456 ]]");
    }

    [Test]
    public void Deserialize() {
        Matrix m = Matrix.Deserialize(
            "0 4x4\n\n\n\n");
        checkString(m.ToString(),
            "[[ -, -, -, - ],",
            " [ -, -, -, - ],",
            " [ -, -, -, - ],",
            " [ -, -, -, - ]]");

        m = Matrix.Deserialize(
            "0 4x4\n0:1.2\n\n\n");
        checkString(m.ToString(),
            "[[ 1.2, -, -, - ],",
            " [  -,  -, -, - ],",
            " [  -,  -, -, - ],",
            " [  -,  -, -, - ]]");

        m = Matrix.Deserialize(
            "0 4x4\n0:1.2\n2:3.4\n1:5.6\n");
        checkString(m.ToString(),
            "[[ 1.2,  -,   -,  - ],",
            " [  -,   -,  3.4, - ],",
            " [  -,  5.6,  -,  - ],",
            " [  -,   -,   -,  - ]]");

        m = Matrix.Deserialize(
            "0 4x4\n0:1.2\n1:7.8 2:3.4\n1:5.6 2:9.1\n");
        checkString(m.ToString(),
            "[[ 1.2,  -,   -,  - ],",
            " [  -,  7.8, 3.4, - ],",
            " [  -,  5.6, 9.1, - ],",
            " [  -,   -,   -,  - ]]");

        m = Matrix.Deserialize(
            "0 4x4\n0:1.2 3:-0.33\n1:11.1 2:3.4\n1:5.6 2:22.2\n");
        checkString(m.ToString(),
            "[[ 1.2,   -,    -,  -0.33 ],",
            " [  -,  11.1,  3.4,   -   ],",
            " [  -,   5.6, 22.2,   -   ],",
            " [  -,    -,    -,    -   ]]");

        m = Matrix.Deserialize(
            "0 4x4\n0:1.2 3:-0.33\n\n1:5.6\n");
        checkString(m.ToString(),
            "[[ 1.2,  -,  -, -0.33 ],",
            " [  -,   -,  -,   -   ],",
            " [  -,  5.6, -,   -   ],",
            " [  -,   -,  -,   -   ]]");
    }

    [Test]
    public void MatrixMatrixMultiplication() {
        Matrix m1 = new(new double[,] {
            { 1.0, 0.0, 0.0, 0.0 },
            { 0.0, 0.0, 1.0, 0.0 },
            { 0.0, 1.0, 0.5, 0.5 },
            { 0.0, 0.0, 0.0, 0.0 } });

        checkString((m1 * m1).ToString(),
            "[[ 1.0,  -,   -,    -   ],",
            " [  -,  1.0, 0.5,  0.5  ],",
            " [  -,  0.5, 1.25, 0.25 ],",
            " [  -,   -,   -,    -   ]]");

        Matrix m2 = new(new double[,] {
            { 0.0, 0.0, 0.0, 1.0 },
            { 0.0, 0.0, 1.0, 0.0 },
            { 0.0, 1.0, 0.0, 0.0 },
            { 1.0, 0.0, 0.0, 0.0 } });

        checkString((m2 * m2).ToString(),
            "[[ 1.0,  -,   -,   -  ],",
            " [  -,  1.0,  -,   -  ],",
            " [  -,   -,  1.0,  -  ],",
            " [  -,   -,   -,  1.0 ]]");

        checkString((m1 * m2).ToString(),
            "[[  -,   -,   -,  1.0 ],",
            " [  -,  1.0,  -,   -  ],",
            " [ 0.5, 0.5, 1.0,  -  ],",
            " [  -,   -,   -,   -  ]]");

        checkString((m2 * m1).ToString(),
            "[[  -,   -,   -,   -  ],",
            " [  -,  1.0, 0.5, 0.5 ],",
            " [  -,   -,  1.0,  -  ],",
            " [ 1.0,  -,   -,   -  ]]");

        Matrix m3 = new(new double[,] {
            { 1.0, 5.0,  9.0, 13.0 },
            { 2.0, 6.0, 10.0, 14.0 },
            { 3.0, 7.0, 11.0, 15.0 },
            { 4.0, 8.0, 12.0, 16.0 } });

        checkString((m2 * m3).ToString(),
            "[[ 4.0, 8.0, 12.0, 16.0 ],",
            " [ 3.0, 7.0, 11.0, 15.0 ],",
            " [ 2.0, 6.0, 10.0, 14.0 ],",
            " [ 1.0, 5.0,  9.0, 13.0 ]]");

        checkString((m3 * m3).ToString(),
            "[[  90.0, 202.0, 314.0, 426.0 ],",
            " [ 100.0, 228.0, 356.0, 484.0 ],",
            " [ 110.0, 254.0, 398.0, 542.0 ],",
            " [ 120.0, 280.0, 440.0, 600.0 ]]");

        checkString(m3.ToString(),
            "[[ 1.0, 5.0,  9.0, 13.0 ],",
            " [ 2.0, 6.0, 10.0, 14.0 ],",
            " [ 3.0, 7.0, 11.0, 15.0 ],",
            " [ 4.0, 8.0, 12.0, 16.0 ]]");
    }

    [Test]
    public void MatrixVectorMultiplication() {
        Matrix m1 = new(new double[,] {
            { 1.0, 5.0,  9.0, 13.0 },
            { 2.0, 6.0, 10.0, 14.0 },
            { 3.0, 7.0, 11.0, 15.0 },
            { 4.0, 8.0, 12.0, 16.0 } });

        Vector v1 = new([1.0, 1.0, 1.0, 1.0]);
        checkString(v1.ToString(),
            "[[ 1.0 ],",
            " [ 1.0 ],",
            " [ 1.0 ],",
            " [ 1.0 ]]");
        checkString((m1 * v1).ToString(),
            "[[ 28.0 ],",
            " [ 32.0 ],",
            " [ 36.0 ],",
            " [ 40.0 ]]");

        Vector v2 = new([2.0, 0.0, 0.0, 0.0]);
        checkString(v2.ToString(),
            "[[ 2.0 ],",
            " [  -  ],",
            " [  -  ],",
            " [  -  ]]");
        checkString((m1 * v2).ToString(),
            "[[ 2.0 ],",
            " [ 4.0 ],",
            " [ 6.0 ],",
            " [ 8.0 ]]");

        Vector v3 = new([0.0, 2.0, 0.0, 0.0]);
        checkString(v3.ToString(),
            "[[  -  ],",
            " [ 2.0 ],",
            " [  -  ],",
            " [  -  ]]");
        checkString((m1 * v3).ToString(),
            "[[ 10.0 ],",
            " [ 12.0 ],",
            " [ 14.0 ],",
            " [ 16.0 ]]");

        Matrix m2 = new(new double[,] {
            { 0.0, 0.0, 0.0, 0.0 },
            { 0.0, 0.0, 1.0, 0.0 },
            { 0.0, 1.0, 1.0, 0.0 },
            { 0.0, 0.0, 0.0, 0.0 } });
        checkString((m2 * v1).ToString(),
            "[[  -  ],",
            " [ 1.0 ],",
            " [ 2.0 ],",
            " [  -  ]]");
        checkString((m2 * v2).ToString(),
            "[[ - ],",
            " [ - ],",
            " [ - ],",
            " [ - ]]");
        checkString((m2 * v3).ToString(),
            "[[  -  ],",
            " [  -  ],",
            " [ 2.0 ],",
            " [  -  ]]");
    }

    [Test]
    public void MatrixAdditionNegationSubtractionScale() {
        Matrix m1 = new(new double[,] {
            { 1.0, 5.0,  9.0, 13.0 },
            { 2.0, 6.0, 10.0, 14.0 },
            { 3.0, 7.0, 11.0, 15.0 },
            { 4.0, 8.0, 12.0, 16.0 } });
        checkString((m1 + m1).ToString(),
            "[[ 2.0, 10.0, 18.0, 26.0 ],",
            " [ 4.0, 12.0, 20.0, 28.0 ],",
            " [ 6.0, 14.0, 22.0, 30.0 ],",
            " [ 8.0, 16.0, 24.0, 32.0 ]]");
        checkString((m1 - m1).ToString(),
            "[[ -, -, -, - ],",
            " [ -, -, -, - ],",
            " [ -, -, -, - ],",
            " [ -, -, -, - ]]");
        checkString((-m1).ToString(),
            "[[ -1.0, -5.0,  -9.0, -13.0 ],",
            " [ -2.0, -6.0, -10.0, -14.0 ],",
            " [ -3.0, -7.0, -11.0, -15.0 ],",
            " [ -4.0, -8.0, -12.0, -16.0 ]]");
        checkString((m1*2.0).ToString(),
            "[[ 2.0, 10.0, 18.0, 26.0 ],",
            " [ 4.0, 12.0, 20.0, 28.0 ],",
            " [ 6.0, 14.0, 22.0, 30.0 ],",
            " [ 8.0, 16.0, 24.0, 32.0 ]]");
        checkString((2.0*m1).ToString(),
            "[[ 2.0, 10.0, 18.0, 26.0 ],",
            " [ 4.0, 12.0, 20.0, 28.0 ],",
            " [ 6.0, 14.0, 22.0, 30.0 ],",
            " [ 8.0, 16.0, 24.0, 32.0 ]]");

        Matrix m2 = new(new double[,] {
            { 1.0, 1.0, 1.0, 1.0 },
            { 1.0, 1.0, 1.0, 1.0 },
            { 1.0, 1.0, 1.0, 1.0 },
            { 1.0, 1.0, 1.0, 1.0 } });
        checkString((m1 + m2).ToString(),
            "[[ 2.0, 6.0, 10.0, 14.0 ],",
            " [ 3.0, 7.0, 11.0, 15.0 ],",
            " [ 4.0, 8.0, 12.0, 16.0 ],",
            " [ 5.0, 9.0, 13.0, 17.0 ]]");
        checkString((m1 - m2).ToString(),
            "[[  -,  4.0,  8.0, 12.0 ],",
            " [ 1.0, 5.0,  9.0, 13.0 ],",
            " [ 2.0, 6.0, 10.0, 14.0 ],",
            " [ 3.0, 7.0, 11.0, 15.0 ]]");
    }

    private static void checkString(string got, params string[] expLines) {
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
