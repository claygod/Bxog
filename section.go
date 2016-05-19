package bxog

// Section
// Each route includes a section
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

type Section struct {
	id       string
	type_sec int // 0 - TYPE_STAT, 1 - TYPE_ARG
}

func newSection(sec string, type_s int) *Section {
	return &Section{id: sec, type_sec: type_s}
}
