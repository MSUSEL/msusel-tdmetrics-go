using Runner.Commands;
using System;
using System.Collections.Generic;
using System.Linq;

namespace Runner;

internal class Program {

    static private readonly IReadOnlyDictionary<string, ICommand> Commands =
        new Dictionary<string, ICommand> {
            ["gtcheck"] = new GtCheckCommand(),
        };

    static int Main(string[] args) {
        if (args.Length < 1) {
            printUsage();
            return 1;
        }

        string name = args[0];
        if (!Commands.TryGetValue(name, out ICommand? cmd)) {
            Console.Error.WriteLine("Runner: unknown command '" + name + "'.");
            printUsage();
            return 1;
        }

        cmd.Run(args.Skip(1).ToArray());
        return 0;
    }

    static private void printUsage() {
        Console.Error.WriteLine("usage: Runner <command> [args...]");
        Console.Error.WriteLine("commands:");
        foreach (ICommand cmd in Commands.Values.OrderBy(c => c.Name))
            Console.Error.WriteLine("  " + cmd.Description);
    }
}
