// Package transaction represents an internal implementation of an xdr
// transaction envelope
package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/stellar/go/build"
	shash "github.com/stellar/go/hash"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/xdr"
)

var Network build.Network

func init() {
	if os.Getenv("MOBIUS_NETWORK") == "" {
		Network = build.Network{Passphrase: network.TestNetworkPassphrase}
	} else if os.Getenv("MOBIUS_NETWORK") == "public" {
		Network = build.Network{Passphrase: network.PublicNetworkPassphrase}
	}
}

// Transaction creates a transaction envelope. Once a Transaction has been created from an enveloper, it is immutable. Signers must be added before submitting to the network or
// forwarding on to additional signers.
type Transaction struct {
	Tx          *xdr.TransactionEnvelope
	envelopeRaw interface{}
}

func New(envelope interface{}) (*Transaction, error) {
	var (
		tx xdr.TransactionEnvelope
	)
	if reflect.TypeOf(envelope).String() == "string" {
		rawr := strings.NewReader(envelope.(string))
		b64r := base64.NewDecoder(base64.StdEncoding, rawr)
		_, err := xdr.Unmarshal(b64r, &tx)
		if err != nil {
			return nil, err
		}
	} else if strings.Split(reflect.TypeOf(envelope).String(), ".")[1] == "xdr.TransactionEnvelope" {
		tx = envelope.(xdr.TransactionEnvelope)
	}
	return &Transaction{
		&tx,
		envelope,
	}, nil
}

// Memo returns the memo attached to the transaction
func (t *Transaction) Memo() xdr.Memo {
	return t.Tx.Tx.Memo
}

// Sign signs the transaction with the given keypair
func (t *Transaction) Sign(kps ...*keypair.Full) {
	for _, kp := range kps {
		h := t.Hash()
		sig, err := kp.SignDecorated(h[:])
		if err != nil {
			log.Fatalf("failed to sign transaction, %v tx: %+v", err, t.Tx)
		}
		t.Tx.Signatures = append(t.Tx.Signatures, sig)
	}
}

// SignHashX adds hashX signer preimage as signature
func (t *Transaction) SignHashX(preimage interface{}) {
	if reflect.TypeOf(preimage).String() == "string" {
		src := []byte(preimage.(string))
		dst := make(xdr.Signature, hex.DecodedLen(len(src)))
		n, err := hex.Decode(dst, src)
		if err != nil {
			log.Fatalf("failed to SignHashX, %v tx: %+v", err, t.Tx)
		}
		if len(dst[:n]) > 64 {
			log.Fatalf("preiamge cannot be longer than 64 bytes")
		}
		h := sha256.New()
		h.Write(dst[:n])
		var hint xdr.SignatureHint
		copy(hint[:], h.Sum(nil)[:4])
		t.Tx.Signatures = append(t.Tx.Signatures, xdr.DecoratedSignature{Hint: hint, Signature: dst})
	}
}

// Hash returns a hash for the current transcation, suitable for signing.
func (t *Transaction) Hash() [32]byte { return shash.Hash(t.SignatureBase()) }

// SignatureBase returns the "signature base" of ths current transcation, which
// is the value that, when hashed, should be signed to create a signature that
// validators on the Stellar Network will accept.
func (t *Transaction) SignatureBase() []byte {
	if Network.Passphrase == "" {
		log.Fatalf("No network selected. Set the `MOBIUS_NETWORK` environment variable to `test` or `public`")
	}
	netId := Network.ID()
	txXDR, err := t.Tx.Tx.MarshalBinary()
	if err != nil {
		log.Fatalf("failed to marhal txt binary, %v", err)
	}
	b := [][]byte{
		netId[:],
		[]byte(xdr.EnvelopeTypeEnvelopeTypeTx.String()),
		txXDR,
	}
	return bytes.Join(b, nil)
}
