package utils

import (
	"log"
	"time"
)

func TimeFunc[P any, R any, F func(P) R](name string, fn F) F {
	return func(args P) R {
		start := time.Now()
		defer func() {
			log.Printf("[Timer] %s took %s", name, time.Since(start))
		}()
		return fn(args)
	}
}

func TimeFuncWithErr[P any, R any, F func(P) (R, error)](name string, fn F) F {
	return func(args P) (R, error) {
		start := time.Now()
		defer func() {
			log.Printf("[Timer] %s took %s", name, time.Since(start))
		}()
		return fn(args)
	}
}

func TimeFunc2WithErr[P1 any, P2, R any, F func(P1, P2) (R, error)](name string, fn F) F {
	return func(arg1 P1, arg2 P2) (R, error) {
		log.Printf("[Timer] Timing %s", name)
		start := time.Now()
		defer func() {
			log.Printf("[Timer] %s took %s", name, time.Since(start))
		}()
		return fn(arg1, arg2)
	}
}
