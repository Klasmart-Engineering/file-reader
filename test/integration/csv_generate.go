package util

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func RandomStringGenerator(allowedChars []rune, min int, max int) func() string {
	// Return generator which makes random strings using supplied characters with length between boundaries min and max
	return func() string {
		length := rand.Intn(max-min) + min
		randomString := make([]rune, length)
		for i := range randomString {
			randomString[i] = allowedChars[rand.Intn(len(allowedChars))]
		}
		return string(randomString)
	}
}

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

func HumanNameFieldGenerator(min int, max int) func() string {
	// Supply min and max bounds for how long names can be, makes random "name" of that length
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ漢字नागरीж ")
	rand.Seed(time.Now().UnixNano())
	return RandomStringGenerator(letters, min, max)
}

func DateGenerator(minYear int, maxYear int, layout string) func() string {
	// Makes a random date between the specified years.
	// The layout string is an example date in the desired format (e.g 2006-01-02 or January 2, 2006)
	rand.Seed(time.Now().UnixNano())
	min := time.Date(minYear, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(maxYear, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min
	return func() string {
		sec := rand.Int63n(delta) + min
		return time.Unix(sec, 0).Format(layout)
	}
}

func GenderGenerator() func() string {
	rand.Seed(time.Now().UnixNano())
	genders := []string{"male", "female"} // Or M, f, other? Needs review.
	return func() string {
		return genders[rand.Intn(len(genders))]
	}
}

func EmptyFieldGenerator() func() string {
	return func() string {
		return ""
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

func MakeEmptyFile() *strings.Reader {
	return strings.NewReader("")
}
