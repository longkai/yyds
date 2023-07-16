package fieldmask_test

import (
	"testing"

	"github.com/longkai/yyds/fieldmask"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/apipb"
	"google.golang.org/protobuf/types/known/sourcecontextpb"
)

func TestFieldMask(t *testing.T) {
	// test both top level and nested level
	api := &apipb.Api{
		Name: "hello api",
		SourceContext: &sourcecontextpb.SourceContext{
			FileName: "api.proto",
		},
		Version: "v1",
	}

	if err := fieldmask.Mask(api, "version", "source_context.file_name"); err != nil {
		t.Fatal(err)
	}

	b, err := protojson.MarshalOptions{Indent: "  "}.Marshal(api)
	if err != nil {
		t.Fatal(err)
	}
	want := `{
  "version":  "v1",
  "sourceContext":  {
    "fileName":  "api.proto"
  }
}`
	if str := string(b); str != want {
		t.Fatalf("want %s, got %s", want, str)
	}
}
