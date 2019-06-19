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

func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// 为指定域名设置解析字符串值
func (k Keeper) Setwhois(ctx sdk.Context, name string, whois Whois) {
	if whois.Owner.Empty() {
		return
	}

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(whois))
}

// 解析域名
func (k Keeper) getWhois(ctx sdk.Context, name string) Whois {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(name)) {
		return NewWhois()
	}
	bz := store.Get([]byte(name))
	var whois Whois
	k.cdc.MustUnmarshalBinaryBare(bz, &whois)
	return whois
}

func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	return k.getWhois(ctx, name).Value
}

func (k Keeper) SetName(ctx sdk.Context, name, value string) {
	whois := k.getWhois(ctx, name)
	whois.Value = value
	k.Setwhois(ctx, name, whois)
}

func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
	return !k.getWhois(ctx, name).Owner.Empty()
}

func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
	return k.getWhois(ctx, name).Owner
}

func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
	whois := k.getWhois(ctx, name)
	whois.Owner = owner
	k.Setwhois(ctx, name, whois)
}

func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	return k.getWhois(ctx, name).Price
}

func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	whois := k.getWhois(ctx, name)
	whois.Price = price
	k.Setwhois(ctx, name, whois)
}

// 遍历获取所有已经域名
func (k Keeper) GetNameIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}
