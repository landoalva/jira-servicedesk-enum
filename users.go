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

type ServiceDeskResponse struct {
	Values []ServiceDesk `json:"values"`
}

type ServiceDesk struct {
	ID          string `json:"id"`
	ProjectID   string `json:"projectId"`
	ProjectName string `json:"projectName"`
	ProjectKey  string `json:"projectKey"`
}

type User struct {
	ID           string `json:"id"`
	AccountID    string `json:"accountId"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
	Avatar       string `json:"avatar"`
}

func enumerateUsers(baseURL, cookie string) error {
	client := newClient(baseURL, cookie)

	resp, err := client.get("/rest/servicedeskapi/servicedesk")
	if err != nil {
		return fmt.Errorf("get service desks: %w", err)
	}

	var desksResp ServiceDeskResponse
	if err := unmarshalJSON(resp, &desksResp); err != nil {
		return fmt.Errorf("parse service desks: %w", err)
	}

	fmt.Printf("\nFound %d service desk(s)\n\n", len(desksResp.Values))

	userMap := make(map[string]User)

	for _, desk := range desksResp.Values {
		fmt.Printf("Service Desk: %s (%s)\n", desk.ProjectName, desk.ProjectKey)

		resp, err := client.get(fmt.Sprintf("/rest/servicedesk/1/customer/portal/%s/user-search/proforma", desk.ID))
		if err != nil {
			fmt.Printf("  Error fetching users: %v\n", err)
			continue
		}

		var users []User
		if err := unmarshalJSON(resp, &users); err != nil {
			fmt.Printf("  Error parsing users: %v\n", err)
			continue
		}

		fmt.Printf("  Found %d user(s)\n", len(users))

		for _, user := range users {
			if _, exists := userMap[user.ID]; !exists {
				userMap[user.ID] = user
			}
		}
	}

	fmt.Printf("\n\nUnique Users (%d):\n", len(userMap))
	fmt.Println("----------------------------------------")

	for _, user := range userMap {
		fmt.Printf("%-40s %-30s %s\n", user.DisplayName, user.EmailAddress, user.AccountID)
	}

	return nil
}
