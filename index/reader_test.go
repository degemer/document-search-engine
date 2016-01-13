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

func BenchmarkCacmDoc(b *testing.B) {
	unparsed_doc := `.I 3204
.T
   An On-Line Program for Non-Numerical Algebra
.W
   The goal of this program is to make a step toward te design of an automated
mathematical assistant. Some requirements for such a program are: it must be
easy to access, and that the result must be obtained in a reasonably short
time. Accordingly the program is written for a time-shared computer. The Q-32
computer as System Development Corporation, Santa Monica, California, was
chosen because it also had a LISP 1.5 compiler. Programming and debugging was
done from a remote teletype console at Stanford University.
.B
CACM August, 1966
.A
Korsvold, K.
.N
CA660818 ES March 17, 1982 10:10 AM
.X
1396	5	3204
3204	5	3204
3204	5	3204
3204	5	3204
964	6	3204
1028	6	3204
1029	6	3204
1083	6	3204
1132	6	3204
1214	6	3204
1278	6	3204
1334	6	3204
1365	6	3204
1386	6	3204
1387	6	3204
1388	6	3204
1392	6	3204
1393	6	3204
1394	6	3204
1395	6	3204
1396	6	3204
1397	6	3204
1496	6	3204
284	6	3204
407	6	3204
3199	6	3204
3200	6	3204
3201	6	3204
3202	6	3204
3203	6	3204
3204	6	3204
561	6	3204
730	6	3204`
	for i := 0; i < b.N; i++ {
		cacmDoc(unparsed_doc)
	}
}
