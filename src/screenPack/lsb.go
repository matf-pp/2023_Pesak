package screenPack

import (
	stego "github.com/sajberk/steganography"
)

// njanja: ne pitajte me ništa razdvojiću ovo u dve fje kad se naspavam
func encdec(encode bool, text, inputpath string) (msg string) {

	if encode {
		err := stego.Encode(&stego.TextCarrier{CarrierFileName: inputpath, TextContent: text})
		if err != nil {
			panic(err)
		}
		return ""
	}

	msg, err := stego.Decode(&stego.TextCarrier{CarrierFileName: inputpath})
	if err != nil {
		panic(err)
	}
	return msg

}
