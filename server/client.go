package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"

	"git.garena.com/kwads-relevance/rcommon"
)

const (
	exampleNum = 100
)

func main() {
	url := "http://localhost:8082/test"
	// prepare requests
	datas, err := loadReqData()
	if err != nil {
		log.Fatal(err)
		return
	}

	// post requests
	for _, jVal := range datas {
		func() {
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(jVal))
			if err != nil {
				log.Fatal(err)
				return
			}

			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("response Body:", string(body))
		}()
	}
}

// your own data preparation
func loadReqData() ([][]byte, error) {
	filePath := "data/train.json"
	kws, titles, err := loadFile(filePath)
	if err != nil {
		return nil, err
	}

	reqBytesAll := make([][]byte, 0, exampleNum)
	for i := 0; i < exampleNum; i++ {
		input := getRandInput(kws, titles)

		req := &rcommon.Request{
			Country: input.Country,
			Version: input.Version,
			Query: &rcommon.QueryFeatures{
				Emb:           input.QueryEmb,
				Keyword:       input.Keyword,
				CatDistribMap: input.CatDistribMap,
			},
			Items: []*rcommon.ItemFeatures{
				{
					Emb:    input.TitleEmb,
					Title:  input.Title,
					CatIDs: input.CatIDs,
				},
			},
		}
		reqBytes, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}
		reqBytesAll = append(reqBytesAll, reqBytes)
	}
	return reqBytesAll, nil
}

type dataLine struct {
	Title string
	Query string
}

func loadFile(filePath string) ([]string, []string, error) {
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return nil, nil, err
	}

	kws := make([]string, 0, 200000)
	titles := make([]string, 0, 200000)
	scanner := bufio.NewScanner(f)
	i := 0
	for scanner.Scan() {
		var line dataLine
		err := json.Unmarshal([]byte(scanner.Text()), &line)
		if err != nil {
			return nil, nil, err
		}

		kws = append(kws, line.Query)
		titles = append(titles, line.Title)

		i++
		if i >= exampleNum {
			break
		}

	}
	if err := scanner.Err(); err != nil {
		log.Print(err)
	}
	return kws, titles, nil
}

type randInput struct {
	Country       string
	Version       string
	Keyword       string
	Title         string
	QueryEmb      []float32
	TitleEmb      []float32
	CatDistribMap map[int32]float32
	CatIDs        []int32
}

func getRandInput(kws, titles []string) *randInput {
	country := "TW"
	version := "1"
	keyword := getRandItem(kws)
	title := getRandItem(titles)
	queryEmb := getRandEmb()
	titleEmb := getRandEmb()
	catDistribMap := getRandCatDistrib()
	catIDs := getRandCat(catDistribMap)

	return &randInput{
		Country:       country,
		Version:       version,
		Keyword:       keyword,
		Title:         title,
		QueryEmb:      queryEmb,
		TitleEmb:      titleEmb,
		CatDistribMap: catDistribMap,
		CatIDs:        catIDs,
	}
}

func getRandItem(items []string) string {
	i := rand.Intn(len(items))
	return items[i]
}

func getRandEmb() []float32 {
	size := 768
	emb := make([]float32, 0, size)
	for i := 0; i < size; i++ {
		emb = append(emb, float32((rand.Float64()-0.5)*0.07))
	}
	return emb
}

func getRandCatDistrib() map[int32]float32 {
	catDistrib := map[int32]float32{}
	sum := float32(0)
	for i := 0; i < 2; i++ {
		k := rand.Int31n(10000)
		v := float32(rand.Float64() * 0.5)
		catDistrib[k] = v
		sum += v
	}
	k := rand.Int31n(10000)
	catDistrib[k] = 1 - sum
	return catDistrib
}

func getRandCat(catDistrib map[int32]float32) []int32 {
	catIDs := make([]int32, 0, 3)
	for key := range catDistrib {
		catIDs = append(catIDs, key)
	}
	return catIDs
}
