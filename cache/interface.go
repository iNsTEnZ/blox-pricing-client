package cache

type ICache interface {
	Set(key interface{}, data interface{})
	Get(key interface{}) (interface{}, bool)
}
