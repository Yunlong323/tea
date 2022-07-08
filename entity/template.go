package entity

type TemplateEntity struct {
	Name           string
	Classification []Classification
}

type Classification struct {
	Type    TemplateType
	Task    Task
	Service Service
	Nodes   []Nodes
}

type TemplateType string

const (
	Client TemplateType = "client"
	Server TemplateType = "server"
)

type Template struct {
	T *TemplateEntity
}

func NewTemplate(entity *TemplateEntity) *Template {
	t := Template{T: entity}
	return &t
}

func (t Template) getTemplateEntity() *TemplateEntity {
	return t.T
}

func (t *Template) setTemplateEntity(templateEntity *TemplateEntity) {
	t.T = templateEntity
}

type Service struct {
	Entry []string
	Api   []Api
}

type Api struct {
	Name   string
	Url    string
	Args   []Parameter
	Return []Parameter
}

type Task struct {
	Cmd  []string
	Args []Parameter
}
