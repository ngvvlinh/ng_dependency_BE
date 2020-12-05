package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"cloud.google.com/go/storage"
	"go.uber.org/atomic"
	"google.golang.org/api/iterator"

	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/storage/filestorage"
	"o.o/backend/pkg/common/storage/gcloudstorage"
	"o.o/common/l"
)

var ll = l.New()
var flDownloadDir = flag.String("dd", "", "directory to download")

type Config struct {
	File   *filestorage.Config   `yaml:"file"`
	Gcloud *gcloudstorage.Config `yaml:"gcloud"`
}

func main() {
	cc.InitFlags()
	cc.ParseFlags()

	var cfg Config
	err := cc.LoadWithDefault(&cfg, Config{})
	ll.Must(err, "load config")

	ctx := context.Background()
	gdriver, err := gcloudstorage.Connect(ctx, *cfg.Gcloud, []string{gcloudstorage.ReadOnly})
	ll.Must(err, "connect gcloud")

	bucket := gdriver.Bucket()
	objIter := bucket.Objects(ctx, &storage.Query{})

	var objs []*storage.ObjectAttrs
	cOK, cErr := 0, 0
loop:
	for {
		attrs, err := objIter.Next()
		switch err {
		case nil:

		case iterator.Done:
			fmt.Printf("scanned %v items (OK: %v, Err: %v)\n", cOK+cErr, cOK, cErr)
			break loop

		default:
			cErr++
			fmt.Printf("scan object: %v\n", err)
			continue
		}

		cOK++
		objs = append(objs, attrs)
		fmt.Printf("[%5v] %v\n      %v\n", cOK+cErr, attrs.Name, attrs.MediaLink)
	}

	if *flDownloadDir == "" {
		return
	}

	const N = 16
	var c atomic.Int64
	var wg sync.WaitGroup
	ch := make(chan *storage.ObjectAttrs, 32)
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			for item := range ch {
				download(c.Inc(), item)
			}
			wg.Done()
		}()
	}

	if *flDownloadDir != "" {
		fmt.Printf("\n\ndownloading\n")
		for _, obj := range objs {
			ch <- obj
		}
		close(ch)
		wg.Wait()
		fmt.Printf("\n\nDownloaded %v items", c.Load())
	}
}

func download(i int64, obj *storage.ObjectAttrs) {
	dirname := filepath.Join(*flDownloadDir, filepath.Dir(obj.Name))
	err := os.MkdirAll(dirname, 0755)
	ll.Must(err, "mkdir")

	saveTo := filepath.Join(*flDownloadDir, obj.Name)
	cmd := exec.Command("wget", "-O", saveTo, obj.MediaLink)
	output, err := cmd.CombinedOutput()
	if err == nil && cmd.ProcessState.ExitCode() == 0 {
		fmt.Printf("[%5v] %v OK\n", i, obj.Name)
		return
	}
	fmt.Printf("[%5v] %v ERROR\n%s\n", i, obj.Name, output)
}
