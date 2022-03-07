/*
 * @Description: concurrent dictionary
 * @Autor: HTmonster
 * @Date: 2022-02-14 21:14:12
 */
package dict

import (
	"math"
	"sync"
	"sync/atomic"
)

const prime32 = uint32(16777619)

// Segment locking strategy

//Segment dict
type Segment struct {
	m      map[string]interface{}
	mutext sync.RWMutex
}

// concurrent dict
type ConcurrentDict struct {
	table    []*Segment
	count    int32
	segCount int
}

/**
 * @description: FNV hash
 * @param {string} key
 * @return {*} hash value
 */
func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

/**
 * @description: Compute ConcurrentDict Capacity
 * @param {*}
 * @return {*} capacity
 */
func computeCapacity(param int) int {
	//default capacity:16
	if param <= 16 {
		return 16
	}
	n := param - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16

	if n < 0 {
		return math.MaxInt32
	} else {
		return int(n + 1)
	}
}

/**
 * @description: make a new concurrent dictionary
 * @param {int} segCount
 * @return {*} concurrent dictionary
 */
func MakeConcurrentDict(segCount int) *ConcurrentDict {
	//0. compute capacity
	capacity := computeCapacity(segCount)
	//1. initial table
	table := make([]*Segment, capacity)
	//2. initial every Segment
	for i := 0; i < capacity; i++ {
		table[i] = &Segment{
			m: make(map[string]interface{}),
		}
	}
	//3. build the concurrent dictionary
	return &ConcurrentDict{
		table:    table,
		count:    0,
		segCount: segCount,
	}
}

/**
 * @description: spread
 * @event: h%n==(n-1)&h
 * @param {uint32} hashCode
 * @return {*}
 */
func (dict *ConcurrentDict) spread(hashCode uint32) uint32 {
	if dict == nil {
		panic("Error, dictionary is nil")
	}
	h := uint32(len(dict.table)) - 1
	return (h - 1) & uint32(hashCode)
}

/**
 * @description: add count atomicly
 */
func (dict *ConcurrentDict) addCount() int32 {
	return atomic.AddInt32(&dict.count, 1)
}

/**
 * @description: decrease count atomicly
 */
func (dict *ConcurrentDict) decreaseCount() int32 {
	return atomic.AddInt32(&dict.count, -1)
}

/**
 * @description: Given a index, return the Segment with check
 * @param {uint32} idx
 * @return {*}
 */
func (dict *ConcurrentDict) getSegment(idx uint32) *Segment {
	if dict == nil {
		panic("Error, dictionary is nil")
	}
	return dict.table[idx]
}

/**
 * @description: Given a key, get the value if exists.
 * @param {string} key
 * @return {*} value, exists or not
 */
func (dict *ConcurrentDict) Get(key string) (value interface{}, exists bool) {
	if dict == nil {
		panic("Error, dictionary is nil")
	}
	//0. get key hash value
	hashCode := fnv32(key)
	//1. get segment index
	index := dict.spread(hashCode)
	//2. get segment by index
	segment := dict.getSegment(index)
	//3. lock when read value
	segment.mutext.RLock()
	defer segment.mutext.RUnlock()
	//4. read value
	value, exists = segment.m[key]
	return value, exists
}

/**
 * @description: Put the key-value into the segment
 * @param {string} key
 * @param {interface{}} value
 * @return {*} 1 if inserted new one, else 0
 */
func (dict *ConcurrentDict) Put(key string, value interface{}) int {
	if dict == nil {
		panic("Error, dictionary is nil")
	}
	//0. get key hash value
	hashCode := fnv32(key)
	//1. get segment index
	index := dict.spread(hashCode)
	//2. get segment by index
	segment := dict.getSegment(index)
	//3. lock when write
	segment.mutext.Lock()
	defer segment.mutext.Unlock()
	//4. write segment
	//4.1 exists or not
	if _, ok := segment.m[key]; ok {
		segment.m[key] = value
		return 0
	} else {
		segment.m[key] = value
		dict.addCount()
		return 1 //add new one
	}
}

/**
 * @description: Remove a key from the dictionary
 * @param {string} key
 * @return {*} 1 if exists, 0 otherwise
 */
func (dict *ConcurrentDict) Remove(key string) int {
	if dict == nil {
		panic("Error, dictionary is nil")
	}
	//0. get key hash value
	hashCode := fnv32(key)
	//1. get segment index
	index := dict.spread(hashCode)
	//2. get segment by index
	segment := dict.getSegment(index)
	//3. lock when write
	segment.mutext.Lock()
	defer segment.mutext.Unlock()
	//4. remove
	if _, ok := segment.m[key]; ok {
		delete(segment.m, key)
		dict.decreaseCount()
		return 1
	}
	return 0
}

/**
 * @description: get the length of the dictionary in atomic way
 * @return {*} length
 */
func (dict *ConcurrentDict) Len() int {
	if dict == nil {
		panic("Error, dictionary is nil")
	}

	return int(atomic.LoadInt32(&dict.count))
}

/**
 * @description: traversal the dictonary
 * @param {*} consumer (operator)
 * @return {*}
 */
func (dict *ConcurrentDict) ForEach(consumer Consumer) {
	if dict == nil {
		panic("Error, dictionary is nil")
	}

	for _, segment := range dict.table {
		// lock when read
		segment.mutext.RLock()
		// operation depended on consumer
		func() {
			defer segment.mutext.RUnlock()
			for key, value := range segment.m {
				if c := consumer(key, value); !c {
					return
				}
			}

		}()
	}
}

/**
 * @description: return all keys
 */
func (dict *ConcurrentDict) Keys() []string {
	keys := make([]string, dict.Len())

	//traversal to get keys
	i := 0
	dict.ForEach(func(key string, value interface{}) bool {
		if i < len(keys) {
			keys[i] = key
			i++
		} else {
			keys = append(keys, key)
		}
		return true
	})

	return keys
}

/**
 * @description: clear all key and value
 * @param {*}
 * @return {*}
 */
func (dict *ConcurrentDict) Clear() {
	*dict = *MakeConcurrentDict(int(dict.segCount))
}
