package mock

import (
	"github.com/ltcsuite/lnd/keychain"
	"github.com/ltcsuite/ltcd/btcec/v2"
	"github.com/ltcsuite/ltcd/btcec/v2/ecdsa"
	"github.com/ltcsuite/ltcd/btcec/v2/schnorr"
	"github.com/ltcsuite/ltcd/chaincfg/chainhash"
	"github.com/ltcsuite/ltcd/txscript"
)

// SecretKeyRing is a mock implementation of the SecretKeyRing interface.
type SecretKeyRing struct {
	RootKey *btcec.PrivateKey
}

// DeriveNextKey currently returns dummy values.
func (s *SecretKeyRing) DeriveNextKey(
	_ keychain.KeyFamily) (keychain.KeyDescriptor, error) {

	return keychain.KeyDescriptor{
		PubKey: s.RootKey.PubKey(),
	}, nil
}

// DeriveKey currently returns dummy values.
func (s *SecretKeyRing) DeriveKey(
	_ keychain.KeyLocator) (keychain.KeyDescriptor, error) {

	return keychain.KeyDescriptor{
		PubKey: s.RootKey.PubKey(),
	}, nil
}

// DerivePrivKey currently returns dummy values.
func (s *SecretKeyRing) DerivePrivKey(
	_ keychain.KeyDescriptor) (*btcec.PrivateKey, error) {

	return s.RootKey, nil
}

// ECDH currently returns dummy values.
func (s *SecretKeyRing) ECDH(_ keychain.KeyDescriptor,
	_ *btcec.PublicKey) ([32]byte, error) {

	return [32]byte{}, nil
}

// SignMessage signs the passed message and ignores the KeyDescriptor.
func (s *SecretKeyRing) SignMessage(_ keychain.KeyLocator,
	msg []byte, doubleHash bool) (*ecdsa.Signature, error) {

	var digest []byte
	if doubleHash {
		digest = chainhash.DoubleHashB(msg)
	} else {
		digest = chainhash.HashB(msg)
	}
	return ecdsa.Sign(s.RootKey, digest), nil
}

// SignMessageCompact signs the passed message.
func (s *SecretKeyRing) SignMessageCompact(_ keychain.KeyLocator,
	msg []byte, doubleHash bool) ([]byte, error) {

	var digest []byte
	if doubleHash {
		digest = chainhash.DoubleHashB(msg)
	} else {
		digest = chainhash.HashB(msg)
	}
	return ecdsa.SignCompact(s.RootKey, digest, true)
}

// SignMessageSchnorr signs the passed message and ignores the KeyDescriptor.
func (s *SecretKeyRing) SignMessageSchnorr(_ keychain.KeyLocator,
	msg []byte, doubleHash bool, taprootTweak []byte) (*schnorr.Signature,
	error) {

	var digest []byte
	if doubleHash {
		digest = chainhash.DoubleHashB(msg)
	} else {
		digest = chainhash.HashB(msg)
	}

	privKey := s.RootKey
	if len(taprootTweak) > 0 {
		privKey = txscript.TweakTaprootPrivKey(*privKey, taprootTweak)
	}

	return schnorr.Sign(privKey, digest)
}
