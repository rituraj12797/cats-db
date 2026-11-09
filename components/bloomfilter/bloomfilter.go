package bloomfilter

type Key interface {
	int | float32 | float64
}

type Bloomfilter[K Key] struct {
	BitArray []K
}

func NewBloomFilter[K Key]() *Bloomfilter[K] {
	x := &Bloomfilter[K]{
		BitArray: make([]K, 10000),
	}
}

func (b *Bloomfilter[K]) Insert() struct {
}
