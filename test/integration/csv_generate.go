package util

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	UUID                = "uuid"
	ORGANIZATION_NAME   = "organization_name"
	OWNER_USER_ID       = "owner_user_id"
	TYPE_OF_UUID        = 0
	TYPE_OF_ID_LIST     = 1
	TYPE_DEFAULT        = 2
	TYPE_OPERATION_NAME = 3
)

// Define func for different data field types
type DataGenerator func(string, int) string
type operationConfig struct {
	name           string
	prefix         string
	opNum          int
	headers        []string
	fieldFunc      DataGenerator
	fieldType      map[string]int
	fieldIdListNum map[string]int
}

func (c operationConfig) generateData(col string, i int) string {
	switch c.fieldType[col] {
	case TYPE_OPERATION_NAME:
		return c.prefix + strconv.Itoa(i)
	case TYPE_OF_UUID:
		return uuid.NewString()
	case TYPE_DEFAULT:
		return uuid.NewString()
	case TYPE_OF_ID_LIST:
		return generateUuidListAsString(c.fieldIdListNum[col])
	default:
		return uuid.NewString()
	}

}

func NewOperationConfig(name string, prefix string, opNum int, headers []string, fieldType map[string]int, fieldIdListNum map[string]int) *operationConfig {
	return &operationConfig{
		name:           name,
		prefix:         prefix,
		opNum:          opNum,
		headers:        headers,
		fieldType:      fieldType,
		fieldIdListNum: fieldIdListNum,
	}
}

func generateUuidListAsString(idNum int) string {
	res := make([]string, idNum)
	for i := 0; i < idNum; i++ {
		res[i] = uuid.NewString()
	}
	return strings.Join(res, ";")
}
func (c operationConfig) MakeCsv() (csv *strings.Reader, op []map[string]string) {
	columnNames := c.headers
	c.fieldFunc = c.generateData
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
			op[col] = c.fieldFunc(col, i)

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
