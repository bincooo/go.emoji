// Package emoji is designed to recognize and parse
// every individual Unicode Emoji characters from a string.
//
// Unicode Emoji Documents: https://www.unicode.org/Public/emoji/
package emoji

import (
	"bytes"

	"github.com/bincooo/go.emoji/internal/official"
)

func HasEmoji(s string) bool {
	for i := range s {
		match, _ := official.AllSequences.HasEmojiPrefix(s[i:])
		if match {
			return true
		}
	}
	return false
}

func FilterEmoji(s string) string {
	return replaceAllEmojiFunc(s, func(_ int, emoji string) string {
		return ""
	})
}

func ReplaceEmoji(s string, f func(index int, emoji string) string) string {
	return replaceAllEmojiFunc(s, f)
}

// ReplaceAllEmojiFunc searches string and find all emojis.
func replaceAllEmojiFunc(s string, f func(index int, emoji string) string) string {
	buff := bytes.Buffer{}
	nextIndex := 0

	for i, r := range s {
		if i < nextIndex {
			continue
		}

		match, length := official.AllSequences.HasEmojiPrefix(s[i:])
		if !match {
			buff.WriteRune(r)
			continue
		}

		nextIndex = i + length
		if f != nil {
			buff.WriteString(f(i, s[i:i+length]))
		}
	}

	return buff.String()
}

// IterateChars iterates a string, and extract every characters.
func IterateChars(s string) CharIterator {
	it := &charIteratorImpl{
		orig: []rune(s),
	}
	return it
}

// CharIterator is used in RangeChars
type CharIterator interface {
	Current() string
	CurrentIsEmoji() bool
	Next() bool
}

type charIteratorImpl struct {
	orig []rune

	curr struct {
		index int
		r     string
		emoji bool
	}
}

func (it *charIteratorImpl) Current() string {
	return it.curr.r
}

func (it *charIteratorImpl) Next() bool {
	if it.curr.index >= len(it.orig) {
		return false
	}

	match, length := official.AllSequences.HasEmojiPrefixRunes(it.orig[it.curr.index:])
	if !match {
		it.curr.r = string(it.orig[it.curr.index])
		it.curr.emoji = false
		it.curr.index++
		return true
	}

	it.curr.r = string(it.orig[it.curr.index : it.curr.index+length])
	it.curr.index += length
	it.curr.emoji = true
	return true
}

func (it *charIteratorImpl) CurrentIsEmoji() bool {
	return it.curr.emoji
}
