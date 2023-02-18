package janis.classFileParser

import com.sun.tools.javac.jvm.ClassFile.*
import java.io.*

/**
 * The constant pool read from a single class file.
 * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.4>
 */
class JConstantPool {
    val constants : Array<JConstant>

    /**
     * Reads the pool count and constants from the given data input stream.
     * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.1>
     */
    constructor(dat: DataInputStream) {
        val size: Int = dat.readUnsignedShort()
        var pool = Array<JConstant>(size) { JConstantUnusable() }
        var i = 1 // indices are [1 .. size)
        while (i < size) {
            val constant = readConstant(dat)
            pool[i] = constant
            if (constant is JConstantLong || constant is JConstantDouble) i++
            i++
        }
        this.constants = pool
    }

    /** Reads a single content from the given input. */
    private fun readConstant(dat: DataInputStream): JConstant {
        return when (val tag = dat.readByte() as Int) {
            CONSTANT_Utf8               -> JConstantUtf8(dat.readUTF())
            CONSTANT_Integer            -> JConstantInteger(dat.readInt())
            CONSTANT_Float              -> JConstantFloat(dat.readFloat())
            CONSTANT_Long               -> JConstantLong(dat.readLong())
            CONSTANT_Double             -> JConstantDouble(dat.readDouble())
            CONSTANT_Class              -> JConstantClass(dat.readUnsignedShort() as UShort)
            CONSTANT_String             -> JConstantString(dat.readUnsignedShort() as UShort)
            CONSTANT_Fieldref           -> JConstantFieldRef(dat.readUnsignedShort() as UShort, dat.readUnsignedShort() as UShort)
            CONSTANT_Methodref          -> JConstantMethodRef(dat.readUnsignedShort() as UShort, dat.readUnsignedShort() as UShort)
            CONSTANT_InterfaceMethodref -> JConstantInterfaceMethodRef(dat.readUnsignedShort() as UShort, dat.readUnsignedShort() as UShort)
            CONSTANT_NameandType        -> JConstantNameAndType(dat.readUnsignedShort() as UShort, dat.readUnsignedShort() as UShort)
            CONSTANT_MethodHandle       -> JConstantMethodHandle(dat.readByte(), dat.readUnsignedShort() as UShort)
            CONSTANT_MethodType         -> JConstantMethodType(dat.readUnsignedShort() as UShort)
            CONSTANT_Dynamic            -> JConstantDynamic(dat.readUnsignedShort() as UShort, dat.readUnsignedShort() as UShort)
            CONSTANT_InvokeDynamic      -> JConstantInvokeDynamic(dat.readUnsignedShort() as UShort, dat.readUnsignedShort() as UShort)
            CONSTANT_Module             -> JConstantModule(dat.readUnsignedShort() as UShort)
            CONSTANT_Package            -> JConstantPackage(dat.readUnsignedShort() as UShort)
            else -> throw IOException("Unknown constant tag: $tag")
        }
    }
}

/** A constant entry in the constant pool. */
interface JConstant

/** A constant entry for unassigned or unusable entries. */
class JConstantUnusable: JConstant

/**
 * A constant Utf8 encoded string.
 * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.4.7>
 */
class JConstantUtf8(val value: String): JConstant

/**
 * A constant integer (32 bit).
 * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.4.4>
 */
class JConstantInteger(val value: Int): JConstant

/**
 * A constant float (32 bit).
 * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.4.4>
 */
class JConstantFloat(val value: Float): JConstant

/**
 * A constant long integer (64 bit).
 * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.4.5>
 */
class JConstantLong(val value: Long): JConstant

/**
 * A constant double float (64 bit).
 * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.4.5>
 */
class JConstantDouble(val value: Double): JConstant

/**
 * A constant used to represent a class or interface.
 * @param nameIndex A valid index into the constant pool for an Utf-8 constant with the class or interface name.
 * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.4.1>
 */
class JConstantClass(val nameIndex: UShort): JConstant

