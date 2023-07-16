package fieldmask_test

import (
	"encoding/json"
	"testing"

	"github.com/longkai/yyds/fieldmask"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/apipb"
	"google.golang.org/protobuf/types/known/sourcecontextpb"
)

func TestFieldMask(t *testing.T) {
	// test both top level and nested level
	var msg = &apipb.Api{
		Name: "hello api",
		SourceContext: &sourcecontextpb.SourceContext{
			FileName: "api.proto",
		},
		Version: "v1",
	}
	cases := []struct {
		desc  string
		paths []string
		want  string
	}{
		{
			desc: "empty path",
			want: `{
  "name": "hello api",
  "version": "v1",
  "sourceContext": {
    "fileName": "api.proto"
  }
}`,
		},
		{
			desc:  "single top level",
			paths: []string{"version"},
			want: `{
  "version": "v1"
}`,
		},
		{
			desc:  "single second level",
			paths: []string{"source_context.file_name"},
			want: `{
  "sourceContext": {
    "fileName": "api.proto"
  }
}`,
		},
		{
			desc:  "both top second level",
			paths: []string{"source_context.file_name", "version"},
			want: `{
  "version": "v1",
  "sourceContext": {
    "fileName": "api.proto"
  }
}`,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			msg := proto.Clone(msg)
			if err := fieldmask.Mask(msg, c.paths...); err != nil {
				t.Fatal(err)
			}
			b, err := stableMarshal(msg)
			if err != nil {
				t.Fatal(err)
			}
			if str := string(b); str != c.want {
				t.Fatalf("want %s, got %s", c.want, str)
			}
		})
	}
}

// protojson will add random space... https://github.com/golang/protobuf/issues/1082
func stableMarshal(msg proto.Message) ([]byte, error) {
	b, err := protojson.MarshalOptions{}.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(json.RawMessage(b), "", "  ")
}
