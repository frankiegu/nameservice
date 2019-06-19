package nameservice

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"strings"
)

const (
	QueryResolve = "resolve" // 传入一个域名,返回nameservice给定的解析值,类似dns查询
	QueryWhois   = "whois"   // 传入一个域名,返回价格,用于确定想要购买名称的成本
	QueryNames   = "names"
)

type QueryResResolve struct {
	Value string `json:"value"`
}

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryResolve:
			return queryResolve(ctx, path[1:], req, keeper)
		case QueryWhois:
			return queryWhois(ctx, path[1:], req, keeper)
		case QueryNames:
			return queryNames(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

func queryResolve(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	name := path[0]

	value := keeper.ResolveName(ctx, name)
	if value == "" {
		return []byte{}, sdk.ErrUnknownRequest("could not resolve name")
	}
	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, QueryResResolve{value})
	if err2 != nil {
		panic("could not marshal result to json")
	}
	return bz, nil

}

func queryWhois(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	name := path[0]

	whois := keeper.getWhois(ctx, name)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, whois)
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

func queryNames(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var nameList QueryResNames

	iterator := keeper.GetNameIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		name := string(iterator.Key())
		nameList = append(nameList, name)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, nameList)
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

func (r QueryResResolve) String() string {
	return r.Value
}

func (w Whois) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Value:%s
Price:%s`, w.Owner, w.Value, w.Price))
}

type QueryResNames []string

func (n QueryResNames) String() string {
	return strings.Join(n[:], "\n")
}
