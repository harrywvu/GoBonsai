package main

import (
	"math"
	"math/bits"
)

// pythonRandom implements the subset of Python's random.Random used by
// PyBonsai. Keeping the same generator and sampling algorithms makes --seed
// produce the same tree in both implementations.
type pythonRandom struct {
	state [624]uint32
	index int
}

func newPythonRandom(seed int64) *pythonRandom {
	r := &pythonRandom{}
	r.Seed(seed)
	return r
}

func (r *pythonRandom) Seed(seed int64) {
	var value uint64
	if seed < 0 {
		value = uint64(-(seed + 1)) + 1
	} else {
		value = uint64(seed)
	}

	key := []uint32{uint32(value)}
	if value>>32 != 0 {
		key = append(key, uint32(value>>32))
	}

	r.initByArray(key)
}

func (r *pythonRandom) initByArray(key []uint32) {
	r.initGenRand(19650218)

	i, j := 1, 0
	k := len(r.state)
	if len(key) > k {
		k = len(key)
	}
	for ; k > 0; k-- {
		prev := r.state[i-1]
		r.state[i] = (r.state[i] ^ ((prev ^ (prev >> 30)) * 1664525)) + key[j] + uint32(j)
		i++
		j++
		if i >= len(r.state) {
			r.state[0] = r.state[len(r.state)-1]
			i = 1
		}
		if j >= len(key) {
			j = 0
		}
	}

	for k = len(r.state) - 1; k > 0; k-- {
		prev := r.state[i-1]
		r.state[i] = (r.state[i] ^ ((prev ^ (prev >> 30)) * 1566083941)) - uint32(i)
		i++
		if i >= len(r.state) {
			r.state[0] = r.state[len(r.state)-1]
			i = 1
		}
	}

	r.state[0] = 0x80000000
	r.index = len(r.state)
}

func (r *pythonRandom) initGenRand(seed uint32) {
	r.state[0] = seed
	for i := 1; i < len(r.state); i++ {
		prev := r.state[i-1]
		r.state[i] = 1812433253*(prev^(prev>>30)) + uint32(i)
	}
	r.index = len(r.state)
}

func (r *pythonRandom) uint32() uint32 {
	const (
		upperMask = uint32(0x80000000)
		lowerMask = uint32(0x7fffffff)
		matrixA   = uint32(0x9908b0df)
	)

	if r.index >= len(r.state) {
		for i := range r.state {
			y := (r.state[i] & upperMask) | (r.state[(i+1)%len(r.state)] & lowerMask)
			r.state[i] = r.state[(i+397)%len(r.state)] ^ (y >> 1)
			if y&1 != 0 {
				r.state[i] ^= matrixA
			}
		}
		r.index = 0
	}

	y := r.state[r.index]
	r.index++

	y ^= y >> 11
	y ^= (y << 7) & 0x9d2c5680
	y ^= (y << 15) & 0xefc60000
	y ^= y >> 18
	return y
}

func (r *pythonRandom) Float64() float64 {
	a := uint64(r.uint32() >> 5)
	b := uint64(r.uint32() >> 6)
	return float64(a*67108864+b) / 9007199254740992
}

func (r *pythonRandom) getRandBits(k int) uint32 {
	if k <= 0 || k > 32 {
		panic("getRandBits supports between 1 and 32 bits")
	}
	return r.uint32() >> (32 - k)
}

func (r *pythonRandom) IntN(n int) int {
	if n <= 0 {
		panic("invalid argument to IntN")
	}

	k := bits.Len(uint(n))
	for {
		value := int(r.getRandBits(k))
		if value < n {
			return value
		}
	}
}

func (r *pythonRandom) NormFloat64() float64 {
	const normalMagic = 1.7155277699214135

	for {
		u1 := r.Float64()
		u2 := 1 - r.Float64()
		z := normalMagic * (u1 - 0.5) / u2
		zz := z * z / 4
		if zz <= -math.Log(u2) {
			return z
		}
	}
}
