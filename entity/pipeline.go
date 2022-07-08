package entity

type Pipeline struct {
	P *PipelineEntity
}

type PipelineEntity struct {
	Name  string
	Steps []Step
	Nodes []Nodes
}

type Nodes struct {
	Name   string
	Socket string
}

func NewPipeline(entity *PipelineEntity) *Pipeline {
	p := Pipeline{P: entity}
	return &p
}

type Step struct {
	Name  string
	Works []Work
}

type FinallyStatusEnum bool

const (
	FinallyOnline FinallyStatusEnum = true
	FinallyDone   FinallyStatusEnum = false
)

type Work struct {
	Template     string
	Name         string
	Description  string
	Parameters   []Parameter
	Deployment   Deployment
	Dependencies []Dependencies
	Serving      FinallyStatusEnum
}

type Deployment struct {
	Client string
	Server string
}

type Dependencies struct {
	Dependency string
}
