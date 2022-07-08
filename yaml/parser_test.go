package yaml

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"tea/entity"
	"testing"
)

const (
	templateYaml = `
name: evaluate # 模板名
kind: template

classification:
  - type: client
  # running阶段依次执行命令
    task:
      args:
        - name: EmployeeInfoPath
          type: string
      cmd:
        - python getEmployeeInfo.py EmployeeInfoPath # 获取模型的操作放在service中，因为service的启动很严格。必须等待server端把model传过来才可启动service
    service:
      entry:                      # 使用命令行起一个serviceA
        - java -jar evaluation.jar     # 起一个Java的service
      api:                      # service提供的具体接口列表
        - name: getQualifiedEmployee        # 接口名
          url: "/api/getQualifiedEmployee"  # 接口url
          args:
            - name: database # 员工的信息都存在一个DB中
              type: string
            - name: modelFileName # 训练模型的文件名
              type: string
            - name: outputFormat
              type: string # csv,txt,xlsx,sql
          return:               # 返回值
            - name: employeeId
              type: list[int] # dapr 类型
    nodes:
      - name: XingYe
        socket: 121.192.123.66:8083

  - type: server
    task:
      args:
        - name: ModelPath
          type: string
        - name: TargetPath
          type: string
      cmd:
        - python getEvaluateModel.py ModelPath TargetPath # 把模型从modelpath传到targetPath
    nodes:
      - name: Company
        socket: 121.192.13.66:8083

`

	templateErrorYaml = `
name: evaluate # 模板名
kind: template
type: nothetype   # 这里我想知道枚举只局限一个新类型值为“server” “client” 如何实现
# running阶段依次执行命令
task:
  args:
    - name: ModelPath
      type: string
    - name: TargetPath #这里错就是格式错误
      type: string
  cmdDisturbed:  # 解析出错 下面一行即便加上注释也会报错
    - python getEvaluateModel.py ModelPath TargetPath # 把模型从modelpath传到targetPath
`
	pipelineYaml = `
# yaml模板示例
name: default   # 名称
kind: pipeline  # 类型
nodes:
  - name: xingye bank
    socket: 1.2.3.4:7788
  - name: electricity company
    socket: 1.2.5.6:7788
# 任务编排
steps:
  # 步骤名
  - name: step1
    # 一个步骤所包含的任务群
    works:
      # 任务所使用的模板
      - template: tmp1 name1
        # 任务名
        name: task1/1
        # 任务描述
        description: a simple task
        # 任务所需的参数
        parameters:
          - name: v1
            type: int
          - name: v2
            type: string
        deployment:
          client: electricity company
          server: xingye bank
        serving: true
      - template: tmp1 name2
        name: task1/2
        description: a simple task
        parameters:
          - name: v3
            type: int
        deployment:
          client: xingye bank
        # 当前任务所依赖的其他任务
        dependencies:
          # 外部依赖的任务
          - dependency: another_task # 只支持internal
        serving: true
  - name: step2
    works:
      - template: tmp1 name3
        name: task2/1
        description: a simple task
        parameters:
          - name: v4
            type: string
        # 仍在线的任务，其余任务全部自动下线
        serving: true
      - template: tmp1 name4
        name: task2/2
        description: a simple task
        parameters:
          - name: v5
            type: int
        dependencies:
          # 内部依赖的任务
          - dependency: task1/1
`

	errorPipelineYaml = `
# yaml模板示例
name: default   # 名称
kind: pipeline  # 类型错误

# 任务编排
steps_error:
  # 步骤名
  - extra: extra
  - name: step1
    # 一个步骤所包含的任务群
    works:
      # 任务所使用的模板
      - template: tmp1 name1
        # 任务名
        name: task1/1
        # 任务描述
        description: a simple task
        # 任务所需的参数
        parameters:
          - name: v1
            type: int
          - name: v2
            type: string
        serving: true
      - template: tmp1 name2
        name: task1/2
        description: a simple task
        parameters:
          - name: v3
            type: int
        # 当前任务所依赖的其他任务
        dependencies:
          # 外部依赖的任务
          - dependency: another_task # 只支持internal
        serving: true
  - name: step2
    works:
      - template: tmp1 name3
        name: task2/1
        description: a simple task
        parameters:
          - name: v4
            type: string
        # 仍在线的任务，其余任务全部自动下线
        serving: true
      - template: tmp1 name4
        name: task2/2
        description: a simple task
        parameters:
          - name: v5
            type: int
        dependencies:
          # 内部依赖的任务
          - dependency: task1/1
`
)

