package config

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

const middle = "========="

var Instance *Config

type Config struct {
	dataMap map[string]string
	secret  string
}

func (this *Config) InitConfig(path string) {
	this.dataMap = make(map[string]string)

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		s := strings.TrimSpace(string(b))
		//fmt.Println(s)
		if strings.Index(s, "#") == 0 {
			continue
		}

		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			this.secret = strings.TrimSpace(s[n1+1 : n2])
			continue
		}

		if len(this.secret) == 0 {
			continue
		}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		first := strings.TrimSpace(s[:index])
		if len(first) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}

		key := this.secret + middle + first
		this.dataMap[key] = strings.TrimSpace(second)
	}
}

func (this *Config) Get(node, key string) string {
	key = node + middle + key
	v, found := this.dataMap[key]
	if !found {
		return ""
	}
	return v
}

func (this *Config) GetInt(node, key string) int {
	key = node + middle + key
	v, found := this.dataMap[key]
	if !found {
		return 0
	}
	conv, _ := strconv.Atoi(v)
	return conv
}
