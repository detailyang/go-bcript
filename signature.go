package bscript

type SignatureVersion uint32

const (
	SignatureVersionBase = 1 << iota
	SignatureVersionWitnessV0
	SignatureVersionForkId
)
