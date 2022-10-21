// Code generated by goyacc -v y.output -o parser.go -p mtail parser.y. DO NOT EDIT.

//line parser.y:5
/* #nosec G104 generated code, errors reported do not make sense */
package parser

import __yyfmt__ "fmt"

//line parser.y:6

import (
	"time"

	"github.com/golang/glog"
	"github.com/google/mtail/internal/metrics"
	"github.com/google/mtail/internal/runtime/compiler/ast"
	"github.com/google/mtail/internal/runtime/compiler/position"
)

//line parser.y:19
type mtailSymType struct {
	yys      int
	intVal   int64
	floatVal float64
	floats   []float64
	op       int
	text     string
	texts    []string
	flag     bool
	n        ast.Node
	kind     metrics.Kind
	duration time.Duration
}

const INVALID = 57346
const COUNTER = 57347
const GAUGE = 57348
const TIMER = 57349
const TEXT = 57350
const HISTOGRAM = 57351
const AFTER = 57352
const AS = 57353
const BY = 57354
const CONST = 57355
const HIDDEN = 57356
const DEF = 57357
const DEL = 57358
const NEXT = 57359
const OTHERWISE = 57360
const ELSE = 57361
const STOP = 57362
const BUCKETS = 57363
const LIMIT = 57364
const BUILTIN = 57365
const REGEX = 57366
const STRING = 57367
const CAPREF = 57368
const CAPREF_NAMED = 57369
const ID = 57370
const DECO = 57371
const INTLITERAL = 57372
const FLOATLITERAL = 57373
const DURATIONLITERAL = 57374
const INC = 57375
const DEC = 57376
const DIV = 57377
const MOD = 57378
const MUL = 57379
const MINUS = 57380
const PLUS = 57381
const POW = 57382
const SHL = 57383
const SHR = 57384
const LT = 57385
const GT = 57386
const LE = 57387
const GE = 57388
const EQ = 57389
const NE = 57390
const BITAND = 57391
const XOR = 57392
const BITOR = 57393
const NOT = 57394
const AND = 57395
const OR = 57396
const ADD_ASSIGN = 57397
const ASSIGN = 57398
const MATCH = 57399
const NOT_MATCH = 57400
const LCURLY = 57401
const RCURLY = 57402
const LPAREN = 57403
const RPAREN = 57404
const LSQUARE = 57405
const RSQUARE = 57406
const COMMA = 57407
const NL = 57408

var mtailToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"INVALID",
	"COUNTER",
	"GAUGE",
	"TIMER",
	"TEXT",
	"HISTOGRAM",
	"AFTER",
	"AS",
	"BY",
	"CONST",
	"HIDDEN",
	"DEF",
	"DEL",
	"NEXT",
	"OTHERWISE",
	"ELSE",
	"STOP",
	"BUCKETS",
	"LIMIT",
	"BUILTIN",
	"REGEX",
	"STRING",
	"CAPREF",
	"CAPREF_NAMED",
	"ID",
	"DECO",
	"INTLITERAL",
	"FLOATLITERAL",
	"DURATIONLITERAL",
	"INC",
	"DEC",
	"DIV",
	"MOD",
	"MUL",
	"MINUS",
	"PLUS",
	"POW",
	"SHL",
	"SHR",
	"LT",
	"GT",
	"LE",
	"GE",
	"EQ",
	"NE",
	"BITAND",
	"XOR",
	"BITOR",
	"NOT",
	"AND",
	"OR",
	"ADD_ASSIGN",
	"ASSIGN",
	"MATCH",
	"NOT_MATCH",
	"LCURLY",
	"RCURLY",
	"LPAREN",
	"RPAREN",
	"LSQUARE",
	"RSQUARE",
	"COMMA",
	"NL",
}

var mtailStatenames = [...]string{}

const mtailEofCode = 1
const mtailErrCode = 2
const mtailInitialStackSize = 16

//line parser.y:733

//  tokenpos returns the position of the current token.
func tokenpos(mtaillex mtailLexer) position.Position {
	return mtaillex.(*parser).t.Pos
}

// markedpos returns the position recorded from the most recent mark_pos
// production.
func markedpos(mtaillex mtailLexer) position.Position {
	return mtaillex.(*parser).pos
}

// positionFromMark returns a position spanning from the last mark to the current position.
func positionFromMark(mtaillex mtailLexer) position.Position {
	tp := tokenpos(mtaillex)
	mp := markedpos(mtaillex)
	return *position.Merge(&mp, &tp)
}

//line yacctab:1
var mtailExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 2,
	1, 1,
	5, 93,
	6, 93,
	7, 93,
	8, 93,
	9, 93,
	-2, 124,
	-1, 22,
	66, 24,
	-2, 68,
	-1, 106,
	5, 93,
	6, 93,
	7, 93,
	8, 93,
	9, 93,
	-2, 124,
}

