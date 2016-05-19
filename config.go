package bxog

// Gonfig
// Multiplexer configuration
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

type type_hash uint32

const (
	TYPE_STAT = iota
	TYPE_ARG
)

// Editable parameters
const (
	// The method used by default when you add a Route.
	// Example: GET, POST etc.
	HTTP_METHOD_DEFAULT = "GET"

	// Maximum number of sections in the URL
	// Example: /abc/:par - 2 sections, /a/:b/:c - 3 sections
	HTTP_SECTION_COUNT = 32

	// The maximum length of URL (characters)
	HTTP_PATTERN_COUNT = 512

	// The maximum wait time during a read operation
	READ_TIME_OUT = 100

	// The maximum wait time during a write operation
	WRITE_TIME_OUT = 100

	// Address directory for files on the website URL
	FILE_PREF = "/file/"

	// Address directory for files on your computer
	FILE_PATH = "./file/"
)

// Non-editable parameters
const (
	DELIMITER_STRING           = "/"
	DELIMITER_BYTE   byte      = 47
	DELIMITER_UINT   type_hash = 47
	SLASH_HASH       type_hash = 1
)
