package hemis

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"

	"github.com/juju/errors"
	"github.com/martinboehm/btcd/wire"
	"github.com/martinboehm/btcutil/chaincfg"
	"github.com/trezor/blockbook/bchain"
	"github.com/trezor/blockbook/bchain/coins/btc"
)

// magic numbers
const (
	MainnetMagic wire.BitcoinNet = 0x99dd14a0
	TestnetMagic wire.BitcoinNet = 0xcad5e6f5
)

// chain parameters
var (
	MainNetParams chaincfg.Params
	TestNetParams chaincfg.Params
)

func init() {
	// HEMIS mainnet Address encoding magics
	MainNetParams = chaincfg.MainNetParams
	MainNetParams.Net = MainnetMagic
	MainNetParams.PubKeyHashAddrID = []byte{40} // starting with 'H'
	MainNetParams.ScriptHashAddrID = []byte{13}
	MainNetParams.PrivateKeyID = []byte{212}

	// HEMIS testnet Address encoding magics
	TestNetParams = chaincfg.TestNet3Params
	TestNetParams.Net = TestnetMagic
	TestNetParams.PubKeyHashAddrID = []byte{139} // starting with 'x' or 'y'
	TestNetParams.ScriptHashAddrID = []byte{19}
	TestNetParams.PrivateKeyID = []byte{239}
}

// HemisParser handle
type HemisParser struct {
	*btc.BitcoinLikeParser
	baseparser                         *bchain.BaseParser
	BitcoinOutputScriptToAddressesFunc btc.OutputScriptToAddressesFunc
}

// NewHemisParser returns new HemisParser instance
func NewHemisParser(params *chaincfg.Params, c *btc.Configuration) *HemisParser {
	p := &HemisParser{
		BitcoinLikeParser: btc.NewBitcoinLikeParser(params, c),
		baseparser:        &bchain.BaseParser{},
	}
	p.BitcoinOutputScriptToAddressesFunc = p.OutputScriptToAddressesFunc
	p.OutputScriptToAddressesFunc = p.outputScriptToAddresses
	return p
}

// GetChainParams contains network parameters for the main Hemis network
func GetChainParams(chain string) *chaincfg.Params {
	if !chaincfg.IsRegistered(&MainNetParams) {
		err := chaincfg.Register(&MainNetParams)
		if err == nil {
			err = chaincfg.Register(&TestNetParams)
		}
		if err != nil {
			panic(err)
		}
	}
	switch chain {
	case "test":
		return &TestNetParams
	default:
		return &MainNetParams
	}
}

// ParseBlock parses raw block to our Block struct
func (p *HemisParser) ParseBlock(b []byte) (*bchain.Block, error) {
	r := bytes.NewReader(b)
	w := wire.MsgBlock{}
	h := wire.BlockHeader{}
	err := h.Deserialize(r)
	if err != nil {
		return nil, errors.Annotatef(err, "Deserialize")
	}

	if h.Version > 7 {
		// Skip new hashFinalSaplingRoot (block version 8 or newer)
		r.Seek(32, io.SeekCurrent)
	}

	err = p.HemisDecodeTransactions(r, 0, &w)
	if err != nil {
		return nil, errors.Annotatef(err, "DecodeTransactions")
	}

	txs := make([]bchain.Tx, len(w.Transactions))
	for ti, t := range w.Transactions {
		txs[ti] = p.TxFromMsgTx(t, false)
	}

	return &bchain.Block{
		BlockHeader: bchain.BlockHeader{
			Size: len(b),
			Time: h.Timestamp.Unix(),
		},
		Txs: txs,
	}, nil
}

// PackTx packs transaction to byte array using protobuf
func (p *HemisParser) PackTx(tx *bchain.Tx, height uint32, blockTime int64) ([]byte, error) {
	return p.baseparser.PackTx(tx, height, blockTime)
}

// UnpackTx unpacks transaction from protobuf byte array
func (p *HemisParser) UnpackTx(buf []byte) (*bchain.Tx, uint32, error) {
	return p.baseparser.UnpackTx(buf)
}

// ParseTx parses byte array containing transaction and returns Tx struct
func (p *HemisParser) ParseTx(b []byte) (*bchain.Tx, error) {
	t := wire.MsgTx{}
	r := bytes.NewReader(b)
	if err := t.Deserialize(r); err != nil {
		return nil, err
	}
	tx := p.TxFromMsgTx(&t, true)
	tx.Hex = hex.EncodeToString(b)
	return &tx, nil
}

