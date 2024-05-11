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
	testTxPacked1 = " 0a2052b116d26f7c8b633c284f8998a431e106d837c0c5888f9ea5273d36c4556bec12f501010000000188557c816acd0a61579b701278c7dde85ea25d57877f9dbc65d3b2df2feacc42320000006b483045022100f5d0e98d064d5256852e420a4a3779527fb182c5edbfecf6143fc70eeba8eeef02202f0b2445185fbf846cca07c56c317733a9a4e46f960615f541da7aa27c33cfa201210251c5555ff3c684aebfca92f5329e2f660da54856299da067060a1bcf5e8fae73ffffffff03000000000000000000f06832fa0100000023210251c5555ff3c684aebfca92f5329e2f660da54856299da067060a1bcf5e8fae73aca038370e000000001976a914b4aa56c103b398f875bb8d15c3bb4136aa62725f88ac0000000018d080f7b10628a5da09329701122042ccea2fdfb2d365bc9d7f87575da25ee8ddc77812709b57610acd6a817c55881832226b483045022100f5d0e98d064d5256852e420a4a3779527fb182c5edbfecf6143fc70eeba8eeef02202f0b2445185fbf846cca07c56c317733a9a4e46f960615f541da7aa27c33cfa201210251c5555ff3c684aebfca92f5329e2f660da54856299da067060a1bcf5e8fae7328ffffffff0f3a003a520a0501fa3268f010011a23210251c5555ff3c684aebfca92f5329e2f660da54856299da067060a1bcf5e8fae73ac2222484c693547354255516556337a6b6932593761326a34566a576b4b7153636f626b4d3a470a040e3738a010021a1976a914b4aa56c103b398f875bb8d15c3bb4136aa62725f88ac222248536a414a6668523877554e547735574536394b6556426977724163536a7331314b"
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
