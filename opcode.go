package bscript

import "errors"

var (
	ErrOPCodeUnknow = errors.New("opcode: unknow")
)

type OPCode uint8

var (
	// push value
	OP_0            OPCode = 0x00
	OP_PUSHBYTES_1  OPCode = 0x01
	OP_PUSHBYTES_2  OPCode = 0x02
	OP_PUSHBYTES_3  OPCode = 0x03
	OP_PUSHBYTES_4  OPCode = 0x04
	OP_PUSHBYTES_5  OPCode = 0x05
	OP_PUSHBYTES_6  OPCode = 0x06
	OP_PUSHBYTES_7  OPCode = 0x07
	OP_PUSHBYTES_8  OPCode = 0x08
	OP_PUSHBYTES_9  OPCode = 0x09
	OP_PUSHBYTES_10 OPCode = 0x0a
	OP_PUSHBYTES_11 OPCode = 0x0b
	OP_PUSHBYTES_12 OPCode = 0x0c
	OP_PUSHBYTES_13 OPCode = 0x0d
	OP_PUSHBYTES_14 OPCode = 0x0e
	OP_PUSHBYTES_15 OPCode = 0x0f
	OP_PUSHBYTES_16 OPCode = 0x10
	OP_PUSHBYTES_17 OPCode = 0x11
	OP_PUSHBYTES_18 OPCode = 0x12
	OP_PUSHBYTES_19 OPCode = 0x13
	OP_PUSHBYTES_20 OPCode = 0x14
	OP_PUSHBYTES_21 OPCode = 0x15
	OP_PUSHBYTES_22 OPCode = 0x16
	OP_PUSHBYTES_23 OPCode = 0x17
	OP_PUSHBYTES_24 OPCode = 0x18
	OP_PUSHBYTES_25 OPCode = 0x19
	OP_PUSHBYTES_26 OPCode = 0x1a
	OP_PUSHBYTES_27 OPCode = 0x1b
	OP_PUSHBYTES_28 OPCode = 0x1c
	OP_PUSHBYTES_29 OPCode = 0x1d
	OP_PUSHBYTES_30 OPCode = 0x1e
	OP_PUSHBYTES_31 OPCode = 0x1f
	OP_PUSHBYTES_32 OPCode = 0x20
	OP_PUSHBYTES_33 OPCode = 0x21
	OP_PUSHBYTES_34 OPCode = 0x22
	OP_PUSHBYTES_35 OPCode = 0x23
	OP_PUSHBYTES_36 OPCode = 0x24
	OP_PUSHBYTES_37 OPCode = 0x25
	OP_PUSHBYTES_38 OPCode = 0x26
	OP_PUSHBYTES_39 OPCode = 0x27
	OP_PUSHBYTES_40 OPCode = 0x28
	OP_PUSHBYTES_41 OPCode = 0x29
	OP_PUSHBYTES_42 OPCode = 0x2a
	OP_PUSHBYTES_43 OPCode = 0x2b
	OP_PUSHBYTES_44 OPCode = 0x2c
	OP_PUSHBYTES_45 OPCode = 0x2d
	OP_PUSHBYTES_46 OPCode = 0x2e
	OP_PUSHBYTES_47 OPCode = 0x2f
	OP_PUSHBYTES_48 OPCode = 0x30
	OP_PUSHBYTES_49 OPCode = 0x31
	OP_PUSHBYTES_50 OPCode = 0x32
	OP_PUSHBYTES_51 OPCode = 0x33
	OP_PUSHBYTES_52 OPCode = 0x34
	OP_PUSHBYTES_53 OPCode = 0x35
	OP_PUSHBYTES_54 OPCode = 0x36
	OP_PUSHBYTES_55 OPCode = 0x37
	OP_PUSHBYTES_56 OPCode = 0x38
	OP_PUSHBYTES_57 OPCode = 0x39
	OP_PUSHBYTES_58 OPCode = 0x3a
	OP_PUSHBYTES_59 OPCode = 0x3b
	OP_PUSHBYTES_60 OPCode = 0x3c
	OP_PUSHBYTES_61 OPCode = 0x3d
	OP_PUSHBYTES_62 OPCode = 0x3e
	OP_PUSHBYTES_63 OPCode = 0x3f
	OP_PUSHBYTES_64 OPCode = 0x40
	OP_PUSHBYTES_65 OPCode = 0x41
	OP_PUSHBYTES_66 OPCode = 0x42
	OP_PUSHBYTES_67 OPCode = 0x43
	OP_PUSHBYTES_68 OPCode = 0x44
	OP_PUSHBYTES_69 OPCode = 0x45
	OP_PUSHBYTES_70 OPCode = 0x46
	OP_PUSHBYTES_71 OPCode = 0x47
	OP_PUSHBYTES_72 OPCode = 0x48
	OP_PUSHBYTES_73 OPCode = 0x49
	OP_PUSHBYTES_74 OPCode = 0x4a
	OP_PUSHBYTES_75 OPCode = 0x4b
	OP_PUSHDATA1    OPCode = 0x4c
	OP_PUSHDATA2    OPCode = 0x4d
	OP_PUSHDATA4    OPCode = 0x4e
	OP_1NEGATE      OPCode = 0x4f
	OP_RESERVED     OPCode = 0x50
	OP_1            OPCode = 0x51
	OP_2            OPCode = 0x52
	OP_3            OPCode = 0x53
	OP_4            OPCode = 0x54
	OP_5            OPCode = 0x55
	OP_6            OPCode = 0x56
	OP_7            OPCode = 0x57
	OP_8            OPCode = 0x58
	OP_9            OPCode = 0x59
	OP_10           OPCode = 0x5a
	OP_11           OPCode = 0x5b
	OP_12           OPCode = 0x5c
	OP_13           OPCode = 0x5d
	OP_14           OPCode = 0x5e
	OP_15           OPCode = 0x5f
	OP_16           OPCode = 0x60

	// control
	OP_NOP      OPCode = 0x61
	OP_VER      OPCode = 0x62
	OP_IF       OPCode = 0x63
	OP_NOTIF    OPCode = 0x64
	OP_VERIF    OPCode = 0x65
	OP_VERNOTIF OPCode = 0x66
	OP_ELSE     OPCode = 0x67
	OP_ENDIF    OPCode = 0x68
	OP_VERIFY   OPCode = 0x69
	OP_RETURN   OPCode = 0x6a

	// stack ops
	OP_TOALTSTACK   OPCode = 0x6b
	OP_FROMALTSTACK OPCode = 0x6c
	OP_2DROP        OPCode = 0x6d
	OP_2DUP         OPCode = 0x6e
	OP_3DUP         OPCode = 0x6f
	OP_2OVER        OPCode = 0x70
	OP_2ROT         OPCode = 0x71
	OP_2SWAP        OPCode = 0x72
	OP_IFDUP        OPCode = 0x73
	OP_DEPTH        OPCode = 0x74
	OP_DROP         OPCode = 0x75
	OP_DUP          OPCode = 0x76
	OP_NIP          OPCode = 0x77
	OP_OVER         OPCode = 0x78
	OP_PICK         OPCode = 0x79
	OP_ROLL         OPCode = 0x7a
	OP_ROT          OPCode = 0x7b
	OP_SWAP         OPCode = 0x7c
	OP_TUCK         OPCode = 0x7d

	// splice ops
	OP_CAT    OPCode = 0x7e
	OP_SUBSTR OPCode = 0x7f
	OP_LEFT   OPCode = 0x80
	OP_RIGHT  OPCode = 0x81
	OP_SIZE   OPCode = 0x82

	// bit logic
	OP_INVERT      OPCode = 0x83
	OP_AND         OPCode = 0x84
	OP_OR          OPCode = 0x85
	OP_XOR         OPCode = 0x86
	OP_EQUAL       OPCode = 0x87
	OP_EQUALVERIFY OPCode = 0x88
	OP_RESERVED1   OPCode = 0x89
	OP_RESERVED2   OPCode = 0x8a

	// numeric
	OP_1ADD      OPCode = 0x8b
	OP_1SUB      OPCode = 0x8c
	OP_2MUL      OPCode = 0x8d
	OP_2DIV      OPCode = 0x8e
	OP_NEGATE    OPCode = 0x8f
	OP_ABS       OPCode = 0x90
	OP_NOT       OPCode = 0x91
	OP_0NOTEQUAL OPCode = 0x92

	OP_ADD    OPCode = 0x93
	OP_SUB    OPCode = 0x94
	OP_MUL    OPCode = 0x95
	OP_DIV    OPCode = 0x96
	OP_MOD    OPCode = 0x97
	OP_LSHIFT OPCode = 0x98
	OP_RSHIFT OPCode = 0x99

	OP_BOOLAND            OPCode = 0x9a
	OP_BOOLOR             OPCode = 0x9b
	OP_NUMEQUAL           OPCode = 0x9c
	OP_NUMEQUALVERIFY     OPCode = 0x9d
	OP_NUMNOTEQUAL        OPCode = 0x9e
	OP_LESSTHAN           OPCode = 0x9f
	OP_GREATERTHAN        OPCode = 0xa0
	OP_LESSTHANOREQUAL    OPCode = 0xa1
	OP_GREATERTHANOREQUAL OPCode = 0xa2
	OP_MIN                OPCode = 0xa3
	OP_MAX                OPCode = 0xa4

	OP_WITHIN OPCode = 0xa5

	// crypto
	OP_RIPEMD160           OPCode = 0xa6
	OP_SHA1                OPCode = 0xa7
	OP_SHA256              OPCode = 0xa8
	OP_HASH160             OPCode = 0xa9
	OP_HASH256             OPCode = 0xaa
	OP_CODESEPARATOR       OPCode = 0xab
	OP_CHECKSIG            OPCode = 0xac
	OP_CHECKSIGVERIFY      OPCode = 0xad
	OP_CHECKMULTISIG       OPCode = 0xae
	OP_CHECKMULTISIGVERIFY OPCode = 0xaf

	// expansion
	OP_NOP1                OPCode = 0xb0
	OP_CHECKLOCKTIMEVERIFY OPCode = 0xb1
	//OP_NOP2 = OP_CHECKLOCKTIMEVERIFY
	OP_CHECKSEQUENCEVERIFY OPCode = 0xb2
	//OP_NOP3 = OP_CHECKSEQUENCEVERIFY
	OP_NOP4  OPCode = 0xb3
	OP_NOP5  OPCode = 0xb4
	OP_NOP6  OPCode = 0xb5
	OP_NOP7  OPCode = 0xb6
	OP_NOP8  OPCode = 0xb7
	OP_NOP9  OPCode = 0xb8
	OP_NOP10 OPCode = 0xb9
)

