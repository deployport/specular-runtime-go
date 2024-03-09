package client

import (
	"fmt"
	"time"

	b64 "encoding/base64"

	"golang.org/x/exp/constraints"
)

// ContentNumericConstraint is a constraint for a string
type ContentNumericConstraint interface {
	constraints.Integer | constraints.Float
}

// ContentStringConstraint is a constraint for a string
type ContentStringConstraint interface {
	~string
}

// ContentRequireNumericProperty gets a numeric property from a content
func ContentRequireNumericProperty[V ContentNumericConstraint](content Content, name string, field *V) error {
	v := content.GetProperty(name)
	if v == nil {
		return fmt.Errorf("missing property %s", name)
	}
	fv, ok := v.(float64)
	if ok {
		vv := V(fv)
		*field = vv
		return nil
	}
	vv := v.(V)
	*field = vv
	return nil
}

// ContentOptionalNumericProperty gets a numeric property from a content
func ContentOptionalNumericProperty[V ContentNumericConstraint](content Content, name string, field **V) error {
	v := content.GetProperty(name)
	if v == nil {
		return nil
	}
	fv, ok := v.(float64)
	if ok {
		vv := V(fv)
		*field = &vv
		return nil
	}
	vv := v.(V)
	*field = &vv
	return nil
}

// ContentOptionalStringProperty gets an optional string property from a content
func ContentOptionalStringProperty[V ContentStringConstraint](content Content, name string, field **V) error {
	v := content.GetProperty(name)
	if v == nil {
		return nil
	}
	fv, ok := v.(string)
	if !ok {
		return nil
	}
	vv := V(fv)
	*field = &vv
	return nil
}

// ContentRequireStringProperty gets a string property from a content, if the property does not exist it or or is not an string or is nil/null, it returns an error
func ContentRequireStringProperty[V ContentStringConstraint](content Content, name string, field *V) error {
	v := content.GetProperty(name)
	if v == nil {
		return fmt.Errorf("missing property %s", name)
	}
	vs, ok := v.(string)
	if !ok {
		return fmt.Errorf("string expected in property %s", name)
	}
	vv := V(vs)
	*field = vv
	return nil
}

// ContentOptionalTimeProperty gets an optional time property from a content or leaves it unassigned
func ContentOptionalTimeProperty(content Content, name string, field **time.Time) error {
	var ost *string
	if err := ContentOptionalStringProperty(content, name, &ost); err != nil {
		return err
	}
	if ost == nil {
		return nil
	}
	st := *ost
	t, err := time.Parse(time.RFC3339, st)
	if err != nil {
		return err
	}
	*field = &t
	return nil
}

// ContentRequireTimeProperty gets a time property from a content or returns an error
func ContentRequireTimeProperty(content Content, name string, field *time.Time) error {
	var ost string
	if err := ContentRequireStringProperty(content, name, &ost); err != nil {
		return fmt.Errorf("failed to require time property, %w", err)
	}
	t, err := time.Parse(time.RFC3339, ost)
	if err != nil {
		return err
	}
	*field = t
	return nil
}

// ContentOptionalBinaryProperty gets an []byte property from a content, if not found it leaves the field unassigned
func ContentOptionalBinaryProperty(content Content, name string, field *[]byte) error {
	var osb *string
	if err := ContentOptionalStringProperty(content, name, &osb); err != nil {
		return fmt.Errorf("failed to require binary property, %w", err)
	}
	if osb == nil {
		return nil
	}
	sb, err := b64.StdEncoding.DecodeString(*osb)
	if err != nil {
		return err
	}
	*field = sb
	return nil
}

// ContentRequireBinaryProperty gets an []byte property from a content or returns an error if missing
func ContentRequireBinaryProperty(content Content, name string, field *[]byte) error {
	var lf *[]byte
	if err := ContentOptionalBinaryProperty(content, name, lf); err != nil {
		return fmt.Errorf("failed to require binary property, %w", err)
	}
	if lf == nil {
		return fmt.Errorf("missing property %s", name)
	}
	*field = *lf
	return nil
}

// ContentRequireBoolProperty gets a bool property from a content or returns an error
func ContentRequireBoolProperty(content Content, name string, field *bool) error {
	v := content.GetProperty(name)
	if v == nil {
		return fmt.Errorf("missing property %s", name)
	}
	vb, ok := v.(bool)
	if !ok {
		return fmt.Errorf("bool expected in property %s", name)
	}
	*field = vb
	return nil
}

// ContentOptionalBoolProperty gets a bool property from a content or leaves it unassigned, returns an error if the property is not bool
func ContentOptionalBoolProperty(content Content, name string, field **bool) error {
	v := content.GetProperty(name)
	if v == nil {
		return nil
	}
	vb, ok := v.(bool)
	if !ok {
		return fmt.Errorf("bool expected in property %s", name)
	}
	*field = &vb
	return nil
}

// BuildintConstraint is a constraint for a string or a float64
type BuildintConstraint interface {
	~string | ContentNumericConstraint
}

// ContentOptionalBuiltinArrayProperty gets an []T property from a content, if not found it leaves the field unassigned
func ContentOptionalBuiltinArrayProperty[T BuildintConstraint](content Content, name string, field *[]T) error {
	err := ContentRequireBuiltinArrayProperty(content, name, field)
	if IsPropertyRequiredError(err) {
		return nil
	} else if err != nil {
		return err
	}
	return nil
}

// ContentRequireBuiltinArrayProperty gets an []T property from a content, if not found it leaves the field unassigned
func ContentRequireBuiltinArrayProperty[T BuildintConstraint](content Content, name string, field *[]T) error {
	v := content.GetProperty(name)
	if v == nil {
		return NewPropertyRequiredError(name)
	}
	fv, ok := v.([]T)
	if !ok {
		// allocate and convert one by one
		vi, ok := v.([]any)
		if !ok {
			return fmt.Errorf("array expected in property %s, %T found instead", name, v)
		}
		fv = make([]T, len(vi))
		for ix, v := range vi {
			vs, ok := v.(T)
			if !ok {
				return fmt.Errorf("array item expected in property %s, %T found instead", name, v)
			}
			fv[ix] = vs
		}
	}
	*field = fv
	return nil
}

// ContentOptionalBuiltinNumericArrayProperty gets an []T property from a content, if not found it leaves the field unassigned
func ContentOptionalBuiltinNumericArrayProperty[T ContentNumericConstraint](content Content, name string, field *[]T) error {
	err := ContentRequireBuiltinNumericArrayProperty(content, name, field)
	if IsPropertyRequiredError(err) {
		return nil
	} else if err != nil {
		return err
	}
	return nil
}

// ContentRequireBuiltinNumericArrayProperty gets an []T property from a content, if not found it leaves the field unassigned
func ContentRequireBuiltinNumericArrayProperty[T ContentNumericConstraint](content Content, name string, field *[]T) error {
	v := content.GetProperty(name)
	if v == nil {
		return NewPropertyRequiredError(name)
	}
	fv, ok := v.([]float64)
	if !ok {
		// allocate and convert one by one
		vi, ok := v.([]any)
		if !ok {
			return fmt.Errorf("array expected in property %s, %T found instead", name, v)
		}
		fv = make([]float64, len(vi))
		for ix, v := range vi {
			vs, ok := v.(float64)
			if !ok {
				return fmt.Errorf("array item expected in property %s, %T found instead", name, v)
			}
			fv[ix] = vs
		}
	}
	res := make([]T, len(fv))
	for ix, v := range fv {
		res[ix] = T(v)
	}
	*field = res
	return nil
}
