package client

// Annotation is the interface for annotation instances
type Annotation interface {
	// Type
}

// AnnotationContainer is the structure for annotation container
type AnnotationContainer struct {
	annotations []Annotation
}

// AddAnnotation adds an annotation
func (op *AnnotationContainer) AddAnnotation(ann Annotation) {
	op.annotations = append(op.annotations, ann)
}

// Annotations returns all annotations
func (op *AnnotationContainer) Annotations() []Annotation {
	return op.annotations
}
