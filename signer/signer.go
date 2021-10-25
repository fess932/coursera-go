package main

import "log"

// сюда писать код
func ExecutePipeline(jobs ...job) {

	in := make(chan interface{}, 1)
	out := make(chan interface{}, 1)

	for _, j := range jobs {
		go j(in, out)
	}

}

func SingleHash(in, out chan interface{}) {
	for v := range in {
		hash := DataSignerCrc32(v.(string)) + "~" + DataSignerCrc32(DataSignerMd5(v.(string)))

		log.Println("kek", hash)
		out <- hash
	}
}

func MultiHash(in, out chan interface{}) {
	for v := range in {
		hash := DataSignerCrc32(v.(string)) + "~" + DataSignerCrc32(DataSignerMd5(v.(string)))

		log.Println("kek", hash)
		out <- hash
	}
}

func CombineResults(in, out chan interface{}) {
	for v := range in {
		hash := DataSignerCrc32(v.(string)) + "~" + DataSignerCrc32(DataSignerMd5(v.(string)))

		log.Println("kek", hash)
		out <- hash
	}
}
