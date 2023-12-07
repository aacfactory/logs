package logs

var (
	lb           = []byte{'{'}
	rb           = []byte{'}'}
	lqb          = []byte{'['}
	rqb          = []byte{']'}
	dqm          = []byte{'"'}
	comma        = []byte{','}
	colon        = []byte{':'}
	space        = []byte(" ")
	lab          = []byte(">")
	equal        = []byte("=")
	newline      = []byte("\n")
	levelIdent   = []byte("level")
	occurIdent   = []byte("occur")
	messageIdent = []byte("message")
	fieldsIdent  = []byte("fields")
	keyIdent     = []byte("key")
	valueIdent   = []byte("value")
	callerIdent  = []byte("caller")
	fnIdent      = []byte("fn")
	fileIdent    = []byte("file")
	lineIdent    = []byte("line")
	causeIdent   = []byte("cause")
)
