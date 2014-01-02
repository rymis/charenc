/* goconv package */
/* tables */

package charenc

import (
	"strings"
)

type pair struct {
	chr byte
	uchr rune
}

type tbls struct {
	name string
	to_ucs [256]rune
	from_ucs [256]pair
}

