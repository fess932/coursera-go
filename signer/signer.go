package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

func main() {
	inputData := []int{0, 1, 1, 2, 3, 5, 8}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			<-in
		}),
	}

	start := time.Now()

	ExecutePipeline(hashSignJobs...)

	end := time.Since(start)

	log.Println("time:", end)
}

func ExecutePipeline(jobs ...job) {

	in := make(chan interface{})
	out := make(chan interface{})

	go jobs[0](make(chan interface{}), in)

	var wg sync.WaitGroup
	go runProcessJobs(&wg, out, jobs[1:]...)

	func() {
	loop:
		for {
			select {
			case v := <-in:
				print("ok")
				out <- v
			case <-time.After(time.Millisecond * 10):
				fmt.Println("No input, close pipeline")
				close(in)
				break loop
			}
		}
	}()

	wg.Wait()
}

func runProcessJobs(wg *sync.WaitGroup, in chan interface{}, jobs ...job) {
	for _, j := range jobs {
		in = run(wg, in, j)
	}
}

func run(wg *sync.WaitGroup, in chan interface{}, j job) (out chan interface{}) {
	out = make(chan interface{})

	wg.Add(1)
	go func(j job, in, out chan interface{}, wg *sync.WaitGroup) {
		j(in, out)
		wg.Done()
	}(j, in, out, wg)

	return out
}

func SingleHash(in, out chan interface{}) {
	for {
		i, ok := <-in
		if !ok {
			log.Println("closed channel")
			break
		}

		fmt.Printf("Value %d was read.\n", i)

		var data string

		switch d := i.(type) {
		case string:
			data = d
		case int:
			data = strconv.Itoa(d)
		default:
			log.Println("wrong type assertion single hash:", i)
		}

		log.Println("data:", data)

		var wg sync.WaitGroup

		var p1, p2 string

		wg.Add(1)
		go func() {
			p1 = DataSignerCrc32(data) // gorutine
			wg.Done()
		}()

		wg.Add(1)

		time.Sleep(time.Millisecond * 10)
		h2 := DataSignerMd5(data) // gorutine

		go func() {
			p2 = DataSignerCrc32(h2) // gorutine
			wg.Done()
		}()

		wg.Wait()

		hash := p1 + "~" + p2
		out <- hash
		log.Println("single", hash)
	}

	close(out)
}

func MultiHash(in, out chan interface{}) {
	for {
		i, ok := <-in
		if !ok {
			log.Println("closed channel")
			break
		}

		data, ok := i.(string)
		if !ok {
			log.Println("wrong type assertion")
		}

		var wg sync.WaitGroup

		var total [6]string

		for i := 0; i < 6; i++ {
			wg.Add(1)
			go func(id int) {
				hash := DataSignerCrc32(strconv.Itoa(id) + data)
				total[id] = hash
				wg.Done()
			}(i)
		}

		wg.Wait()

		var td string

		for _, v := range total {
			td += v
		}

		log.Println("multi", td)
		out <- td
	}

	close(out)
}

func CombineResults(in, out chan interface{}) {
	total := ""
	for v := range in {
		data, ok := v.(string)
		if !ok {
			log.Println("wrong type assertion combine results")
		}

		if total == "" {
			total += data
		} else {
			total += "_" + data
		}

	}
	log.Println("combine", total)
	out <- total
}
