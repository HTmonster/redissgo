/*
 * @Description: dict
 * @Autor: HTmonster
 * @Date: 2022-02-15 11:13:56
 */
package dict

// used in traversal dict
type Consumer func(key string, value interface{}) bool

// dictionary interface
type dict interface {
	Get(key string) (value interface{}, exists bool)
	Put(key string, value interface{}) int
	Remove(key string) int
	Len() int
	ForEach(consumer Consumer)
	Keys() []string
	Clear()
}
