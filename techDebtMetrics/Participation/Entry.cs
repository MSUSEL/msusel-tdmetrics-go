namespace Participation;

/// <summary>An entry in the matrix that is returned while enumerating.</summary>
/// <param name="Row">The zero based row of the entry.</param>
/// <param name="Column">The zero based column of the entry.</param>
/// <param name="Value">The value at this entry.</param>
public readonly record struct Entry(int Row, int Column, double Value);