func NewOPCodeFromString(s string) (OPCode, error) {
	switch s {
	case "OP_0":
		return OP_0, nil
	case "OP_PUSHBYTES_1":
		return OP_PUSHBYTES_1, nil
	case "OP_PUSHBYTES_2":
		return OP_PUSHBYTES_2, nil
	case "OP_PUSHBYTES_3":
		return OP_PUSHBYTES_3, nil
	case "OP_PUSHBYTES_4":
		return OP_PUSHBYTES_4, nil
	case "OP_PUSHBYTES_5":
		return OP_PUSHBYTES_5, nil
	case "OP_PUSHBYTES_6":
		return OP_PUSHBYTES_6, nil
	case "OP_PUSHBYTES_7":
		return OP_PUSHBYTES_7, nil
	case "OP_PUSHBYTES_8":
		return OP_PUSHBYTES_8, nil
	case "OP_PUSHBYTES_9":
		return OP_PUSHBYTES_9, nil
	case "OP_PUSHBYTES_10":
		return OP_PUSHBYTES_10, nil
	case "OP_PUSHBYTES_11":
		return OP_PUSHBYTES_11, nil
	case "OP_PUSHBYTES_12":
		return OP_PUSHBYTES_12, nil
	case "OP_PUSHBYTES_13":
		return OP_PUSHBYTES_13, nil
	case "OP_PUSHBYTES_14":
		return OP_PUSHBYTES_14, nil
	case "OP_PUSHBYTES_15":
		return OP_PUSHBYTES_15, nil
	case "OP_PUSHBYTES_16":
		return OP_PUSHBYTES_16, nil
	case "OP_PUSHBYTES_17":
		return OP_PUSHBYTES_17, nil
	case "OP_PUSHBYTES_18":
		return OP_PUSHBYTES_18, nil
	case "OP_PUSHBYTES_19":
		return OP_PUSHBYTES_19, nil
	case "OP_PUSHBYTES_20":
		return OP_PUSHBYTES_20, nil
	case "OP_PUSHBYTES_21":
		return OP_PUSHBYTES_21, nil
	case "OP_PUSHBYTES_22":
		return OP_PUSHBYTES_22, nil
	case "OP_PUSHBYTES_23":
		return OP_PUSHBYTES_23, nil
	case "OP_PUSHBYTES_24":
		return OP_PUSHBYTES_24, nil
	case "OP_PUSHBYTES_25":
		return OP_PUSHBYTES_25, nil
	case "OP_PUSHBYTES_26":
		return OP_PUSHBYTES_26, nil
	case "OP_PUSHBYTES_27":
		return OP_PUSHBYTES_27, nil
	case "OP_PUSHBYTES_28":
		return OP_PUSHBYTES_28, nil
	case "OP_PUSHBYTES_29":
		return OP_PUSHBYTES_29, nil
	case "OP_PUSHBYTES_30":
		return OP_PUSHBYTES_30, nil
	case "OP_PUSHBYTES_31":
		return OP_PUSHBYTES_31, nil
	case "OP_PUSHBYTES_32":
		return OP_PUSHBYTES_32, nil
	case "OP_PUSHBYTES_33":
		return OP_PUSHBYTES_33, nil
	case "OP_PUSHBYTES_34":
		return OP_PUSHBYTES_34, nil
	case "OP_PUSHBYTES_35":
		return OP_PUSHBYTES_35, nil
	case "OP_PUSHBYTES_36":
		return OP_PUSHBYTES_36, nil
	case "OP_PUSHBYTES_37":
		return OP_PUSHBYTES_37, nil
	case "OP_PUSHBYTES_38":
		return OP_PUSHBYTES_38, nil
	case "OP_PUSHBYTES_39":
		return OP_PUSHBYTES_39, nil
	case "OP_PUSHBYTES_40":
		return OP_PUSHBYTES_40, nil
	case "OP_PUSHBYTES_41":
		return OP_PUSHBYTES_41, nil
	case "OP_PUSHBYTES_42":
		return OP_PUSHBYTES_42, nil
	case "OP_PUSHBYTES_43":
		return OP_PUSHBYTES_43, nil
	case "OP_PUSHBYTES_44":
		return OP_PUSHBYTES_44, nil
	case "OP_PUSHBYTES_45":
		return OP_PUSHBYTES_45, nil
	case "OP_PUSHBYTES_46":
		return OP_PUSHBYTES_46, nil
	case "OP_PUSHBYTES_47":
		return OP_PUSHBYTES_47, nil
	case "OP_PUSHBYTES_48":
		return OP_PUSHBYTES_48, nil
	case "OP_PUSHBYTES_49":
		return OP_PUSHBYTES_49, nil
	case "OP_PUSHBYTES_50":
		return OP_PUSHBYTES_50, nil
	case "OP_PUSHBYTES_51":
		return OP_PUSHBYTES_51, nil
	case "OP_PUSHBYTES_52":
		return OP_PUSHBYTES_52, nil
	case "OP_PUSHBYTES_53":
		return OP_PUSHBYTES_53, nil
	case "OP_PUSHBYTES_54":
		return OP_PUSHBYTES_54, nil
	case "OP_PUSHBYTES_55":
		return OP_PUSHBYTES_55, nil
	case "OP_PUSHBYTES_56":
		return OP_PUSHBYTES_56, nil
	case "OP_PUSHBYTES_57":
		return OP_PUSHBYTES_57, nil
	case "OP_PUSHBYTES_58":
		return OP_PUSHBYTES_58, nil
	case "OP_PUSHBYTES_59":
		return OP_PUSHBYTES_59, nil
	case "OP_PUSHBYTES_60":
		return OP_PUSHBYTES_60, nil
	case "OP_PUSHBYTES_61":
		return OP_PUSHBYTES_61, nil
	case "OP_PUSHBYTES_62":
		return OP_PUSHBYTES_62, nil
	case "OP_PUSHBYTES_63":
		return OP_PUSHBYTES_63, nil
	case "OP_PUSHBYTES_64":
		return OP_PUSHBYTES_64, nil
	case "OP_PUSHBYTES_65":
		return OP_PUSHBYTES_65, nil
	case "OP_PUSHBYTES_66":
		return OP_PUSHBYTES_66, nil
	case "OP_PUSHBYTES_67":
		return OP_PUSHBYTES_67, nil
	case "OP_PUSHBYTES_68":
		return OP_PUSHBYTES_68, nil
	case "OP_PUSHBYTES_69":
		return OP_PUSHBYTES_69, nil
	case "OP_PUSHBYTES_70":
		return OP_PUSHBYTES_70, nil
	case "OP_PUSHBYTES_71":
		return OP_PUSHBYTES_71, nil
	case "OP_PUSHBYTES_72":
		return OP_PUSHBYTES_72, nil
	case "OP_PUSHBYTES_73":
		return OP_PUSHBYTES_73, nil
	case "OP_PUSHBYTES_74":
		return OP_PUSHBYTES_74, nil
	case "OP_PUSHBYTES_75":
		return OP_PUSHBYTES_75, nil
	case "OP_PUSHDATA1":
		return OP_PUSHDATA1, nil
	case "OP_PUSHDATA2":
		return OP_PUSHDATA2, nil
	case "OP_PUSHDATA4":
		return OP_PUSHDATA4, nil
	case "OP_1NEGATE":
		return OP_1NEGATE, nil
	case "OP_RESERVED":
		return OP_RESERVED, nil
	case "OP_1":
		return OP_1, nil
	case "OP_2":
		return OP_2, nil
	case "OP_3":
		return OP_3, nil
	case "OP_4":
		return OP_4, nil
	case "OP_5":
		return OP_5, nil
	case "OP_6":
		return OP_6, nil
	case "OP_7":
		return OP_7, nil
	case "OP_8":
		return OP_8, nil
	case "OP_9":
		return OP_9, nil
	case "OP_10":
		return OP_10, nil
	case "OP_11":
		return OP_11, nil
	case "OP_12":
		return OP_12, nil
	case "OP_13":
		return OP_13, nil
	case "OP_14":
		return OP_14, nil
	case "OP_15":
		return OP_15, nil
	case "OP_16":
		return OP_16, nil

	// control
	case "OP_NOP":
		return OP_NOP, nil
	case "OP_VER":
		return OP_VER, nil
	case "OP_IF":
		return OP_IF, nil
	case "OP_NOTIF":
		return OP_NOTIF, nil
	case "OP_VERIF":
		return OP_VERIF, nil
	case "OP_VERNOTIF":
		return OP_VERNOTIF, nil
	case "OP_ELSE":
		return OP_ELSE, nil
	case "OP_ENDIF":
		return OP_ENDIF, nil
	case "OP_VERIFY":
		return OP_VERIFY, nil
	case "OP_RETURN":
		return OP_RETURN, nil

	// stack ops
	case "OP_TOALTSTACK":
		return OP_TOALTSTACK, nil
	case "OP_FROMALTSTACK":
		return OP_FROMALTSTACK, nil
	case "OP_2DROP":
		return OP_2DROP, nil
	case "OP_2DUP":
		return OP_2DUP, nil
	case "OP_3DUP":
		return OP_3DUP, nil
	case "OP_2OVER":
		return OP_2OVER, nil
	case "OP_2ROT":
		return OP_2ROT, nil
	case "OP_2SWAP":
		return OP_2SWAP, nil
	case "OP_IFDUP":
		return OP_IFDUP, nil
	case "OP_DEPTH":
		return OP_DEPTH, nil
	case "OP_DROP":
		return OP_DROP, nil
	case "OP_DUP":
		return OP_DUP, nil
	case "OP_NIP":
		return OP_NIP, nil
	case "OP_OVER":
		return OP_OVER, nil
	case "OP_PICK":
		return OP_PICK, nil
	case "OP_ROLL":
		return OP_ROLL, nil
	case "OP_ROT":
		return OP_ROT, nil
	case "OP_SWAP":
		return OP_SWAP, nil
	case "OP_TUCK":
		return OP_TUCK, nil

	// splice ops
	case "OP_CAT":
		return OP_CAT, nil
	case "OP_SUBSTR":
		return OP_SUBSTR, nil
	case "OP_LEFT":
		return OP_LEFT, nil
	case "OP_RIGHT":
		return OP_RIGHT, nil
	case "OP_SIZE":
		return OP_SIZE, nil

	// bit logic
	case "OP_INVERT":
		return OP_INVERT, nil
	case "OP_AND":
		return OP_AND, nil
	case "OP_OR":
		return OP_OR, nil
	case "OP_XOR":
		return OP_XOR, nil
	case "OP_EQUAL":
		return OP_EQUAL, nil
	case "OP_EQUALVERIFY":
		return OP_EQUALVERIFY, nil
	case "OP_RESERVED1":
		return OP_RESERVED1, nil
	case "OP_RESERVED2":
		return OP_RESERVED2, nil

	// numeric
	case "OP_1ADD":
		return OP_1ADD, nil
	case "OP_1SUB":
		return OP_1SUB, nil
	case "OP_2MUL":
		return OP_2MUL, nil
	case "OP_2DIV":
		return OP_2DIV, nil
	case "OP_NEGATE":
		return OP_NEGATE, nil
	case "OP_ABS":
		return OP_ABS, nil
	case "OP_NOT":
		return OP_NOT, nil
	case "OP_0NOTEQUAL":
		return OP_0NOTEQUAL, nil

	case "OP_ADD":
		return OP_ADD, nil
	case "OP_SUB":
		return OP_SUB, nil
	case "OP_MUL":
		return OP_MUL, nil
	case "OP_DIV":
		return OP_DIV, nil
	case "OP_MOD":
		return OP_MOD, nil
	case "OP_LSHIFT":
		return OP_LSHIFT, nil
	case "OP_RSHIFT":
		return OP_RSHIFT, nil

	case "OP_BOOLAND":
		return OP_BOOLAND, nil
	case "OP_BOOLOR":
		return OP_BOOLOR, nil
	case "OP_NUMEQUAL":
		return OP_NUMEQUAL, nil
	case "OP_NUMEQUALVERIFY":
		return OP_NUMEQUALVERIFY, nil
	case "OP_NUMNOTEQUAL":
		return OP_NUMNOTEQUAL, nil
	case "OP_LESSTHAN":
		return OP_LESSTHAN, nil
	case "OP_GREATERTHAN":
		return OP_GREATERTHAN, nil
	case "OP_LESSTHANOREQUAL":
		return OP_LESSTHANOREQUAL, nil
	case "OP_GREATERTHANOREQUAL":
		return OP_GREATERTHANOREQUAL, nil
	case "OP_MIN":
		return OP_MIN, nil
	case "OP_MAX":
		return OP_MAX, nil

	case "OP_WITHIN":
		return OP_WITHIN, nil

	// crypto
	case "OP_RIPEMD160":
		return OP_RIPEMD160, nil
	case "OP_SHA1":
		return OP_SHA1, nil
	case "OP_SHA256":
		return OP_SHA256, nil
	case "OP_HASH160":
		return OP_HASH160, nil
	case "OP_HASH256":
		return OP_HASH256, nil
	case "OP_CODESEPARATOR":
		return OP_CODESEPARATOR, nil
	case "OP_CHECKSIG":
		return OP_CHECKSIG, nil
	case "OP_CHECKSIGVERIFY":
		return OP_CHECKSIGVERIFY, nil
	case "OP_CHECKMULTISIG":
		return OP_CHECKMULTISIG, nil
	case "OP_CHECKMULTISIGVERIFY":
		return OP_CHECKMULTISIGVERIFY, nil

	// expansion
	case "OP_NOP1":
		return OP_NOP1, nil
	case "OP_CHECKLOCKTIMEVERIFY":
		return OP_CHECKLOCKTIMEVERIFY, nil
	//OP_NOP2 = OP_CHECKLOCKTIMEVERIFY
	case "OP_CHECKSEQUENCEVERIFY":
		return OP_CHECKSEQUENCEVERIFY, nil
	//OP_NOP3 = OP_CHECKSEQUENCEVERIFY
	case "OP_NOP4":
		return OP_NOP4, nil
	case "OP_NOP5":
		return OP_NOP5, nil
	case "OP_NOP6":
		return OP_NOP6, nil
	case "OP_NOP7":
		return OP_NOP7, nil
	case "OP_NOP8":
		return OP_NOP8, nil
	case "OP_NOP9":
		return OP_NOP9, nil
	case "OP_NOP10":
		return OP_NOP10, nil
	}

	return 0, errors.New("unknow opcode")
}

