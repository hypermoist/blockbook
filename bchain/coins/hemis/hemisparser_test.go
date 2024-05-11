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
	// Mint transaction
	testTx1       bchain.Tx
	testTxPacked1 = "0a2052b116d26f7c8b633c284f8998a431e106d837c0c5888f9ea5273d36c4556bec12f501010000000188557c816acd0a61579b701278c7dde85ea25d57877f9dbc65d3b2df2feacc42320000006b483045022100f5d0e98d064d5256852e420a4a3779527fb182c5edbfecf6143fc70eeba8eeef02202f0b2445185fbf846cca07c56c317733a9a4e46f960615f541da7aa27c33cfa201210251c5555ff3c684aebfca92f5329e2f660da54856299da067060a1bcf5e8fae73ffffffff03000000000000000000f06832fa0100000023210251c5555ff3c684aebfca92f5329e2f660da54856299da067060a1bcf5e8fae73aca038370e000000001976a914b4aa56c103b398f875bb8d15c3bb4136aa62725f88ac000000001883a8aacd052880ea30329701122042ccea2fdfb2d365bc9d7f87575da25ee8ddc77812709b57610acd6a817c55881832226b483045022100f5d0e98d064d5256852e420a4a3779527fb182c5edbfecf6143fc70eeba8eeef02202f0b2445185fbf846cca07c56c317733a9a4e46f960615f541da7aa27c33cfa201210251c5555ff3c684aebfca92f5329e2f660da54856299da067060a1bcf5e8fae7328ffffffff0f3a003a520a0501fa3268f010011a23210251c5555ff3c684aebfca92f5329e2f660da54856299da067060a1bcf5e8fae73ac2222444b4c33517a43624a71724870524b4148764571736f6d7344686b515076567a5a673a470a040e3738a010021a1976a914b4aa56c103b398f875bb8d15c3bb4136aa62725f88ac2222444d634e45393855667571454b32674746664b423459705741577162774852415448"

	// Normal transaction
	testTx2       bchain.Tx
	testTxPacked2 = "0a20599d5d797a4575eb25e1c291c0e7630bd6fdc0e6ec5fa9b14147f929a4e41bf212ae020100000001b56a0fe242a8de7dcbb58ae1009e44e7f2ec25a65eeb8b815cf53393309741ca0100000049483045022100cc208a59341dca98207ec8a4a42c014d435192694a77c69d40e51467800c0a0802205ac1782d4ecefa260b33340d92c2ab2396b43c1073a67b4180aa8ef2aede8af801ffffffff0200e876481700000087c10281004c816f5ce1eeda911203319a256e8560c8dbfd47b569ff32c27559bda78854e63e49718ce43036e5120dce357b5630afd745d399f91e675a921adbb45224a6661656217fcfe32396fb25609b724646759116326964f2f1f7ddb7c340dc24be2b75a0a9dc05ca2fdf805c03c7a04d972456beb82a51de73d8842b39a553919dfa5d8e003e98dab7210000001976a914dda91c0396050d660f9c0e38f78064486bbfcb2c88ac00000000189dab96cf0528b8dc3432751220ca4197309333f55c818beb5ea625ecf2e7449e00e18ab5cb7ddea842e20f6ab518012249483045022100cc208a59341dca98207ec8a4a42c014d435192694a77c69d40e51467800c0a0802205ac1782d4ecefa260b33340d92c2ab2396b43c1073a67b4180aa8ef2aede8af80128ffffffff0f3a91010a05174876e8001a8701c10281004c816f5ce1eeda911203319a256e8560c8dbfd47b569ff32c27559bda78854e63e49718ce43036e5120dce357b5630afd745d399f91e675a921adbb45224a6661656217fcfe32396fb25609b724646759116326964f2f1f7ddb7c340dc24be2b75a0a9dc05ca2fdf805c03c7a04d972456beb82a51de73d8842b39a553919dfa5d8e003a480a0521b7da983e10011a1976a914dda91c0396050d660f9c0e38f78064486bbfcb2c88ac222244524d38546169593338716348626764797470386f455472656f62424c4874706545"
)

func init() {
	testTx1 = bchain.Tx{
		Hex:      "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff050385700200ffffffff0100000000000000000000000000",
		Txid:     "19c4eb500005e2a65fc0ddd3a97f95bb6e8e2c23752d27bbb84805f099dd65f1",
		LockTime: 0,
		Vin: []bchain.Vin{
			{
				Coinbase: "0385700200",
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
		},
		Blocktime: 1715377125,
		Time:      1715377125,
	}

	testTx2 = bchain.Tx{
		Hex:      "010000000135c95c18a72e7c2f68fd64947fd9b3e838c2c692d08fefeb90a5d7197ead890b010000006d483045022100bef5efa8a4cef888281f74713ca8c1e8e1313456a1751a35a74261a5e20486300220195fd7f88e634c0538b337bf4b0226f69fd86ad2b25e7fe5d83605023520715801015121028822c07871ae35acc2029159c9564115a0367d1e053ebef3eccc9a69ed595346ffffffff02000000000000000000481391d4020000003376a97b63d114f1eb7778da943aed6d1879765292ace8756d25b1671498833c37aabb0d983c94370ed5450ed55d41b6786888ac00000000",
		Txid:     "77640dafac13374e8afd2f63343e9114ab3ddd9c437d4fcc92b0b8cc2c4cdf5a",
		LockTime: 0,
		Vin: []bchain.Vin{
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "483045022100bef5efa8a4cef888281f74713ca8c1e8e1313456a1751a35a74261a5e20486300220195fd7f88e634c0538b337bf4b0226f69fd86ad2b25e7fe5d83605023520715801015121028822c07871ae35acc2029159c9564115a0367d1e053ebef3eccc9a69ed595346",
				},
				Txid:     "0b89ad7e19d7a590ebef8fd092c6c238e8b3d97f9464fd682f7c2ea7185cc935",
				Vout:     1,
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
				ValueSat: *big.NewInt(121562120400),
				N:        1,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "76a9145b1d583a4c270f2f14be77b298f0a9c6df97471388ac",
					Addresses: []string{
						"HLRYJtL33LLp6bDHohdZecYnuavvdq7cfT",
					},
				},
			},
		},
		Blocktime: 1715377125,
		Time:      1715377125,
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
				height:    409054,
				blockTime: 1562853038,
				parser:    NewHemisParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    testTxPacked1,
			wantErr: false,
		},
		{
			name: "hemis-2",
			args: args{
				tx:        testTx2,
				height:    409054,
				blockTime: 1562853038,
				parser:    NewHemisParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    testTxPacked2,
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
		time: 1715322960,
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
