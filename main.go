package main

import "sync"

type Data struct {
	number int
	square float64
}

func ProcessDataWithNRoutines(numbers []int) []*Data {
	outputCh := make(chan *Data)
	var wg sync.WaitGroup
	for _, v := range numbers {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			outputCh <- &Data{
				number: v,
				square: float64(v) * float64(v),
			}
		}(v)
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()

	output := make([]*Data, 0)

	for d := range outputCh {
		output = append(output, d)
	}

	return output
}

func ProcessDataWithWorkerPool(numbers []int) []*Data {
	output := make([]*Data, 0)
	outputCh := make(chan *Data)
	workerCount := 5
	var wg sync.WaitGroup

	workerCh := publishData(numbers)

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for v := range workerCh {
				outputCh <- &Data{
					number: v,
					square: float64(v) * float64(v),
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()

	for d := range outputCh {
		output = append(output, d)
	}

	return output
}

func publishData(numbers []int) <-chan int {
	workerCh := make(chan int)
	go func() {
		for _, v := range numbers {
			workerCh <- v
		}
		close(workerCh)
	}()
	return workerCh
}
