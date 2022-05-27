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
	return uuid.NewString
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
func getKeys(fieldGenMap map[string]func() string) []string {
	keys := []string{}
	for k, _ := range fieldGenMap {
		keys = append(keys, k)
	}
	return keys
}
func MakeCsv(numRows int, fieldGenMap map[string]func() string) (csv *strings.Reader, op []map[string]string) {
	// Get headers from map
	headers := getKeys(fieldGenMap)
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
			op[col] = fieldGenMap[col]()

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

func MakeSchoolsCsv(numSchools int) (*strings.Reader, []map[string]string) {
	columnNames := []string{"uuid", "school_name", "organization_uuid", "program_ids"}
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
	schools := []map[string]string{}
	for i := 0; i < numSchools; i++ {
		school := map[string]string{}
		school["uuid"] = uuid.NewString()
		school["school_name"] = "school" + strconv.Itoa(i)
		school["organization_uuid"] = uuid.NewString()
		programIds := []string{}
		for i := 0; i < 10; i++ {
			programIds = append(programIds, uuid.NewString())
		}
		school["program_ids"] = strings.Join(programIds, ";")
		schools = append(schools, school)
	}

	// Create file for test
	lines := []string{strings.Join(columnNames, ",")}
	for _, school := range schools {
		cols := make([]string, len(columnNames))
		for _, col := range columnNames {
			cols[colIndexMap[col]] = school[col]
		}

		line := strings.Join(cols, ",")
		lines = append(lines, line)
	}

	file := strings.NewReader(strings.Join(lines, "\n"))
	return file, schools
}
