//go:build unittest

package eth

import "testing"

func Test_parseSimpleStringProperty(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "1",
			args: "0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000758504c4f44444500000000000000000000000000000000000000000000000000",
			want: "XPLODDE",
		},
		{
			name: "2",
			args: "0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000022426974436c617665202d20436f6e73756d657220416374697669747920546f6b656e00000000000000",
			want: "BitClave - Consumer Activity Token",
		},
		{
			name: "short",
			args: "0x44616920537461626c65636f696e2076312e3000000000000000000000000000",
			want: "Dai Stablecoin v1.0",
		},
		{
			name: "short2",
			args: "0x44616920537461626c65636f696e2076312e3020444444444444444444444444",
			want: "Dai Stablecoin v1.0 DDDDDDDDDDDD",
		},
		{
			name: "long",
			args: "0x556e6973776170205631000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			want: "Uniswap V1",
		},
		{
			name: "garbage",
			args: "0x2234880850896048596206002535425366538144616734015984380565810000",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseSimpleStringProperty(tt.args)
			// the addresses could have different case
			if got != tt.want {
				t.Errorf("parseSimpleStringProperty = %v, want %v", got, tt.want)
			}
		})
	}
}
