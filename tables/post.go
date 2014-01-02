/* end of codepage tables */

func isalpha(c byte) bool {
	return (c >= 97 && c <= 122) || (c >= 65 && c <= 90)
}

func isnumber(c byte) bool {
	return (c >= 48 && c <= 58)
}

func isalnum(c byte) bool {
	return isalpha(c) || isnumber(c)
}

// Get internal number of this encoding:
func Open8bit(encoding string) int {
	// We must convert name special way:
	benc := make([]byte, len(encoding))
	copy(benc, encoding)
	for i := range(benc) {
		if !isalnum(benc[i]) {
			benc[i] = '_'
		}
	}
	enc := strings.ToUpper(string(benc))

	for i := range(names) {
		if strings.ToUpper(names[i].name) == enc {
			return i
		}
	}

	return -1
}

// Convert one byte to rune in specific encoding (0 is error):
func ByteToRune(codec int, b byte) rune {
	if codec < 0 || codec >= len(names) {
		return 0
	}

	return names[codec].to_ucs[b]
}

// Convert one rune to byte in specific encoding (0 is error):
func RuneToByte(codec int, ch rune) byte {
	var c int

	if codec < 0 || codec >= len(names) {
		return 0
	}

	a := 0
	b := len(names) - 1

	if names[codec].from_ucs[a].uchr > ch || names[codec].from_ucs[b].uchr < ch {
		return 0
	}

	for b - a > 1 {
		c = (a + b) / 2
		if names[codec].from_ucs[c].uchr < ch {
			a = c
		} else {
			b = c
		}
	}

	if names[codec].from_ucs[a].uchr == ch {
		return names[codec].from_ucs[a].chr
	}

	if names[codec].from_ucs[b].uchr == ch {
		return names[codec].from_ucs[b].chr
	}

	return 0
}

