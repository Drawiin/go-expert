package events

import "errors"

var ErrHandlerAlreadyRegistered = errors.New("handler already registered")

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	if handlers, ok := ed.handlers[eventName]; ok {
		for _, h := range handlers {
			if h.Id() == handler.Id() {
				return ErrHandlerAlreadyRegistered
			}
		}
	}
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

func (ed *EventDispatcher) Clear() {
	ed.handlers = map[string][]EventHandlerInterface{}
}


func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool { 
	if handlers, ok := ed.handlers[eventName]; ok {
		for _, h := range handlers {
			if h.Id() == handler.Id() {
				return true
			}
		}
	}

	return false
}