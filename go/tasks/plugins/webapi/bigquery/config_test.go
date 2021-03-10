package bigquery

import (
	"github.com/magiconair/properties/assert"
	"google.golang.org/protobuf/types/known/structpb"
	"testing"
)

func TestUnmarshalBigQueryQueryConfig(t *testing.T) {
	custom := structpb.Struct{
		Fields: map[string]*structpb.Value{
			"projectId": structpb.NewStringValue("project-id"),
			"location": structpb.NewStringValue("EU"),
			"query": structpb.NewStringValue("SELECT 1"),
		},
	}

	config, err := unmarshalBigQueryQueryConfig(&custom)

	assert.Equal(t, err, nil)

	assert.Equal(t, config, &QueryJobConfig{
		ProjectID: "project-id",
		Location: "EU",
		Query: "SELECT 1",
	})

}
