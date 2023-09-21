package tg

import (
	//"reflect"
)

// The argument for handling in channenl behaviours.
type ChannelContext struct {
}
type CC = ChannelContext
type ChannelAction struct {
	Act (*ChannelContext)
}
