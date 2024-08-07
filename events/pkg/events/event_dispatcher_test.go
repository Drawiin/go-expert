package events

import (
	"testing"
	"time"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload interface{}
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

type TestEventHandler struct {
	id string
}

func (h *TestEventHandler) Handle(event EventInterface) {
	// do nothing
}

func (h *TestEventHandler) Id() string {
	return h.id
}

type EventDispatcherTestSuite struct {
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher *EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.event = TestEvent{Name: "test", Payload: "test"}
	suite.event2 = TestEvent{Name: "test2", Payload: "test2"}
	suite.handler = TestEventHandler{id: "1"}
	suite.handler2 = TestEventHandler{id: "2"}
	suite.handler3 = TestEventHandler{id: "3"}
	suite.eventDispatcher = NewEventDispatcher()
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	suite.Equal(&suite.handler, suite.eventDispatcher.handlers[suite.event.GetName()][0])
	suite.Equal(&suite.handler2, suite.eventDispatcher.handlers[suite.event.GetName()][1])
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register_WithSameHandler() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Equal(ErrHandlerAlreadyRegistered, err)
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler)
	suite.Nil(err)

	suite.eventDispatcher.Clear()

	suite.Equal(0, len(suite.eventDispatcher.handlers))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Has() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)

	suite.True(suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler))
	suite.False(suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler2))
	suite.False(suite.eventDispatcher.Has(suite.event2.GetName(), &suite.handler))
}

type MockEventHandler struct {
	mock.Mock
}

func (m *MockEventHandler) Handle(event EventInterface) {
	m.Called(event)
}

func (m *MockEventHandler) Id() string {
	return "mock"
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Dispatch() {
	eh := &MockEventHandler{}
	eh.On("Handle", &suite.event)

	suite.eventDispatcher.Register(suite.event.GetName(), eh)
	suite.eventDispatcher.Dispatch(&suite.event)
	eh.AssertExpectations(suite.T())
	eh.AssertNumberOfCalls(suite.T(), "Handle", 1)
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Remove() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)

	err = suite.eventDispatcher.Remove(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))
	suite.Equal(&suite.handler2, suite.eventDispatcher.handlers[suite.event.GetName()][0])
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