func TestParseTemplateSuccess(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name string
		args args
		want *entity.Template
	}{
		// TODO: Add test cases.
		{name: "template case 1", args: args{reader: strings.NewReader(templateYaml)}, want: &entity.Template{
			T: &entity.TemplateEntity{
				Name: "evaluate",
				Classification: []entity.Classification{
					{
						Type: "client",
						Task: entity.Task{
							Cmd: []string{
								"python getEmployeeInfo.py EmployeeInfoPath",
							},
							Args: []entity.Parameter{
								{
									Name: "EmployeeInfoPath",
									Type: "string",
								},
							},
						},
						Service: entity.Service{
							Entry: []string{
								"java -jar evaluation.jar",
							},
							Api: []entity.Api{
								{
									Name: "getQualifiedEmployee",
									Url:  "/api/getQualifiedEmployee",
									Args: []entity.Parameter{
										{
											Name: "database",
											Type: "string",
										}, {
											Name: "modelFileName",
											Type: "string",
										}, {
											Name: "outputFormat",
											Type: "string"},
									},
									Return: []entity.Parameter{
										{
											Name: "employeeId",
											Type: "list[int]",
										}},
								},
							},
						},
						Nodes: []entity.Nodes{
							{
								Name:   "XingYe",
								Socket: "121.192.123.66:8083",
							},
						},
					}, {
						Type: "server",
						Task: entity.Task{
							Cmd: []string{
								"python getEvaluateModel.py ModelPath TargetPath",
							},
							Args: []entity.Parameter{
								{
									Name: "ModelPath",
									Type: "string",
								},
								{
									Name: "TargetPath",
									Type: "string",
								},
							},
						},
						Nodes: []entity.Nodes{
							{
								Name:   "Company",
								Socket: "121.192.13.66:8083",
							},
						},
					},
				},
			}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseTemplate(tt.args.reader); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// pipeline 表格驱动测试 失败案例
// 1. 类型错误，传入的并非pipeline类型的yaml
// 2. 格式错误，传入的yaml无法被正确解析
func TestParseTemplateFail(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "template case 3", args: args{reader: strings.NewReader(pipelineYaml)}, want: "get kind pipeline instead of template"},
		{name: "template case 4", args: args{reader: strings.NewReader(templateErrorYaml)}, want: "failed to parse template"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					errMsg := fmt.Sprintf("%v", err)
					fmt.Println(errMsg)
					begin := strings.Index(errMsg, "map[error:")
					end := strings.Index(errMsg, "]")
					errMsg = errMsg[begin+len("map[error:") : end]
					fmt.Println(errMsg)
					if !reflect.DeepEqual(errMsg, tt.want) {
						t.Errorf("error = %v, want %v", errMsg, tt.want)
					}
				}
			}()
			_ = ParseTemplate(tt.args.reader)
		})
	}
}

// pipeline 表格驱动测试 成功案例
func TestParsePipelineSuccess(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name string
		args args
		want *entity.Pipeline
	}{
		{name: "pipeline case 1", args: args{reader: strings.NewReader(pipelineYaml)}, want: &entity.Pipeline{
			P: &entity.PipelineEntity{
				Name: "default",
				Nodes: []entity.Nodes{
					{
						Name:   "xingye bank",
						Socket: "1.2.3.4:7788",
					}, {
						Name:   "electricity company",
						Socket: "1.2.5.6:7788",
					},
				},
				Steps: []entity.Step{
					{
						Name: "step1",
						Works: []entity.Work{
							{
								Template:    "tmp1 name1",
								Name:        "task1/1",
								Description: "a simple task",
								Parameters: []entity.Parameter{
									{
										Name: "v1",
										Type: "int",
									}, {
										Name: "v2",
										Type: "string",
									},
								},
								Deployment: entity.Deployment{
									Client: "electricity company",
									Server: "xingye bank",
								},
								Serving: entity.FinallyOnline,
							},
							{
								Template:    "tmp1 name2",
								Name:        "task1/2",
								Description: "a simple task",
								Parameters: []entity.Parameter{
									{
										Name: "v3",
										Type: "int",
									},
								},
								Deployment: entity.Deployment{
									Client: "xingye bank",
								},
								Dependencies: []entity.Dependencies{
									{
										Dependency: "another_task",
									},
								},
								Serving: entity.FinallyOnline,
							},
						},
					},
					{
						Name: "step2",
						Works: []entity.Work{
							{
								Template:    "tmp1 name3",
								Name:        "task2/1",
								Description: "a simple task",
								Parameters: []entity.Parameter{
									{
										Name: "v4",
										Type: "string",
									},
								},
								Serving: entity.FinallyOnline,
							},
							{
								Template:    "tmp1 name4",
								Name:        "task2/2",
								Description: "a simple task",
								Parameters: []entity.Parameter{
									{
										Name: "v5",
										Type: "int",
									},
								},
								Dependencies: []entity.Dependencies{
									{
										Dependency: "task1/1",
									},
								},
							},
						},
					},
				},
			}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParsePipeline(tt.args.reader); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePipeline() = %v, want %v", got, tt.want)
			}
		})
	}
}

// pipeline 表格驱动测试 失败案例
// 1. 类型错误，传入的并非pipeline类型的yaml
// 2. 格式错误，传入的yaml无法被正确解析
func TestParsePipelineFailure(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "pipeline case 2", args: args{reader: strings.NewReader(templateYaml)}, want: "get kind template instead of pipeline"},
		{name: "pipeline case 3", args: args{reader: strings.NewReader(errorPipelineYaml)}, want: "failed to parse pipeline"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					errMsg := fmt.Sprintf("%v", err)
					fmt.Println(errMsg)
					begin := strings.Index(errMsg, "map[error:")
					end := strings.Index(errMsg, "]")
					errMsg = errMsg[begin+len("map[error:") : end]
					fmt.Println(errMsg)
					if !reflect.DeepEqual(errMsg, tt.want) {
						t.Errorf("error = %v, want %v", errMsg, tt.want)
					}
				}
			}()
			_ = ParsePipeline(tt.args.reader)
		})
	}
}
