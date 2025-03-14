package blindindexstore

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"unicode"
)

type TransformerInterface interface {
	Transform(string) string
}

// Example transformer not doing anything (do not use in production)
type NoChangeTransformer struct{}

var _ TransformerInterface = new(NoChangeTransformer)

func (t *NoChangeTransformer) Transform(v string) string {
	return v
}

// Example transformer using ROT13 (do not use in production)
type Rot13Transformer struct{}

var _ TransformerInterface = new(Rot13Transformer)

func (t *Rot13Transformer) Transform(v string) string {
	return rot13(v)
}

// Example transformer using SHA256
type Sha256Transformer struct{}

var _ TransformerInterface = new(Sha256Transformer)

func (t *Sha256Transformer) Transform(v string) string {
	return sha256Transform(v)
}

// Example custom transformer (do not use in production)
type UniTransformer struct{}

var _ TransformerInterface = new(UniTransformer)

func (t *UniTransformer) Transform(v string) string {
	return unicodeToCharNumber(v)
}

func rot13(s string) string {
	var result string
	for _, r := range s {
		if unicode.IsLetter(r) {
			r = unicode.ToLower(r)
			offset := int(r) - 'a'
			newOffset := (offset + 13) % 26
			newChar := rune('a' + newOffset)
			if unicode.IsUpper(r) {
				newChar = unicode.ToUpper(newChar)
			}
			result += string(newChar)
		} else {
			result += string(r)
		}
	}
	return result
}

func sha256Transform(inputString string) string {
	hash := sha256.Sum256([]byte(inputString))
	return hex.EncodeToString(hash[:])
}

func unicodeToCharNumber(unicodeString string) string {
	var charNumbers []string
	for _, r := range unicodeString {
		charNumbers = append(charNumbers, strconv.Itoa(int(r)))
	}
	return strings.Join(charNumbers, ";")
}
