﻿using Participation;
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
    public void Multiplication() {
        Matrix m1 = new(new double[,] {
            { 1.0, 0.0, 0.0, 0.0 },
            { 0.0, 0.0, 1.0, 0.0 },
            { 0.0, 1.0, 0.5, 0.5 },
            { 0.0, 0.0, 0.0, 0.0 } });

        Matrix m2 = m1*m1;
        checkString(m2.ToString(),
            "[[ 1.0,  -,   -,    -   ],",
            " [  -,  1.0, 0.5,  0.5  ],",
            " [  -,  0.5, 1.25, 0.25 ],",
            " [  -,   -,   -,    -   ]]");

    }

    [Test]
    public void Addition() {
        // TODO: Finish
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
