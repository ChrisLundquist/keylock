package keylock

import (
	"fmt"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	_ = New(128)
	_ = New(256)
	_ = New(512)
}

func TestLocking(t *testing.T) {
	locker := New(128)
	locker.Lock("foo")
	locker.Lock("bar")
	locker.Lock("baz")

	locker.Unlock("baz")
	locker.Lock("baz")

	locker.Unlock("baz")
	locker.Unlock("bar")
	locker.Unlock("foo")
}

func TestGoLocking(t *testing.T) {
	var wg sync.WaitGroup
	locker := New(256)
	character_likes := map[string]int{
		"Jean-Luc Picard":  0,
		"William Riker":    0,
		"Deanna Troi":      0,
		"Beverly Crusher":  0,
		"Data":             0,
		"Geordi La Forge":  0,
		"Worf":             0,
		"Miles O'Brien":    0,
		"Ro Laren":         0,
		"James T. Kirk":    0,
		"Spock":            0,
		"Leonard McCoy":    0,
		"Montgomery Scott": 0,
		"Nyota Uhura":      0,
		"Hikaru Sulu":      0,
		"Pavel Chekov":     0,
	}

	wg.Add(1000)
	for r := 0; r < 1000; r++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				for name, _ := range character_likes {
					locker.Lock(name)
					character_likes[name] += 1
					locker.Unlock(name)
				}
			}
		}()
	}
	wg.Wait()

	wg.Add(100 * len(character_likes))
	for name, _ := range character_likes {
		for i := 0; i < 100; i++ {
			go func(name string) {
				defer wg.Done()
				for i := 0; i < 1000; i++ {
					locker.Lock(name)
					character_likes[name] += 1
					locker.Unlock(name)
				}
			}(name)
		}
	}
	wg.Wait()
	for name, value := range character_likes {
		locker.Lock(name)
		fmt.Println(name, value)
		locker.Unlock(name)
	}
}
