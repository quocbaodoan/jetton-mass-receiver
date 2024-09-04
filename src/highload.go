package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/log"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/jetton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

func massSender(seedPhrase string, jettonMasterAddress, commentary, receiverAddress string) {
	client := liteclient.NewConnectionPool()

	// connect to testnet lite server
	err := client.AddConnectionsFromConfigUrl(context.Background(), "https://ton.org/global.config.json")
	if err != nil {
		log.Error(err)
		log.Error("Skipped!")
		return
	}

	ctx := client.StickyContext(context.Background())

	// initialize ton api lite connection wrapper
	api := ton.NewAPIClient(client)

	// seed words of account, you can generate them with any wallet or using wallet.NewSeed() method
	words := strings.Split(seedPhrase, " ")

	w, err := wallet.FromSeed(api, words, wallet.V4R2)
	if err != nil {
		log.Error("FromSeed err:", err.Error())
		log.Error("Skipped!")
		return
	}

	walletAddress := w.WalletAddress()

	token := jetton.NewJettonMasterClient(api, address.MustParseAddr(jettonMasterAddress))

	// find our jetton wallet
	tokenWallet, err := token.GetJettonWallet(ctx, w.WalletAddress())
	if err != nil {
		log.Error(err)
		log.Error("Skipped!")
		return
	}

	tokenBalance, err := tokenWallet.GetBalance(ctx)
	if err != nil {
		log.Error(err)
		log.Error("Skipped!")
		return
	}

	sendBalance := float64(tokenBalance.Int64()) / 1000000000.0

	log.Printf("Wallet address: %s | Jetton balance: %f", walletAddress, sendBalance)

	amountTokens := tlb.MustFromDecimal(fmt.Sprintf("%.9f", sendBalance), 9)

	comment, err := wallet.CreateCommentCell(commentary)
	if err != nil {
		log.Fatal(err)
	}

	if float64(amountTokens.Nano().Int64()) > 0 {
		// address of receiver's wallet (not token wallet, just usual)
		to := address.MustParseAddr(receiverAddress)
		transferPayload, err := tokenWallet.BuildTransferPayloadV2(to, to, amountTokens, tlb.ZeroCoins, comment, nil)
		if err != nil {
			log.Fatal(err)
		}

		// your TON balance must be > 0.05 to send
		msg := wallet.SimpleMessage(tokenWallet.Address(), tlb.MustFromTON("0.05"), transferPayload)

		log.Printf("Sending transaction with amount %f", sendBalance)
		tx, _, err := w.SendWaitTransaction(ctx, msg)
		if err != nil {
			panic(err)
		}
		log.Printf("Transaction sent, hash: %s", base64.StdEncoding.EncodeToString(tx.Hash))
		log.Printf("Explorer link: https://tonscan.org/tx/%s", base64.URLEncoding.EncodeToString(tx.Hash))

		log.Print("Sleeping for 5 seconds...")
		time.Sleep(5 * time.Second)
	} else {
		log.Printf("Wallet address: %s don't have Jetton balance to send!", walletAddress)
		time.Sleep(5 * time.Second)
	}
}
