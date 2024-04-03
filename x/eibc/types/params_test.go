package types_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension/v3/x/eibc/types"
	"github.com/stretchr/testify/require"
)

func TestValidateErrAckFee(t *testing.T) {
	testCases := []struct {
		name        string
		input       any
		expectedErr error
	}{
		{
			name:        "valid fee",
			input:       sdk.NewDecWithPrec(5, 2), // 0.05
			expectedErr: nil,
		},
		{
			name:        "wrong type",
			input:       123,
			expectedErr: fmt.Errorf("invalid parameter type: %T", 123),
		},
		{
			name:        "nil value",
			input:       sdk.Dec{},
			expectedErr: fmt.Errorf("invalid global pool params: %+v", sdk.Dec{}),
		},
		{
			name:        "negative fee",
			input:       sdk.NewDec(-1),
			expectedErr: types.ErrNegativeErrAckFee,
		},
		{
			name:        "too much fee",
			input:       sdk.OneDec(), // 1
			expectedErr: types.ErrTooMuchErrAckFee,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := types.ValidateErrAckFee(tc.input)
			if tc.expectedErr != nil {
				require.EqualError(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
