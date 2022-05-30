package util

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func NameFieldGenerator(prefix string, n int) func() string {
	// Return a generator which adds a random number up to n to the supplied prefix
	rand.Seed(time.Now().UnixNano())
	return func() string {
		i := rand.Intn(n)
		return prefix + strconv.Itoa(i)
	}
}

func UuidFieldGenerator() func() string {
	return func() string {
		return uuid.NewString()
	}
}

func RepeatedFieldGenerator(gen func() string, min int, max int) func() string {
	// Supply a field generator function and min and max number of times to repeat it.
	// Returns a generator which generates a ; delimited string of those fields
	rand.Seed(time.Now().UnixNano())
	return func() string {
		numFields := rand.Intn(max-min) + min
		fields := []string{}
		for i := 0; i < numFields; i++ {
			fields = append(fields, gen())
		}
		return strings.Join(fields, ";")
	}
}

func MakeCsv(headers []string, numRows int, fieldGenMap map[string]func() string) (csv *strings.Reader, op []map[string]string) {

	// Reorder columns to random order and map column names to their random index
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(headers), func(i, j int) {
		headers[i], headers[j] = headers[j], headers[i]
	})
	colIndexMap := map[string]int{}
	for i, col := range headers {
		colIndexMap[col] = i
	}

	// Create operation (gets returned for use in assertions)
	ops := []map[string]string{}
	for i := 0; i < numRows; i++ {
		op := map[string]string{}
		for _, col := range headers {
			// check existence so we can generate default value for extra headers
			_, exists := fieldGenMap[col]
			if exists {
				op[col] = fieldGenMap[col]()
			} else {
				op[col] = uuid.NewString()
			}

		}

		ops = append(ops, op)
	}

	// Create file for test
	lines := []string{strings.Join(headers, ",")}
	for _, op := range ops {
		cols := make([]string, len(headers))
		for _, col := range headers {

			cols[colIndexMap[col]] = op[col]

		}
		line := strings.Join(cols, ",")
		lines = append(lines, line)
	}

	file := strings.NewReader(strings.Join(lines, "\n"))
	return file, ops

}

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
		for _, col := range columnNames {
			cols[colIndexMap[col]] = org[col]
		}

		line := strings.Join(cols, ",")
		lines = append(lines, line)
	}

	file := strings.NewReader(strings.Join(lines, "\n"))
	return file, organizations
}
