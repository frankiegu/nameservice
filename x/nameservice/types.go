package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}

type Whois struct {
	Value string         `json:"value"` // 域名解析出为的值
	Owner sdk.AccAddress `json:"owner"` // 该域名当前所有者的地址
	Price sdk.Coins      `json:"price"` // 你需要为购买域名支付的费用
}

func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
	}
}
