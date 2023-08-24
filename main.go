package main

import (
	"context"
	"fmt"
	ttypes "github.com/tendermint/tendermint/types"
	"log"

	// Importing the general purpose Cosmos blockchain client
	"github.com/ignite/cli/ignite/pkg/cosmosclient"

	// Importing the types package of your blog blockchain
	"github.com/Murphy-hub/hello/x/hello/types"
)

func main() {
	ctx := context.Background()
	// Create a Cosmos client instance
	client, err := cosmosclient.New(ctx,
		cosmosclient.WithNodeAddress("http://localhost:26657"),
	)

	// 定义一个账户名称
	accountName := "jerry"
	account, err := client.Account(accountName)
	if err != nil {
		log.Fatal(err)
	}
	addr, err := client.Address(accountName)
	if err != nil {
		log.Fatal(err)
	}

	// 创建 kv 消息
	msg := &types.MsgCreateKv{
		Creator: addr,
		Index:   "username",
		Value:   "I'm Jerry",
	}

	// 广播交易
	txResp, err := client.BroadcastTx(ctx, account, msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("MsgCreateKv:\n\n")
	fmt.Println(txResp)

	// 查询Kv
	queryClient := types.NewQueryClient(client.Context())
	queryResp, err := queryClient.KvAll(ctx, &types.QueryAllKvRequest{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("All Kv:")
	fmt.Println(queryResp)
}

// CreateAccount 创建一个客户端账号
func CreateAccount(client *cosmosclient.Client, accountName string) {
	_, mnemonic, err := client.AccountRegistry.Create(accountName)
	if err != nil {
		log.Fatal(err)
	}
	addr, err := client.Address(accountName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("addr:")
	fmt.Println(addr)
	fmt.Println("mnemonic:")
	fmt.Println(mnemonic)

	// TODO 创建完账号后，需要使用水龙头给这个账号发送token,才可使用
	// curl -X POST "http://localhost:4500/" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"address\": \"cosmos164avpya9w5lvva5fyrp708lpluvzszhy45g80j\",  \"coins\": [    \"10token\"  ]}"
}

// Subscribe 订阅事件
func Subscribe(client *cosmosclient.Client) {
	_ = client.RPC.Start()
	const subscriber = "TestCreateKvEvents"
	eventCh, err := client.RPC.Subscribe(context.Background(), subscriber, ttypes.QueryForEvent(ttypes.EventTx).String())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.RPC.UnsubscribeAll(context.Background(), subscriber); err != nil {
			log.Fatal(err)
		}
		_ = client.RPC.Stop()
	}()
	for {
		event := <-eventCh
		txEvent, ok := event.Data.(ttypes.EventDataTx)
		if ok {
			// txEvent.Result.Events
			fmt.Println(txEvent.String())
		}
	}
}
