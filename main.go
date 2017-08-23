package main

import (
	"fmt"
	"github.com/urfave/cli"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

const RANDOM_BYTES = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Fetcher struct {
	Url       string
	Limit     int
	MinWidth  int
	MinHeight int
	MaxWidth  int
	MaxHeight int
	Random    bool
	Dest      string
}

func (f *Fetcher) Download() {
	var wg sync.WaitGroup
	wg.Add(f.Limit)
	for i := 1; i <= f.Limit; i++ {
		go func() {
			defer wg.Done()
			saved := f.SaveFile()
			if saved {
				fmt.Println("* Download successfull")
			} else {
				fmt.Println("- Download failed")
			}
		}()
	}
	wg.Wait()
}

func (f *Fetcher) SaveFile() bool {
	var url string
	if f.Random {
		url = f.getUrlRandom()
	} else {
		url = f.getUrl()
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Could not download image")
	}
	defer resp.Body.Close()

	os.Mkdir(f.Dest, os.FileMode(0777))

	file, err := os.Create(path.Join(f.Dest, generateFileName(25)+".jpg"))
	if err != nil {
		fmt.Println("Unable to write on disk.")
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Unable to save image on disk")
	}

	file.Close()

	return true
}

func (f *Fetcher) getUrl() string {
	rand.Seed(time.Now().Unix())
	width := strconv.Itoa(rand.Intn(f.MaxWidth-f.MinWidth) + f.MinWidth)
	height := strconv.Itoa(rand.Intn(f.MaxHeight-f.MinHeight) + f.MinHeight)
	return f.Url + width + "x" + height
}

func (f *Fetcher) getUrlRandom() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	width := strconv.Itoa(r.Intn(1000))
	height := strconv.Itoa(r.Intn(1000))
	return f.Url + width + "x" + height
}

func generateFileName(size int) string {
	buf := make([]byte, size)
	for i := range buf {
		buf[i] = RANDOM_BYTES[rand.Int63()%int64(len(RANDOM_BYTES))]
	}

	return string(buf)
}

func main() {
	app := cli.NewApp()
	app.Name = "IMGSEED"
	app.Usage = "download random images"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url",
			Value: "https://source.unsplash.com/random/",
			Usage: "source url to fetch images from",
		},
		cli.StringFlag{
			Name:  "dest",
			Value: "img",
			Usage: "dir to save downloaded images",
		},
		cli.IntFlag{
			Name:  "limit",
			Value: 10,
			Usage: "number of images to download",
		},
		cli.IntFlag{
			Name:  "min-width",
			Value: 100,
			Usage: "min width of image",
		},
		cli.IntFlag{
			Name:  "max-width",
			Value: 1000,
			Usage: "max width of image",
		},
		cli.IntFlag{
			Name:  "min-height",
			Value: 100,
			Usage: "min height of image",
		},
		cli.IntFlag{
			Name:  "max-height",
			Value: 1000,
			Usage: "max height of image",
		},
		cli.BoolFlag{
			Name:  "random",
			Usage: "download random size images",
		},
	}

	f := &Fetcher{}

	app.Action = func(c *cli.Context) error {
		f.Url = c.String("url")
		f.Limit = c.Int("limit")
		f.Random = c.Bool("random")
		f.Dest = c.String("dest")
		f.MinWidth = c.Int("min-width")
		f.MaxWidth = c.Int("max-width")
		f.MinHeight = c.Int("min-height")
		f.MaxHeight = c.Int("max-height")

		return nil
	}

	app.Run(os.Args)
	f.Download()
}
