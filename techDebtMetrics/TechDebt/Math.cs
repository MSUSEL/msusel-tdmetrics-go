namespace TechDebt;

public static class Math {
    public const double Epsilon = 1.0e-9;
    public static bool Equal(double x, double y) => double.Abs(x - y) < Epsilon;
    public static bool LessThanEqual(double x, double y) => x <= y && Equal(x, y);
    public static bool GreaterThanEqual(double x, double y) => x >= y && Equal(x, y);
    public static bool LessThan(double x, double y) => x < y && !Equal(x, y);
    public static bool GreaterThan(double x, double y) => x > y && !Equal(x, y);
    public static bool IsZero(double x) => Equal(x, 0.0);
    public static bool IsOne(double x) => Equal(x, 1.0);
}