func NewOPCode(n uint8) (OPCode, error) {
	switch n {
	case 0x00:
		return OP_0, nil
	case 0x01:
		return OP_PUSHBYTES_1, nil
	case 0x02:
		return OP_PUSHBYTES_2, nil
	case 0x03:
		return OP_PUSHBYTES_3, nil
	case 0x04:
		return OP_PUSHBYTES_4, nil
	case 0x05:
		return OP_PUSHBYTES_5, nil
	case 0x06:
		return OP_PUSHBYTES_6, nil
	case 0x07:
		return OP_PUSHBYTES_7, nil
	case 0x08:
		return OP_PUSHBYTES_8, nil
	case 0x09:
		return OP_PUSHBYTES_9, nil
	case 0x0a:
		return OP_PUSHBYTES_10, nil
	case 0x0b:
		return OP_PUSHBYTES_11, nil
	case 0x0c:
		return OP_PUSHBYTES_12, nil
	case 0x0d:
		return OP_PUSHBYTES_13, nil
	case 0x0e:
		return OP_PUSHBYTES_14, nil
	case 0x0f:
		return OP_PUSHBYTES_15, nil
	case 0x10:
		return OP_PUSHBYTES_16, nil
	case 0x11:
		return OP_PUSHBYTES_17, nil
	case 0x12:
		return OP_PUSHBYTES_18, nil
	case 0x13:
		return OP_PUSHBYTES_19, nil
	case 0x14:
		return OP_PUSHBYTES_20, nil
	case 0x15:
		return OP_PUSHBYTES_21, nil
	case 0x16:
		return OP_PUSHBYTES_22, nil
	case 0x17:
		return OP_PUSHBYTES_23, nil
	case 0x18:
		return OP_PUSHBYTES_24, nil
	case 0x19:
		return OP_PUSHBYTES_25, nil
	case 0x1a:
		return OP_PUSHBYTES_26, nil
	case 0x1b:
		return OP_PUSHBYTES_27, nil
	case 0x1c:
		return OP_PUSHBYTES_28, nil
	case 0x1d:
		return OP_PUSHBYTES_29, nil
	case 0x1e:
		return OP_PUSHBYTES_30, nil
	case 0x1f:
		return OP_PUSHBYTES_31, nil
	case 0x20:
		return OP_PUSHBYTES_32, nil
	case 0x21:
		return OP_PUSHBYTES_33, nil
	case 0x22:
		return OP_PUSHBYTES_34, nil
	case 0x23:
		return OP_PUSHBYTES_35, nil
	case 0x24:
		return OP_PUSHBYTES_36, nil
	case 0x25:
		return OP_PUSHBYTES_37, nil
	case 0x26:
		return OP_PUSHBYTES_38, nil
	case 0x27:
		return OP_PUSHBYTES_39, nil
	case 0x28:
		return OP_PUSHBYTES_40, nil
	case 0x29:
		return OP_PUSHBYTES_41, nil
	case 0x2a:
		return OP_PUSHBYTES_42, nil
	case 0x2b:
		return OP_PUSHBYTES_43, nil
	case 0x2c:
		return OP_PUSHBYTES_44, nil
	case 0x2d:
		return OP_PUSHBYTES_45, nil
	case 0x2e:
		return OP_PUSHBYTES_46, nil
	case 0x2f:
		return OP_PUSHBYTES_47, nil
	case 0x30:
		return OP_PUSHBYTES_48, nil
	case 0x31:
		return OP_PUSHBYTES_49, nil
	case 0x32:
		return OP_PUSHBYTES_50, nil
	case 0x33:
		return OP_PUSHBYTES_51, nil
	case 0x34:
		return OP_PUSHBYTES_52, nil
	case 0x35:
		return OP_PUSHBYTES_53, nil
	case 0x36:
		return OP_PUSHBYTES_54, nil
	case 0x37:
		return OP_PUSHBYTES_55, nil
	case 0x38:
		return OP_PUSHBYTES_56, nil
	case 0x39:
		return OP_PUSHBYTES_57, nil
	case 0x3a:
		return OP_PUSHBYTES_58, nil
	case 0x3b:
		return OP_PUSHBYTES_59, nil
	case 0x3c:
		return OP_PUSHBYTES_60, nil
	case 0x3d:
		return OP_PUSHBYTES_61, nil
	case 0x3e:
		return OP_PUSHBYTES_62, nil
	case 0x3f:
		return OP_PUSHBYTES_63, nil
	case 0x40:
		return OP_PUSHBYTES_64, nil
	case 0x41:
		return OP_PUSHBYTES_65, nil
	case 0x42:
		return OP_PUSHBYTES_66, nil
	case 0x43:
		return OP_PUSHBYTES_67, nil
	case 0x44:
		return OP_PUSHBYTES_68, nil
	case 0x45:
		return OP_PUSHBYTES_69, nil
	case 0x46:
		return OP_PUSHBYTES_70, nil
	case 0x47:
		return OP_PUSHBYTES_71, nil
	case 0x48:
		return OP_PUSHBYTES_72, nil
	case 0x49:
		return OP_PUSHBYTES_73, nil
	case 0x4a:
		return OP_PUSHBYTES_74, nil
	case 0x4b:
		return OP_PUSHBYTES_75, nil
	case 0x4c:
		return OP_PUSHDATA1, nil
	case 0x4d:
		return OP_PUSHDATA2, nil
	case 0x4e:
		return OP_PUSHDATA4, nil
	case 0x4f:
		return OP_1NEGATE, nil
	case 0x50:
		return OP_RESERVED, nil
	case 0x51:
		return OP_1, nil
	case 0x52:
		return OP_2, nil
	case 0x53:
		return OP_3, nil
	case 0x54:
		return OP_4, nil
	case 0x55:
		return OP_5, nil
	case 0x56:
		return OP_6, nil
	case 0x57:
		return OP_7, nil
	case 0x58:
		return OP_8, nil
	case 0x59:
		return OP_9, nil
	case 0x5a:
		return OP_10, nil
	case 0x5b:
		return OP_11, nil
	case 0x5c:
		return OP_12, nil
	case 0x5d:
		return OP_13, nil
	case 0x5e:
		return OP_14, nil
	case 0x5f:
		return OP_15, nil
	case 0x60:
		return OP_16, nil

	case 0x61:
		return OP_NOP, nil
	case 0x62:
		return OP_VER, nil
	case 0x63:
		return OP_IF, nil
	case 0x64:
		return OP_NOTIF, nil
	case 0x65:
		return OP_VERIF, nil
	case 0x66:
		return OP_VERNOTIF, nil
	case 0x67:
		return OP_ELSE, nil
	case 0x68:
		return OP_ENDIF, nil
	case 0x69:
		return OP_VERIFY, nil
	case 0x6a:
		return OP_RETURN, nil

	case 0x6b:
		return OP_TOALTSTACK, nil
	case 0x6c:
		return OP_FROMALTSTACK, nil
	case 0x6d:
		return OP_2DROP, nil
	case 0x6e:
		return OP_2DUP, nil
	case 0x6f:
		return OP_3DUP, nil
	case 0x70:
		return OP_2OVER, nil
	case 0x71:
		return OP_2ROT, nil
	case 0x72:
		return OP_2SWAP, nil
	case 0x73:
		return OP_IFDUP, nil
	case 0x74:
		return OP_DEPTH, nil
	case 0x75:
		return OP_DROP, nil
	case 0x76:
		return OP_DUP, nil
	case 0x77:
		return OP_NIP, nil
	case 0x78:
		return OP_OVER, nil
	case 0x79:
		return OP_PICK, nil
	case 0x7a:
		return OP_ROLL, nil
	case 0x7b:
		return OP_ROT, nil
	case 0x7c:
		return OP_SWAP, nil
	case 0x7d:
		return OP_TUCK, nil

	case 0x7e:
		return OP_CAT, nil
	case 0x7f:
		return OP_SUBSTR, nil
	case 0x80:
		return OP_LEFT, nil
	case 0x81:
		return OP_RIGHT, nil
	case 0x82:
		return OP_SIZE, nil

	case 0x83:
		return OP_INVERT, nil
	case 0x84:
		return OP_AND, nil
	case 0x85:
		return OP_OR, nil
	case 0x86:
		return OP_XOR, nil
	case 0x87:
		return OP_EQUAL, nil
	case 0x88:
		return OP_EQUALVERIFY, nil
	case 0x89:
		return OP_RESERVED1, nil
	case 0x8a:
		return OP_RESERVED2, nil

	case 0x8b:
		return OP_1ADD, nil
	case 0x8c:
		return OP_1SUB, nil
	case 0x8d:
		return OP_2MUL, nil
	case 0x8e:
		return OP_2DIV, nil
	case 0x8f:
		return OP_NEGATE, nil
	case 0x90:
		return OP_ABS, nil
	case 0x91:
		return OP_NOT, nil
	case 0x92:
		return OP_0NOTEQUAL, nil

	case 0x93:
		return OP_ADD, nil
	case 0x94:
		return OP_SUB, nil
	case 0x95:
		return OP_MUL, nil
	case 0x96:
		return OP_DIV, nil
	case 0x97:
		return OP_MOD, nil
	case 0x98:
		return OP_LSHIFT, nil
	case 0x99:
		return OP_RSHIFT, nil

	case 0x9a:
		return OP_BOOLAND, nil
	case 0x9b:
		return OP_BOOLOR, nil
	case 0x9c:
		return OP_NUMEQUAL, nil
	case 0x9d:
		return OP_NUMEQUALVERIFY, nil
	case 0x9e:
		return OP_NUMNOTEQUAL, nil
	case 0x9f:
		return OP_LESSTHAN, nil
	case 0xa0:
		return OP_GREATERTHAN, nil
	case 0xa1:
		return OP_LESSTHANOREQUAL, nil
	case 0xa2:
		return OP_GREATERTHANOREQUAL, nil
	case 0xa3:
		return OP_MIN, nil
	case 0xa4:
		return OP_MAX, nil
	case 0xa5:
		return OP_WITHIN, nil

	case 0xa6:
		return OP_RIPEMD160, nil
	case 0xa7:
		return OP_SHA1, nil
	case 0xa8:
		return OP_SHA256, nil
	case 0xa9:
		return OP_HASH160, nil
	case 0xaa:
		return OP_HASH256, nil
	case 0xab:
		return OP_CODESEPARATOR, nil
	case 0xac:
		return OP_CHECKSIG, nil
	case 0xad:
		return OP_CHECKSIGVERIFY, nil
	case 0xae:
		return OP_CHECKMULTISIG, nil
	case 0xaf:
		return OP_CHECKMULTISIGVERIFY, nil

	case 0xb0:
		return OP_NOP1, nil
	case 0xb1:
		return OP_CHECKLOCKTIMEVERIFY, nil
	case 0xb2:
		return OP_CHECKSEQUENCEVERIFY, nil
	case 0xb3:
		return OP_NOP4, nil
	case 0xb4:
		return OP_NOP5, nil
	case 0xb5:
		return OP_NOP6, nil
	case 0xb6:
		return OP_NOP7, nil
	case 0xb7:
		return OP_NOP8, nil
	case 0xb8:
		return OP_NOP9, nil
	case 0xb9:
		return OP_NOP10, nil
	default:
		return 0, ErrOPCodeUnknow
	}
}

