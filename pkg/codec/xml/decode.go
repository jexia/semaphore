package xml

// import (
// 	"encoding/xml"
// 	"fmt"
//
// 	"github.com/jexia/semaphore/pkg/references"
// 	"github.com/jexia/semaphore/pkg/specs"
// 	"github.com/jexia/semaphore/pkg/specs/types"
// )
//
// func decodeValue(prop *specs.Property, resource string, raw xml.CharData, store references.Store) error {
// 	var ref = &references.Reference{
// 		Path: prop.Path,
// 	}
//
// 	if prop.Type == types.Enum {
// 		enum, ok := prop.Enum.Keys[string(raw)]
// 		if !ok {
// 			return fmt.Errorf("unknown enum %s", raw)
// 		}
//
// 		ref.Enum = &enum.Position
// 		store.StoreReference(resource, ref)
//
// 		return nil
// 	}
//
// 	value, err := DecodeType(string(raw), prop.Type)
// 	if err != nil {
// 		return err
// 	}
//
// 	ref.Value = value
// 	store.StoreReference(resource, ref)
//
// 	return nil
// }
