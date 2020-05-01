package cache

import "container/list"

// lruCache is a LRU cache.
type lruCache struct {
	cache *cache
	cap   int
	ls    list.List
}

// init initializes cache list.
func (l *lruCache) init(c *cache, cap int) {
	l.cache = c
	l.cap = cap
	l.ls.Init()
}

// add addes new entry to the cache and returns evicted entry if necessary.
func (l *lruCache) add(en *entry) *entry {
	l.cache.mu.Lock()
	defer l.cache.mu.Unlock()

	el := l.cache.data[en.key]
	if el != nil {
		// Entry had been added
		setEntry(el, en)
		l.ls.MoveToFront(el)
		return nil
	}
	if l.cap <= 0 || l.ls.Len() < l.cap {
		// Add this entry
		el = l.ls.PushFront(en)
		l.cache.data[en.key] = el
		return nil
	}
	// Replace with the last one
	el = l.ls.Back()
	if el == nil {
		// Can happen if cap is zero
		return en
	}
	remEn := getEntry(el)
	setEntry(el, en)
	l.ls.MoveToFront(el)

	delete(l.cache.data, remEn.key)
	l.cache.data[en.key] = el
	return remEn
}

// hit updates cache entry for a get.
func (l *lruCache) hit(el *list.Element) {
	l.ls.MoveToFront(el)
}

// remove removes an entry from the cache.
func (l *lruCache) remove(el *list.Element) *entry {
	en := getEntry(el)
	l.cache.mu.Lock()
	defer l.cache.mu.Unlock()

	if _, ok := l.cache.data[en.key]; !ok {
		return nil
	}
	l.ls.Remove(el)
	delete(l.cache.data, en.key)
	return en
}

// walk walks thourgh all lists.
func (l *lruCache) walk(f func(list *list.List)) {
	f(&l.ls)
}

type listID int

const (
	admissionWindow listID = iota
	probationSegment
	protectedSegment
)

const (
	protectedRatio = 0.8
)

// slruCache is a segmented LRU.
// See http://highscalability.com/blog/2016/1/25/design-of-a-modern-cache.html
type slruCache struct {
	cache *cache

	probationCap int
	probationLs  list.List

	protectedCap int
	protectedLs  list.List
}

// init initializes the cache list.
func (l *slruCache) init(c *cache, cap int) {
	l.cache = c
	l.protectedCap = int(float64(cap) * protectedRatio)
	l.probationCap = cap - l.protectedCap
	l.probationLs.Init()
	l.protectedLs.Init()
}

// length returns total number of entries in the cache.
func (l *slruCache) length() int {
	return l.probationLs.Len() + l.protectedLs.Len()
}

// add addes new entry to the cache and returns evicted entry if necessary.
func (l *slruCache) add(en *entry) *entry {
	l.cache.mu.Lock()
	defer l.cache.mu.Unlock()

	el, ok := l.cache.data[en.key]
	if ok {
		// Copy value but preserve listID
		curEn := getEntry(el)
		en.listID = curEn.listID
		setEntry(el, en)
		l.accessNoLock(el)
		return nil
	}

	en.listID = probationSegment
	if l.probationCap <= 0 || l.probationLs.Len() < l.probationCap ||
		l.length() < (l.probationCap+l.protectedCap) {
		// probation list can exceed its capacity if number of entries
		// is still under total capacity
		el := l.probationLs.PushFront(en)
		l.cache.data[en.key] = el
		return nil
	}

	// Reuse the last entry in probation list
	el = l.probationLs.Back()
	if el == nil {
		// Can happen if cap is zero
		return en
	}
	remEn := getEntry(el)
	setEntry(el, en)
	l.probationLs.MoveToFront(el)
	delete(l.cache.data, remEn.key)
	l.cache.data[en.key] = el
	return remEn
}

// hit updates cache entry for a get.
func (l *slruCache) hit(el *list.Element) {
	en := getEntry(el)

	// Already in the protected segment
	if en.listID == protectedSegment {
		l.protectedLs.MoveToFront(el)
		return
	}

	// Promote this entry to the protected segment
	if l.protectedLs.Len() < l.protectedCap {
		l.cache.mu.Lock()
		l.probationLs.Remove(el)
		en.listID = protectedSegment
		l.cache.data[en.key] = l.protectedLs.PushFront(en)
		l.cache.mu.Unlock()
		return
	}

	// Swap with the last entry in the protected segment
	el2 := l.protectedLs.Back()
	if el2 == nil {
		return
	}
	en2 := getEntry(el2)

	l.cache.mu.Lock()
	setEntry(el, en2)
	setEntry(el2, en)

	// Update entry details
	en.listID = protectedSegment
	l.protectedLs.MoveToFront(el2)

	en2.listID = probationSegment
	l.probationLs.MoveToFront(el)

	l.cache.data[en.key] = el2
	l.cache.data[en2.key] = el
	l.cache.mu.Unlock()
}

// accessNoLock is similar to access but does not lock cache data.
func (l *slruCache) accessNoLock(el *list.Element) {
	en := getEntry(el)

	// Already in the protected segment
	if en.listID == protectedSegment {
		l.protectedLs.MoveToFront(el)
		return
	}

	// Promote this entry to the protected segment
	if l.protectedLs.Len() < l.protectedCap {
		l.probationLs.Remove(el)
		en.listID = protectedSegment
		l.cache.data[en.key] = l.protectedLs.PushFront(en)
		return
	}

	// Swap with the last entry in the protected segment
	el2 := l.protectedLs.Back()
	if el2 == nil {
		return
	}
	en2 := getEntry(el2)

	setEntry(el, en2)
	setEntry(el2, en)

	// Update entry details
	en.listID = protectedSegment
	l.protectedLs.MoveToFront(el2)

	en2.listID = probationSegment
	l.probationLs.MoveToFront(el)

	l.cache.data[en.key] = el2
	l.cache.data[en2.key] = el
}

// remove removes an entry from the cache and returns the removed entry or nil
// if it is not found.
func (l *slruCache) remove(el *list.Element) *entry {
	en := getEntry(el)
	l.cache.mu.Lock()
	defer l.cache.mu.Unlock()

	if _, ok := l.cache.data[en.key]; !ok {
		return nil
	}
	if en.listID == protectedSegment {
		l.protectedLs.Remove(el)
	} else {
		l.probationLs.Remove(el)
	}
	delete(l.cache.data, en.key)
	return en
}

// victim returns the last entry in probation list if total entries reached the limit.
func (l *slruCache) victim() *entry {
	if l.probationCap <= 0 || l.length() < (l.probationCap+l.protectedCap) {
		return nil
	}
	el := l.probationLs.Back()
	if el == nil {
		return nil
	}
	return getEntry(el)
}

// walk walks thourgh all lists.
func (l *slruCache) walk(f func(list *list.List)) {
	f(&l.protectedLs)
	f(&l.probationLs)
}
