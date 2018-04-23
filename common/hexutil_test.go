// Copyright 2017. box.la authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"github.com/ethereum/go-ethereum/common"
	"strings"
	"testing"
)

type testCase struct {
	input  string
	expect interface{}
	result bool
}

var sTestCases = []testCase{
	{
		input:  "AllowFlow(bytes32)",
		expect: "0x19d544d2bbf9f77ed0ab86140926d49d992616ab5fc07309b84eb6f1d3576d7e",
		result: true,
	},
}

func TestSignEvent(t *testing.T) {
	for _, tCase := range sTestCases {
		hash := SignEvent(tCase.input)
		assertEquals(t, strings.ToLower(hash.String()), tCase.expect.(string), tCase.result)
	}
}

func assertEquals(t *testing.T, a string, b string, result bool) {
	if (a == b) != result {
		t.Errorf("The result not equals to be expected value.\na: %s\nb: %s", a, b)
		t.FailNow()
	}
}

var addrTestCases = []testCase{
	{
		input:  "0x599d7abdb0a289f85aaca706b55d1b96cc07f349",
		expect: []byte{89, 157, 122, 189, 176, 162, 137, 248, 90, 172, 167, 6, 181, 93, 27, 150, 204, 7, 243, 73},
		result: true,
	},

	{
		input:  "0x599d7abdb0a289f85aaca706b55d1b96cc07f349",
		expect: []byte{0},
		result: false,
	},
}

func TestAddressEquals(t *testing.T) {
	for _, tCase := range addrTestCases {
		a := common.HexToAddress(tCase.input)
		b := common.BytesToAddress(tCase.expect.([]byte))
		//t.Logf("address: %v", b.Bytes())
		result := AddressEquals(a, b)
		if result != tCase.result {
			t.Errorf("unexpected result. address: %s", tCase.input)
			t.FailNow()
		}
	}
}
