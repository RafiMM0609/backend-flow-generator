package main

import "fmt"

type Manusia interface {
	Bicara() string
}

type Guru struct {
	Materi string
}

type Programmer struct {
	Aplikasi string
}

func (g Guru) Bicara() string {
	return "Saya mengajar materi: " + g.Materi
}

func (p Programmer) Bicara() string {
	return "Saya sedang mengembangkan aplikasi: " + p.Aplikasi
}

func Kegiatan(m Manusia) string {
	fmt.Printf("%s\n", m.Bicara())
	return m.Bicara()
}

func main() {
	guru := Guru{Materi: "Sejarah"}
	programmer := Programmer{Aplikasi: "Pertamina"}
	Kegiatan(guru)
	Kegiatan(programmer)
}