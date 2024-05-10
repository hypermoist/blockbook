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
			args:    args{address: "HLRYJtL33LLp6bDHohdZecYnuavvdq7cfT"},
			want:    "76a91498833c37aabb0d983c94370ed5450ed55d41b67888ac",
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
			args:    args{script: "76a91498833c37aabb0d983c94370ed5450ed55d41b67888ac"},
			want:    []string{"HLRYJtL33LLp6bDHohdZecYnuavvdq7cfT"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "pubkey",
			args:    args{script: "03a8bb91411c991859d89415db9748f03ca2efc526004fcf52171225650745e31a"},
			want:    []string{"HUMhGKdTEG5JasTo4UqHGXFXoWMRfM5Ux2"},
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
	testTxPacked1 = "0100000001aec51b66adfa98cc3c40a651083de2a0d0525e2e1e7eb8cf8ae530833c9f972b050000006a47304402205f4bfb6701b6f4060b426df29e74dd7838fbb6aaa34847c56fc170335a3e296502203b48c46910bf2bad4a6aaa2e364938d436dfdee317155329f0309c0435a5faf9012102bed9fe429195cca195c0cfc0ae13fcf6a9985ee636f289a2fa319b4bfc00a940ffffffff02000000000000000000fa4d3eb4020000001976a91434d5cd5bc462727c3f7ed1b28a0e3c65b43e5cd588ac00000000"
)

func init() {
	testTx1 = bchain.Tx{
		Hex:      "0100000001aec51b66adfa98cc3c40a651083de2a0d0525e2e1e7eb8cf8ae530833c9f972b050000006a47304402205f4bfb6701b6f4060b426df29e74dd7838fbb6aaa34847c56fc170335a3e296502203b48c46910bf2bad4a6aaa2e364938d436dfdee317155329f0309c0435a5faf9012102bed9fe429195cca195c0cfc0ae13fcf6a9985ee636f289a2fa319b4bfc00a940ffffffff02000000000000000000fa4d3eb4020000001976a91434d5cd5bc462727c3f7ed1b28a0e3c65b43e5cd588ac00000000",
		Txid:     "64f626e9f7ff3d37c1332f138f4abbb3624678e63f23111539cc4e385ac38655",
		LockTime: 0,
		Vin: []bchain.Vin{
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "47304402205f4bfb6701b6f4060b426df29e74dd7838fbb6aaa34847c56fc170335a3e296502203b48c46910bf2bad4a6aaa2e364938d436dfdee317155329f0309c0435a5faf9012102bed9fe429195cca195c0cfc0ae13fcf6a9985ee636f289a2fa319b4bfc00a940",
				},
				Txid:     "400522dec366e254b337f230794c4e824c7d23a01f8c1782c42a956624af9d85",
				Vout:     5,
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
					Hex: "76a91434d5cd5bc462727c3f7ed1b28a0e3c65b43e5cd588ac",
					Addresses: []string{
						"HBLVeZAwf7BNqUZLXcqFzUF3GLazvFfuMc",
					},
				},
			},
		},
		Blocktime: 1715323245,
		Time:      1715323245,
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
				blockTime: 1715323245,
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
			"62b30b5b09da27ee57631b142fe9a0e87e098438710776ab01b617ee4939bee2",
			"64f626e9f7ff3d37c1332f138f4abbb3624678e63f23111539cc4e385ac38655",
		},
	},
	159017: {
		size: 450,
		time: 1715323455,
		txs: []string{
			"004da915937db5ac6d2837293e94764251dd4be7ab0245f0f6483e33d9bdfa32",
			"fe435b7d3ea87169e45799c1c0f28904cb9f2669551d11d0c62af7927acfa484",
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
