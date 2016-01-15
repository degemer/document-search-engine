package measure

import (
	"bufio"
	"github.com/degemer/document-search-engine/index"
	"github.com/degemer/document-search-engine/search"
	"log"
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
	temp.queries = loadQueries(options["cacm_path"])
	temp.results = loadResults(options["cacm_path"])
	return temp
}

func (cm *CacmMeasurer) Measure(searcher search.Searcher) {
	meanAveragePrecision := 0.0
	for query_id, request := range cm.queries {
		result := searcher.Search(request)
		// sortedIdResult := make([]index.DocScore, len(result))
		// copy(sortedIdResult, result)
		// sort.Sort(index.ById(sortedIdResult))
		// precision, rappel := precisionRappel(sortedIdResult, cm.results[query_id])
		// log.Println(query_id,
		// 			precision,
		// 			rappel,
		// 			eMeasure(precision, rappel, 0.5),
		// 			fMeasure(precision, rappel, 1),
		// 			rPrecision(sortedIdResult, cm.results[query_id]))
		meanAveragePrecision += averagePrecision(result, cm.results[query_id])
	}
	meanAveragePrecision /= float64(len(cm.queries))
	log.Println("MAP:", meanAveragePrecision)
}

func precisionRappel(result, expectedResult []index.DocScore) (precision, rappel float64) {
	pertinentDocsFound := search.Intersect(result, expectedResult)
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

func eMeasure(precision, rappel, alpha float64) float64 {
	return 1 - 1/(alpha/precision+(1-alpha)/rappel)
}

func fMeasure(precision, rappel, beta float64) float64 {
	return 1 - eMeasure(precision, rappel, 1/(1+beta*beta))
}

func rPrecision(result, expectedResult []index.DocScore) (rPrecision float64) {
	return precisionK(result, expectedResult, len(expectedResult))
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
	reader := index.NewReader(map[string]string{"cacm_path": filepath.Join(cacmPath, "query.text")})
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
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		query_id, err := strconv.Atoi(text[0:2])
		if err != nil {
			log.Fatalln("Unable to convert query id ", text[0:2], "to int: ", err)
		}
		doc_id, err := strconv.Atoi(text[3:7])
		if err != nil {
			log.Fatalln("Unable to convert doc id ", text[3:7], "to int: ", err)
		}
		results[query_id] = append(results[query_id], index.DocScore{Id: doc_id})
	}
	return results
}
