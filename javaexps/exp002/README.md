# Experiment 002

_Janis = Java Analysis_

This uses [INRIA's Spoon](https://github.com/INRIA/spoon).

## Getting Test Projects

From this same folder as this README file, run:

    git clone https://github.com/apache/commons-bcel.git ./testProjects/bcel

## Schema

````
Main{
    packages []package;
    data []data;
    methods []method;
}

package {
    name string:
    subs []package.name;
}

data {
    name string;
    
    static bool;
    public bool;
    interface bool;
    external bool;
    
    supers []data.name;
    fields []data.name;
    members []method.name;
}

method {
    name string;
    
    public bool;
    static bool;
    external bool;
    
    receivers []data.name;
    parameters []data.name;
}
```
