package screenPack

import (
	stego "github.com/zhcppy/steganography"
)

// ovo inače ne radi pa ću forkovati repo da ga sredim i dodaću kasnije šta fali idemo po 00 00 release - njanja
func encdec(encode bool, text, inputpath string) {
	if encode {
		err := stego.Encode(&stego.TextCarrier{CarrierFileName: inputpath, TextContent: text})
		if err != nil {
			panic(err)
		}
	} else {
		err := stego.Decode(&stego.TextCarrier{CarrierFileName: inputpath})
		if err != nil {
			panic(err)
		}
	}
}
