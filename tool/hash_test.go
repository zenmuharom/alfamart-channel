package tool

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckHash(t *testing.T) {
	type TestCase struct {
		Signature string
		IsMD5     bool
		IsSha256  bool
	}

	testCases := []TestCase{
		{
			Signature: "",
			IsMD5:     false,
			IsSha256:  false,
		},
		{
			Signature: "c0af658015ade859afce40a0c7b6ef3d",
			IsMD5:     true,
			IsSha256:  false,
		},
		{
			Signature: "954a994e941a28a8faa0ba27699ebdc6992fb4b7bdcac313a8d7513052ba6301",
			IsMD5:     false,
			IsSha256:  true,
		},
	}

	for _, tc := range testCases {
		md5Res := IsMD5Hash(tc.Signature)
		sha256Res := IsSHA256Hash(tc.Signature)
		require.Equal(t, tc.IsMD5, md5Res, fmt.Sprintf("MD5 expectation:%v but actual:%v", tc.IsMD5, md5Res))
		require.Equal(t, tc.IsSha256, sha256Res, fmt.Sprintf("Sha256 expectation:%v but actual:%v", tc.IsSha256, sha256Res))
	}
}
