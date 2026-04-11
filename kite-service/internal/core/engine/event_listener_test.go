package engine

import (
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/kitecloud/kite/kite-service/internal/model"
	"github.com/stretchr/testify/assert"
)

func newTestEventListener(filter *model.EventListenerFilter) *EventListener {
	return &EventListener{
		listener: &model.EventListener{
			Filter: filter,
		},
		// flow is not used by shouldHandleEvent
	}
}

func TestShouldHandleEvent_NoFilter(t *testing.T) {
	l := newTestEventListener(nil)

	assert.True(t, l.shouldHandleEvent(&gateway.MessageCreateEvent{Message: discord.Message{Author: discord.User{Bot: false}}}),
		"non-bot message create should be handled")
	assert.False(t, l.shouldHandleEvent(&gateway.MessageCreateEvent{Message: discord.Message{Author: discord.User{Bot: true}}}),
		"bot message create should be ignored")
	assert.True(t, l.shouldHandleEvent(&gateway.MessageDeleteEvent{}), "message delete should be handled")
	assert.True(t, l.shouldHandleEvent(&gateway.GuildMemberAddEvent{}), "guild member add should be handled")
	assert.True(t, l.shouldHandleEvent(&gateway.GuildMemberRemoveEvent{}), "guild member remove should be handled")
	assert.True(t, l.shouldHandleEvent(&gateway.MessageReactionAddEvent{}), "reaction add should be handled")
	assert.True(t, l.shouldHandleEvent(&gateway.MessageReactionRemoveEvent{}), "reaction remove should be handled")
}

func TestShouldHandleEvent_MessageReactionFilter_Unicode(t *testing.T) {
	emoji := "👍"
	l := newTestEventListener(&model.EventListenerFilter{
		MessageReaction: &model.EventListenerFilterMessageReaction{
			Emoji: emoji,
		},
	})

	matchingEmoji := discord.Emoji{Name: "👍"}
	nonMatchingEmoji := discord.Emoji{Name: "👎"}

	assert.True(t, l.shouldHandleEvent(&gateway.MessageReactionAddEvent{
		Emoji: matchingEmoji,
	}), "matching unicode emoji should be accepted for ReactionAdd")

	assert.False(t, l.shouldHandleEvent(&gateway.MessageReactionAddEvent{
		Emoji: nonMatchingEmoji,
	}), "non-matching unicode emoji should be rejected for ReactionAdd")

	assert.True(t, l.shouldHandleEvent(&gateway.MessageReactionRemoveEvent{
		Emoji: matchingEmoji,
	}), "matching unicode emoji should be accepted for ReactionRemove")

	assert.False(t, l.shouldHandleEvent(&gateway.MessageReactionRemoveEvent{
		Emoji: nonMatchingEmoji,
	}), "non-matching unicode emoji should be rejected for ReactionRemove")
}

func TestShouldHandleEvent_MessageReactionFilter_CustomEmoji(t *testing.T) {
	// Custom emoji APIString format is "name:id"
	filterEmoji := "myemoji:123456789"
	l := newTestEventListener(&model.EventListenerFilter{
		MessageReaction: &model.EventListenerFilterMessageReaction{
			Emoji: filterEmoji,
		},
	})

	matchingEmoji := discord.Emoji{Name: "myemoji", ID: discord.EmojiID(123456789)}
	nonMatchingEmoji := discord.Emoji{Name: "other", ID: discord.EmojiID(999)}

	assert.True(t, l.shouldHandleEvent(&gateway.MessageReactionAddEvent{
		Emoji: matchingEmoji,
	}), "matching custom emoji should be accepted for ReactionAdd")

	assert.False(t, l.shouldHandleEvent(&gateway.MessageReactionAddEvent{
		Emoji: nonMatchingEmoji,
	}), "non-matching custom emoji should be rejected for ReactionAdd")

	assert.True(t, l.shouldHandleEvent(&gateway.MessageReactionRemoveEvent{
		Emoji: matchingEmoji,
	}), "matching custom emoji should be accepted for ReactionRemove")

	assert.False(t, l.shouldHandleEvent(&gateway.MessageReactionRemoveEvent{
		Emoji: nonMatchingEmoji,
	}), "non-matching custom emoji should be rejected for ReactionRemove")
}

func TestShouldHandleEvent_MessageReactionFilter_NonReactionEventFallsThrough(t *testing.T) {
	// When a reaction filter is set but a non-reaction event arrives,
	// shouldHandleEvent should fall through to the default event-type switch,
	// which handles non-reaction events independently of the filter.
	l := newTestEventListener(&model.EventListenerFilter{
		MessageReaction: &model.EventListenerFilterMessageReaction{
			Emoji: "👍",
		},
	})

	// A non-bot message create falls through to the switch and returns true.
	assert.True(t, l.shouldHandleEvent(&gateway.MessageCreateEvent{
		Message: discord.Message{Author: discord.User{Bot: false}},
	}), "non-bot message create event should still be handled even when reaction filter is set")

	// A bot message create falls through to the switch and returns false.
	assert.False(t, l.shouldHandleEvent(&gateway.MessageCreateEvent{
		Message: discord.Message{Author: discord.User{Bot: true}},
	}), "bot message create event should still be rejected even when reaction filter is set")
}
