package iputil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestIsTrusted(t *testing.T) {
	type want struct {
		err    error
		result bool
	}

	type args struct {
		ip            string
		trustedSubnet string
	}

	testCases := []struct {
		name string
		args args
		want want
	}{
		{
			name: "ip matches",
			args: args{
				ip:            "127.0.0.23",
				trustedSubnet: "127.0.0.1/24",
			},
			want: want{
				err:    nil,
				result: true,
			},
		},
		{
			name: "ip doesn't match",
			args: args{
				ip:            "123.123.123.123",
				trustedSubnet: "127.0.0.1/24",
			},
			want: want{
				err:    nil,
				result: false,
			},
		},
		{
			name: "invalid ip address",
			args: args{
				ip:            "invalid ip address",
				trustedSubnet: "127.0.0.1",
			},
			want: want{
				err:    fmt.Errorf("invalid ip address: %s", "invalid ip address"),
				result: false,
			},
		},
		{
			name: "empty ip address",
			args: args{
				ip:            "",
				trustedSubnet: "127.0.0.1",
			},
			want: want{
				err:    fmt.Errorf("invalid ip address: %s", ""),
				result: false,
			},
		},
		{
			name: "invalid trusted subnet",
			args: args{
				ip:            "127.0.0.1",
				trustedSubnet: "invalid trusted subnet",
			},
			want: want{
				err:    fmt.Errorf("invalid trusted subnet: %s, err: %s", "invalid trusted subnet", "invalid CIDR address: invalid trusted subnet"),
				result: false,
			},
		},
		{
			name: "trusted subnet is empty",
			args: args{
				ip:            "127.0.0.1",
				trustedSubnet: "",
			},
			want: want{
				err:    nil,
				result: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := IsTrusted(tc.args.ip, tc.args.trustedSubnet)

			if tc.want.err == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.want.err.Error())
			}

			assert.Equal(t, tc.want.result, res)
		})
	}
}

func TestIPFromRequest(t *testing.T) {
	ip := "127.0.0.1"

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)
	req.Header.Set(IPHeader, ip)

	assert.Equal(t, ip, IPFromRequest(req))
}
