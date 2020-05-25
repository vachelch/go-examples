package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"git.garena.com/kwads-relevance/rcommon"
	"git.garena.com/kwads-relevance/rel"
)

// init model
var scorer rel.Scorer

func test(w http.ResponseWriter, r *http.Request) {
	// get request
	req := &rcommon.Request{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Print(err)
		return
	}
	// caculate
	resp, err := scorer.ItemsScore(req)
	if err != nil {
		log.Print(err)
		return
	}
	// response
	scores := make([]float64, 0, len(resp.Items))
	for i := 0; i < len(resp.Items); i++ {
		scores = append(scores, resp.Items[i].Score)
	}
	fmt.Fprintf(w, "scores: %v", scores)
}

func main() {
	var err error
	scorer, err = rel.NewScorer("data/config.yaml")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/test", test)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