const mtailPrivate = 57344

const mtailLast = 249

var mtailAct = [...]int{
	171, 88, 126, 28, 15, 91, 42, 44, 27, 30,
	103, 127, 41, 24, 20, 86, 40, 167, 22, 128,
	163, 182, 19, 26, 45, 29, 104, 25, 36, 34,
	35, 43, 54, 38, 39, 87, 181, 85, 89, 125,
	108, 46, 36, 34, 35, 43, 47, 38, 39, 90,
	162, 163, 62, 63, 2, 31, 49, 87, 76, 77,
	68, 130, 74, 73, 37, 138, 62, 63, 50, 112,
	93, 94, 117, 97, 96, 118, 168, 50, 37, 119,
	120, 70, 72, 71, 121, 122, 123, 66, 67, 124,
	107, 129, 185, 184, 111, 79, 80, 81, 82, 83,
	84, 169, 106, 131, 179, 135, 132, 175, 15, 133,
	129, 43, 27, 110, 100, 101, 99, 134, 20, 102,
	174, 135, 22, 173, 87, 129, 19, 160, 87, 151,
	156, 142, 155, 157, 158, 87, 87, 87, 164, 166,
	165, 161, 153, 159, 140, 154, 152, 136, 139, 178,
	177, 13, 141, 116, 66, 67, 115, 49, 105, 109,
	11, 23, 1, 176, 10, 129, 180, 12, 64, 145,
	13, 65, 36, 34, 35, 43, 75, 38, 39, 11,
	23, 98, 183, 10, 95, 69, 12, 92, 61, 78,
	18, 36, 34, 35, 43, 170, 38, 39, 143, 31,
	56, 57, 58, 59, 60, 172, 144, 137, 37, 36,
	34, 35, 43, 16, 38, 39, 146, 55, 31, 33,
	51, 53, 114, 48, 9, 8, 7, 37, 49, 113,
	6, 32, 16, 21, 52, 17, 31, 148, 147, 5,
	50, 14, 4, 3, 0, 37, 0, 149, 150,
}

var mtailPact = [...]int{
	-1000, -1000, 166, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 83, -1000, -1000, -13, 205, -1000, -34, 195, 13,
	13, -1000, 54, -1000, 21, 32, -1000, 7, 1, -1000,
	52, 184, -25, -1000, -1000, -1000, -1000, 184, -1000, -1000,
	29, -1000, 35, -1000, 79, -40, 139, -1000, -13, -21,
	-1000, 85, -13, 17, -1000, 128, -1000, -1000, -1000, -1000,
	-1000, -40, -1000, -1000, -40, -1000, -1000, -1000, -40, -40,
	-1000, -1000, -1000, -40, -40, -40, -1000, -1000, -40, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 54, -1000, 134, 184,
	-1, -1000, -40, -1000, -1000, -40, -1000, -1000, -40, -1000,
	-1000, -1000, -1000, -1000, -1000, -13, 147, -1000, 3, 120,
	-13, -1000, 121, 226, -1000, -1000, -1000, 184, 184, 83,
	184, 184, 184, 17, 184, -14, -1000, 13, -1000, 33,
	-1000, 184, 184, 184, 21, 42, -1000, -1000, -1000, -45,
	41, -1000, 69, -1000, -1000, -1000, -1000, 95, 82, 119,
	74, 13, 32, -1000, -1000, -1000, 52, 13, 13, -1000,
	-1000, 29, -1000, 184, 35, 79, -1000, -1000, -1000, -1000,
	-29, -1000, -1000, -1000, -1000, -1000, -44, -1000, -1000, -1000,
	-1000, 95, 62, -1000, -1000, -1000,
}

var mtailPgo = [...]int{
	0, 54, 243, 39, 41, 242, 241, 239, 235, 3,
	7, 6, 15, 5, 233, 9, 16, 27, 11, 231,
	12, 13, 19, 230, 229, 226, 225, 25, 23, 224,
	222, 219, 2, 217, 216, 206, 205, 0, 198, 195,
	190, 189, 187, 185, 168, 184, 181, 176, 171, 169,
	163, 162, 10, 1, 159,
}

var mtailR1 = [...]int{
	0, 51, 1, 1, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 5, 5, 5, 6, 6, 6,
	7, 7, 4, 8, 8, 14, 14, 18, 18, 18,
	18, 44, 44, 17, 17, 43, 43, 43, 15, 15,
	41, 41, 41, 41, 41, 41, 16, 16, 42, 42,
	11, 11, 45, 45, 28, 28, 47, 47, 22, 21,
	21, 21, 10, 10, 46, 46, 46, 46, 13, 13,
	12, 12, 48, 48, 9, 9, 9, 9, 9, 9,
	9, 9, 19, 19, 20, 31, 31, 3, 3, 32,
	32, 27, 23, 40, 40, 24, 24, 24, 24, 24,
	30, 30, 33, 33, 33, 33, 33, 38, 39, 39,
	37, 35, 34, 49, 50, 50, 50, 50, 25, 26,
	29, 29, 36, 36, 53, 54, 52, 52,
}

