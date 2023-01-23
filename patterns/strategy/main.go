package main

import "fmt"

// Стратегия -Паттерн Стратегия предлагает определить семейство схожих алгоритмов, которые часто изменяются или расширяются, и вынести их в собственные классы, называемые стратегиями.
/*Вместо того, чтобы изначальный класс сам выполнял тот или иной алгоритм, он будет играть роль контекста, ссылаясь на одну из стратегий и делегируя ей выполнение работы. Чтобы сменить алгоритм, вам будет достаточно подставить в контекст другой объект-стратегию.*/
/*Плюсы:  Горячая замена алгоритмов на лету.
Изолирует код и данные алгоритмов от остальных классов.
Уход от наследования к делегированию.
Реализует принцип открытости/закрытости.*/
/*Минусы:  Усложняет программу за счёт дополнительных классов.
Клиент должен знать, в чём состоит разница между стратегиями, чтобы выбрать подходящую.*/
// Интерфейс стратегии
type EvictionAlgo interface {
	evict(c *Cache)
}

// Конкретная стратегия
type Fifo struct {
}

func (l *Fifo) evict(c *Cache) {
	fmt.Println("Evicting by fifo strаtegy")
}

// Конкретная стратегия
type Lru struct {
}

func (l *Lru) evict(c *Cache) {
	fmt.Println("Evicting by lru strаtegy")
}

type Lfu struct {
}

func (l *Lfu) evict(c *Cache) {
	fmt.Println("Evicting by lfu strаtegy")
}

// Кэш
type Cache struct {
	storage      map[string]string
	evictionAlgo EvictionAlgo
	capacity     int
	maxCapacity  int
}

func initCache(e EvictionAlgo) *Cache {
	storage := make(map[string]string)
	return &Cache{
		storage:      storage,
		evictionAlgo: e,
		capacity:     0,
		maxCapacity:  2,
	}
}
func (c *Cache) setEvictionAlgo(e EvictionAlgo) {
	c.evictionAlgo = e
}
func (c *Cache) add(key, value string) {
	if c.capacity == c.maxCapacity {
		c.evict()
	}
	c.capacity++
	c.storage[key] = value
}
func (c *Cache) get(key string) {
	delete(c.storage, key)
}
func (c *Cache) evict() {
	c.evictionAlgo.evict(c)
	c.capacity--
}
func main() {
	lfu := &Lfu{}
	cache := initCache(lfu)

	cache.add("a", "1")
	cache.add("b", "2")

	cache.add("c", "3")

	lru := &Lru{}
	cache.setEvictionAlgo(lru)

	cache.add("d", "4")

	fifo := &Fifo{}
	cache.setEvictionAlgo(fifo)

	cache.add("e", "5")

}
