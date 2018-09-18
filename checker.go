package bscript

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

func CheckSignatureEncoding(sig []byte, flag Flag) error {
	return nil
}
