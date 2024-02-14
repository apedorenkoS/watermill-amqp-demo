package event

import "strings"

const _concatenatesStrDelim = "||"

func concatenate(topic string, routingKey string) string {
	return topic + _concatenatesStrDelim + routingKey
}

func parseTopic(concatenatedStr string) string {
	entries := strings.Split(concatenatedStr, _concatenatesStrDelim)
	return entries[0]
}

func parseRoutingKey(concatenatedStr string) string {
	entries := strings.Split(concatenatedStr, _concatenatesStrDelim)
	if len(entries) > 1 {
		return entries[1]
	}
	return entries[0]
}
