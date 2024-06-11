package client

// import "fmt"

// // ContentObjectProperty gets an object property from a content
// func ContentObjectProperty[S any](content Content, pk *Package, name string, field *S) error {
// 	v := content.GetProperty(name)
// 	if v == nil {
// 		return nil
// 	}
// 	vmap, ok := v.(map[string]any)
// 	if !ok {
// 		return fmt.Errorf("expected map in field %s", name)
// 	}
// 	vc := Content(vmap)
// 	fqdn, err := TypeFQTNFromString(vc.GetStruct())
// 	if err != nil {
// 		return err
// 	}
// 	tp, err := pk.TypeByPath(fqdn)
// 	if err != nil {
// 		// TODO: polymorphic hydratation
// 		return err
// 	}
// 	object := tp.TypeBuilder()()
// 	if err := StructFromContent(vc, pk, object); err != nil {
// 		return err
// 	}
// 	sobj := object.(S)
// 	*field = sobj
// 	return nil
// }

// // ContentObjectArrayProperty gets an object array property from a content
// func ContentObjectArrayProperty[S any, A ~[]S](content Content, pk *Package, name string, sp *StructPath, field *A) error {
// 	v := content.GetProperty(name)
// 	if v == nil {
// 		return nil
// 	}
// 	vs, ok := v.([]any)
// 	if !ok {
// 		return fmt.Errorf("expected array in field %s", name)
// 	}
// 	structs := make(A, 0, len(vs))
// 	for _, vmap := range vs {
// 		vc := Content(vmap.(map[string]any))
// 		fqdn, err := TypeFQTNFromString(vc.GetStruct())
// 		if err != nil {
// 			return err
// 		}
// 		tp, err := pk.TypeByPath(fqdn)
// 		if err != nil {
// 			// TODO: polymorphic hydratation
// 			return err
// 		}
// 		object := tp.TypeBuilder()()
// 		if err := StructFromContent(vc, pk, object); err != nil {
// 			return err
// 		}
// 		sobj := object.(S)
// 		structs = append(structs, sobj)
// 	}
// 	*field = structs
// 	return nil
// }
