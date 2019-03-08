package main

import (
	"errors"
)

// ErrNoAvatar is the error that is returned when the instance of Avatar
// is unable to provide an Avatar URL.
var ErrNoAvatarURL = errors.New("chat: unabe to get an avatar URL.")

// Avatar represents types capable of representing user profile pictures
type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}
