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
- `document-search-engine -i tf-idf-stem measure vectorial-overlap` will load/create a tf-idf with stemming index, create a `vectorial` with overlap function search, and then run the test queries on it, and output the results.
- `document-search-engine --cacm ~/custom_directory/cacm index` will create a tf-idf index, using CACM from `~/custom_directory/cacm`

##
