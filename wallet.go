package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/lucasmbaia/blockchain/crypto"
	"github.com/lucasmbaia/blockchain/utils"
	"golang.org/x/crypto/ripemd160"
	"hash"
)

const (
	VERSION_WALLET       = byte(0x00)
	ADDRESS_CHECKSUM_LEN = 4
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
	Address    []byte
}

type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallet() (*Wallet, error) {
	var (
		wallet    = &Wallet{}
		err       error
		pubHash   []byte
		version   []byte
		d2version utils.Hash
		cksum     []byte
	)

	if wallet.PrivateKey, wallet.PublicKey, err = newKeyPair(); err != nil {
		return wallet, err
	}

	if pubHash, err = HashPubKey(wallet.PublicKey); err != nil {
		return wallet, err
	}

	version = append([]byte{VERSION_WALLET}, pubHash...)
	d2version = utils.CalcDoubleHash(version)
	cksum = d2version[:][:ADDRESS_CHECKSUM_LEN]
	wallet.Address = utils.Base58Encode(append(version, cksum...))

	return wallet, nil
}

func (w *Wallet) PrivateToHex() string {
	return crypto.ECDSToHEX(w.PrivateKey.D.Bytes())
}

func HexToPrivate(input string) (*ecdsa.PrivateKey, error) {
	return crypto.ToECDS(input)
}

func AddressToPubHash(address []byte) []byte {
	var decoded = utils.Base58Decode(address)

	return decoded[1 : len(decoded)-4]
}

func UnlockWallet(private string, address string) (bool, *Wallet, error) {
	var (
		pk        *ecdsa.PrivateKey
		pubHash   []byte
		version   []byte
		cksum     []byte
		d2version utils.Hash
		err       error
		w         *Wallet
	)

	if pk, err = crypto.ToECDS(private); err != nil {
		return false, w, err
	}

	if pubHash, err = HashPubKey(append(pk.PublicKey.X.Bytes(), pk.PublicKey.Y.Bytes()...)); err != nil {
		return false, w, err
	}

	version = append([]byte{VERSION_WALLET}, pubHash...)
	d2version = utils.CalcDoubleHash(version)
	cksum = d2version[:][:ADDRESS_CHECKSUM_LEN]

	w = &Wallet{
		PrivateKey: *pk,
		PublicKey:  append(pk.PublicKey.X.Bytes(), pk.PublicKey.Y.Bytes()...),
		Address:    utils.Base58Encode(append(version, cksum...))[:],
	}

	if string(utils.Base58Encode(append(version, cksum...))[:]) == address {
		return true, w, nil
	}

	return false, w, nil
}

func HashPubKey(key []byte) ([]byte, error) {
	var (
		hashSHA256 utils.Hash
		hashRIP160 hash.Hash
		err        error
	)

	hashSHA256 = utils.CalcHash(key)
	hashRIP160 = ripemd160.New()

	if _, err = hashRIP160.Write(hashSHA256[:]); err != nil {
		return []byte{}, err
	}

	return hashRIP160.Sum(nil), nil
}

func CheckValidAddress(address []byte) bool {
	var (
		decoded []byte
		cksum   [4]byte
	)

	decoded = utils.Base58Decode(address)
	copy(cksum[:], decoded[len(decoded)-4:])

	if checksum(decoded[:len(decoded)-4]) != cksum {
		return false
	}

	return true
}

func checksum(input []byte) (cksum [4]byte) {
	first := sha256.Sum256(input)
	second := sha256.Sum256(first[:])
	copy(cksum[:], second[:4])
	return
}

func newKeyPair() (ecdsa.PrivateKey, []byte, error) {
	var (
		err        error
		privateKey *ecdsa.PrivateKey
	)

	var curve = elliptic.P256()
	if privateKey, err = ecdsa.GenerateKey(curve, rand.Reader); err != nil {
		return *privateKey, []byte{}, err
	}

	var publicKey = append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey, nil
}
