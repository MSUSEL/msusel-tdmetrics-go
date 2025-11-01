using SCG = System.Collections.Generic;

namespace TechDebt;

public class Project {

    public readonly SCG.SortedSet<Class> Classes = [];

    public readonly SCG.SortedSet<Method> Methods = [];


}
