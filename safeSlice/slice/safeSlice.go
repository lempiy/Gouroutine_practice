package safeSlice

type UpdateFunc func(interface{}) interface{}

type SafeSlice interface {
	Append(interface{})
	At(int) interface{}
	Close() []interface{}
	Delete(int)
	Len() int
	Update(int, UpdateFunc)
}

type commandData struct {
	action  commandAction
	index   int
	value   interface{}
	result  chan<- interface{}
	data    chan<- []interface{}
	updater UpdateFunc
}

type commandAction int

const (
	remove commandAction = iota
	end
	at
	add
	length
	update
)

type safeSlice chan commandData

func (sfSlice safeSlice) Append(value interface{}) {
	sfSlice <- commandData{action: add, value: value}
}

func (sfSlice safeSlice) At(index int) interface{} {
	reply := make(chan interface{})
	sfSlice <- commandData{action: at, index: index, result: reply}
	return <-reply
}

func (sfSlice safeSlice) Delete(index int) {
	sfSlice <- commandData{action: remove, index: index}
}

func (sfSlice safeSlice) Len() int {
	reply := make(chan interface{})
	sfSlice <- commandData{action: length, result: reply}
	result := (<-reply).(int)
	return result
}

func (sfSlice safeSlice) Update(index int, updater UpdateFunc) {
	sfSlice <- commandData{action: update, index: index, updater: updater}
}

func (sfSlice safeSlice) Close() []interface{} {
	reply := make(chan []interface{})
	sfSlice <- commandData{action: end, data: reply}
	return <-reply
}

func (sfSlice safeSlice) run() {
	var store []interface{}
	for command := range sfSlice {
		switch command.action {
		case add:
			store = append(store, command.value)
		case remove:
			store = append(store[:command.index], store[command.index+1:]...)
		case at:
			command.result <- store[command.index]
		case length:
			command.result <- len(store)
		case update:
			store[command.index] = command.updater(store[command.index])
		case end:
			close(sfSlice)
			command.data <- store
		}
	}
}

func New() SafeSlice {
	sfSlice := make(safeSlice)
	go sfSlice.run()
	return sfSlice
}
