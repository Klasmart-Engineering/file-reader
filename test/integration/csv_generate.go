package util

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func MakeOrgsCsv(numOrgs int) (csv *strings.Reader, orgs []map[string]string) {
	columnNames := []string{"uuid", "organization_name", "owner_user_id"}
	// Reorder columns to random order and map column names to their random index
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(columnNames), func(i, j int) {
		columnNames[i], columnNames[j] = columnNames[j], columnNames[i]
	})
	colIndexMap := map[string]int{}
	for i, col := range columnNames {
		colIndexMap[col] = i
	}

	// Create organizations (gets returned for use in assertions)
	organizations := []map[string]string{}
	for i := 0; i < numOrgs; i++ {
		org := map[string]string{}
		org["uuid"] = uuid.NewString()
		org["organization_name"] = "org" + strconv.Itoa(i)
		org["owner_user_id"] = uuid.NewString()
		organizations = append(organizations, org)
	}

	// Create file for test
	lines := []string{strings.Join(columnNames, ",")}
	for _, org := range organizations {
		cols := make([]string, len(columnNames))
		cols[colIndexMap["uuid"]] = org["uuid"]
		cols[colIndexMap["organization_name"]] = org["organization_name"]
		cols[colIndexMap["owner_user_id"]] = org["owner_user_id"]

		line := strings.Join(cols, ",")
		lines = append(lines, line)
	}

	file := strings.NewReader(strings.Join(lines, "\n"))
	return file, organizations
}
