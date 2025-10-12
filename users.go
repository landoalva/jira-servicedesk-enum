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
	"encoding/csv"
	"fmt"
	"os"
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

const defaultAvatar = "/default-avatar.png"

func enumerateUsers(baseURL, cookie string, maxUsers int, deskID string, customQuery string, alphabet string, selfAccountID string, outputPath string) error {
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
				fmt.Fprintf(os.Stderr, "Error: fetching users (query=%s): %v\n", query, err)
				continue
			}

			var users []User
			if err := unmarshalJSON(resp, &users); err != nil {
				fmt.Fprintf(os.Stderr, "Error: parsing users (query=%s): %v\n", query, err)
				continue
			}

			newUsersThisBatch := 0
			for _, user := range users {
				if seenAccountIDs[user.AccountID] {
					continue
				}

				if user.AccountID == selfAccountID {
					continue
				}

				if maxUsers > 0 && totalFetched >= maxUsers {
					capped = true
					break
				}

				seenAccountIDs[user.AccountID] = true
				newUsersThisBatch++
				totalFetched++

				if _, exists := userMap[user.AccountID]; !exists {
					userMap[user.AccountID] = user
				}
			}

			if len(users) == 50 && newUsersThisBatch > 0 && customQuery == "" {
				if (maxUsers == 0 || maxUsers > 50) && totalFetched < maxUsers {
					for _, char := range alphabet {
						queries = append(queries, query+string(char))
					}
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

	if len(userMap) == 0 {
		fmt.Println("\nNo users found")
		return nil
	}

	if outputPath != "" {
		return writeUsersToCSV(userMap, outputPath)
	}

	printUsers(userMap)
	return nil
}

func writeUsersToCSV(userMap map[string]User, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"AccountID", "DisplayName", "Email", "Avatar"}); err != nil {
		return fmt.Errorf("write CSV header: %w", err)
	}

	for _, user := range userMap {
		avatar := user.Avatar
		if strings.Contains(avatar, defaultAvatar) {
			avatar = ""
		}
		if err := writer.Write([]string{user.AccountID, user.DisplayName, user.EmailAddress, avatar}); err != nil {
			return fmt.Errorf("write CSV row: %w", err)
		}
	}

	fmt.Printf("\nWrote %d users to %s\n", len(userMap), outputPath)
	return nil
}

func printUsers(userMap map[string]User) {
	fmt.Printf("\n\nUnique Users (%d):\n", len(userMap))

	for _, user := range userMap {
		fmt.Printf("AccountID: %s\n", user.AccountID)
		fmt.Printf("  Name: %s\n", user.DisplayName)
		if user.EmailAddress != "" {
			fmt.Printf("  Email: %s\n", user.EmailAddress)
		}
		if user.Avatar != "" && !strings.Contains(user.Avatar, defaultAvatar) {
			fmt.Printf("  Avatar: %s\n", user.Avatar)
		}
		fmt.Println()
	}
}
