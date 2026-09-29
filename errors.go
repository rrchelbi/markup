package markup

import "errors"

// Errors returned while rendering. Use errors.Is to test for them.
var (
	// ErrInvalidName reports a tag or attribute name with unsafe characters.
	ErrInvalidName = errors.New("markup: invalid tag or attribute name")
	// ErrUnsafeAttr reports an attribute that NewAttr refuses to render,
	// such as an event handler.
	ErrUnsafeAttr = errors.New("markup: unsafe attribute")
	// ErrVoidChildren reports children on a void element such as <br>.
	ErrVoidChildren = errors.New("markup: void element cannot have children")
	// ErrUnsafeRawText reports script or style content that could break
	// out of its element.
	ErrUnsafeRawText = errors.New("markup: unsafe content in raw text element")
)
