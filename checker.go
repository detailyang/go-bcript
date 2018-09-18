package bscript

type Checker interface {
	CheckLockTime(locktime uint32) error
	CheckSequence(sequence uint32) error
	CheckSignature() error
}

type NoopChecker struct{}

func (n *NoopChecker) CheckLockTime(locktime uint32) error { return nil }
func (n *NoopChecker) CheckSequence(sequence uint32) error { return nil }
func (n *NoopChecker) CheckSignature() error               { return nil }

func CheckHashTypeEncoding(hashtype byte, flag Flag) error {
	return nil
}

func CheckPubkeyEncoding(pubkey []byte, flag Flag) error {
	return nil
}

func CheckSignatureEncoding(sig []byte, flag Flag) error {
	return nil
}
