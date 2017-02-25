package model

import (
	"fmt"
	"hash"
)

type HashBytes []byte

func (h *HashBytes) Trim(bitlen uint) {
	size := bitlen / 8
	if bitlen%8 != 0 {
		size++
		(*h)[size-1] &= 0xFF << (8 - bitlen%8)
	}
	for i := size; i < uint(len(*h)); i++ {
		(*h)[i] = 0x00
	}
}

func (h HashBytes) Print() {
	fmt.Printf("%x\n", h)
}

func (h HashBytes) Sum(ctx hash.Hash) HashBytes {
	ctx.Reset()
	ctx.Write(h)
	return ctx.Sum(nil)
}
