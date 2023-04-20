package mat

type Materijal int

const (
    Zid Materijal = -1
    Prazno Materijal = 0
    Pesak Materijal = 1
    Voda Materijal = 2
    Kamen Materijal = 3
    Metal Materijal = 4
)

var boja = map[Materijal]uint32{
    Zid : 0xff77ff,
    Prazno : 0x000000,
    Pesak : 0xffff66,
    Voda : 0x3333ff,
    Kamen : 0x666666,
    Metal : 0x33334b,
}