/**
 * A constant used to represent a string.
 * @param stringIndex A valid index into the constant pool for an Utf-8 constant with the sequence of Unicode for the string.
 * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.4.3>
 */
class JConstantString(val stringIndex: UShort): JConstant

/**
 * A constant used to reference a field.
 * @param classIndex A valid index into the constant pool for a class or interface constant.
 * @param nameAndTypeIndex A valid index into the constant pool for a Name-and-type for a field descriptor.
 * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.4.2>
 */
class JConstantFieldRef(val classIndex: UShort, val nameAndTypeIndex: UShort): JConstant

/**
 * A constant used to reference a method.
 * @param classIndex A valid index into the constant pool for a class constant.
 * @param nameAndTypeIndex A valid index into the constant pool for a Name-and-type for a method descriptor.
 * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.4.2>
 */
class JConstantMethodRef(val classIndex: UShort, val nameAndTypeIndex: UShort): JConstant

/**
 * A constant used to reference an interface method.
 * @param classIndex A valid index into the constant pool for a interface constant.
 * @param nameAndTypeIndex A valid index into the constant pool for a Name-and-type for a method descriptor.
 * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.4.2>
 */
class JConstantInterfaceMethodRef(val classIndex: UShort, val nameAndTypeIndex: UShort): JConstant

/**
 * A constant used to represent a field or method without indicating which class or interface it belongs to.
 * @param nameIndex A valid index into the constant pool for an Utf-8 constant with a field or method name.
 * @param descIndex A valid index into the constant pool for an Utf-8 constant with a field or method descriptor.
 * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.4.6>
 */
class JConstantNameAndType(val nameIndex: UShort, val descIndex: UShort): JConstant

/**
 * A constant used to represent a handle to a method.
 * @param refKind A value to indicate the reference kind:
 *   {1: getField, 2: getStatic, 3: putField, 4: putStatic, 5: invokeVirtual,
 *    6: invokeStatic, 7: invokeSpecial, 8: newInvokeSpecial, 9: invokeInterface}
 * @param refIndex A valid index into the constant pool for a constant type depending on the refKind.
 * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.4.8>
 */
class JConstantMethodHandle(val refKind: Byte, val refIndex: UShort): JConstant

/**
 * A constant representing a method type.
 * @param descIndex A valid index into the constant pool for an Utf-8 constant with a method descriptor.
 * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.4.9>
 */
class JConstantMethodType(val descIndex: UShort): JConstant

/**
 * A constant representing a dynamically-computed constant, an arbitrary value that is produced
 *   by invocation of a bootstrap method in the course of an instruction.
 * @param attrIndex A valid index into the bootstrap methods array in the bootstrap method table.
 * @param nameAndTypeIndex A valid index into the constant pool for a Name-and-type for a descriptor.
 * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.4.10>
 */
class JConstantDynamic(val attrIndex: UShort, val nameAndTypeIndex: UShort): JConstant

/**
 * A constant representing a dynamically-computed call site, an instance of java.lang.invoke.CallSite
 *   that is produced by invocation of a bootstrap method in the course of an invokedynamic instruction.
 * @param attrIndex A valid index into the bootstrap methods array in the bootstrap method table.
 * @param nameAndTypeIndex A valid index into the constant pool for a Name-and-type for a descriptor.
 * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.4.10>
 */
class JConstantInvokeDynamic(val attrIndex: UShort, val nameAndTypeIndex: UShort): JConstant

/**
 * A constant representing a module.
 * @param nameIndex A valid index into the constant pool for an Utf-8 constant with the module name.
 * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.4.11>
 */
class JConstantModule(val nameIndex: UShort): JConstant

/**
 * A constant representing a package.
 * @param nameIndex A valid index into the constant pool for an Utf-8 constant with an encoded package name.
 * @see <https://docs.oracle.com/javase/specs/jvms/se19/html/jvms-4.html#jvms-4.4.12>
 */
class JConstantPackage(val nameIndex: UShort): JConstant
