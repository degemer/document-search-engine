# document-search-engine

Search engine on the [CACM collection](http://ir.dcs.gla.ac.uk/resources/test_collections/cacm/).

## Installation

```
go get github.com/degemer/document-search-engine
```

Or download a pre-compiled binary (Mac-only) on https://github.com/degemer/document-search-engine/releases.

If you want to build it, it has 2 dependencies on go packages :
- https://github.com/codegangsta/cli (for a beautiful CLI)
- https://github.com/reiver/go-porterstemmer (for the Porter stemmer implementation)

## Usage

Before using `document-search-engine`, you need to download and extract the [CACM collection](http://ir.dcs.gla.ac.uk/resources/test_collections/cacm/cacm.tar.gz), preferably in the same directory as `document-search-engine`, so that in the end you have :

```
├── document-search-engine
└── cacm
    ├── README
    ├── cacm.all
    ├── cite.info
    ├── common_words
    ├── qrels.text
    └── query.text

```

Otherwise you can use the `--cacm` option to specify the CACM directory.

You can run `document-search-engine --help` to see all the options, but it has mainly three uses :
- `document-search-engine -i INDEX_TYPE index` : create and save a `INDEX_TYPE` index of the CACM database
- `document-search-engine -i INDEX_TYPE search SEARCH_TYPE` : query the `INDEX_TYPE` index using a SEARCH_TYPE search
- `document-search-engine -i INDEX_TYPE measure SEARCH_TYPE` : return the main measures (minimum, maximum, average of precision, recall, E-Measure, ...) on `SEARCH_TYPE` search using `INDEX_TYPE` index

`INDEX_TYPE` can take the values : `tf-idf`, `tf-idf-norm`, `tf-norm`, `tf-idf-stem`, `tf-idf-norm-stem`, `tf-norm-stem` (when `-stem` is present, Porter stemming will be used)

`SEARCH_TYPE` can take the values : `vectorial`, `vectorial-dice`, `vectorial-jaccard`, `vectorial-overlap`, `boolean` and `probabilistic`.

If you specify an incorrect or empty `INDEX_TYPE` and `SEARCH_TYPE`, they will default to `tf-idf` and `vectorial`.

A few examples :
- `document-search-engine --save ~/custom_directory/ search` will load from `~/custom_directory/` a Tf-Idf index or create it (and save it in `~/custom_directory/`), create a `vectorial` search, and then wait for an input.
- `document-search-engine -i tf-idf-stem measure vectorial-overlap` will load/create a Tf-Idf with stemming index, create a `vectorial` with overlap function search, and then run the test queries on it, and output the results.
- `document-search-engine --cacm ~/custom_directory/cacm index` will create a Tf-Idf index, using CACM from `~/custom_directory/cacm`


## Implementation

```
.
├── index
├── measure
└── search
```

### index

`index` contains all code relative to indices, i.e. the different implementation, and the pipeline to build an index :
```
Reader -> Tokenizer -> Filter -> Stemmer -> Counter
```

`Counter` returns the word count per document (for Tf) and global word count (for Idf). Each index then deals with them in its own way.


### search

`search` has an index as argument, and then creates the needed search.

### measure

`measure` creates an index and a search, and then test all the `cacm/query.text` requests.

## Performance

### Time performance

Performance is tested using Golang benchmarks (in `*_test.go` files). They are also output when creating an index, querying an index, or running a `measure`.

To run benchmarks : `go test -bench . ./...`.
Test output has the following format : `TestName numberOfIterations timePerIteration`.

Results on a i5 Macbook Pro 2013 (with comments) :

```
################################
# Benchmarks of index creation #
################################
# Benchmark of Tf-Idf creation, 149ms
BenchmarkTfIdfCreate-4        	      10	 149468796 ns/op
# Benchmark of Tf-Idf with stemming creation, 191ms
BenchmarkTfIdfStemCreate-4    	      10	 191585950 ns/op
# Benchmark of Tf-Idf normalized, 167ms
BenchmarkTfIdfNormCreate-4    	      10	 167351149 ns/op
# Benchmark of Tf-Idf with stemming normalized, 203ms
BenchmarkTfIdfNormStemCreate-4	       5	 203404449 ns/op
# Benchmark of Tf normalized, 108ms
BenchmarkTfNormCreate-4       	      10	 108301465 ns/op
# Benchmark of Tf with stemming normalized, 140ms
BenchmarkTfNormStemCreate-4   	      10	 140150647 ns/op

############################
# Benchmarks of query time #
############################
# Benchmark of Boolean Search with Tf-Idf, <1ms
BenchmarkBooleanSearch-4          	    2000	    899328 ns/op
# Benchmark of Probabilistic Search with Tf-Idf, <0.5ms
BenchmarkProbabilisticSearch-4    	    3000	    451336 ns/op
# Benchmark of Probabilistic Search with Tf-Idf stemmed, <1ms
BenchmarkProbabilisticSearchStem-4	    2000	    885313 ns/op
# Benchmark of Vectorial Search with Tf-Idf, <0.5ms
BenchmarkVectorialSearch-4        	    3000	    451840 ns/op
# Benchmark of Vectorial Search with Tf-Idf stemmed, <1ms
BenchmarkVectorialSearchStem-4    	    2000	    853147 ns/op
# Benchmark of Vectorial Overlap  Search with Tf-Idf, <0.5ms
BenchmarkVectorialSearchSum-4     	    3000	    482698 ns/op
# Benchmark of Vectorial Overlap Search with Tf-Idf stemmed, ~1ms
BenchmarkVectorialSearchSumStem-4 	    2000	   1008302 ns/op
```

### Precision performance

Everything can be tested using the `measure` command.

A few remarks :
- best MAP without stemming is obtained by using `vectorial-overlap` with `tf-idf` index.
- best MAP with stemming is still obtained by using `vectorial-overlap` with `tf-idf-stem` index.
- worst index is `tf-norm` (with stemming or without)
- stemming improves MAP
- `vectorial` and `probabilistic` have nearly the same MAP

