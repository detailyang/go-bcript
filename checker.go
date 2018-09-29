package bscript

import "math/bits"

type Checker interface {
	CheckLockTime(locktime uint32) error
	CheckSequence(sequence uint32) error
	CheckSignature(sig, pubkey []byte, script *Script, version SignatureVersion) error
}

type NoopChecker struct{}

func NewNoopChecker() *NoopChecker                         { return &NoopChecker{} }
func (n *NoopChecker) CheckLockTime(locktime uint32) error { return nil }
func (n *NoopChecker) CheckSequence(sequence uint32) error { return nil }
func (n *NoopChecker) CheckSignature(sig, pubkey []byte, script *Script, version SignatureVersion) error {
	return nil
}

func CheckHashTypeEncoding(hashtype byte, flag Flag) error {
	return nil
}

func CheckPubkeyEncoding(pubkey []byte, flag Flag) error {
	return nil
}

func CheckSignatureEncoding(sig []byte, flag Flag, sigver SignatureVersion) error {
	if len(sig) == 0 {
		return nil
	}

	if flag.Has(ScriptVerifyDERSignatures) ||
		flag.Has(ScriptVerifyLowS) ||
		flag.Has(ScriptVerifyStrictEncoding) ||
		!isValidSignatureEncoding(sig) {
		return ErrInterpreterBadSignatureDer
	}

	if flag.Has(ScriptVerifyLowS) {
		if !isLowDerSignature(sig) {
			return ErrInterpreterSigantureHighS
		}
	}

	if flag.Has(ScriptVerifyStrictEncoding) &&
		!isDefinedHashtypeSiganture(sigver, sig) {
		return ErrInterpreterBadSignatureHashType
	}

	if flag.Has(ScriptVerifyStrictEncoding) {
		// TODO: check fork id
	}

	return nil
}

/// A canonical signature exists of: <30> <total len> <02> <len R> <R> <02> <len S> <S> <hashtype>
/// Where R and S are not negative (their first byte has its highest bit not set), and not
/// excessively padded (do not start with a 0 byte, unless an otherwise negative number follows,
/// in which case a single 0 byte is necessary and even required).
///
/// See https://bitcointalk.org/index.php?topic=8392.msg127623#msg127623
///
/// This function is consensus-critical since BIP66.
func isValidSignatureEncoding(sig []byte) bool {
	// Format: 0x30 [total-length] 0x02 [R-length] [R] 0x02 [S-length] [S] [sighash]
	// * total-length: 1-byte length descriptor of everything that follows,
	//   excluding the sighash byte.
	// * R-length: 1-byte length descriptor of the R value that follows.
	// * R: arbitrary-length big-endian encoded R value. It must use the shortest
	//   possible encoding for a positive integers (which means no null bytes at
	//   the start, except a single one when the next byte has its highest bit set).
	// * S-length: 1-byte length descriptor of the S value that follows.
	// * S: arbitrary-length big-endian encoded S value. The same rules apply.
	// * sighash: 1-byte value indicating what data is hashed (not part of the DER
	//   signature)
	if len(sig) < 9 || len(sig) > 73 {
		return false
	}

	if sig[0] != 0x30 {
		return false
	}

	if int(sig[1]) != len(sig)-3 {
		return false
	}

	nr := int(sig[3])
	if nr+5 >= len(sig) {
		return false
	}

	ns := int(sig[nr+5])

	if nr+ns+7 != len(sig) {
		return false
	}

	if sig[2] != 2 {
		return false
	}

	if nr == 0 {
		return false
	}

	if sig[4]&0x80 != 0 {
		return false
	}

	if nr > 1 && sig[4] == 0 && (sig[5]&0x80 == 0x80) {
		return false
	}

	if sig[nr+4] != 2 {
		return false
	}

	if ns == 0 {
		return false
	}

	if sig[nr+6]&0x80 != 0 {
		return false
	}

	if ns > 1 && (sig[nr+6] == 0) && (sig[nr+7]&0x80) == 0 {
		return false
	}

	return true
}

func isLowDerSignature(sig []byte) bool {
	if !isValidSignatureEncoding(sig) {
		return false
	}

	// TODO: check low s

	return true
}

func isDefinedHashtypeSiganture(sigver SignatureVersion, sig []byte) bool {
	if len(sig) == 0 {
		return false
	}

	u := uint(len(sig))
	switch sigver {
	case SignatureVersionForkId:
		u = u & bits.Reverse(0x40|0x80)
	default:
		u = u & bits.Reverse(0x80)
	}

	switch u {
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		return true
	default:
		return false
	}

	return false
}
