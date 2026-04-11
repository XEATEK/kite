package wire

import (
	"testing"

	"github.com/kitecloud/kite/kite-service/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestEventListenerToWire_NilInput(t *testing.T) {
	assert.Nil(t, EventListenerToWire(nil))
}

func TestEventListenerToWire_NilFilter(t *testing.T) {
	el := &model.EventListener{}
	result := EventListenerToWire(el)
	assert.NotNil(t, result)
	assert.Nil(t, result.Filter)
}

func TestEventListenerToWire_FilterNilMessageReaction(t *testing.T) {
	el := &model.EventListener{
		Filter: &model.EventListenerFilter{},
	}
	result := EventListenerToWire(el)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Filter)
	assert.Nil(t, result.Filter.MessageReaction)
}

func TestEventListenerToWire_FilterWithMessageReaction(t *testing.T) {
	el := &model.EventListener{
		Filter: &model.EventListenerFilter{
			MessageReaction: &model.EventListenerFilterMessageReaction{
				Emoji: "👍",
			},
		},
	}
	result := EventListenerToWire(el)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Filter)
	assert.NotNil(t, result.Filter.MessageReaction)
	assert.Equal(t, "👍", result.Filter.MessageReaction.Emoji)
}

func TestEventListenerToWire_FilterWithCustomEmoji(t *testing.T) {
	el := &model.EventListener{
		Filter: &model.EventListenerFilter{
			MessageReaction: &model.EventListenerFilterMessageReaction{
				Emoji: "myemoji:123456789",
			},
		},
	}
	result := EventListenerToWire(el)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Filter)
	assert.NotNil(t, result.Filter.MessageReaction)
	assert.Equal(t, "myemoji:123456789", result.Filter.MessageReaction.Emoji)
}
