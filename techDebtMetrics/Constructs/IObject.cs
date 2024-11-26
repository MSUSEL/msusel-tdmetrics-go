using Constructs.Exceptions;

namespace Constructs;

public interface IObject: ITypeDesc {
    public string Name { get; }
    public InterfaceDesc Interface { get; }
    public StructDesc Data { get; }
}
