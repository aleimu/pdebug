package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	coreV1 "k8s.io/api/core/v1"
	"strings"

	//"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"log"
	"os"
)

type tool struct {
	Data    *[]byte  // 源数据
	Kind    string   // 数据源的格式
	Type    string   // 指定输出格式 json|yaml
	Command []string // debug 命令行,用于阻塞pod退出
	Suffix  string   // 后缀
	transformer
}

type Deploy struct {
	Meta *v1.Deployment
}

type StatefulSet struct {
	Meta *v1.StatefulSet
}
type Job struct {
	Meta *batchv1.Job
}

type Pod struct {
	Meta *coreV1.Pod
}

type transformer interface {
	Read(b *[]byte)
	Replace(cmd []string, suffix string) (v interface{})
}

func (d *Deploy) Read(b *[]byte) {
	var data v1.Deployment
	err := json.Unmarshal(*b, &data)
	if err != nil {
		log.Fatal(err)
	}
	d.Meta = &data
}

func (d *Deploy) Replace(cmd []string, suffix string) (v interface{}) {
	d.Meta.Name = d.Meta.Name + suffix
	d.Meta.Status = v1.DeploymentStatus{}
	d.Meta.Spec.Template.Name = d.Meta.Spec.Template.Name + suffix
	containers := d.Meta.Spec.Template.Spec.Containers
	if len(containers) == 0 {
		log.Fatal("Containers 长度异常")
	}
	container := containers[0]
	container.Command = cmd
	container.ReadinessProbe = nil
	container.LivenessProbe = nil
	container.Lifecycle = nil
	d.Meta.Spec.Template.Spec.Containers[0] = container
	return d.Meta
}

func (s *StatefulSet) Read(b *[]byte) {
	var data v1.StatefulSet
	err := json.Unmarshal(*b, &data)
	if err != nil {
		log.Fatal(err)
	}
	s.Meta = &data
}
func (s *StatefulSet) Replace(cmd []string, suffix string) (v interface{}) {
	s.Meta.Name = s.Meta.Name + suffix
	s.Meta.Status = v1.StatefulSetStatus{}
	s.Meta.Spec.Template.Name = s.Meta.Spec.Template.Name + suffix
	containers := s.Meta.Spec.Template.Spec.Containers
	if len(containers) == 0 {
		log.Fatal("Containers 长度异常")
	}
	container := containers[0]
	container.Command = cmd
	container.ReadinessProbe = nil
	container.LivenessProbe = nil
	container.Lifecycle = nil
	s.Meta.Spec.Template.Spec.Containers[0] = container
	return s.Meta
}

func (j *Job) Read(b *[]byte) {
	var data batchv1.Job
	err := json.Unmarshal(*b, &data)
	if err != nil {
		log.Fatal(err)
	}
	j.Meta = &data
}
func (j *Job) Replace(cmd []string, suffix string) (v interface{}) {
	j.Meta.Name = j.Meta.Name + suffix
	j.Meta.Status = batchv1.JobStatus{}
	j.Meta.Spec.Template.Name = j.Meta.Spec.Template.Name + suffix
	containers := j.Meta.Spec.Template.Spec.Containers
	if len(containers) == 0 {
		log.Fatal("Containers 长度异常")
	}
	container := containers[0]
	container.Command = cmd
	container.ReadinessProbe = nil
	container.LivenessProbe = nil
	container.Lifecycle = nil
	j.Meta.Spec.Template.Spec.Containers[0] = container
	return j.Meta
}

func (p *Pod) Read(b *[]byte) {
	var data coreV1.Pod
	err := json.Unmarshal(*b, &data)
	if err != nil {
		log.Fatal(err)
	}
	p.Meta = &data
}
func (p *Pod) Replace(cmd []string, suffix string) (v interface{}) {
	p.Meta.Name = p.Meta.Name + suffix
	p.Meta.Status = coreV1.PodStatus{}
	p.Meta.OwnerReferences = nil
	p.Meta.Spec.NodeName = ""
	containers := p.Meta.Spec.Containers
	if len(containers) == 0 {
		log.Fatal("Containers 长度异常")
	}
	container := containers[0]
	container.Command = cmd
	container.ReadinessProbe = nil
	container.LivenessProbe = nil
	container.Lifecycle = nil
	p.Meta.Spec.Containers[0] = container
	return p.Meta
}

