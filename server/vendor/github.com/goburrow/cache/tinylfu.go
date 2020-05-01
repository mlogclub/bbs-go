package cache

import "container/list"

const (
	samplesMultiplier        = 8
	insertionsMultiplier     = 2
	countersMultiplier       = 1
	falsePositiveProbability = 0.1
	admissionRatio           = 0.01
)

// tinyLFU is an implementation of TinyLFU. It utilizing 4bit Count Min Sketch
// and Bloom Filter as a Doorkeeper and Segmented LRU for long term retention.
// See https://arxiv.org/pdf/1512.00727v2.pdf
type tinyLFU struct {
	filter  bloomFilter    // 1bit counter
	counter countMinSketch // 4bit counter

	additions int
	samples   int

	lru  lruCache
	slru slruCache
}

func (l *tinyLFU) init(c *cache, cap int) {
	if cap > 0 {
		// Only enable doorkeeper when capacity is finite.
		l.samples = samplesMultiplier * cap
		l.filter.init(insertionsMultiplier*cap, falsePositiveProbability)
		l.counter.init(countersMultiplier * cap)
	}
	lruCap := int(float64(cap) * admissionRatio)
	l.lru.init(c, lruCap)
	l.slru.init(c, cap-lruCap)
}

func (l *tinyLFU) add(en *entry) *entry {
	if l.lru.cap <= 0 {
		return l.slru.add(en)
	}
	l.increase(en.hash)
	candidate := l.lru.add(en)
	if candidate == nil {
		return nil
	}
	victim := l.slru.victim()
	if victim == nil {
		return l.slru.add(candidate)
	}
	// Determine one going to be evicted
	candidateFreq := l.estimate(candidate.hash)
	victimFreq := l.estimate(victim.hash)
	if candidateFreq > victimFreq {
		return l.slru.add(candidate)
	}
	return candidate
}

func (l *tinyLFU) hit(el *list.Element) {
	en := getEntry(el)
	l.increase(en.hash)
	if en.listID == admissionWindow {
		l.lru.hit(el)
	} else {
		l.slru.hit(el)
	}
}

func (l *tinyLFU) remove(el *list.Element) *entry {
	en := getEntry(el)
	if en.listID == admissionWindow {
		return l.lru.remove(el)
	}
	return l.slru.remove(el)
}

// increase adds the given hash to the filter and counter.
func (l *tinyLFU) increase(h uint64) {
	if l.samples <= 0 {
		return
	}
	l.additions++
	if l.additions >= l.samples {
		l.filter.reset()
		l.counter.reset()
		l.additions = 0
	}
	if l.filter.put(h) {
		l.counter.add(h)
	}
}

// estimate estimates frequency of the given hash value.
func (l *tinyLFU) estimate(h uint64) uint8 {
	freq := l.counter.estimate(h)
	if l.filter.contains(h) {
		freq++
	}
	return freq
}

func (l *tinyLFU) walk(f func(list *list.List)) {
	l.slru.walk(f)
	l.lru.walk(f)
}
