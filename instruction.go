package bscript

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInstructionReserved = errors.New("instruction: reserved")
)

type Instruction struct {
	OPCode OPCode
	Step   int
	Data   []byte
}

func (ins *Instruction) IsConditional() bool {
	if ins.OPCode == OP_IF || ins.OPCode == OP_NOTIF || ins.OPCode == OP_ELSE ||
		ins.OPCode == OP_ENDIF {
		return true
	}

	return false
}

func (ins *Instruction) String() string {
	rv := make([]string, 0, 36)
	rv = append(rv, fmt.Sprintf("%-16s", ins.OPCode))

	switch ins.OPCode {
	case OP_PUSHDATA1:
		rv = append(rv, fmt.Sprintf("0x%02x", len(ins.Data)))
	case OP_PUSHDATA2:
		buf := make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, uint16(len(ins.Data)))
		rv = append(rv, fmt.Sprintf("0x%02x", buf[0]))
		rv = append(rv, fmt.Sprintf("0x%02x", buf[1]))
	case OP_PUSHDATA4:
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(len(ins.Data)))
		rv = append(rv, fmt.Sprintf("0x%02x", buf[0]))
		rv = append(rv, fmt.Sprintf("0x%02x", buf[1]))
		rv = append(rv, fmt.Sprintf("0x%02x", buf[2]))
		rv = append(rv, fmt.Sprintf("0x%02x", buf[3]))
	}

	for i := 0; i < len(ins.Data); i++ {
		rv = append(rv, fmt.Sprintf("0x%02x", ins.Data[i]))
	}

	return strings.Join(rv, " ")
}

type Operator func(*InterpreterContext) error

