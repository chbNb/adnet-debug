package expired_map

import "sync"

type ConccurentLock struct {
	m sync.Map
}

func CreateCoccurrentLock() *ConccurentLock {
	return new(ConccurentLock)
}

func (k *ConccurentLock) TryLock(key interface{}) bool {
	_, ok := k.m.LoadOrStore(key, struct{}{})
	return !ok
}

func (k *ConccurentLock) UnLock(key interface{}) {
	k.m.Delete(key)
}
