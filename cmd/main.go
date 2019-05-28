package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"internal"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	configPath = "./config/"
	defaultConfig = "development.yaml"
)

var (
	file *os.File
	nodes []Node
	cfg *config.Config
)

type Node struct {
	Time	int64
}




func init()  {
	//в зависимости от контура выбираем файл конфига
	cfgFileName := flag.String("c", "development.json", "config file name")
	flag.Parse()

	//cfg = LoadConfig(configPath, cfgFileName)
	cfg = New(configPath, cfgFileName)
	cfg.Load()
	fmt.Println("cfg: ...", cfg)
	//проверяем при старте наличие истории запросов
	if _, err := os.Stat(cfg.StorageName); !os.IsNotExist(err) {
		var data []Node
		storage, err := ioutil.ReadFile(cfg.StorageName)
		if err != nil {
			log.Fatal(err)
		}
		if len(storage) != 0 {
			dec := gob.NewDecoder(bytes.NewReader(storage))
			if err := dec.Decode(&data); err != nil {
				log.Fatal(err)
			}
			nodes = data
			log.Printf("история запросов на момент запуска %v", nodes)
		}
	}
}

func main()  {
	http.HandleFunc(cfg.URLPattern, HTTPCounter)
	Notify()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil))
}

func HTTPCounter(w http.ResponseWriter, _ *http.Request)  {
	node := Node{
		Time: time.Now().Unix(),
	}
	cutOff := node.Time - 10
	nodes = append(nodes, node)
	from := Counter(nodes, cutOff)
	nodes = nodes[from:]
	res := len(nodes)

	if _, err := fmt.Fprintf(w, "%s", strconv.Itoa(res)); err != nil {
		log.Println(err)
	}
}

func Counter(nodes []Node, cutOff int64) int {
	var from int

	for i, n := range nodes {
		if from == 0 {
			if n.Time >= cutOff {
				from = i
				break
			}
		}
	}

	return from
}

func Notify()  {
	var err error
	pid := os.Getpid()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func(pid int) {
		select {
		case _ = <-sigc:
			//save file
			var data bytes.Buffer
			enc := gob.NewEncoder(&data)
			if err := enc.Encode(nodes); err != nil {
				fmt.Println("err: ...", err)
			}
			if _, err = os.Stat(cfg.StorageName); os.IsNotExist(err) {
				file, err = os.Create(cfg.StorageName)
			}
			if err = ioutil.WriteFile(cfg.StorageName, data.Bytes(), 777); err != nil {
				log.Fatal(err)
			}

			defer func() {
				if err = file.Close(); err != nil {
					log.Println(err)
				}
			}()

			if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
				log.Fatal(err)
			}
		}
	}(pid)
}
