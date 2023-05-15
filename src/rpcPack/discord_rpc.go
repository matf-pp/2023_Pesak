// Package rpcPack sluzi za Discord rich presence u Pesku
package rpcPack

import (
	"strings"
	"time"

	"github.com/sajberk/rich-go/client"
)

// ConnectToDiscord sluzi za povezivanje na Discord i postavljanje odredjenih labela u Discordu
func ConnectToDiscord() {
	client.Login("1100118057147437207") // ovo zapravo pravi probleme i ne radi kad si oflajn
	UpdateRPC("pesak")
}

func UpdateRPC(materijal string) {
	now := time.Now()
	client.SetActivity(client.Activity{
		State:      "bleja u pesku",
		Details:    "sipa se " + strings.ToLower(materijal),
		LargeImage: "bleja",
		LargeText:  "je l se učitalo ovo", //xDDD -s //;D -nj //:3 -l
		Timestamps: &client.Timestamps{
			Start: &now,
		},
		Buttons: []*client.Button{
			{
				Label: "priključi se",
				Url:   "https://github.com/matf-pp/2023_Pesak",
			},
		},
	})
}
