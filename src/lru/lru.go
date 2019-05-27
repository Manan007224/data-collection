// Implementation of LRU Cache which supports O(1) insertion and O(1) removal

package lru

import (
	"container/list"
)

type LRU struct {
	cap 				int 													// The max no of items LRU can hold
	cache 			map[interface{}]*list.Element 	// The cache for our items 
	evictList  *list.List 										// The acutal list holding our data
}


// An unexported field which we actually store in our cache
type entry struct {
	key, value interface{}
}

func New(cap int) *LRU {	
	lru_cache := &LRU {
		cap: cap,
		evictList: list.New(),
		cache: make(map[interface{}]*list.Element, cap+1),
	}
	return lru_cache
} 

// Used to automatically initialize cache without the New method for eg:
// var L LRU
// L.Add("a", 5)

func (this* LRU) lazyInit() {
	if this.evictList == nil {
		this.evictList = list.New()
		this.cache = make(map[interface{}]*list.Element, this.cap+1)
	}
}

func (this *LRU) Len() int {
	this.lazyInit()
	return this.evictList.Len()
} 

func (this *LRU) Add(k, v interface{}) {
	this.lazyInit()

	// If the item already exists
	if ent, ok := this.cache[k]; ok {
		ent.Value.(*entry).value = v
		this.evictList.MoveToFront(ent)
		return
	}

	this.cache[k] = this.evictList.PushFront(&entry{key: k, value: v})

	// If the capacity is full
	// Get the element which was least recently used from the evictList
	if this.cap > 0 && this.evictList.Len() > this.cap {
		this.removeOldest()
	}
	return
}


func (this *LRU) Get(k interface{}) (value interface{}, ok bool) {
	this.lazyInit()

	// Move the item at the head of the evictList
	if ent, ok := this.cache[k]; ok {
		this.evictList.MoveToFront(ent)
		return ent.Value.(*entry).value, true
	} else {
		return nil, false
	}
}

func (this *LRU) GetLatest() (k, v interface{}) {
	this.lazyInit()

	if this.evictList.Len() == 0 {
		return nil, nil
	}
	ent := this.evictList.Front().Value.(*entry)
	return ent.key, ent.value
}

func (this *LRU) remove(le *list.Element) (k, v interface{}) {
	k_v := le.Value.(*entry)
	this.evictList.Remove(le)
	delete(this.cache, k_v.key)
	return k_v.key, k_v.value	
}

func (this *LRU) Remove(k, v interface{}) {
	this.lazyInit()

	ent, ok := this.cache[k]
	if !ok {
		return
	}
	this.remove(ent)
}

func (this *LRU) removeOldest() (k, v interface{}) {
	le := this.evictList.Back()

	if le == nil {
		return
	}
	return this.remove(le)
}

