// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type File struct {
	InputStream antlr.CharStream
	Parser      *dmlParser
}

func ParseStream(input antlr.CharStream, args ...any) *File {
	lexer := NewdmlLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := NewdmlParser(stream)

	ignoreErrorListeners := false
	for _, argv := range args {
		switch v := argv.(type) {
		case antlr.ErrorListener:
			if !ignoreErrorListeners {
				p.RemoveErrorListeners()
				lexer.RemoveErrorListeners()
				ignoreErrorListeners = true
			}
			p.AddErrorListener(v)
			lexer.AddErrorListener(v)
		case antlr.ParseTreeListener:
			p.AddParseListener(v)
		}
	}

	return &File{
		InputStream: input,
		Parser:      p,
	}
}

func ParseFile(filepath string, args ...any) (*File, error) {
	input, err := antlr.NewFileStream(filepath)
	if err != nil {
		return nil, err
	}

	return ParseStream(input, args...), nil
}
