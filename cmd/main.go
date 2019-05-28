package main

import (
	"bytes"
	"encoding/gob"
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

	"github.com/iluffka/test-task/internal/app/counter"
	"github.com/iluffka/test-task/internal/config"
)

var (
	file  *os.File
	nodes []counter.Node
	cfg   *config.Config
)

func init() {
	//в зависимости от контура выбираем файл конфига
	envType := flag.String("env", config.EnvDev, "config file name")
	cfgFileName := config.SetEnv(*envType)
	flag.Parse()

	cfg = config.New(config.ConfigPath, cfgFileName)
	cfg.Load()

	//проверяем при старте наличие истории запросов
	if _, err := os.Stat(cfg.StorageName); !os.IsNotExist(err) {
		var data []counter.Node
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

func main() {
	http.HandleFunc(cfg.URLPattern, HTTPServe)
	shutdown()

	log.Printf("запуск на порту %s", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil))
}

func HTTPServe(w http.ResponseWriter, _ *http.Request) {
	node := counter.Node{
		Time: time.Now().Unix(),
	}
	cutOff := node.Time - cfg.Period
	nodes = append(nodes, node)
	from := counter.Counter(nodes, cutOff)
	nodes = nodes[from:]
	res := len(nodes)

	if _, err := fmt.Fprintf(w, "%s", strconv.Itoa(res)); err != nil {
		log.Println(err)
	}
}

func shutdown() {
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
		case sig := <-sigc:
			log.Printf("последний запуск приложения в %v", cfg.Start)
			log.Printf("внешний вызов %s", sig.String())
			log.Printf("история последних запросов %v", nodes)

			var data bytes.Buffer
			enc := gob.NewEncoder(&data)
			if err := enc.Encode(nodes); err != nil {
				log.Println(err)
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
