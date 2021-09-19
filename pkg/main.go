package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

type Worker interface {
	Process()
}

type Job struct {
	root string
}

type Result struct {
	file string
	isDir bool
	error error
}

type worker struct {
	Jobs chan Job
	Results chan Result
	Done chan struct{}
}

func main(){
 w := worker{
	 Jobs:    make(chan Job),
	 Results: make(chan Result),
	 Done:    nil,
 }
 go w.Generate()
 go w.Consume(5)
 time.Sleep(30 * time.Second)
}

func(w *worker) Generate(){
	fmt.Println("\ninit Generate")
	w.Jobs <- Job{"/Users/bhakiyarajkalimuthu/Documents/Notes/"}

	go func(){
		for r :=  range w.Results{
			if r.error!=nil {
				fmt.Printf("error %s\n", r.error)
				continue
			}
			fmt.Printf("file name : %s, isDir : %t\n ", r.file,r.isDir)

		}
		defer close(w.Results)
	}()
	//defer close(w.Jobs)
	fmt.Println("\ndone Generate")
}

func(w *worker) Consume( workload int){
	fmt.Println("\ninit Consume")
	var wg sync.WaitGroup
	for i:=0;i< workload;i++{
		wg.Add(1)
		go func(wgp *sync.WaitGroup){
			w.Process(context.Background())
			wgp.Done()
			}(&wg)
	}
	wg.Wait()
	fmt.Println("\ndone Consume")
}

func(w *worker) Process(ctx context.Context){
	fmt.Println("\ninit Process")
	for {
		select{
			case job,ok := <- w.Jobs:
				if !ok {
					fmt.Println("\nnot ok")
					return
				}
				w.ListFiles(job.root)
				fmt.Println("\nlist files called")
				return
			case <-ctx.Done():
				fmt.Println("\nnothing available")
				return
		}
	}
	fmt.Println("done Generate")
}

func(w *worker) ListFiles(root string){
	// TODO: Validate root is dir ?
	fmt.Printf("ListFiles root %s\n",root)
	f,err := os.Open(root)
	if err!=nil {
		w.Results <- Result{
			error : fmt.Errorf("failed to open %s reason %s\n",root,err),
		}
	}

	dirs,err := f.ReadDir(-1)
	f.Close()
	if err!=nil {
		w.Results <- Result{
			error : fmt.Errorf("failed to read directory %s reason %s\n",root,err),
		}
	}
	for _,i := range dirs {
		w.Results <- Result{
			file:  i.Name(),
			isDir: i.IsDir(),
		}
	}
	//for _,i := range dirs {
	//	if i.IsDir(){
	//		go func(){
	//			//fmt.Printf("isDir %s",i)
	//			w.Jobs <- Job{root+i.Name()}
	//		}()
	//		continue
	//	}
	//	go func() {
	//		w.Results <- Result{
	//			file: i.Name(),
	//		}
	//	}()
	//}
}