package janis

import json.*
import java.io.*
import java.util.*

enum class ConstantTag(val value: Byte) {
    UTF8(1),
    UNICODE(2),
    INTEGER(3),
    FLOAT(4),
    LONG(5),
    DOUBLE(6),
    CLASS(7),
    STRING(8),
    FIELD_REF(9),
    METHOD_REF(10),
    INTERFACE_METHOD_REF(11),
    NAME_AND_TYPE(12),
    METHOD_HANDLE(15),
    METHOD_TYPE(16),
    INVOKE_DYNAMIC(18);
}

/**
 *
 * @see <https://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html#jvms-4.4>
 */
open class JConstant(val tag: ConstantTag): Jsonable {
    companion object {
        fun readConstantPool(dat: DataInputStream): List<JConstant> {
            val size: Int = dat.readUnsignedShort()
            val constants = LinkedList<JConstant>()
            var i = 0
            while (i < size) {
                val constant = readConstant(dat)
                constants.add(constant)
                if (isLarge(constant.tag)) i++
                i++
            }
            return constants
        }

        private fun isLarge(tag: ConstantTag) = tag == ConstantTag.LONG || tag == ConstantTag.DOUBLE

        private fun findConstantTag(value: Byte): ConstantTag = ConstantTag.values().find { it.value == value } ?:
            throw IOException("Unknown constant value: $value")

        private fun readConstant(dat: DataInputStream): JConstant {
            return when (val tag = findConstantTag(dat.readByte())) {
                ConstantTag.UTF8    -> JConstantValue(tag, dat.readUTF())
                ConstantTag.UNICODE -> TODO()
                ConstantTag.INTEGER -> JConstantValue(tag, dat.readInt())
                ConstantTag.FLOAT   -> JConstantValue(tag, dat.readFloat())
                ConstantTag.LONG    -> JConstantValue(tag, dat.readLong())
                ConstantTag.DOUBLE  -> JConstantValue(tag, dat.readDouble())
                ConstantTag.CLASS   -> JConstantClass(dat.readUnsignedShort() as UShort)
                ConstantTag.STRING  -> JConstantString(dat.readUnsignedShort() as UShort)
                ConstantTag.FIELD_REF, ConstantTag.METHOD_REF, ConstantTag.INTERFACE_METHOD_REF ->
                    JConstantRef(tag, dat.readUnsignedShort() as UShort, dat.readUnsignedShort() as UShort)

                ConstantTag.METHOD_HANDLE ->TODO()
                ConstantTag.METHOD_TYPE ->TODO()
                ConstantTag.NAME_AND_TYPE->TODO()
                ConstantTag.INVOKE_DYNAMIC ->TODO()
            }
        }
    }

    override fun toJson(): JsonObj? {
        TODO("Not yet implemented")
    }
}

/**
 * A constant used to represent a class or interface.
 * @param nameIndex A valid index into the constant pool for an UTF8 constant with the class or interface name.
 * @see <https://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html#jvms-4.4.1>
 */
class JConstantClass(val nameIndex: UShort): JConstant(ConstantTag.CLASS)

/**
 * A constant used to reference a field, method, or interface method.
 * @param tag Must be FIELD_REF, METHOD_REF, or INTERFACE_METHOD_REF
 * @param classIndex A valid index into the constant pool for a class or interface constant.
 *   FIELD_REF can be a class or interface, METHOD_REF must be a class, INTERFACE_METHOD_REF method must be an interface.
 * @param nameAndTypeIndex A valid index into the constant pool for a NAME_AND_TYPE.
 *   FIELD_REF must be a field descriptor, METHOD_REF and INTERFACE_METHOD_REF must be a method descriptor.
 * @see <https://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html#jvms-4.4.2>
 */
class JConstantRef(tag: ConstantTag, val classIndex: UShort, val nameAndTypeIndex: UShort): JConstant(tag)

/**
 * A constant used to represent a string.
 * @param stringIndex a valid index into the constant pool for an UTF8 constant
 *   with the sequence of Unicode for the string.
 * @see <https://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html#jvms-4.4.3>
 */
class JConstantString(val stringIndex: UShort): JConstant(ConstantTag.STRING)

/**
 * A constant used to represent a numerical value or a UTF8 string.
 * @param tag Must be INTEGER, FLOAT, LONG, DOUBLE, or UTF8.
 * @param value The numerical value of this constant.
 * @see <https://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html#jvms-4.4.4>
 * @see <https://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html#jvms-4.4.5>
 * @see <https://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html#jvms-4.4.7>
 */
class JConstantValue(tag: ConstantTag, val value: Any): JConstant(tag)

/**
 *
 *
 * @see <https://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html#jvms-4.4.6>
 */
class JConstantClass(val nameIndex: UShort, val descIndex: UShort): JConstant(ConstantTag.CLASS)

