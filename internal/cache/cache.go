package cache

type Cache interface {
	Set(key interface{}, value interface{}) (interface{}, bool)
	Get(key interface{}) (interface{}, bool)
	Len() int
	Clear()
}
