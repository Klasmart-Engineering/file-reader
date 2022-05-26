package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRowToOrganizationAvro(t *testing.T) {
	type TestCases struct {
		description   string
		org           []string
		headerIndexes map[string]int
		trackingId    string
		schemaId      int
		expected      []byte
	}

	for _, scenario := range []TestCases{
		{
			description:   "should return byte encoding which starts with schemaId",
			org:           []string{"38ad4d35-0a29-436a-97bd-df26adbea6b4", "org0", "e893b974-dc12-11ec-9d64-0242ac120002"},
			headerIndexes: map[string]int{"uuid": 0, "organization_name": 1, "owner_user_id": 2},
			trackingId:    "trackingid1",
			schemaId:      3,
			expected:      []byte{0, 0, 0, 0, 3},
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			bytes, err := RowToOrganizationAvro(scenario.org, scenario.trackingId, scenario.schemaId, scenario.headerIndexes)
			assert.Equal(t, scenario.expected, bytes[:5])
			assert.Equal(t, err, nil)
		})
	}
}
