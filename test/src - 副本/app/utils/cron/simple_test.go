package cron

import (
	"log"
	"testing"
)

func TestCron(t *testing.T) {
	i := 0
	c := New()
	spec := "*/5 * * * * ?"
	c.AddFunc(spec, func() {
		i++
		log.Println("cron running:", i)
	})
	c.Start()

	select {}
}
