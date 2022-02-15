/*
 * @Description:
 * @Autor: HTmonster
 * @Date: 2022-02-15 16:55:45
 */

package dict

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestConcurrentDictPut(t *testing.T) {
	var wg sync.WaitGroup

	dict := MakeConcurrentDict(16)
	max := 100
	wg.Add(max)
	for i := 0; i < max; i++ {
		go func(i int) {
			// key name
			key := "key" + strconv.Itoa(i)
			// try to put
			if ret := dict.Put(key, i); ret != 1 {
				t.Errorf("test to put key %s fail.", key)
			}
			// try to get
			val, ok := dict.Get(key)
			if !ok {
				t.Errorf("test to get key %s fail.", key)
			} else {
				intval := val.(int)
				if intval != i {
					t.Errorf("test fail to get different value: %d  excepted:%d", intval, i)
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestConcurrentDictRemove(t *testing.T) {
	var wg sync.WaitGroup

	dict := MakeConcurrentDict(16)
	max := 100
	wg.Add(max)
	for i := 0; i < max; i++ {
		go func(i int) {
			// key name
			key := "key" + strconv.Itoa(i)
			// try to put
			if ret := dict.Put(key, i); ret != 1 {
				t.Errorf("test to put key %s fail.", key)
			}
			// try to get
			val, ok := dict.Get(key)
			if !ok {
				t.Errorf("test to get key %s fail.", key)
			} else {
				intval := val.(int)
				if intval != i {
					t.Errorf("test fail to get different value: %d  excepted:%d", intval, i)
				}
			}
			// tyr to remove
			if ret := dict.Remove(key); ret != 1 {
				t.Errorf("test fail to remove key %s", key)
			}
			// try to get agian
			if _, ok := dict.Get(key); ok {
				t.Errorf("test fail: remove fail %s", key)
			}

			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestConcurrentDictLen(t *testing.T) {
	dict := MakeConcurrentDict(16)
	count := rand.Intn(1000)
	for i := 0; i < count; i++ {
		key := "key" + strconv.Itoa(i)
		// try to put
		if ret := dict.Put(key, i); ret != 1 {
			t.Errorf("test to put key %s fail.", key)
		}
	}

	if dict.Len() != count {
		t.Errorf("test len error, not same: %d %d", dict.Len(), count)
	}
}

func TestConcurrentDictKeys(t *testing.T) {
	dict := MakeConcurrentDict(16)
	count := rand.Intn(1000)

	keys := make([]string, count)
	for i := 0; i < count; i++ {
		key := "key" + strconv.Itoa(i)
		keys[i] = key
		// try to put
		if ret := dict.Put(key, i); ret != 1 {
			t.Errorf("test to put key %s fail.", key)
		}
	}

	keys2 := dict.Keys()
	for i := range keys2 {
		flag := false
		for j := range keys {
			if keys2[i] == keys[j] {
				keys = append(keys[:j], keys[j+1:]...)
				flag = true
				break
			}
		}
		if !flag {
			t.Errorf("test to get error key: %s", keys2[i])
		}
	}
}
