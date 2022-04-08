package types

import (
	"bytes"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// GravityDenomPrefix indicates the prefix for all assests minted by this module
	GravityDenomPrefix = ModuleName

	// GravityDenomSeparator is the separator for gravity denoms
	GravityDenomSeparator = ""

	// EthereumContractAddressLen is the length of contract address strings
	EthereumContractAddressLen = 42

	// GravityDenomLen is the length of the denoms generated by the gravity module
	GravityDenomLen = len(GravityDenomPrefix) + len(GravityDenomSeparator) + EthereumContractAddressLen
)

// EthereumAddrLessThan migrates the Ethereum address less than function
func EthereumAddrLessThan(e, o string) bool {
	return bytes.Compare([]byte(e)[:], []byte(o)[:]) == -1
}

// ValidateEthereumAddress validates the ethereum address strings
// func ValidateEthereumAddress(a string) error {
// 	if a == "" {
// 		return fmt.Errorf("empty")
// 	}
// 	if !regexp.MustCompile("^0x[0-9a-fA-F]{40}$").MatchString(a) {
// 		return fmt.Errorf("address(%s) doesn't pass regex", a)
// 	}
// 	if len(a) != EthereumContractAddressLen {
// 		return fmt.Errorf("address(%s) of the wrong length exp(%d) actual(%d)", a, len(a), EthereumContractAddressLen)
// 	}
// 	return nil
// }

/////////////////////////
//     ERC20Token      //
/////////////////////////

// NewERC20Token returns a new instance of an ERC20
func NewERC20Token(amount uint64, contract common.Address) ERC20Token {
	return ERC20Token{
		Amount:   sdk.NewIntFromUint64(amount),
		Contract: contract.Hex(),
	}
}

func NewSDKIntERC20Token(amount sdk.Int, contract common.Address) ERC20Token {
	return ERC20Token{
		Amount:   amount,
		Contract: contract.Hex(),
	}
}

func GravityDenom(contract common.Address) string {
	return strings.Join([]string{GravityDenomPrefix, contract.Hex()}, GravityDenomSeparator)
}

// GravityCoin returns the gravity representation of the ERC20
func (e ERC20Token) GravityCoin() sdk.Coin {
	return sdk.Coin{Amount: e.Amount, Denom: GravityDenom(common.HexToAddress(e.Contract))}
}

func GravityDenomToERC20(denom string) (string, error) {
	fullPrefix := GravityDenomPrefix + GravityDenomSeparator
	if !strings.HasPrefix(denom, fullPrefix) {
		return "", fmt.Errorf("denom prefix(%s) not equal to expected(%s)", denom, fullPrefix)
	}
	contract := strings.TrimPrefix(denom, fullPrefix)
	switch {
	case !common.IsHexAddress(contract):
		return "", fmt.Errorf("error validating ethereum contract address")
	case len(denom) != GravityDenomLen:
		return "", fmt.Errorf("len(denom)(%d) not equal to GravityDenomLen(%d)", len(denom), GravityDenomLen)
	default:
		return common.HexToAddress(contract).Hex(), nil
	}
}

func NormalizeCoinDenom(coin *sdk.Coin) {
	coin.Denom = NormalizeDenom(coin.Denom)
}

func NormalizeDenom(denom string) string {
	if contract, err := GravityDenomToERC20(denom); err == nil {
		return GravityDenom(common.HexToAddress(contract))
	}

	return denom
}

func NewSendToEthereumTx(id uint64, tokenContract common.Address, sender sdk.AccAddress, recipient common.Address, amount, feeAmount uint64) *SendToEthereum {
	return &SendToEthereum{
		Id:                id,
		Erc20Fee:          NewERC20Token(feeAmount, tokenContract),
		Sender:            sender.String(),
		EthereumRecipient: recipient.Hex(),
		Erc20Token:        NewERC20Token(amount, tokenContract),
	}
}

// Id:                2,
// Erc20Fee:          types.NewERC20Token(3, myTokenContractAddr),
// Sender:            mySender.String(),
// EthereumRecipient: myReceiver,
// Erc20Token:        types.NewERC20Token(101, myTokenContractAddr),
