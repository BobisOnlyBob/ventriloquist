package proxytag

import (
	"unicode"
	"unicode/utf8"
)

// Shuck removes the first and last character of a string, analogous to
// shucking off the husk of an ear of corn.
func Shuck(victim string) string {
	return victim[1 : len(victim)-1]
}

func isSigil(inp rune) bool {
	switch inp {
	case ';', '.', '?', '!':
		return false
	}

	return unicode.IsSymbol(inp) || unicode.IsPunct(inp)
}

func firstRune(inp string) rune {
	for _, rn := range inp {
		return rn
	}

	return rune(0)
}

func lastRune(inp string) rune {
	var result rune
	for _, rn := range inp {
		result = rn
	}

	return result
}

// HalfSigilStart parses the "half sigil at the start" method of proxy tagging.
//
// Given a message of the form:
//
//     foo]
//
// This returns
//
//     Match{EndSigil:"]", Method: "HalfSigilEnd", Body: "foo"}
func HalfSigilEnd(message string) (Match, error) {
	if len(message) < 2 {
		return Match{}, ErrNoMatch
	}

	lst := lastRune(message)
	body := message[:len(message)-utf8.RuneLen(lst)]
	if !isSigil(lst) {
		return Match{}, ErrNoMatch
	}

	return Match{
		EndSigil: string(lst),
		Method:   "HalfSigilEnd",
		Body:     body,
	}, nil
}

// HalfSigilStart parses the "half sigil at the start" method of proxy tagging.
//
// Given a message of the form:
//
//     [foo
//
// This returns
//
//     Match{InitialSigil:"[", Method: "HalfSigils", Body: "foo"}
func HalfSigilStart(message string) (Match, error) {
	if len(message) < 2 {
		return Match{}, ErrNoMatch
	}

	fst := firstRune(message)
	body := message[utf8.RuneLen(fst):]
	if !isSigil(fst) {
		return Match{}, ErrNoMatch
	}

	return Match{
		InitialSigil: string(fst),
		Method:       "HalfSigilStart",
		Body:         body,
	}, nil
}

// Sigils parses the "sigils" method of proxy tagging.
//
// Given a message of the form:
//
//     [foo]
//
// This returns
//
//     Match{InitialSigil:"[", EndSigil: "]", Method: "Sigils", Body: "foo"}
func Sigils(message string) (Match, error) {
	if len(message) < 3 {
		return Match{}, ErrNoMatch
	}

	fst := firstRune(message)
	lst := lastRune(message)
	body := Shuck(message)

	// prevent mistakes like `[ <@72838115944828928>` being mis-read
	if fst != '<' && lst == '>' {
		return Match{}, ErrNoMatch
	}

	if !isSigil(fst) {
		return Match{}, ErrNoMatch
	}

	if !isSigil(lst) {
		return Match{}, ErrNoMatch
	}

	return Match{
		InitialSigil: string(fst),
		EndSigil:     string(lst),
		Method:       "Sigils",
		Body:         body,
	}, nil
}
