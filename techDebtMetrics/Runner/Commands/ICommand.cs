namespace Runner.Commands;

/// <summary>A single runnable subcommand of the Runner tool.</summary>
public interface ICommand {
    /// <summary>The name used to select this command from the CLI args.</summary>
    string Name { get; }

    /// <summary>A short description printed by the usage message.</summary>
    string Description { get; }

    /// <summary>Executes the command with its subcommand-specific arguments (the first CLI arg has been stripped).</summary>
    void Run(string[] args);
}
