package util

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	UUID              = "uuid"
	ORGANIZATION_NAME = "organization_name"
	OWNER_USER_ID     = "owner_user_id"
)

const (
	TYPE_OF_UUID    = 0
	TYPE_OF_ID_LIST = 1
	TYPE_DEFAULT    = 2
	ORGANIZATION    = "organization"
)

type operationConfig struct {
	name           string
	opNum          int
	headers        []string
	fieldType      map[string]int
	fieldIdListNum map[string]int
}

var FieldType = map[string]int{
	"uuid": TYPE_OF_UUID,
}

func NewOperationConfig(name string, opNum int, headers []string, fieldType map[string]int, fieldIdListNum map[string]int) *operationConfig {
	return &operationConfig{
		name:           name,
		headers:        headers,
		fieldType:      fieldType,
		fieldIdListNum: fieldIdListNum,
	}
}
func generateUuidList(num int) []string {
	res := make([]string, num)
	for i := 0; i < num; i++ {
		res[i] = uuid.NewString()
	}
	return res
}
func (c operationConfig) MakeCsv() (csv *strings.Reader, op []map[string]string) {
	columnNames := c.headers

	// Reorder columns to random order and map column names to their random index
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(columnNames), func(i, j int) {
		columnNames[i], columnNames[j] = columnNames[j], columnNames[i]
	})
	colIndexMap := map[string]int{}
	for i, col := range columnNames {
		colIndexMap[col] = i
	}

	// Create operation (gets returned for use in assertions)
	ops := []map[string]string{}
	for i := 0; i < c.opNum; i++ {
		op := map[string]string{}
		for _, col := range columnNames {
			switch c.fieldType[col] {
			case TYPE_OF_UUID:
				op[col] = uuid.NewString()
			case TYPE_OF_ID_LIST:
				op[col] = strings.Join(generateUuidList(c.fieldIdListNum[col]), ",")
			default:
				switch col {
				case ORGANIZATION_NAME:
					op[col] = "org" + strconv.Itoa(i)
				default:
					op[col] = uuid.NewString()
				}
			}
		}
		ops = append(ops, op)
	}

	// Create file for test
	lines := []string{strings.Join(columnNames, ",")}
	for _, op := range ops {
		cols := make([]string, len(columnNames))
		for _, col := range columnNames {
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
