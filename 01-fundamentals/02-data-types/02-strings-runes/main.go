// Strings and runes: a string is bytes, and a character can be several.
//
// Go stores a string as UTF-8 bytes. An English letter is one byte. A
// Bangla letter is three. len(s) counts bytes, not characters. A rune is
// one Unicode character, and utf8.RuneCountInString counts those.
//
// s[i] gives you one byte. "for i, r := range s" gives you one rune at a
// time, and i is the byte where it starts.
//
// A profile page shows the first few characters of a user's name.
//
// Problem: shortByBytes cuts with name[:maxShown], so it cuts bytes. That
// works for "Asif", so the bug hides in tests. For "আসিফ" it cuts a letter
// in half, and the page shows a broken character.
//
// Fix: shortByRunes turns the string into []rune first, then cuts. main
// runs both, so you can compare.
//
// Note: some letters on screen are more than one rune, like "সি", which is
// স plus a vowel sign. Cutting by runes never breaks a rune, but it can
// split those.
//
// Real use:
//   - a "max 20 characters" check on a signup form: count runes with
//     utf8.RuneCountInString, not len.
//   - cutting display names, titles or SMS text to a length.
//
// Run: go run ./01-fundamentals/02-data-types/02-strings-runes
package main

import (
	"fmt"
	"unicode/utf8"
)

// maxShown is how many characters of a name the page shows.
const maxShown = 2

func main() {
	english := "Asif"
	name := "আসিফ"
	fmt.Printf("%s: %d bytes, %d runes\n", english, len(english), utf8.RuneCountInString(english))
	fmt.Printf("%s: %d bytes, %d runes\n", name, len(name), utf8.RuneCountInString(name))

	fmt.Println("first byte of", name+":", name[0])
	for i, r := range name {
		fmt.Printf("rune %c (%U) starts at byte %d\n", r, r, i)
	}

	fmt.Println("problem, cut by bytes:")
	show(english, shortByBytes(english))
	show(name, shortByBytes(name))

	fmt.Println("fix, cut by runes:")
	show(english, shortByRunes(english))
	show(name, shortByRunes(name))
}

// show prints a name next to its short form, and whether the short form
// is still valid text.
func show(name, short string) {
	fmt.Printf("  %s -> %q, valid UTF-8: %t\n", name, short+"...", utf8.ValidString(short))
}

// shortByBytes keeps the first maxShown bytes. It can split a letter.
func shortByBytes(name string) string {
	return name[:maxShown]
}

// shortByRunes keeps the first maxShown runes, so no rune is broken.
func shortByRunes(name string) string {
	return string([]rune(name)[:maxShown])
}
