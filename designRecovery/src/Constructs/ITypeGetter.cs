namespace designRecovery.src.Constructs;

internal interface ITypeGetter {
    T Get<T>(uint index) where T : ITypeDesc;
}