var instructionOperator = map[OPCode]Operator{
	OP_0:            instructionPushOP0,
	OP_PUSHBYTES_1:  instructionPushOPBytes,
	OP_PUSHBYTES_2:  instructionPushOPBytes,
	OP_PUSHBYTES_3:  instructionPushOPBytes,
	OP_PUSHBYTES_4:  instructionPushOPBytes,
	OP_PUSHBYTES_5:  instructionPushOPBytes,
	OP_PUSHBYTES_6:  instructionPushOPBytes,
	OP_PUSHBYTES_7:  instructionPushOPBytes,
	OP_PUSHBYTES_8:  instructionPushOPBytes,
	OP_PUSHBYTES_9:  instructionPushOPBytes,
	OP_PUSHBYTES_10: instructionPushOPBytes,
	OP_PUSHBYTES_11: instructionPushOPBytes,
	OP_PUSHBYTES_12: instructionPushOPBytes,
	OP_PUSHBYTES_13: instructionPushOPBytes,
	OP_PUSHBYTES_14: instructionPushOPBytes,
	OP_PUSHBYTES_15: instructionPushOPBytes,
	OP_PUSHBYTES_16: instructionPushOPBytes,
	OP_PUSHBYTES_17: instructionPushOPBytes,
	OP_PUSHBYTES_18: instructionPushOPBytes,
	OP_PUSHBYTES_19: instructionPushOPBytes,
	OP_PUSHBYTES_20: instructionPushOPBytes,
	OP_PUSHBYTES_21: instructionPushOPBytes,
	OP_PUSHBYTES_22: instructionPushOPBytes,
	OP_PUSHBYTES_23: instructionPushOPBytes,
	OP_PUSHBYTES_24: instructionPushOPBytes,
	OP_PUSHBYTES_25: instructionPushOPBytes,
	OP_PUSHBYTES_26: instructionPushOPBytes,
	OP_PUSHBYTES_27: instructionPushOPBytes,
	OP_PUSHBYTES_28: instructionPushOPBytes,
	OP_PUSHBYTES_29: instructionPushOPBytes,
	OP_PUSHBYTES_30: instructionPushOPBytes,
	OP_PUSHBYTES_31: instructionPushOPBytes,
	OP_PUSHBYTES_32: instructionPushOPBytes,
	OP_PUSHBYTES_33: instructionPushOPBytes,
	OP_PUSHBYTES_34: instructionPushOPBytes,
	OP_PUSHBYTES_35: instructionPushOPBytes,
	OP_PUSHBYTES_36: instructionPushOPBytes,
	OP_PUSHBYTES_37: instructionPushOPBytes,
	OP_PUSHBYTES_38: instructionPushOPBytes,
	OP_PUSHBYTES_39: instructionPushOPBytes,
	OP_PUSHBYTES_40: instructionPushOPBytes,
	OP_PUSHBYTES_41: instructionPushOPBytes,
	OP_PUSHBYTES_42: instructionPushOPBytes,
	OP_PUSHBYTES_43: instructionPushOPBytes,
	OP_PUSHBYTES_44: instructionPushOPBytes,
	OP_PUSHBYTES_45: instructionPushOPBytes,
	OP_PUSHBYTES_46: instructionPushOPBytes,
	OP_PUSHBYTES_47: instructionPushOPBytes,
	OP_PUSHBYTES_48: instructionPushOPBytes,
	OP_PUSHBYTES_49: instructionPushOPBytes,
	OP_PUSHBYTES_50: instructionPushOPBytes,
	OP_PUSHBYTES_51: instructionPushOPBytes,
	OP_PUSHBYTES_52: instructionPushOPBytes,
	OP_PUSHBYTES_53: instructionPushOPBytes,
	OP_PUSHBYTES_54: instructionPushOPBytes,
	OP_PUSHBYTES_55: instructionPushOPBytes,
	OP_PUSHBYTES_56: instructionPushOPBytes,
	OP_PUSHBYTES_57: instructionPushOPBytes,
	OP_PUSHBYTES_58: instructionPushOPBytes,
	OP_PUSHBYTES_59: instructionPushOPBytes,
	OP_PUSHBYTES_60: instructionPushOPBytes,
	OP_PUSHBYTES_61: instructionPushOPBytes,
	OP_PUSHBYTES_62: instructionPushOPBytes,
	OP_PUSHBYTES_63: instructionPushOPBytes,
	OP_PUSHBYTES_64: instructionPushOPBytes,
	OP_PUSHBYTES_65: instructionPushOPBytes,
	OP_PUSHBYTES_66: instructionPushOPBytes,
	OP_PUSHBYTES_67: instructionPushOPBytes,
	OP_PUSHBYTES_68: instructionPushOPBytes,
	OP_PUSHBYTES_69: instructionPushOPBytes,
	OP_PUSHBYTES_70: instructionPushOPBytes,
	OP_PUSHBYTES_71: instructionPushOPBytes,
	OP_PUSHBYTES_72: instructionPushOPBytes,
	OP_PUSHBYTES_73: instructionPushOPBytes,
	OP_PUSHBYTES_74: instructionPushOPBytes,
	OP_PUSHBYTES_75: instructionPushOPBytes,
	OP_PUSHDATA1:    instructionPushOPBytes,
	OP_PUSHDATA2:    instructionPushOPBytes,
	OP_PUSHDATA4:    instructionPushOPBytes,
	OP_1NEGATE:      instructionPushOPN,
	OP_RESERVED:     instructionRESERVED,
	OP_1:            instructionPushOPN,
	OP_2:            instructionPushOPN,
	OP_3:            instructionPushOPN,
	OP_4:            instructionPushOPN,
	OP_5:            instructionPushOPN,
	OP_6:            instructionPushOPN,
	OP_7:            instructionPushOPN,
	OP_8:            instructionPushOPN,
	OP_9:            instructionPushOPN,
	OP_10:           instructionPushOPN,
	OP_11:           instructionPushOPN,
	OP_12:           instructionPushOPN,
	OP_13:           instructionPushOPN,
	OP_14:           instructionPushOPN,
	OP_15:           instructionPushOPN,
	OP_16:           instructionPushOPN,

	// control
	OP_NOP:      instructionRESERVED,
	OP_VER:      instructionRESERVED,
	OP_IF:       instructionIF,
	OP_NOTIF:    instructionIF,
	OP_VERIF:    instructionRESERVED,
	OP_VERNOTIF: instructionRESERVED,
	OP_ELSE:     instructionELSE,
	OP_ENDIF:    instructionENDIF,
	OP_VERIFY:   instructionVERIFY,
	OP_RETURN:   instructionRESERVED,

	// stack ops
	OP_TOALTSTACK:   instructionTOTALSTACK,
	OP_FROMALTSTACK: instructionFROMALTSTACK,
	OP_2DROP:        instruction2DROP,
	OP_2DUP:         instruction2DUP,
	OP_3DUP:         instruction3DUP,
	OP_2OVER:        instruction2OVER,
	OP_2ROT:         instruction2ROT,
	OP_2SWAP:        instruction2SWAP,
	OP_IFDUP:        instructionIFDUP,
	OP_DEPTH:        instructionDEPTH,
	OP_DROP:         instructionDROP,
	OP_DUP:          instructionDUP,
	OP_NIP:          instructionNIP,
	OP_OVER:         instructionOVER,
	OP_PICK:         instructionROLL,
	OP_ROLL:         instructionROLL,
	OP_ROT:          instructionROT,
	OP_SWAP:         instructionSWAP,
	OP_TUCK:         instructionTUCK,

	// splice ops
	OP_CAT:    instructionCAT,
	OP_SUBSTR: instructionSUBSTR,
	OP_LEFT:   instructionLEFT,
	OP_RIGHT:  instructionRIGHT,
	OP_SIZE:   instructionSIZE,

	// bit logic
	OP_INVERT:      instructionINVERT,
	OP_AND:         instructionBITOP,
	OP_OR:          instructionBITOP,
	OP_XOR:         instructionBITOP,
	OP_EQUAL:       instructionEQUAL,
	OP_EQUALVERIFY: instructionEQUALVERIFY,
	OP_RESERVED1:   instructionRESERVED,
	OP_RESERVED2:   instructionRESERVED,

	// numeric
	OP_1ADD:      instructionUNARY,
	OP_1SUB:      instructionUNARY,
	OP_2MUL:      instructionBINARY,
	OP_2DIV:      instructionBINARY,
	OP_NEGATE:    instructionUNARY,
	OP_ABS:       instructionUNARY,
	OP_NOT:       instructionUNARY,
	OP_0NOTEQUAL: instructionUNARY,

	OP_ADD:    instructionBINARY,
	OP_SUB:    instructionBINARY,
	OP_MUL:    instructionBINARY,
	OP_DIV:    instructionBINARY,
	OP_MOD:    instructionBINARY,
	OP_LSHIFT: instructionSFHIT,
	OP_RSHIFT: instructionSFHIT,

	OP_BOOLAND:            instructionBINARY,
	OP_BOOLOR:             instructionBINARY,
	OP_NUMEQUAL:           instructionBINARY,
	OP_NUMEQUALVERIFY:     instructionBINARY,
	OP_NUMNOTEQUAL:        instructionBINARY,
	OP_LESSTHAN:           instructionBINARY,
	OP_GREATERTHAN:        instructionBINARY,
	OP_LESSTHANOREQUAL:    instructionBINARY,
	OP_GREATERTHANOREQUAL: instructionBINARY,
	OP_MIN:                instructionBINARY,
	OP_MAX:                instructionBINARY,
	OP_WITHIN:             instructionWITHIN,

	// crypto
	OP_RIPEMD160:           instructionRIPEMD160,
	OP_SHA1:                instructionSHA1,
	OP_SHA256:              instructionSHA256,
	OP_HASH160:             instructionHASH160,
	OP_HASH256:             instructionHASH256,
	OP_CODESEPARATOR:       instructionCODESEPARATOR,
	OP_CHECKSIG:            instructionCHECKSIG,
	OP_CHECKSIGVERIFY:      instructionCHECKSIG,
	OP_CHECKMULTISIG:       instructionCHECKMULTISIG,
	OP_CHECKMULTISIGVERIFY: instructionCHECKMULTISIGVERIFY,

	// expansion
	OP_NOP1:                instructionRESERVED,
	OP_CHECKLOCKTIMEVERIFY: instructionCHECKLOCKTIMEVERIFY,
	//OP_NOP2 = OP_CHECKLOCKTIMEVERIFY
	OP_CHECKSEQUENCEVERIFY: instructionCHECKSEQUENCEVERIFY,

	//OP_NOP3 = OP_CHECKSEQUENCEVERIFY
	OP_NOP4:  instructionRESERVED,
	OP_NOP5:  instructionRESERVED,
	OP_NOP6:  instructionRESERVED,
	OP_NOP7:  instructionRESERVED,
	OP_NOP8:  instructionRESERVED,
	OP_NOP9:  instructionRESERVED,
	OP_NOP10: instructionRESERVED,
}

func instructionRESERVED(ctx *InterpreterContext) error {
	return ErrInstructionReserved
}
