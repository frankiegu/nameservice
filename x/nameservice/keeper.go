package nameservice

// 处理同存储的交互
import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type Keeper struct {
	coinKeeper bank.Keeper // 控制账户转账
	storeKey   sdk.StoreKey
	cdc        *codec.Codec // 提供负责cosmos编码格式的工具
}

// 为指定域名设置解析字符串值
func (k Keeper) Setwhois(ctx sdk.Context, name string, whois Whois) {
	if whois.Owner.Empty() {
		return
	}

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(name),k.cdc.MustMarshalBinaryBare(whois))
}


// 解析域名
func (k Keeper) getWhois(ctx sdk.Context,name string) Whois  {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(name)) {
		return NewWhois()
	}
	bz := store.Get([]byte(name))
	var whois Whois
	k.cdc.MustUnmarshalBinaryBare(bz,&whois)
	return whois
}