// Alex Ray 2011 <ajray@ncsu.edu>
// Modifications by Matt Layher, 2013 <mdlayher@gmail.com>

// Reference:
// http://en.wikipedia.org/wiki/Bencode

package bencode

import "fmt"

// Bencode an Integer
func EncInt(i int) []byte {
	return []byte(fmt.Sprintf("i%de", i))
}

// Bencode a byte string
func EncBytes(a []byte) []byte {
	return []byte(fmt.Sprintf("%d:%s", len(a), a))
}

// Bencode a string
func EncString(s string) []byte {
	return EncBytes([]byte(s))
}

// Bencode a list of bencoded values
func EncList(list [][]byte) []byte {
	b := make([]byte, 0, 2)
	b = append(b, 'l')
	for _, value := range list {
		b = append(b, value...)
	}
	b = append(b, 'e')
	return b
}

// Bencode a dictionary of key:value pairs given as a list
func EncDict(dict [][]byte) []byte {
	b := make([]byte, 0, 2)
	b = append(b, 'd')
	for _, value := range dict {
		b = append(b, value...)
	}
	b = append(b, 'e')
	return b
}

// Bencode a dictionary of key:value pairs given as a map
func EncDictMap(dict map[string][]byte) []byte {
	b := make([]byte, 0, 2)
	b = append(b, 'd')
	for key, value := range dict {
		b = append(b, EncBytes([]byte(key))...)
		b = append(b, value...)
	}
	b = append(b, 'e')
	return b
}
