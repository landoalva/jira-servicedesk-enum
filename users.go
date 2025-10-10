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
	"strings"
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

func enumerateUsers(baseURL, cookie string, maxUsers int, deskID string, customQuery string) error {
	client := newClient(baseURL, cookie)

	var desks []ServiceDesk
	targetingSingleDesk := deskID != ""

	if targetingSingleDesk {
		desks = []ServiceDesk{{ID: deskID}}
	} else {
		resp, err := client.get("/rest/servicedeskapi/servicedesk")
		if err != nil {
			return fmt.Errorf("get service desks: %w", err)
		}

		var desksResp ServiceDeskResponse
		if err := unmarshalJSON(resp, &desksResp); err != nil {
			return fmt.Errorf("parse service desks: %w", err)
		}

		desks = desksResp.Values
		fmt.Printf("\nFound %d service desk(s)\n\n", len(desks))
	}

	userMap := make(map[string]User)

	for _, desk := range desks {
		if desk.ProjectName != "" {
			fmt.Printf("Service Desk: %s (%s) [ID: %s]\n", desk.ProjectName, desk.ProjectKey, desk.ID)
		} else {
			fmt.Printf("Service Desk: [ID: %s]\n", desk.ID)
		}

		totalFetched := 0
		capped := false
		seenAccountIDs := make(map[string]bool)

		var queries []string
		if customQuery != "" {
			queries = []string{customQuery}
		} else {
			queries = []string{""}
			alphabet := "abcdefghijklmnopqrstuvwxyz0123456789"
			for _, char := range alphabet {
				queries = append(queries, string(char))
			}
		}

		for len(queries) > 0 {
			if maxUsers > 0 && totalFetched >= maxUsers {
				capped = true
				break
			}

			query := queries[0]
			queries = queries[1:]

			var url string
			if query == "" {
				url = fmt.Sprintf("/rest/servicedesk/1/customer/portal/%s/user-search/proforma", desk.ID)
			} else {
				url = fmt.Sprintf("/rest/servicedesk/1/customer/portal/%s/user-search/proforma?query=%s", desk.ID, query)
			}

			resp, err := client.get(url)
			if err != nil {
				fmt.Printf("  Error fetching users (query=%s): %v\n", query, err)
				continue
			}

			var users []User
			if err := unmarshalJSON(resp, &users); err != nil {
				fmt.Printf("  Error parsing users (query=%s): %v\n", query, err)
				continue
			}

			newUsersThisBatch := 0
			for _, user := range users {
				if seenAccountIDs[user.AccountID] {
					continue
				}
				seenAccountIDs[user.AccountID] = true
				newUsersThisBatch++
				totalFetched++

				if maxUsers > 0 && totalFetched > maxUsers {
					capped = true
					break
				}

				if targetingSingleDesk {
					fmt.Printf("AccountID: %s\n", user.AccountID)
					fmt.Printf("  Name: %s\n", user.DisplayName)
					if user.EmailAddress != "" {
						fmt.Printf("  Email: %s\n", user.EmailAddress)
					}
					if user.Avatar != "" && !strings.Contains(user.Avatar, "default-avatar.png") {
						fmt.Printf("  Avatar: %s\n", user.Avatar)
					}
					fmt.Println()
				} else {
					if _, exists := userMap[user.AccountID]; !exists {
						userMap[user.AccountID] = user
					}
				}
			}

			if len(users) == 50 && newUsersThisBatch > 0 && customQuery == "" {
				alphabet := "abcdefghijklmnopqrstuvwxyz0123456789"
				for _, char := range alphabet {
					queries = append(queries, query+string(char))
				}
			}

			if capped {
				break
			}
		}

		if capped {
			fmt.Printf("  Found %d user(s) [CAPPED at max=%d]\n", totalFetched, maxUsers)
		} else {
			fmt.Printf("  Found %d user(s)\n", totalFetched)
		}
	}

	if !targetingSingleDesk {
		fmt.Printf("\n\nUnique Users (%d):\n", len(userMap))

		for _, user := range userMap {
			fmt.Printf("AccountID: %s\n", user.AccountID)
			fmt.Printf("  Name: %s\n", user.DisplayName)
			if user.EmailAddress != "" {
				fmt.Printf("  Email: %s\n", user.EmailAddress)
			}
			if user.Avatar != "" && !strings.Contains(user.Avatar, "/default-avatar.png") {
				fmt.Printf("  Avatar: %s\n", user.Avatar)
			}
			fmt.Println()
		}
	}

	return nil
}