var mtailR2 = [...]int{
	0, 1, 0, 2, 1, 1, 1, 1, 1, 1,
	1, 4, 1, 1, 4, 2, 3, 1, 4, 1,
	1, 2, 3, 1, 1, 4, 4, 1, 1, 4,
	4, 1, 1, 1, 4, 1, 1, 1, 1, 4,
	1, 1, 1, 1, 1, 1, 1, 4, 1, 1,
	1, 4, 1, 1, 4, 4, 1, 1, 1, 1,
	4, 4, 1, 4, 1, 1, 1, 1, 1, 2,
	1, 2, 1, 1, 1, 1, 1, 1, 1, 3,
	1, 1, 1, 4, 1, 4, 5, 1, 3, 1,
	1, 5, 3, 0, 1, 2, 2, 2, 2, 1,
	1, 1, 1, 1, 1, 1, 1, 2, 1, 3,
	1, 2, 2, 2, 1, 1, 3, 3, 4, 3,
	5, 3, 1, 1, 0, 0, 0, 1,
}

var mtailChk = [...]int{
	-1000, -51, -1, -2, -5, -7, -23, -25, -26, -29,
	17, 13, 20, 4, -6, -53, 66, -8, -40, -22,
	-18, -14, -12, 14, -21, -17, -28, -13, -9, -27,
	-15, 52, -19, -31, 26, 27, 25, 61, 30, 31,
	-16, -20, -11, 28, -10, -20, -4, 59, 18, 23,
	35, 15, 29, 16, 66, -33, 5, 6, 7, 8,
	9, -44, 53, 54, -44, -48, 33, 34, 39, -43,
	49, 51, 50, 56, 55, -47, 57, 58, -41, 43,
	44, 45, 46, 47, 48, -13, -12, -9, -53, 63,
	-18, -13, -42, 41, 42, -45, 39, 38, -46, 37,
	35, 36, 40, -52, 66, 19, -1, -4, 61, -54,
	28, -4, -12, -24, -30, 28, 25, -52, -52, -52,
	-52, -52, -52, -52, -52, -3, -32, -18, -22, -53,
	62, -52, -52, -52, -21, -53, -4, 60, 62, -3,
	24, -4, 10, -38, -35, -49, -34, 12, 11, 21,
	22, -18, -17, -28, -27, -20, -15, -18, -18, -22,
	-9, -16, 64, 65, -11, -10, -13, 62, 35, 32,
	-39, -37, -36, 28, 25, 25, -50, 31, 30, 30,
	-32, 65, 65, -37, 31, 30,
}

var mtailDef = [...]int{
	2, -2, -2, 3, 4, 5, 6, 7, 8, 9,
	10, 0, 12, 13, 0, 0, 20, 0, 0, 17,
	19, 23, -2, 94, 58, 27, 28, 62, 70, 59,
	33, 124, 74, 75, 76, 77, 78, 124, 80, 81,
	38, 82, 46, 84, 50, 126, 15, 2, 0, 0,
	125, 0, 0, 124, 21, 0, 102, 103, 104, 105,
	106, 126, 31, 32, 126, 71, 72, 73, 126, 126,
	35, 36, 37, 126, 126, 126, 56, 57, 126, 40,
	41, 42, 43, 44, 45, 69, 68, 70, 0, 124,
	0, 62, 126, 48, 49, 126, 52, 53, 126, 64,
	65, 66, 67, 124, 127, 0, -2, 16, 124, 0,
	0, 119, 121, 92, 99, 100, 101, 124, 124, 124,
	124, 124, 124, 124, 124, 0, 87, 89, 90, 0,
	79, 124, 124, 124, 11, 0, 14, 22, 85, 0,
	0, 118, 0, 95, 96, 97, 98, 0, 0, 0,
	0, 18, 29, 30, 60, 61, 34, 25, 26, 54,
	55, 39, 83, 124, 47, 51, 63, 86, 91, 120,
	107, 108, 110, 122, 123, 111, 113, 114, 115, 112,
	88, 0, 0, 109, 116, 117,
}

var mtailTok1 = [...]int{
	1,
}

var mtailTok2 = [...]int{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62, 63, 64, 65, 66,
}

var mtailTok3 = [...]int{
	0,
}

var mtailErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{
	{109, 4, "unexpected end of file, expecting '/' to end regex"},
	{15, 1, "unexpected end of file, expecting '}' to end block"},
	{15, 1, "unexpected end of file, expecting '}' to end block"},
	{15, 1, "unexpected end of file, expecting '}' to end block"},
	{14, 63, "unexpected indexing of an expression"},
	{14, 66, "statement with no effect, missing an assignment, `+' concatenation, or `{}' block?"},
}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	mtailDebug        = 0
	mtailErrorVerbose = false
)

type mtailLexer interface {
	Lex(lval *mtailSymType) int
	Error(s string)
}

type mtailParser interface {
	Parse(mtailLexer) int
	Lookahead() int
}

type mtailParserImpl struct {
	lval  mtailSymType
	stack [mtailInitialStackSize]mtailSymType
	char  int
}

func (p *mtailParserImpl) Lookahead() int {
	return p.char
}

func mtailNewParser() mtailParser {
	return &mtailParserImpl{}
}

const mtailFlag = -1000

func mtailTokname(c int) string {
	if c >= 1 && c-1 < len(mtailToknames) {
		if mtailToknames[c-1] != "" {
			return mtailToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func mtailStatname(s int) string {
	if s >= 0 && s < len(mtailStatenames) {
		if mtailStatenames[s] != "" {
			return mtailStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func mtailErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !mtailErrorVerbose {
		return "syntax error"
	}

	for _, e := range mtailErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + mtailTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := mtailPact[state]
	for tok := TOKSTART; tok-1 < len(mtailToknames); tok++ {
		if n := base + tok; n >= 0 && n < mtailLast && mtailChk[mtailAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if mtailDef[state] == -2 {
		i := 0
		for mtailExca[i] != -1 || mtailExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; mtailExca[i] >= 0; i += 2 {
			tok := mtailExca[i]
			if tok < TOKSTART || mtailExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if mtailExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += mtailTokname(tok)
	}
	return res
}

func mtaillex1(lex mtailLexer, lval *mtailSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = mtailTok1[0]
		goto out
	}
	if char < len(mtailTok1) {
		token = mtailTok1[char]
		goto out
	}
	if char >= mtailPrivate {
		if char < mtailPrivate+len(mtailTok2) {
			token = mtailTok2[char-mtailPrivate]
			goto out
		}
	}
	for i := 0; i < len(mtailTok3); i += 2 {
		token = mtailTok3[i+0]
		if token == char {
			token = mtailTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = mtailTok2[1] /* unknown char */
	}
	if mtailDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", mtailTokname(token), uint(char))
	}
	return char, token
}

func mtailParse(mtaillex mtailLexer) int {
	return mtailNewParser().Parse(mtaillex)
}

func (mtailrcvr *mtailParserImpl) Parse(mtaillex mtailLexer) int {
	var mtailn int
	var mtailVAL mtailSymType
	var mtailDollar []mtailSymType
	_ = mtailDollar // silence set and not used
	mtailS := mtailrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	mtailstate := 0
	mtailrcvr.char = -1
	mtailtoken := -1 // mtailrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		mtailstate = -1
		mtailrcvr.char = -1
		mtailtoken = -1
	}()
	mtailp := -1
	goto mtailstack

ret0:
	return 0

ret1:
	return 1

mtailstack:
	/* put a state and value onto the stack */
	if mtailDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", mtailTokname(mtailtoken), mtailStatname(mtailstate))
	}

	mtailp++
	if mtailp >= len(mtailS) {
		nyys := make([]mtailSymType, len(mtailS)*2)
		copy(nyys, mtailS)
		mtailS = nyys
	}
	mtailS[mtailp] = mtailVAL
	mtailS[mtailp].yys = mtailstate

mtailnewstate:
	mtailn = mtailPact[mtailstate]
	if mtailn <= mtailFlag {
		goto mtaildefault /* simple state */
	}
	if mtailrcvr.char < 0 {
		mtailrcvr.char, mtailtoken = mtaillex1(mtaillex, &mtailrcvr.lval)
	}
	mtailn += mtailtoken
	if mtailn < 0 || mtailn >= mtailLast {
		goto mtaildefault
	}
	mtailn = mtailAct[mtailn]
	if mtailChk[mtailn] == mtailtoken { /* valid shift */
		mtailrcvr.char = -1
		mtailtoken = -1
		mtailVAL = mtailrcvr.lval
		mtailstate = mtailn
		if Errflag > 0 {
			Errflag--
		}
		goto mtailstack
	}

mtaildefault:
	/* default state action */
	mtailn = mtailDef[mtailstate]
	if mtailn == -2 {
		if mtailrcvr.char < 0 {
			mtailrcvr.char, mtailtoken = mtaillex1(mtaillex, &mtailrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if mtailExca[xi+0] == -1 && mtailExca[xi+1] == mtailstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			mtailn = mtailExca[xi+0]
			if mtailn < 0 || mtailn == mtailtoken {
				break
			}
		}
		mtailn = mtailExca[xi+1]
		if mtailn < 0 {
			goto ret0
		}
	}
	if mtailn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			mtaillex.Error(mtailErrorMessage(mtailstate, mtailtoken))
			Nerrs++
			if mtailDebug >= 1 {
				__yyfmt__.Printf("%s", mtailStatname(mtailstate))
				__yyfmt__.Printf(" saw %s\n", mtailTokname(mtailtoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for mtailp >= 0 {
				mtailn = mtailPact[mtailS[mtailp].yys] + mtailErrCode
				if mtailn >= 0 && mtailn < mtailLast {
					mtailstate = mtailAct[mtailn] /* simulate a shift of "error" */
					if mtailChk[mtailstate] == mtailErrCode {
						goto mtailstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if mtailDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", mtailS[mtailp].yys)
				}
				mtailp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if mtailDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", mtailTokname(mtailtoken))
			}
			if mtailtoken == mtailEofCode {
				goto ret1
			}
			mtailrcvr.char = -1
			mtailtoken = -1
			goto mtailnewstate /* try again in the same state */
		}
	}

	/* reduction by production mtailn */
	if mtailDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", mtailn, mtailStatname(mtailstate))
	}

	mtailnt := mtailn
	mtailpt := mtailp
	_ = mtailpt // guard against "declared and not used"

	mtailp -= mtailR2[mtailn]
	// mtailp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if mtailp+1 >= len(mtailS) {
		nyys := make([]mtailSymType, len(mtailS)*2)
		copy(nyys, mtailS)
		mtailS = nyys
	}
	mtailVAL = mtailS[mtailp+1]

	/* consult goto table to find next state */
	mtailn = mtailR1[mtailn]
	mtailg := mtailPgo[mtailn]
	mtailj := mtailg + mtailS[mtailp].yys + 1

	if mtailj >= mtailLast {
		mtailstate = mtailAct[mtailg]
	} else {
		mtailstate = mtailAct[mtailj]
		if mtailChk[mtailstate] != -mtailn {
			mtailstate = mtailAct[mtailg]
		}
	}
	// dummy call; replaced with literal code
	switch mtailnt {

	case 1:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:93
		{
			mtaillex.(*parser).root = mtailDollar[1].n
		}
	case 2:
		mtailDollar = mtailS[mtailpt-0 : mtailpt+1]
//line parser.y:101
		{
			mtailVAL.n = &ast.StmtList{}
		}
	case 3:
		mtailDollar = mtailS[mtailpt-2 : mtailpt+1]
//line parser.y:105
		{
			mtailVAL.n = mtailDollar[1].n
			if mtailDollar[2].n != nil {
				mtailVAL.n.(*ast.StmtList).Children = append(mtailVAL.n.(*ast.StmtList).Children, mtailDollar[2].n)
			}
		}
	case 4:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:116
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 5:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:118
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 6:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:120
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 7:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:122
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 8:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:124
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 9:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:126
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 10:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:128
		{
			mtailVAL.n = &ast.NextStmt{tokenpos(mtaillex)}
		}
	case 11:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:132
		{
			mtailVAL.n = &ast.PatternFragment{ID: mtailDollar[2].n, Expr: mtailDollar[4].n}
		}
	case 12:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:136
		{
			mtailVAL.n = &ast.StopStmt{tokenpos(mtaillex)}
		}
	case 13:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:140
		{
			mtailVAL.n = &ast.Error{tokenpos(mtaillex), mtailDollar[1].text}
		}
	case 14:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:148
		{
			mtailVAL.n = &ast.CondStmt{mtailDollar[1].n, mtailDollar[2].n, mtailDollar[4].n, nil}
		}
	case 15:
		mtailDollar = mtailS[mtailpt-2 : mtailpt+1]
//line parser.y:152
		{
			if mtailDollar[1].n != nil {
				mtailVAL.n = &ast.CondStmt{mtailDollar[1].n, mtailDollar[2].n, nil, nil}
			} else {
				mtailVAL.n = mtailDollar[2].n
			}
		}
	case 16:
		mtailDollar = mtailS[mtailpt-3 : mtailpt+1]
//line parser.y:160
		{
			o := &ast.OtherwiseStmt{positionFromMark(mtaillex)}
			mtailVAL.n = &ast.CondStmt{o, mtailDollar[3].n, nil, nil}
		}
	case 17:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:168
		{
			mtailVAL.n = &ast.UnaryExpr{P: tokenpos(mtaillex), Expr: mtailDollar[1].n, Op: MATCH}
		}
	case 18:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:172
		{
			mtailVAL.n = &ast.BinaryExpr{
				LHS: &ast.UnaryExpr{P: tokenpos(mtaillex), Expr: mtailDollar[1].n, Op: MATCH},
				RHS: mtailDollar[4].n,
				Op:  mtailDollar[2].op,
			}
		}
	case 19:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:180
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 20:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:186
		{
			mtailVAL.n = nil
		}
	case 21:
		mtailDollar = mtailS[mtailpt-2 : mtailpt+1]
//line parser.y:188
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 22:
		mtailDollar = mtailS[mtailpt-3 : mtailpt+1]
//line parser.y:194
		{
			mtailVAL.n = mtailDollar[2].n
		}
	case 23:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:202
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 24:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:204
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 25:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:210
		{
			mtailVAL.n = &ast.BinaryExpr{LHS: mtailDollar[1].n, RHS: mtailDollar[4].n, Op: mtailDollar[2].op}
		}
	case 26:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:214
		{
			mtailVAL.n = &ast.BinaryExpr{LHS: mtailDollar[1].n, RHS: mtailDollar[4].n, Op: mtailDollar[2].op}
		}
	case 27:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:222
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 28:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:224
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 29:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:226
		{
			mtailVAL.n = &ast.BinaryExpr{LHS: mtailDollar[1].n, RHS: mtailDollar[4].n, Op: mtailDollar[2].op}
		}
	case 30:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:230
		{
			mtailVAL.n = &ast.BinaryExpr{LHS: mtailDollar[1].n, RHS: mtailDollar[4].n, Op: mtailDollar[2].op}
		}
	case 31:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:237
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 32:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:239
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 33:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:245
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 34:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:247
		{
			mtailVAL.n = &ast.BinaryExpr{LHS: mtailDollar[1].n, RHS: mtailDollar[4].n, Op: mtailDollar[2].op}
		}
	case 35:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:254
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 36:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:256
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 37:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:258
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 38:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:264
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 39:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:266
		{
			mtailVAL.n = &ast.BinaryExpr{LHS: mtailDollar[1].n, RHS: mtailDollar[4].n, Op: mtailDollar[2].op}
		}
	case 40:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:273
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 41:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:275
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 42:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:277
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 43:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:279
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 44:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:281
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 45:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:283
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 46:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:289
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 47:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:291
		{
			mtailVAL.n = &ast.BinaryExpr{LHS: mtailDollar[1].n, RHS: mtailDollar[4].n, Op: mtailDollar[2].op}
		}
	case 48:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:298
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 49:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:300
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 50:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:306
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 51:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:308
		{
			mtailVAL.n = &ast.BinaryExpr{LHS: mtailDollar[1].n, RHS: mtailDollar[4].n, Op: mtailDollar[2].op}
		}
	case 52:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:315
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 53:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:317
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 54:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:323
		{
			mtailVAL.n = &ast.BinaryExpr{LHS: mtailDollar[1].n, RHS: mtailDollar[4].n, Op: mtailDollar[2].op}
		}
	case 55:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:327
		{
			mtailVAL.n = &ast.BinaryExpr{LHS: mtailDollar[1].n, RHS: mtailDollar[4].n, Op: mtailDollar[2].op}
		}
	case 56:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:334
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 57:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:336
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 58:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:343
		{
			mtailVAL.n = &ast.PatternExpr{Expr: mtailDollar[1].n}
		}
	case 59:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:351
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 60:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:353
		{
			mtailVAL.n = &ast.BinaryExpr{LHS: mtailDollar[1].n, RHS: mtailDollar[4].n, Op: PLUS}
		}
	case 61:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:357
		{
			mtailVAL.n = &ast.BinaryExpr{LHS: mtailDollar[1].n, RHS: mtailDollar[4].n, Op: PLUS}
		}
	case 62:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:365
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 63:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:367
		{
			mtailVAL.n = &ast.BinaryExpr{LHS: mtailDollar[1].n, RHS: mtailDollar[4].n, Op: mtailDollar[2].op}
		}
	case 64:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:374
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 65:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:376
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 66:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:378
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 67:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:380
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 68:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:386
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 69:
		mtailDollar = mtailS[mtailpt-2 : mtailpt+1]
//line parser.y:388
		{
			mtailVAL.n = &ast.UnaryExpr{P: tokenpos(mtaillex), Expr: mtailDollar[2].n, Op: mtailDollar[1].op}
		}
	case 70:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:396
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 71:
		mtailDollar = mtailS[mtailpt-2 : mtailpt+1]
//line parser.y:398
		{
			mtailVAL.n = &ast.UnaryExpr{P: tokenpos(mtaillex), Expr: mtailDollar[1].n, Op: mtailDollar[2].op}
		}
	case 72:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:405
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 73:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:407
		{
			mtailVAL.op = mtailDollar[1].op
		}
	case 74:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:413
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 75:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:415
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 76:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:417
		{
			mtailVAL.n = &ast.CaprefTerm{tokenpos(mtaillex), mtailDollar[1].text, false, nil}
		}
	case 77:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:421
		{
			mtailVAL.n = &ast.CaprefTerm{tokenpos(mtaillex), mtailDollar[1].text, true, nil}
		}
	case 78:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:425
		{
			mtailVAL.n = &ast.StringLit{tokenpos(mtaillex), mtailDollar[1].text}
		}
	case 79:
		mtailDollar = mtailS[mtailpt-3 : mtailpt+1]
//line parser.y:429
		{
			mtailVAL.n = mtailDollar[2].n
		}
	case 80:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:433
		{
			mtailVAL.n = &ast.IntLit{tokenpos(mtaillex), mtailDollar[1].intVal}
		}
	case 81:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:437
		{
			mtailVAL.n = &ast.FloatLit{tokenpos(mtaillex), mtailDollar[1].floatVal}
		}
	case 82:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:445
		{
			// Build an empty IndexedExpr so that the recursive rule below doesn't need to handle the alternative.
			mtailVAL.n = &ast.IndexedExpr{LHS: mtailDollar[1].n, Index: &ast.ExprList{}}
		}
	case 83:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:450
		{
			mtailVAL.n = mtailDollar[1].n
			mtailVAL.n.(*ast.IndexedExpr).Index.(*ast.ExprList).Children = append(
				mtailVAL.n.(*ast.IndexedExpr).Index.(*ast.ExprList).Children,
				mtailDollar[3].n.(*ast.ExprList).Children...)
		}
	case 84:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:461
		{
			mtailVAL.n = &ast.IDTerm{tokenpos(mtaillex), mtailDollar[1].text, nil, false}
		}
	case 85:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:469
		{
			mtailVAL.n = &ast.BuiltinExpr{P: positionFromMark(mtaillex), Name: mtailDollar[2].text, Args: nil}
		}
	case 86:
		mtailDollar = mtailS[mtailpt-5 : mtailpt+1]
//line parser.y:473
		{
			mtailVAL.n = &ast.BuiltinExpr{P: positionFromMark(mtaillex), Name: mtailDollar[2].text, Args: mtailDollar[4].n}
		}
	case 87:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:482
		{
			mtailVAL.n = &ast.ExprList{}
			mtailVAL.n.(*ast.ExprList).Children = append(mtailVAL.n.(*ast.ExprList).Children, mtailDollar[1].n)
		}
	case 88:
		mtailDollar = mtailS[mtailpt-3 : mtailpt+1]
//line parser.y:487
		{
			mtailVAL.n = mtailDollar[1].n
			mtailVAL.n.(*ast.ExprList).Children = append(mtailVAL.n.(*ast.ExprList).Children, mtailDollar[3].n)
		}
	case 89:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:495
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 90:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:497
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 91:
		mtailDollar = mtailS[mtailpt-5 : mtailpt+1]
//line parser.y:503
		{
			mtailVAL.n = &ast.PatternLit{P: positionFromMark(mtaillex), Pattern: mtailDollar[4].text}
		}
	case 92:
		mtailDollar = mtailS[mtailpt-3 : mtailpt+1]
//line parser.y:511
		{
			mtailVAL.n = mtailDollar[3].n
			d := mtailVAL.n.(*ast.VarDecl)
			d.Kind = mtailDollar[2].kind
			d.Hidden = mtailDollar[1].flag
		}
	case 93:
		mtailDollar = mtailS[mtailpt-0 : mtailpt+1]
//line parser.y:522
		{
			mtailVAL.flag = false
		}
	case 94:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:526
		{
			mtailVAL.flag = true
		}
	case 95:
		mtailDollar = mtailS[mtailpt-2 : mtailpt+1]
//line parser.y:534
		{
			mtailVAL.n = mtailDollar[1].n
			mtailVAL.n.(*ast.VarDecl).Keys = mtailDollar[2].texts
		}
	case 96:
		mtailDollar = mtailS[mtailpt-2 : mtailpt+1]
//line parser.y:539
		{
			mtailVAL.n = mtailDollar[1].n
			mtailVAL.n.(*ast.VarDecl).ExportedName = mtailDollar[2].text
		}
	case 97:
		mtailDollar = mtailS[mtailpt-2 : mtailpt+1]
//line parser.y:544
		{
			mtailVAL.n = mtailDollar[1].n
			mtailVAL.n.(*ast.VarDecl).Buckets = mtailDollar[2].floats
		}
	case 98:
		mtailDollar = mtailS[mtailpt-2 : mtailpt+1]
//line parser.y:549
		{
			mtailVAL.n = mtailDollar[1].n
			mtailVAL.n.(*ast.VarDecl).Limit = mtailDollar[2].intVal
		}
	case 99:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:554
		{
			mtailVAL.n = mtailDollar[1].n
		}
	case 100:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:562
		{
			mtailVAL.n = &ast.VarDecl{P: tokenpos(mtaillex), Name: mtailDollar[1].text}
		}
	case 101:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:566
		{
			mtailVAL.n = &ast.VarDecl{P: tokenpos(mtaillex), Name: mtailDollar[1].text}
		}
	case 102:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:574
		{
			mtailVAL.kind = metrics.Counter
		}
	case 103:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:578
		{
			mtailVAL.kind = metrics.Gauge
		}
	case 104:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:582
		{
			mtailVAL.kind = metrics.Timer
		}
	case 105:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:586
		{
			mtailVAL.kind = metrics.Text
		}
	case 106:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:590
		{
			mtailVAL.kind = metrics.Histogram
		}
	case 107:
		mtailDollar = mtailS[mtailpt-2 : mtailpt+1]
//line parser.y:598
		{
			mtailVAL.texts = mtailDollar[2].texts
		}
	case 108:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:605
		{
			mtailVAL.texts = make([]string, 0)
			mtailVAL.texts = append(mtailVAL.texts, mtailDollar[1].text)
		}
	case 109:
		mtailDollar = mtailS[mtailpt-3 : mtailpt+1]
//line parser.y:610
		{
			mtailVAL.texts = mtailDollar[1].texts
			mtailVAL.texts = append(mtailVAL.texts, mtailDollar[3].text)
		}
	case 110:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:618
		{
			mtailVAL.text = mtailDollar[1].text
		}
	case 111:
		mtailDollar = mtailS[mtailpt-2 : mtailpt+1]
//line parser.y:624
		{
			mtailVAL.text = mtailDollar[2].text
		}
	case 112:
		mtailDollar = mtailS[mtailpt-2 : mtailpt+1]
//line parser.y:631
		{
			mtailVAL.intVal = mtailDollar[2].intVal
		}
	case 113:
		mtailDollar = mtailS[mtailpt-2 : mtailpt+1]
//line parser.y:639
		{
			mtailVAL.floats = mtailDollar[2].floats
		}
	case 114:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:645
		{
			mtailVAL.floats = make([]float64, 0)
			mtailVAL.floats = append(mtailVAL.floats, mtailDollar[1].floatVal)
		}
	case 115:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:650
		{
			mtailVAL.floats = make([]float64, 0)
			mtailVAL.floats = append(mtailVAL.floats, float64(mtailDollar[1].intVal))
		}
	case 116:
		mtailDollar = mtailS[mtailpt-3 : mtailpt+1]
//line parser.y:655
		{
			mtailVAL.floats = mtailDollar[1].floats
			mtailVAL.floats = append(mtailVAL.floats, mtailDollar[3].floatVal)
		}
	case 117:
		mtailDollar = mtailS[mtailpt-3 : mtailpt+1]
//line parser.y:660
		{
			mtailVAL.floats = mtailDollar[1].floats
			mtailVAL.floats = append(mtailVAL.floats, float64(mtailDollar[3].intVal))
		}
	case 118:
		mtailDollar = mtailS[mtailpt-4 : mtailpt+1]
//line parser.y:668
		{
			mtailVAL.n = &ast.DecoDecl{P: markedpos(mtaillex), Name: mtailDollar[3].text, Block: mtailDollar[4].n}
		}
	case 119:
		mtailDollar = mtailS[mtailpt-3 : mtailpt+1]
//line parser.y:676
		{
			mtailVAL.n = &ast.DecoStmt{markedpos(mtaillex), mtailDollar[2].text, mtailDollar[3].n, nil, nil}
		}
	case 120:
		mtailDollar = mtailS[mtailpt-5 : mtailpt+1]
//line parser.y:684
		{
			mtailVAL.n = &ast.DelStmt{P: positionFromMark(mtaillex), N: mtailDollar[3].n, Expiry: mtailDollar[5].duration}
		}
	case 121:
		mtailDollar = mtailS[mtailpt-3 : mtailpt+1]
//line parser.y:688
		{
			mtailVAL.n = &ast.DelStmt{P: positionFromMark(mtaillex), N: mtailDollar[3].n}
		}
	case 122:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:695
		{
			mtailVAL.text = mtailDollar[1].text
		}
	case 123:
		mtailDollar = mtailS[mtailpt-1 : mtailpt+1]
//line parser.y:699
		{
			mtailVAL.text = mtailDollar[1].text
		}
	case 124:
		mtailDollar = mtailS[mtailpt-0 : mtailpt+1]
//line parser.y:709
		{
			glog.V(2).Infof("position marked at %v", tokenpos(mtaillex))
			mtaillex.(*parser).pos = tokenpos(mtaillex)
		}
	case 125:
		mtailDollar = mtailS[mtailpt-0 : mtailpt+1]
//line parser.y:719
		{
			mtaillex.(*parser).inRegex()
		}
	}
	goto mtailstack /* stack new state and value */
}