func (o OPCode) IsDisabled() bool {
	switch o {
	case OP_CAT:
		fallthrough
	case OP_SUBSTR:
		fallthrough
	case OP_LEFT:
		fallthrough
	case OP_RIGHT:
		fallthrough
	case OP_INVERT:
		fallthrough
	case OP_AND:
		fallthrough
	case OP_OR:
		fallthrough
	case OP_XOR:
		fallthrough
	case OP_2MUL:
		fallthrough
	case OP_2DIV:
		fallthrough
	case OP_MUL:
		fallthrough
	case OP_DIV:
		fallthrough
	case OP_MOD:
		fallthrough
	case OP_LSHIFT:
		fallthrough
	case OP_RSHIFT:
		return true
	}

	return false
}

func (o OPCode) String() string {
	switch o {
	case OP_0:
		return "OP_0"
	case OP_PUSHBYTES_1:
		return "OP_PUSHBYTES_1"
	case OP_PUSHBYTES_2:
		return "OP_PUSHBYTES_2"
	case OP_PUSHBYTES_3:
		return "OP_PUSHBYTES_3"
	case OP_PUSHBYTES_4:
		return "OP_PUSHBYTES_4"
	case OP_PUSHBYTES_5:
		return "OP_PUSHBYTES_5"
	case OP_PUSHBYTES_6:
		return "OP_PUSHBYTES_6"
	case OP_PUSHBYTES_7:
		return "OP_PUSHBYTES_7"
	case OP_PUSHBYTES_8:
		return "OP_PUSHBYTES_8"
	case OP_PUSHBYTES_9:
		return "OP_PUSHBYTES_9"
	case OP_PUSHBYTES_10:
		return "OP_PUSHBYTES_10"
	case OP_PUSHBYTES_11:
		return "OP_PUSHBYTES_11"
	case OP_PUSHBYTES_12:
		return "OP_PUSHBYTES_12"
	case OP_PUSHBYTES_13:
		return "OP_PUSHBYTES_13"
	case OP_PUSHBYTES_14:
		return "OP_PUSHBYTES_14"
	case OP_PUSHBYTES_15:
		return "OP_PUSHBYTES_15"
	case OP_PUSHBYTES_16:
		return "OP_PUSHBYTES_16"
	case OP_PUSHBYTES_17:
		return "OP_PUSHBYTES_17"
	case OP_PUSHBYTES_18:
		return "OP_PUSHBYTES_18"
	case OP_PUSHBYTES_19:
		return "OP_PUSHBYTES_19"
	case OP_PUSHBYTES_20:
		return "OP_PUSHBYTES_20"
	case OP_PUSHBYTES_21:
		return "OP_PUSHBYTES_21"
	case OP_PUSHBYTES_22:
		return "OP_PUSHBYTES_22"
	case OP_PUSHBYTES_23:
		return "OP_PUSHBYTES_23"
	case OP_PUSHBYTES_24:
		return "OP_PUSHBYTES_24"
	case OP_PUSHBYTES_25:
		return "OP_PUSHBYTES_25"
	case OP_PUSHBYTES_26:
		return "OP_PUSHBYTES_26"
	case OP_PUSHBYTES_27:
		return "OP_PUSHBYTES_27"
	case OP_PUSHBYTES_28:
		return "OP_PUSHBYTES_28"
	case OP_PUSHBYTES_29:
		return "OP_PUSHBYTES_29"
	case OP_PUSHBYTES_30:
		return "OP_PUSHBYTES_30"
	case OP_PUSHBYTES_31:
		return "OP_PUSHBYTES_31"
	case OP_PUSHBYTES_32:
		return "OP_PUSHBYTES_32"
	case OP_PUSHBYTES_33:
		return "OP_PUSHBYTES_33"
	case OP_PUSHBYTES_34:
		return "OP_PUSHBYTES_34"
	case OP_PUSHBYTES_35:
		return "OP_PUSHBYTES_35"
	case OP_PUSHBYTES_36:
		return "OP_PUSHBYTES_36"
	case OP_PUSHBYTES_37:
		return "OP_PUSHBYTES_37"
	case OP_PUSHBYTES_38:
		return "OP_PUSHBYTES_38"
	case OP_PUSHBYTES_39:
		return "OP_PUSHBYTES_39"
	case OP_PUSHBYTES_40:
		return "OP_PUSHBYTES_40"
	case OP_PUSHBYTES_41:
		return "OP_PUSHBYTES_41"
	case OP_PUSHBYTES_42:
		return "OP_PUSHBYTES_42"
	case OP_PUSHBYTES_43:
		return "OP_PUSHBYTES_43"
	case OP_PUSHBYTES_44:
		return "OP_PUSHBYTES_44"
	case OP_PUSHBYTES_45:
		return "OP_PUSHBYTES_45"
	case OP_PUSHBYTES_46:
		return "OP_PUSHBYTES_46"
	case OP_PUSHBYTES_47:
		return "OP_PUSHBYTES_47"
	case OP_PUSHBYTES_48:
		return "OP_PUSHBYTES_48"
	case OP_PUSHBYTES_49:
		return "OP_PUSHBYTES_49"
	case OP_PUSHBYTES_50:
		return "OP_PUSHBYTES_50"
	case OP_PUSHBYTES_51:
		return "OP_PUSHBYTES_51"
	case OP_PUSHBYTES_52:
		return "OP_PUSHBYTES_52"
	case OP_PUSHBYTES_53:
		return "OP_PUSHBYTES_53"
	case OP_PUSHBYTES_54:
		return "OP_PUSHBYTES_54"
	case OP_PUSHBYTES_55:
		return "OP_PUSHBYTES_55"
	case OP_PUSHBYTES_56:
		return "OP_PUSHBYTES_56"
	case OP_PUSHBYTES_57:
		return "OP_PUSHBYTES_57"
	case OP_PUSHBYTES_58:
		return "OP_PUSHBYTES_58"
	case OP_PUSHBYTES_59:
		return "OP_PUSHBYTES_59"
	case OP_PUSHBYTES_60:
		return "OP_PUSHBYTES_60"
	case OP_PUSHBYTES_61:
		return "OP_PUSHBYTES_61"
	case OP_PUSHBYTES_62:
		return "OP_PUSHBYTES_62"
	case OP_PUSHBYTES_63:
		return "OP_PUSHBYTES_63"
	case OP_PUSHBYTES_64:
		return "OP_PUSHBYTES_64"
	case OP_PUSHBYTES_65:
		return "OP_PUSHBYTES_65"
	case OP_PUSHBYTES_66:
		return "OP_PUSHBYTES_66"
	case OP_PUSHBYTES_67:
		return "OP_PUSHBYTES_67"
	case OP_PUSHBYTES_68:
		return "OP_PUSHBYTES_68"
	case OP_PUSHBYTES_69:
		return "OP_PUSHBYTES_69"
	case OP_PUSHBYTES_70:
		return "OP_PUSHBYTES_70"
	case OP_PUSHBYTES_71:
		return "OP_PUSHBYTES_71"
	case OP_PUSHBYTES_72:
		return "OP_PUSHBYTES_72"
	case OP_PUSHBYTES_73:
		return "OP_PUSHBYTES_73"
	case OP_PUSHBYTES_74:
		return "OP_PUSHBYTES_74"
	case OP_PUSHBYTES_75:
		return "OP_PUSHBYTES_75"
	case OP_PUSHDATA1:
		return "OP_PUSHDATA1"
	case OP_PUSHDATA2:
		return "OP_PUSHDATA2"
	case OP_PUSHDATA4:
		return "OP_PUSHDATA4"
	case OP_1NEGATE:
		return "OP_1NEGATE"
	case OP_RESERVED:
		return "OP_RESERVED"
	case OP_1:
		return "OP_1"
	case OP_2:
		return "OP_2"
	case OP_3:
		return "OP_3"
	case OP_4:
		return "OP_4"
	case OP_5:
		return "OP_5"
	case OP_6:
		return "OP_6"
	case OP_7:
		return "OP_7"
	case OP_8:
		return "OP_8"
	case OP_9:
		return "OP_9"
	case OP_10:
		return "OP_10"
	case OP_11:
		return "OP_11"
	case OP_12:
		return "OP_12"
	case OP_13:
		return "OP_13"
	case OP_14:
		return "OP_14"
	case OP_15:
		return "OP_15"
	case OP_16:
		return "OP_16"

	// control
	case OP_NOP:
		return "OP_NOP"
	case OP_VER:
		return "OP_VER"
	case OP_IF:
		return "OP_IF"
	case OP_NOTIF:
		return "OP_NOTIF"
	case OP_VERIF:
		return "OP_VERIF"
	case OP_VERNOTIF:
		return "OP_VERNOTIF"
	case OP_ELSE:
		return "OP_ELSE"
	case OP_ENDIF:
		return "OP_ENDIF"
	case OP_VERIFY:
		return "OP_VERIFY"
	case OP_RETURN:
		return "OP_RETURN"

	// stack ops
	case OP_TOALTSTACK:
		return "OP_TOALTSTACK"
	case OP_FROMALTSTACK:
		return "OP_FROMALTSTACK"
	case OP_2DROP:
		return "OP_2DROP"
	case OP_2DUP:
		return "OP_2DUP"
	case OP_3DUP:
		return "OP_3DUP"
	case OP_2OVER:
		return "OP_2OVER"
	case OP_2ROT:
		return "OP_2ROT"
	case OP_2SWAP:
		return "OP_2SWAP"
	case OP_IFDUP:
		return "OP_IFDUP"
	case OP_DEPTH:
		return "OP_DEPTH"
	case OP_DROP:
		return "OP_DROP"
	case OP_DUP:
		return "OP_DUP"
	case OP_NIP:
		return "OP_NIP"
	case OP_OVER:
		return "OP_OVER"
	case OP_PICK:
		return "OP_PICK"
	case OP_ROLL:
		return "OP_ROLL"
	case OP_ROT:
		return "OP_ROT"
	case OP_SWAP:
		return "OP_SWAP"
	case OP_TUCK:
		return "OP_TUCK"

	// splice ops
	case OP_CAT:
		return "OP_CAT"
	case OP_SUBSTR:
		return "OP_SUBSTR"
	case OP_LEFT:
		return "OP_LEFT"
	case OP_RIGHT:
		return "OP_RIGHT"
	case OP_SIZE:
		return "OP_SIZE"

	// bit logic
	case OP_INVERT:
		return "OP_INVERT"
	case OP_AND:
		return "OP_AND"
	case OP_OR:
		return "OP_OR"
	case OP_XOR:
		return "OP_XOR"
	case OP_EQUAL:
		return "OP_EQUAL"
	case OP_EQUALVERIFY:
		return "OP_EQUALVERIFY"
	case OP_RESERVED1:
		return "OP_RESERVED1"
	case OP_RESERVED2:
		return "OP_RESERVED2"

	// numeric
	case OP_1ADD:
		return "OP_1ADD"
	case OP_1SUB:
		return "OP_1SUB"
	case OP_2MUL:
		return "OP_2MUL"
	case OP_2DIV:
		return "OP_2DIV"
	case OP_NEGATE:
		return "OP_NEGATE"
	case OP_ABS:
		return "OP_ABS"
	case OP_NOT:
		return "OP_NOT"
	case OP_0NOTEQUAL:
		return "OP_0NOTEQUAL"

	case OP_ADD:
		return "OP_ADD"
	case OP_SUB:
		return "OP_SUB"
	case OP_MUL:
		return "OP_MUL"
	case OP_DIV:
		return "OP_DIV"
	case OP_MOD:
		return "OP_MOD"
	case OP_LSHIFT:
		return "OP_LSHIFT"
	case OP_RSHIFT:
		return "OP_RSHIFT"

	case OP_BOOLAND:
		return "OP_BOOLAND"
	case OP_BOOLOR:
		return "OP_BOOLOR"
	case OP_NUMEQUAL:
		return "OP_NUMEQUAL"
	case OP_NUMEQUALVERIFY:
		return "OP_NUMEQUALVERIFY"
	case OP_NUMNOTEQUAL:
		return "OP_NUMNOTEQUAL"
	case OP_LESSTHAN:
		return "OP_LESSTHAN"
	case OP_GREATERTHAN:
		return "OP_GREATERTHAN"
	case OP_LESSTHANOREQUAL:
		return "OP_LESSTHANOREQUAL"
	case OP_GREATERTHANOREQUAL:
		return "OP_GREATERTHANOREQUAL"
	case OP_MIN:
		return "OP_MIN"
	case OP_MAX:
		return "OP_MAX"

	case OP_WITHIN:
		return "OP_WITHIN"

	// crypto
	case OP_RIPEMD160:
		return "OP_RIPEMD160"
	case OP_SHA1:
		return "OP_SHA1"
	case OP_SHA256:
		return "OP_SHA256"
	case OP_HASH160:
		return "OP_HASH160"
	case OP_HASH256:
		return "OP_HASH256"
	case OP_CODESEPARATOR:
		return "OP_CODESEPARATOR"
	case OP_CHECKSIG:
		return "OP_CHECKSIG"
	case OP_CHECKSIGVERIFY:
		return "OP_CHECKSIGVERIFY"
	case OP_CHECKMULTISIG:
		return "OP_CHECKMULTISIG"
	case OP_CHECKMULTISIGVERIFY:
		return "OP_CHECKMULTISIGVERIFY"

	// expansion
	case OP_NOP1:
		return "OP_NOP1"
	case OP_CHECKLOCKTIMEVERIFY:
		return "OP_CHECKLOCKTIMEVERIFY"
	//OP_NOP2 = OP_CHECKLOCKTIMEVERIFY
	case OP_CHECKSEQUENCEVERIFY:
		return "OP_CHECKSEQUENCEVERIFY"
	//OP_NOP3 = OP_CHECKSEQUENCEVERIFY
	case OP_NOP4:
		return "OP_NOP4"
	case OP_NOP5:
		return "OP_NOP5"
	case OP_NOP6:
		return "OP_NOP6"
	case OP_NOP7:
		return "OP_NOP7"
	case OP_NOP8:
		return "OP_NOP8"
	case OP_NOP9:
		return "OP_NOP9"
	case OP_NOP10:
		return "OP_NOP10"
	}

	return "OP_UNKNOW"
}

func (o OPCode) IsCountable() bool {
	return o > OP_16
}

func (o OPCode) IsPushValue() bool {
	return OP_1NEGATE <= o && o <= OP_16
}
