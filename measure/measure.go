package measure

import (
	"bufio"
	"fmt"
	"github.com/degemer/document-search-engine/index"
	"github.com/degemer/document-search-engine/search"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type Measurer interface {
	Measure(search.Searcher)
}

type CacmMeasurer struct {
	queries map[int]string
	results map[int][]index.DocScore
}

func New(name string, options map[string]string) Measurer {
	temp := new(CacmMeasurer)
	temp.queries = loadQueries(options["cacm"])
	temp.results = loadResults(options["cacm"])
	return temp
}

func (cm *CacmMeasurer) Measure(searcher search.Searcher) {
	meanAveragePrecision := 0.0
	mins := []float64{10, 10, 10, 10, 10}
	maxs := make([]float64, 5)
	sums := make([]float64, 5)
	for query_id, request := range cm.queries {
		result := searcher.Search(request)
		precision, rappel := precisionRappel(result, cm.results[query_id])
		minMaxSum(mins, maxs, sums, precision, 0)
		minMaxSum(mins, maxs, sums, rappel, 1)
		minMaxSum(mins, maxs, sums, eMeasure(precision, rappel, 0.5), 2)
		minMaxSum(mins, maxs, sums, fMeasure(precision, rappel, 1), 3)
		minMaxSum(mins, maxs, sums, rPrecision(result, cm.results[query_id]), 4)
		meanAveragePrecision += averagePrecision(result, cm.results[query_id])
	}
	nbTests := float64(len(cm.queries))
	fmt.Println("Precision - min:", mins[0], "max:", maxs[0], "avg:", sums[0]/nbTests)
	fmt.Println("Rappel - min:", mins[1], "max:", maxs[1], "avg:", sums[1]/nbTests)
	fmt.Println("E-Measure - min:", mins[2], "max:", maxs[2], "avg:", sums[2]/nbTests)
	fmt.Println("F-Measure - min:", mins[3], "max:", maxs[3], "avg:", sums[3]/nbTests)
	fmt.Println("R-Measure - min:", mins[4], "max:", maxs[4], "avg:", sums[4]/nbTests)
	fmt.Println("MAP -", meanAveragePrecision/nbTests)
}

func minMaxSum(mins, maxs, sums []float64, measure float64, ind int) {
	if measure < mins[ind] {
		mins[ind] = measure
	} else if measure > maxs[ind] {
		maxs[ind] = measure
	}
	sums[ind] += measure
}

func precisionRappel(result, expectedResult []index.DocScore) (precision, rappel float64) {
	pertinentDocsFound := intersect(result, expectedResult)
	precision = float64(len(pertinentDocsFound)) / float64(len(result))
	rappel = float64(len(pertinentDocsFound)) / float64(len(expectedResult))
	if len(result) == 0 && len(expectedResult) == 0 {
		precision = 1
		rappel = 1
	} else if len(result) == 0 {
		precision = 0
	} else if len(expectedResult) == 0 {
		rappel = 1
	}
	return
}

func intersect(result, expectedResult []index.DocScore) (intersection []int) {
	resultMap := make(map[int]int)
	for _, docScore := range result {
		resultMap[docScore.Id] = 1
	}
	for _, docScore := range expectedResult {
		if resultMap[docScore.Id] == 1 {
			intersection = append(intersection, docScore.Id)
		}
	}
	return
}

func eMeasure(precision, rappel, alpha float64) float64 {
	return 1 - 1/(alpha/precision+(1-alpha)/rappel)
}

func fMeasure(precision, rappel, beta float64) float64 {
	return 1 - eMeasure(precision, rappel, 1/(1+beta*beta))
}

func rPrecision(result, expectedResult []index.DocScore) (rPrecision float64) {
	k := len(expectedResult)
	if k == 0 {
		k = 1
	}
	return precisionK(result, expectedResult, k)
}

func precisionK(result, expectedResult []index.DocScore, k int) (rPrecision float64) {
	rPrecision, _ = precisionRappel(result[:k], expectedResult)
	return
}

func averagePrecision(result, expectedResult []index.DocScore) (aveP float64) {
	if len(expectedResult) == 0 {
		if len(result) == 0 {
			return 1
		}
		return 0
	}
	for i, val := range result {
		if inExpected(val, expectedResult) {
			aveP += precisionK(result, expectedResult, i+1)
		}
	}
	return aveP / float64(len(expectedResult))
}

func inExpected(result index.DocScore, expectedResult []index.DocScore) bool {
	id := sort.Search(len(expectedResult), func(i int) bool { return result.Id <= expectedResult[i].Id })
	return id != len(expectedResult) && expectedResult[id].Id == result.Id
}

func loadQueries(cacmPath string) map[int]string {
	queries := make(map[int]string)
	reader := index.NewReader(map[string]string{"cacm_file": filepath.Join(cacmPath, "query.text")})
	for query := range reader.Read() {
		queries[query.Id] = query.Content
	}
	return queries
}

func loadResults(cacmPath string) map[int][]index.DocScore {
	results := make(map[int][]index.DocScore)
	resultsFilePath := filepath.Join(cacmPath, "qrels.text")
	file, err := os.Open(resultsFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		query_id, err := strconv.Atoi(text[0:2])
		if err != nil {
			fmt.Println("Unable to convert query id ", text[0:2], "to int: ", err)
			os.Exit(1)
		}
		doc_id, err := strconv.Atoi(text[3:7])
		if err != nil {
			fmt.Println("Unable to convert doc id ", text[3:7], "to int: ", err)
			os.Exit(1)
		}
		results[query_id] = append(results[query_id], index.DocScore{Id: doc_id})
	}
	return results
}