func readPipeline() *[]byte {
	stdinInfo, _ := os.Stdin.Stat()
	if (stdinInfo.Mode() & os.ModeNamedPipe) != os.ModeNamedPipe {
		log.Fatal("读取方式不是控制台输入")
	}
	b1 := new(bytes.Buffer)
	s := bufio.NewScanner(os.Stdin)
	s.Bytes()
	for s.Scan() {
		b1.Write(s.Bytes())
	}
	b2 := b1.Bytes()
	return &b2
}

func readFile(path string) *[]byte {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("文件读取异常:", err.Error())
	}
	return &f
}

func handleFlag() *tool {
	var operation string
	var cmd string
	var filepath string
	var suffix string
	flag.StringVar(&filepath, "f", "-", "默认从控制台读入数据,亦可指定文件")
	flag.StringVar(&operation, "t", "json", "输出格式,默认为json,可设定为yaml")
	flag.StringVar(&suffix, "s", "-copy", "后缀名,用于区别原文件中的资源名字")
	flag.StringVar(&cmd, "c", "sh -c /usr/sbin/sshd -D & touch debug && tail -f debug", "阻塞pod退出的命令")
	flag.Parse()
	t := new(tool)
	if filepath == "-" {
		t.Data = readPipeline()
	} else {
		t.Data = readFile(filepath)
	}
	t.Command = cmd2Cmds(cmd)
	t.Type = operation
	t.Suffix = suffix
	t.checkKind()
	return t
}

func (t *tool) checkKind() {
	// 可以通过官方库获取类型
	//us:=unstructured.Unstructured{}
	//us.GetKind()
	var data map[string]interface{}
	t.Kind = getKind(t.Data)
	if t.Kind != "" {
		return
	}
	err := yaml.Unmarshal(*t.Data, &data)
	if err != nil {
		log.Fatal("尝试yaml解析失败...", err)
	}
	b, err := yaml.YAMLToJSON(*t.Data)
	if err != nil {
		log.Fatal("尝试yaml转json失败...", err)
	}
	t.Data = &b
	t.Kind = getKind(&b)
}

func getKind(b *[]byte) string {
	var data map[string]interface{}
	if json.Valid(*b) {
		err := json.Unmarshal(*b, &data)
		if err != nil {
			log.Fatal("尝试json解析失败...", err)
		}
		if kind, ok := data["kind"]; !ok {
			log.Fatal("kind not found")
		} else {
			return kind.(string)
		}
	}
	return ""
}

func cmd2Cmds(cmd string) []string {
	cmds := []string{"/bin/bash", "-c", cmd}
	return cmds
}

func (t *tool) write(v interface{}) {
	var data []byte
	var err error
	if strings.ToLower(t.Type) == "json" {
		data, err = json.Marshal(v)
		if err != nil {
			log.Fatal("结果转换为json格式失败", err)
		}
	} else if strings.ToLower(t.Type) == "json" {
		data, err = yaml.Marshal(v)
		if err != nil {
			log.Fatal("结果转换为yaml格式失败", err)
		}
	}
	if len(data) == 0 {
		log.Fatal("转换失败")
	}
	os.Stdout.Write(data)
}

func main() {
	// 读取cmd
	t := handleFlag()
	// 策略模式
	switch t.Kind {
	case "Pod":
		t.transformer = new(Pod)
	case "StatefulSet":
		t.transformer = new(StatefulSet)
	case "Deployment":
		t.transformer = new(Deploy)
	case "Job":
		t.transformer = new(Job)
	default:
		log.Fatalf("%s 类型无对应实现", t.Kind)
	}
	t.Read(t.Data)
	v := t.Replace(t.Command, t.Suffix)
	// 输出结果和报错
	t.write(v)
}

/*
- 实现从shell管道中读取json/yaml并转换
- 发生错误时直接log.Fatal退出程序
*/
