package events

import (
	"errors"
	"sync"
)

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

func (ed *EventDispatcher) Dispatch(event EventInterface) error {
	if handlers, ok := ed.handlers[event.GetName()]; ok {
		wg := &sync.WaitGroup{}
		for _, handler := range handlers {
			wg.Add(1)
			go func() {
				handler.Handle(event)
				wg.Done()
			}()
		}
		wg.Wait()
	}

	return nil
}

func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	if handlers, ok := ed.handlers[eventName]; ok {
		for i, h := range handlers {
			if h.Id() == handler.Id() {
				ed.handlers[eventName] = append(handlers[:i], handlers[i+1:]...)
				return nil
			}
		}
	}

	return nil
}
