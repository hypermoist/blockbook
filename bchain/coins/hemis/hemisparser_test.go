//go:build unittest

package hemis

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/martinboehm/btcutil/chaincfg"
	"github.com/trezor/blockbook/bchain"
	"github.com/trezor/blockbook/bchain/coins/btc"
)

func TestMain(m *testing.M) {
	c := m.Run()
	chaincfg.ResetParams()
	os.Exit(c)
}

func Test_GetAddrDescFromAddress_Mainnet(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "P2PKH1",
			args:    args{address: "HSjAJfhR8wUNTw5WE69KeVBiwrAcSjs11K"},
			want:    "76a914dda91c0396050d660f9c0e38f78064486bbfcb2c88ac",
			wantErr: false,
		},
	}
	parser := NewHemisParser(GetChainParams("main"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.GetAddrDescFromAddress(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddrDescFromAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("GetAddrDescFromAddress() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_GetAddressesFromAddrDesc(t *testing.T) {
	type args struct {
		script string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		want2   bool
		wantErr bool
	}{
		{
			name:    "P2PKH1",
			args:    args{script: "76a914dda91c0396050d660f9c0e38f78064486bbfcb2c88ac"},
			want:    []string{"HSjAJfhR8wUNTw5WE69KeVBiwrAcSjs11K"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "pubkey",
			args:    args{script: "210251c5555ff3c684aebfca92f5329e2f660da54856299da067060a1bcf5e8fae73ac"},
			want:    []string{"HLi5G5BUQeV3zki2Y7a2j4VjWkKqScobkM"},
			want2:   false,
			wantErr: false,
		},
	}

	parser := NewHemisParser(GetChainParams("main"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.script)
			got, got2, err := parser.GetAddressesFromAddrDesc(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddressesFromAddrDesc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

var (
	// regular transaction
	testTx1       bchain.Tx
	testTxPacked1 = "0b000000a1e1c9b2948a44d1e223ee711331224fca82dac44a4ebd0e6fd2d57ce6d1590a366f41fdac00dddd95b9414a9f8e976d2f38bea648cb66129a384c59685e707e6dc13d669472361b00000000fbc2f4300c01f0b7820d00e3347c8da4ee614674376cbc45359daa54f9b5493e0201000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0503256d0200ffffffff01000000000000000000000000000100000001aec51b66adfa98cc3c40a651083de2a0d0525e2e1e7eb8cf8ae530833c9f972b050000006a47304402205f4bfb6701b6f4060b426df29e74dd7838fbb6aaa34847c56fc170335a3e296502203b48c46910bf2bad4a6aaa2e364938d436dfdee317155329f0309c0435a5faf9012102bed9fe429195cca195c0cfc0ae13fcf6a9985ee636f289a2fa319b4bfc00a940ffffffff02000000000000000000fa4d3eb4020000001976a91434d5cd5bc462727c3f7ed1b28a0e3c65b43e5cd588ac000000004630440220565a495b3fced45a081e97d80c28d3db3258c7dbe3af2711f40d129c843cea82022075cfc82d2c4076e238407e960bc157714d67e3595fd4ccc6b4ec75d6aac40bdd"
)

func init() {
	testTx1 = bchain.Tx{
		Hex:      "010000000188557c816acd0a61579b701278c7dde85ea25d57877f9dbc65d3b2df2feacc42320000006b483045022100f5d0e98d064d5256852e420a4a3779527fb182c5edbfecf6143fc70eeba8eeef02202f0b2445185fbf846cca07c56c317733a9a4e46f960615f541da7aa27c33cfa201210251c5555ff3c684aebfca92f5329e2f660da54856299da067060a1bcf5e8fae73ffffffff03000000000000000000f06832fa0100000023210251c5555ff3c684aebfca92f5329e2f660da54856299da067060a1bcf5e8fae73aca038370e000000001976a914b4aa56c103b398f875bb8d15c3bb4136aa62725f88ac00000000",
		Txid:     "52b116d26f7c8b633c284f8998a431e106d837c0c5888f9ea5273d36c4556bec",
		LockTime: 0,
		Vin: []bchain.Vin{
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "483045022100f5d0e98d064d5256852e420a4a3779527fb182c5edbfecf6143fc70eeba8eeef02202f0b2445185fbf846cca07c56c317733a9a4e46f960615f541da7aa27c33cfa201210251c5555ff3c684aebfca92f5329e2f660da54856299da067060a1bcf5e8fae73",
				},
				Txid:     "42ccea2fdfb2d365bc9d7f87575da25ee8ddc77812709b57610acd6a817c5588",
				Vout:     50,
				Sequence: 4294967295,
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(0),
				N:        0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "",
				},
			},
			{
				ValueSat: *big.NewInt(8492574960),
				N:        1,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "210251c5555ff3c684aebfca92f5329e2f660da54856299da067060a1bcf5e8fae73ac",
					Addresses: []string{
						"HLi5G5BUQeV3zki2Y7a2j4VjWkKqScobkM",
					},
				},
			},
			{
				ValueSat: *big.NewInt(238500000),
				N:        2,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "76a914b4aa56c103b398f875bb8d15c3bb4136aa62725f88ac",
					Addresses: []string{
						"HSjAJfhR8wUNTw5WE69KeVBiwrAcSjs11K",
					},
				},
			},
		},
		Blocktime: 1504351235,
		Time:      1504351235,
	}
}

func Test_PackTx(t *testing.T) {
	type args struct {
		tx        bchain.Tx
		height    uint32
		blockTime int64
		parser    *HemisParser
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "hemis-1",
			args: args{
				tx:        testTx1,
				height:    159013,
				blockTime: 1715322960,
				parser:    NewHemisParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    testTxPacked1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.parser.PackTx(&tt.args.tx, tt.args.height, tt.args.blockTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("packTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("packTx() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_UnpackTx(t *testing.T) {
	type args struct {
		packedTx string
		parser   *HemisParser
	}
	tests := []struct {
		name    string
		args    args
		want    *bchain.Tx
		want1   uint32
		wantErr bool
	}{
		{
			name: "hemis-1",
			args: args{
				packedTx: testTxPacked1,
				parser:   NewHemisParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    &testTx1,
			want1:   159013,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.packedTx)
			got, got1, err := tt.args.parser.UnpackTx(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpackTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unpackTx() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("unpackTx() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

type testBlock struct {
	size int
	time int64
	txs  []string
}

var testParseBlockTxs = map[int]testBlock{
	159013: {
		size: 449,
		time: 1715323245,
		txs: []string{
			"556569e1bd20ae007853d839fda5cbefed4883ac53e6327a0a8a30180d242e24",
			"52b116d26f7c8b633c284f8998a431e106d837c0c5888f9ea5273d36c4556bec",
		},
	},
}

func helperLoadBlock(t *testing.T, height int) []byte {
	name := fmt.Sprintf("block_dump.%d", height)
	path := filepath.Join("testdata", name)

	d, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	d = bytes.TrimSpace(d)

	b := make([]byte, hex.DecodedLen(len(d)))
	_, err = hex.Decode(b, d)
	if err != nil {
		t.Fatal(err)
	}

	return b
}

func TestParseBlock(t *testing.T) {
	p := NewHemisParser(GetChainParams("main"), &btc.Configuration{})

	for height, tb := range testParseBlockTxs {
		b := helperLoadBlock(t, height)

		blk, err := p.ParseBlock(b)
		if err != nil {
			t.Fatal(err)
		}

		if blk.Size != tb.size {
			t.Errorf("ParseBlock() block size: got %d, want %d", blk.Size, tb.size)
		}

		if blk.Time != tb.time {
			t.Errorf("ParseBlock() block time: got %d, want %d", blk.Time, tb.time)
		}

		if len(blk.Txs) != len(tb.txs) {
			t.Errorf("ParseBlock() number of transactions: got %d, want %d", len(blk.Txs), len(tb.txs))
		}

		for ti, tx := range tb.txs {
			if blk.Txs[ti].Txid != tx {
				t.Errorf("ParseBlock() transaction %d: got %s, want %s", ti, blk.Txs[ti].Txid, tx)
			}
		}
	}
}
