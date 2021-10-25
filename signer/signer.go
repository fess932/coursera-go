package main

import "log"

// сюда писать код
func ExecutePipeline(jobs ...job) {

	in := make(chan interface{})
	out := make(chan interface{})

	for _, j := range jobs {
		j(in, out)
	}

}

func SingleHash(in, out chan interface{}) {
	for v := range in {
		hash := DataSignerCrc32(v.(string))
		log.Println("kek", hash)
		out <- hash
	}
}

func MultiHash(in, out chan interface{}) {

}

func CombineResults(in, out chan interface{}) {

}
