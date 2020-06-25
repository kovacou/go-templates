// Copyright © 2020 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package templates

// Copyright © 2020 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

import (
	"encoding/json"
	"html/template"
	"strings"

	"github.com/kovacou/go-convert"
	"github.com/kovacou/go-types"
)

var (
	FuncMap = types.Map{
		// Treat as HTML
		"html": func(v string) template.HTML {
			return template.HTML(v)
		},

		// Treat as JS
		"js": func(v interface{}) template.JS {
			encoded, _ := json.Marshal(v)
			return template.JS(encoded)
		},

		// Strings
		"ToLower": strings.ToLower,
		"ToUpper": strings.ToUpper,

		// Casts & Conversions
		"Int":    convert.Int,
		"String": convert.String,
	}
)
