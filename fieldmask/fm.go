package fieldmask

import (
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// Mask filters the msg to contain only those fields specified in the paths.
//
// Empty paths will mask nothing.
//
// For more information please go to https://protobuf.dev/reference/protobuf/google.protobuf/#field-mask
func Mask(msg proto.Message, paths ...string) error {
	fm, err := fieldmaskpb.New(msg, paths...)
	if err != nil {
		return err
	}
	fm.Normalize()
	newMasker(fm.Paths).mask(msg)
	return nil
}

// masker is like a map based trie tree.
type masker map[string]masker

func newMasker(paths []string) masker {
	root := make(masker)
	for _, path := range paths {
		cur := root
		for _, seg := range strings.Split(path, ".") {
			if node, ok := cur[seg]; ok {
				cur = node
				continue
			}
			node := make(masker)
			cur[seg] = node
			cur = node
		}
	}
	return root
}

func (m masker) mask(msg proto.Message) {
	if len(m) == 0 {
		return
	}
	pr := msg.ProtoReflect()
	pr.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if node, ok := m[string(fd.Name())]; ok {
			if fd.Kind() == protoreflect.MessageKind {
				node.mask(v.Message().Interface())
			}
			return true
		}
		pr.Clear(fd)
		return true
	})
}
