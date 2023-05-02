//discord rich presence u Pesku
package rpcPack

import (
	"time"

	"github.com/sajberk/rich-go/client"
)

func ConnectToDiscord() {
	client.Login("1100118057147437207") // ovo zapravo pravi probleme i ne radi kad si oflajn

	now := time.Now()
	client.SetActivity(client.Activity{
		State:      "bleja u pesku",
		Details:    "sipa pesak",
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
