package client

// Operation is the metadata of operation
type Operation struct {
	AnnotationContainer
	res      *Resource
	name     string
	problems map[StructPathID]*StructDefinition
	input    *StructDefinition
	output   *StructDefinition
	streamed bool
}

func newOperation(res *Resource, name string) *Operation {
	return &Operation{
		res:      res,
		name:     name,
		problems: map[StructPathID]*StructDefinition{},
	}
}

// Name returns the name of operation
func (op *Operation) Name() string {
	return op.name
}

// Resource returns the resource of operation
func (op *Operation) Resource() *Resource {
	return op.res
}

// RegisterProblemType registers a problem type
func (op *Operation) RegisterProblemType(tp *StructDefinition) {
	op.problems[tp.Path().ID()] = tp
}

// IsProblemType returns true if the type is a problem type
func (op *Operation) IsProblemType(tp StructPathID) bool {
	_, ok := op.problems[tp]
	return ok
}

// SetInput sets the input type
func (op *Operation) SetInput(tp *StructDefinition) {
	op.input = tp
}

// Input returns the input type
func (op *Operation) Input() *StructDefinition {
	return op.input
}

// SetOutput sets the output type
func (op *Operation) SetOutput(tp *StructDefinition) {
	op.output = tp
}

// SetStreamed sets the streamed flag
func (op *Operation) SetStreamed() {
	op.streamed = true
}

// IsStreamed returns true if the operation is streamed
func (op *Operation) IsStreamed() bool {
	return op.streamed
}

// Output returns the output type
func (op *Operation) Output() *StructDefinition {
	return op.output
}
