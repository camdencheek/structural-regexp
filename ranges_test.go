package regexp

import "testing"

const source = `
const fetch = require('node-fetch')
globalThis.fetch = fetch
const { Request, Response, Headers } = fetch
Object.assign(globalThis, { Request, Response, Headers })
`

func TestRanges(t *testing.T) {
	res := ParseJavascript([]byte(source))
	t.Fatalf("%+v\n", res)
}
