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
package errors

import (
	"fmt"
)

var (
	NoConfigErr        = &VoucherERR{Message: "no specify config"}
	NoBasedirErr       = &VoucherERR{Message: "the configuration file does not specify the basedir path."}
	NoPrivateNodeErr   = &VoucherERR{Message: "no private node"}
	NoDataErr          = &VoucherERR{Message: "no data"}
	GenServerCertErr   = &VoucherERR{Message: "generate server cert failed."}
	GenClientCertErr   = &VoucherERR{Message: "generate client cert failed."}
	GenServerSecretErr = &VoucherERR{Message: "generate server secret failed."}
	GenRandomPassErr   = &VoucherERR{Message: "generate random password failed."}
	SaveSecretErr      = &VoucherERR{Message: "save the secret to db failed."}
	SavePassErr        = &VoucherERR{Message: "save the password to db failed."}
)

type VoucherERR struct {
	Message string
	Err     error
}

func (e *VoucherERR) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s -> %v", e.Message, e.Err)
	} else {
		return fmt.Sprintf("%s", e.Message)
	}
}
