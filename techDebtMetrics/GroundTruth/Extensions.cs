using System.ComponentModel;
using System.Reflection;

namespace GroundTruth;

public static class DeclTypeExt {
    private static readonly Lazy<Dictionary<string, DeclType>> declTypeByName = new(() => {
        DeclType[] types = Enum.GetValues<DeclType>();
        Dictionary<string, DeclType> result = [];
        foreach(DeclType dt in types) result[dt.ToName()] = dt;
        return result;
    });

    public static DeclType FromName(string name) =>
        declTypeByName.Value.TryGetValue(name, out DeclType value) ? value : //DeclType.Unknown;
        throw new Exception("Failed to find DeclType for " + name + ".");

    public static string ToName(this Enum value) {
        MemberInfo[] memberInfo = value.GetType().GetMember(value.ToString());
        if (memberInfo.Length <= 0) return value.ToString();
        object[] attr = memberInfo[0].GetCustomAttributes(typeof(DescriptionAttribute), false);
        if (attr.Length <= 0) return value.ToString(); 
        return ((DescriptionAttribute)attr[0]).Description;
    }
}