// TxFromMsgTx parses tx
func (p *HemisParser) TxFromMsgTx(t *wire.MsgTx, parseAddresses bool) bchain.Tx {
	vin := make([]bchain.Vin, len(t.TxIn))
	for i, in := range t.TxIn {

		s := bchain.ScriptSig{
			Hex: hex.EncodeToString(in.SignatureScript),
			// missing: Asm,
		}

		txid := in.PreviousOutPoint.Hash.String()

		vin[i] = bchain.Vin{
			Txid:      txid,
			Vout:      in.PreviousOutPoint.Index,
			Sequence:  in.Sequence,
			ScriptSig: s,
		}
	}
	vout := make([]bchain.Vout, len(t.TxOut))
	for i, out := range t.TxOut {
		addrs := []string{}
		if parseAddresses {
			addrs, _, _ = p.OutputScriptToAddressesFunc(out.PkScript)
		}
		s := bchain.ScriptPubKey{
			Hex:       hex.EncodeToString(out.PkScript),
			Addresses: addrs,
			// missing: Asm,
			// missing: Type,
		}
		var vs big.Int
		vs.SetInt64(out.Value)
		vout[i] = bchain.Vout{
			ValueSat:     vs,
			N:            uint32(i),
			ScriptPubKey: s,
		}
	}
	tx := bchain.Tx{
		Txid:     t.TxHash().String(),
		Version:  t.Version,
		LockTime: t.LockTime,
		Vin:      vin,
		Vout:     vout,
		// skip: BlockHash,
		// skip: Confirmations,
		// skip: Time,
		// skip: Blocktime,
	}
	return tx
}

// ParseTxFromJson parses JSON message containing transaction and returns Tx struct
func (p *HemisParser) ParseTxFromJson(msg json.RawMessage) (*bchain.Tx, error) {
	var tx bchain.Tx
	err := json.Unmarshal(msg, &tx)
	if err != nil {
		return nil, err
	}

	for i := range tx.Vout {
		vout := &tx.Vout[i]
		// convert vout.JsonValue to big.Int and clear it, it is only temporary value used for unmarshal
		vout.ValueSat, err = p.AmountToBigInt(vout.JsonValue)
		if err != nil {
			return nil, err
		}
		vout.JsonValue = ""

		if vout.ScriptPubKey.Addresses == nil {
			vout.ScriptPubKey.Addresses = []string{}
		}
	}

	return &tx, nil
}

// outputScriptToAddresses converts ScriptPubKey to bitcoin addresses
func (p *HemisParser) outputScriptToAddresses(script []byte) ([]string, bool, error) {

	rv, s, _ := p.BitcoinOutputScriptToAddressesFunc(script)
	return rv, s, nil
}

func (p *HemisParser) GetAddrDescForUnknownInput(tx *bchain.Tx, input int) bchain.AddressDescriptor {
	if len(tx.Vin) > input {
		scriptHex := tx.Vin[input].ScriptSig.Hex

		if scriptHex != "" {
			script, _ := hex.DecodeString(scriptHex)
			return script
		}
	}

	s := make([]byte, 10)
	return s
}

func (p *HemisParser) HemisDecodeTransactions(r *bytes.Reader, pver uint32, blk *wire.MsgBlock) error {
	maxTxPerBlock := uint64((wire.MaxBlockPayload / 10) + 1)

	txCount, err := wire.ReadVarInt(r, pver)
	if err != nil {
		return err
	}

	// Prevent more transactions than could possibly fit into a block.
	// It would be possible to cause memory exhaustion and panics without
	// a sane upper bound on this count.
	if txCount > maxTxPerBlock {
		str := fmt.Sprintf("too many transactions to fit into a block "+
			"[count %d, max %d]", txCount, maxTxPerBlock)
		return &wire.MessageError{Func: "utils.decodeTransactions", Description: str}
	}

	blk.Transactions = make([]*wire.MsgTx, 0, txCount)
	for i := uint64(0); i < txCount; i++ {
		tx := wire.MsgTx{}

		// read version & seek back to original state
		var version uint32 = 0
		if err = binary.Read(r, binary.LittleEndian, &version); err != nil {
			return err
		}
		if _, err = r.Seek(-4, io.SeekCurrent); err != nil {
			return err
		}

		if err != nil {
			return err
		}
		blk.Transactions = append(blk.Transactions, &tx)
	}

	return nil
}

func (p *HemisParser) HemisDecode(MsgTx *wire.MsgTx, r *bytes.Reader, pver uint32, enc wire.MessageEncoding) error {
	if err := MsgTx.BtcDecode(r, pver, enc); err != nil {
		return err
	}

	// extra
	version := uint32(MsgTx.Version)
	txVersion := version & 0xffff

	if txVersion >= 3 {
		// valueBalance
		r.Seek(9, io.SeekCurrent)

		// bindingSig
		r.Seek(64, io.SeekCurrent)
	}

	return nil
}
