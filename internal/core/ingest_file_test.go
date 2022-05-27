package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	zapLogger "github.com/KL-Engineering/file-reader/internal/log"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestReadRows(t *testing.T) {
	type TestCases struct {
		description  string
		fileBody     string
		fileRows     chan []string
		expectedRows int
	}
	testcases := []TestCases{
		{
			description: "When all rows have same number of columns, all rows should get passed through",
			fileBody: `uuid,organization_name,owner_user_id
0755f5c0-dd08-11ec-9d64-0242ac120002,org1,0755f5c0-dd08-11ec-9d64-0242ac120003
0755fa70-dd08-11ec-9d64-0242ac120005,org2,0755fa70-dd08-11ec-9d64-0242ac120006
0755fb9c-dd08-11ec-9d64-0242ac120008,org3,0755fb9c-dd08-11ec-9d64-0242ac120009`,
			fileRows:     make(chan []string, 4),
			expectedRows: 4,
		},
		{
			description: "When number of columns on row is less than number of headers, that row should skip, and others should be passed through",
			fileBody: `uuid,organization_name,owner_user_id
0755f5c0-dd08-11ec-9d64-0242ac120002,org1,0755f5c0-dd08-11ec-9d64-0242ac120003
0755fa70-dd08-11ec-9d64-0242ac120005,0755fa70-dd08-11ec-9d64-0242ac120006
0755fb9c-dd08-11ec-9d64-0242ac120008,org3,0755fb9c-dd08-11ec-9d64-0242ac120009`,
			fileRows:     make(chan []string, 4),
			expectedRows: 3,
		},
		{
			description: "When number of columns on row is greater than number of headers, that row should skip, and others should be passed through",
			fileBody: `uuid,organization_name,owner_user_id
0755f5c0-dd08-11ec-9d64-0242ac120002,org1,0755f5c0-dd08-11ec-9d64-0242ac120003
0755fa70-dd08-11ec-9d64-0242ac120005,org2,0755fa70-dd08-11ec-9d64-0242ac120006,extracolumn
0755fb9c-dd08-11ec-9d64-0242ac120008,org3,0755fb9c-dd08-11ec-9d64-0242ac120009`,
			fileRows:     make(chan []string, 4),
			expectedRows: 3,
		},
		{
			description: "When there are empty values, the rows should still get passed along successfully",
			fileBody: `uuid,organization_name,slogan,owner_user_id,founder
0755f5c0-dd08-11ec-9d64-0242ac120002,org1,we the best,0755f5c0-dd08-11ec-9d64-0242ac120003,Steve
0755fa70-dd08-11ec-9d64-0242ac120005,,org2,0755fa70-dd08-11ec-9d64-0242ac120006,Jim
0755fb9c-dd08-11ec-9d64-0242ac120008,,org3,0755fb9c-dd08-11ec-9d64-0242ac120009,
0755fb9c-dd08-11ec-9d64-0242ac120012,lorem ipsum,org4,0755fb9c-dd08-11ec-9d64-0242ac120044,`,
			fileRows:     make(chan []string, 5),
			expectedRows: 5,
		},
	}
	for i, scenario := range testcases {
		t.Run(scenario.description, func(t *testing.T) {
			f, _ := ioutil.TempFile("", "ingest-file-test-"+strconv.Itoa(i))
			defer os.Remove(f.Name())
			_, err := f.WriteString(scenario.fileBody)
			if err != nil {
				fmt.Println(err)
				t.FailNow()
			}
			f.Close()
			f, _ = os.Open(f.Name())

			ctx := context.Background()
			l, _ := zap.NewDevelopment()
			logger := zapLogger.Wrap(l)
			ReadRows(ctx, logger, f, "text/csv", scenario.fileRows)
			assert.Equal(t, scenario.expectedRows, len(scenario.fileRows))
		})
	}
}

func TestUpdateHeaderIndexes(t *testing.T) {
	type TestCases struct {
		description     string
		expectedHeaders []string
		headers         []string
		expected        map[string]int
	}
	testcases := []TestCases{
		{
			description:     "When headers are supplied in some order, the index map should match that order",
			expectedHeaders: []string{"uuid", "organization_name", "owner_user_id"},
			headers:         []string{"organization_name", "owner_user_id", "uuid"},
			expected:        map[string]int{"uuid": 2, "organization_name": 0, "owner_user_id": 1},
		},
		{
			description:     "When extra headers are supplied, the index map should still be correct",
			expectedHeaders: []string{"uuid", "organization_name", "owner_user_id"},
			headers:         []string{"uuid", "slogan", "owner_user_id", "organization_name"},
			expected:        map[string]int{"uuid": 0, "organization_name": 3, "owner_user_id": 2},
		},
		{
			description:     "When headers are missing, an error should be raised",
			expectedHeaders: []string{"uuid", "organization_name", "owner_user_id"},
			headers:         []string{"uuid", "owner_user_id"},
			expected:        nil,
		},
	}
	for _, scenario := range testcases {
		t.Run(scenario.description, func(t *testing.T) {
			headerIndexes, err := GetHeaderIndexes(scenario.expectedHeaders, scenario.headers)
			assert.Equal(t, scenario.expected, headerIndexes)
			if headerIndexes == nil {
				assert.ErrorContains(t, err, "missing header")
			}
		})
	}
}
