package eventmgr

type EventType int32

const (
	EventPlayerLogin EventType = iota + 10001
	EventPlayerLogout
)

// Event 定义事件结构
type Event struct {
	EType EventType   // 事件类型
	Data  interface{} // 事件数据
}

// Listener 定义监听者接口
type Listener func(event Event)

// EventBus 定义同步事件总线
type EventBus struct {
	listeners map[EventType][]Listener // 事件名称 -> 监听器列表
}

// NewEventBus 创建一个新的事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		listeners: make(map[EventType][]Listener),
	}
}

// Subscribe 订阅指定事件
func (eb *EventBus) Subscribe(eventName EventType, listener Listener) {
	eb.listeners[eventName] = append(eb.listeners[eventName], listener)
}

// Publish 发布事件（同步执行）
func (eb *EventBus) Publish(event Event) {
	if listeners, exists := eb.listeners[event.EType]; exists {
		for _, listener := range listeners {
			listener(event) // 同步调用监听器
		}
	}
}
