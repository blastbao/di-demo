package di

import (
	"sync"
	"reflect"
	"fmt"
	"strings"
	"errors"
)

var (
	ErrFactoryNotFound = errors.New("factory not found")
)

type factory = func() (interface{}, error)



// 容器
type Container struct {
	sync.Mutex
	singletons map[string]interface{}
	factories  map[string]factory
}
// 容器实例化
func NewContainer() *Container {
	return &Container{
		singletons: make(map[string]interface{}),
		factories:  make(map[string]factory),
	}
}

// 注册单例对象
func (p *Container) SetSingleton(name string, singleton interface{}) {
	p.Lock()
	p.singletons[name] = singleton
	p.Unlock()
}

// 获取单例对象
func (p *Container) GetSingleton(name string) interface{} {
	return p.singletons[name]
}

// 获取实例对象
func (p *Container) GetPrototype(name string) (interface{}, error) {
	factory, ok := p.factories[name]
	if !ok {
		return nil, ErrFactoryNotFound
	}
	return factory()
}

// 设置实例对象工厂
func (p *Container) SetPrototype(name string, factory factory) {
	p.Lock()
	p.factories[name] = factory
	p.Unlock()
}

// 注入依赖

// 最重要的是Ensure方法，该方法通过反射来扫描 instance 实例的所有export字段，并读取di标签，如果有该标签则启动注入。
// 判断di标签的类型来确定注入 singleton 或者 prototype 对象。
func (p *Container) Ensure(instance interface{}) error {

	// 注意，这里的 instance 一定是某种指针类型（golang中指针也是一种interface{})，
	// 因此这里无论获取Type还是Value都需要通过Elem()获取其实际元素类型。

	// 另注：在实际测试中，发现不是指针，也可以正常注入，反射并不强制要求被set的字段一定是指针，这块是个疑问。

	elemType := reflect.TypeOf(instance).Elem()
	ele := reflect.ValueOf(instance).Elem()


	for i := 0; i < elemType.NumField(); i++ { 	// 遍历字段
		fieldType := elemType.Field(i)			// 字段类型
		tag := fieldType.Tag.Get("di") 		// 获取字段tag
		diName := p.injectName(tag)				// 取出tags[0]，为空则忽略
		if diName == "" {
			continue
		}

		var (
			diInstance interface{}
			err        error
		)

		//若是单例类型注入，则通过diName获取被注入对象的单例实例
		if p.isSingleton(tag) {
			diInstance = p.GetSingleton(diName)
		}
		//若是工厂类型注入，则通过diName创建新对象并注入该实例
		if p.isPrototype(tag) {
			diInstance, err = p.GetPrototype(diName)
		}
		//工厂类型注入时，创建新对象出错，则注入失败，报错
		if err != nil {
			return err
		}
		//单例/工厂类型获取被注入对象的实例均出错，报错
		if diInstance == nil {
			return errors.New(diName + " dependency not found")
		}

		//把被注入对象diInstance的值设置为当前字段值（均是指针类型）
		ele.Field(i).Set(reflect.ValueOf(diInstance))
	}
	return nil
}

// 获取需要注入的依赖名称
func (p *Container) injectName(tag string) string {
	tags := strings.Split(tag, ",")
	if len(tags) == 0 {
		return ""
	}
	return tags[0]
}

// 检测是否单例依赖
func (p *Container) isSingleton(tag string) bool {
	tags := strings.Split(tag, ",")
	for _, name := range tags {
		if name == "prototype" {
			return false
		}
	}
	return true
}

// 检测是否实例依赖
func (p *Container) isPrototype(tag string) bool {
	tags := strings.Split(tag, ",")
	for _, name := range tags {
		if name == "prototype" {
			return true
		}
	}
	return false
}

// 打印容器内部实例
func (p *Container) String() string {
	lines := make([]string, 0, len(p.singletons)+len(p.factories)+2)
	lines = append(lines, "singletons:")
	for name, item := range p.singletons {
		line := fmt.Sprintf("  %s: %x %s", name, &item, reflect.TypeOf(item).String())
		lines = append(lines, line)
	}
	lines = append(lines, "factories:")
	for name, item := range p.factories {
		line := fmt.Sprintf("  %s: %x %s", name, &item, reflect.TypeOf(item).String())
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
