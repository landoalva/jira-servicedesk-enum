// Copyright 2025 İrem Kuyucu
// Copyright 2025 Laurynas Četyrkinas
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
)

func signup(baseURL, email string) error {
	client := newClient(baseURL, "")

	body := map[string]string{
		"email":          email,
		"secondaryEmail": "",
	}

	resp, err := client.post("/rest/servicedesk/1/customer/pages/user/signup", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 204:
		return nil
	case 403:
		return fmt.Errorf("signup forbidden (403)")
	default:
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
