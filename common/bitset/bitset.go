package bitset

import (
	"encoding/binary"
	"slices"
)

type BitSet struct {
	words []uint64
}

func (b *BitSet) Clone() *BitSet {
	cloned := BitSet{
		words: slices.Clone(b.words),
	}

	return &cloned
}

func (b *BitSet) Set(flag uint64) {
	word, bit := flag/64, flag%64

	for len(b.words) <= int(word) {
		b.words = append(b.words, 0)
	}

	b.words[word] |= 1 << bit
}

func (b *BitSet) Clear(flag uint64) {
	word, bit := flag/64, flag%64
	if int(word) < len(b.words) {
		b.words[word] &^= 1 << bit
	}
}

func (b *BitSet) Key() string {
	buf := make([]byte, len(b.words)*8)
	for i, word := range b.words {
		binary.LittleEndian.AppendUint64(buf[i*8:], word)
	}

	return string(buf)
}

func (b *BitSet) IsSet(flag uint64) bool {
	word, bit := flag/64, flag%64
	if int(word) >= len(b.words) {
		return false
	}
	return b.words[word]&(1<<bit) != 0
}

func (b *BitSet) HasAll(o *BitSet) bool {
	for i, word := range o.words {
		if i >= len(b.words) {
			return false
		}

		if b.words[i]&word != word {
			return false
		}
	}

	return true
}

func (b *BitSet) HasNone(o *BitSet) bool {
	for i, word := range o.words {
		if i >= len(b.words) {
			break
		}

		if b.words[i]&word != 0 {
			return false
		}
	}

	return true
}
