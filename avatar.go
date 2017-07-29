package main

import "errors"

var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません")

type Avatar interface {
	// 指定されたクライアントアバターのURL
	// 問題が発生した場合はエラーを返す
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct {

}

var UserAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(c * client) (string, error){
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}