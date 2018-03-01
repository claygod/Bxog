// Copyright © 2016-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

package bxog

// Section

// Each route includes a section
type section struct {
	id      string
	typeSec int // 0 - TYPE_STAT, 1 - TYPE_ARG
}

func newSection(sec string, tps int) *section {
	return &section{id: sec, typeSec: tps}
}
