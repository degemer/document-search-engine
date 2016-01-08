package index

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCacmDoc(t *testing.T) {
	unparsed_req := `.I 1
.W
 What articles exist which deal with TSS (Time Sharing System), an
operating system for IBM computers?
.N
 1. Richard Alexander, Comp Serv, Langmuir Lab (TSS)`
	expected_req := `
 What articles exist which deal with TSS (Time Sharing System), an
operating system for IBM computers?`
	parsed_req_id, parsed_req := cacmDoc(unparsed_req)
	assert.Equal(t, parsed_req_id, 1, "Id incorrect")
	assert.Equal(t, parsed_req, expected_req, "Content incorrect")

	unparsed_doc := `.I 1655
.T
Code Extension Procedures for Information
Interchange* (Proposed USA Standard)
.B
CACM December, 1968
.K
standard code, code, information interchange, characters,
shift out, shift in, escape, data link
escape, control functions, standard procedures,
code extension, code table, bit pattern
.C
1.0 2.0 2.43 3.20 3.24 3.50 3.51 3.52 3.53 3.54 3.55 3.56 3.57 3.70 3.71 3.72
3.73, 3.74, 3.75, 3.80, 3.81, 3.82, 3.83, 5.0, 5.1, 6.2, 6.21, 6.22
.N
CA681211 JB February 21, 1978  12:16 PM
.X
1655	5	1655
1655	5	1655
1655	5	1655`
	expected_doc := `
Code Extension Procedures for Information
Interchange* (Proposed USA Standard)
standard code, code, information interchange, characters,
shift out, shift in, escape, data link
escape, control functions, standard procedures,
code extension, code table, bit pattern`
	parsed_doc_id, parsed_doc := cacmDoc(unparsed_doc)
	assert.Equal(t, parsed_doc_id, 1655, "Id incorrect")
	assert.Equal(t, parsed_doc, expected_doc, "Content incorrect")
}
