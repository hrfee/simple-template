export enum BlockType {
	PlainText,                   // Plain text
	LogicOpen,                   // { (or {{)
	LogicClose,                  // } (or }})
	Word,                        // A single word (only parsed like this when within braces)
	String,                      // Plain text within quotes (only parsed like this when within braces). Quotes can be ', ", or `.
	EOF,                         // End of file.
}

const seekBufferSize = 3;

function truthy(val: any): boolean {
	switch (typeof val) {
	case "string":
		return val != "";
	case "boolean":
		return val;
	case "number":
		return val != 0;
	}
	return false
}

class block {
	parent: string;
	a: number; // start/end indices (inclusive)
    b: number;
	Type: BlockType;

    String(): string {
        if (this.Type == BlockType.String) {
            return this.parent.slice(this.a+1,  this.b);
        } else if (this.Type == BlockType.EOF) {
            return ""
        }
        return this.parent.slice(this.a, this.b+1);
    }

    Describe(): string {
        return `${block.typeToString(this.Type)}[${this.a}:${this.b}]="${this.String()}"`;
    }

    static typeToString(b: BlockType): string {
        switch (b) {
        case BlockType.PlainText:
            return "PlainText"
        case BlockType.LogicOpen:
            return "LogicOpen"
        case BlockType.LogicClose:
            return "LogicClose"
        case BlockType.Word:
            return "Word"
        case BlockType.String:
            return "String"
        case BlockType.EOF:
            return "EOF"
        }
        return "?"
    }

    expected(...expected: BlockType[]): Error|null {
        return new ExpectedTypeError(this.b, this.Type, expected)
    }

    expectedWord(expected: string): Error|null {
        return new ExpectedError(this.b, this.String(), expected)
    }

    copy(): block { return new block(this.parent, this.a, this.b, this.Type); }

    constructor(parent?: string, a?: number, b?: number, Type?: BlockType) {
        this.parent = parent;
        this.a = a;
        this.b = b;
        this.Type = Type;
    }
}

// DoubleBraceError indicates double braces were used instead of single braces. This being returned does not indicate that templating failed.
export class DoubleBraceError extends Error {
    constructor(pos: number) {
        super(`double braces ("{{"/"}}") near char ${pos}, use single braces only.`);
        this.name = "DoubleBraceError";
        Object.setPrototypeOf(this, DoubleBraceError.prototype);
    }
}

// SingleEqualsError indicates a single equals sign ("=") was used in a comparison rather than two ("=="). This being returned does not indicate that templating failed.
export class SingleEqualsError extends Error {
    constructor(pos: number) {
        super(`single equals ("=") used in if block near char ${pos}, use double equals ("==").`);
        this.name = "SingleEqualsError";
        Object.setPrototypeOf(this, SingleEqualsError.prototype);
    }
}

// ExpectedTypeError indicates the wrong block type was found at the position.
export class ExpectedTypeError extends Error {
    Pos: number;
    Got: BlockType;
    Expected: BlockType[]; // Expected one of these
    constructor(pos: number, got: BlockType, expected: BlockType[]) {
        let expectedString = "";
        for (let i = 0; i < expected.length; i++) {
            expectedString += block.typeToString(expected[i]);
            if (i != expected.length-1) expectedString += " or ";
        }
        super(`near char ${pos}: got type ${block.typeToString(got)}, expected ${expectedString}`);
        this.name = "ExpectedTypeError";
        this.Pos = pos;
        this.Got = got;
        this.Expected = expected;
        Object.setPrototypeOf(this, ExpectedTypeError.prototype);
    }
}

// ExpectedError indicates the wrong text or character was found at the position.
export class ExpectedError extends Error {
    Pos: number;
    Got: string;
    Expected: string;
    constructor(pos: number, got: string, expected: string) {
        super(`near char %d: got \"${got}\", expected ${expected}`);
        this.name = "ExpectedError";
        this.Pos = pos;
        this.Got = got;
        this.Expected = expected;
        Object.setPrototypeOf(this, ExpectedError.prototype);
    }
}

// Template completes the given template string given the values provided.
// If failed, will return an empty string and an Error.
// If succeeded, will return the templated string and null.
// If succeeded with a warning, will return the templated string and an Error.
export function Template(input: string, vals?: Map<string, any>): [string, Error] {
	const t = new templater(input, vals);
    
	let a: block = new block();
	let err: Error|null | null = null;
    let out: string;
	a.Type = BlockType.EOF;
	while (true) {
		a = t.nextFromBuf()
		if (a.Type == BlockType.EOF) {
			break
		}
		[out, err] = t.process(a)
		if (err != null) {
			break
		}
        t.output += out;
	}
	if (err == null) {
		err = t.warning
	}
	return [t.output, err];
}

class templater {
	input: string;
	output: string;
	len: number;
	// Last read byte (i.e. start at -1)
	pos: number;
	vals: Map<string, any>;
	// Flag set when we're in a { ... } (or {{ ... }}) block, indicating we should tokenize text.
	inLogic: boolean;
	// Flag set when we're in a quoted string with the flag byte, or 0 when not.
	inString: '\'' | '"' | '`' | null;
	buffer: {
		buf: Array<block>; // of length seekBufferSize, but I don't wanna figure out how to "type" that
		pos: number;
	};
	warning: Error|null | null; // Non-fatal Error, returned at completion, rather than terminating early.

    constructor(input: string, vals?: Map<string, any>) {
        this.input = input;
        this.output = "";
        this.len = input.length;
        this.pos = -1;
        this.vals = vals;
        this.inLogic = false;
        this.inString = null;
        this.warning = null;
        if (!(vals)) {
            this.vals = new Map<string, any>();
        }
        this.buffer = {
            buf: new Array<block>(seekBufferSize),
            pos: 0
        };
        for (let i = 0; i < seekBufferSize; i++) {
            this.buffer.buf[i] = new block();
            this.next(i)
        }
    }

    getChar(): string | null {
        if (this.pos+1 == this.len) {
            return null;
        }
        this.pos += 1
        return this.input.charAt(this.pos)
    }

    peekChar(): string | null {
        if (this.pos+1 == this.len) {
            return null;
        }
        return this.input[this.pos+1]
    }

    next = (i: number /* index of this.buffer.buf */) => {
        this.buffer.buf[i].parent = this.input;
        this.buffer.buf[i].Type = BlockType.PlainText;
        this.buffer.buf[i].a = this.pos + 1;
        this.buffer.buf[i].b = -1;
        let c: string | null = null;
        while (true) {
            c = this.getChar();
            if (c == null) {
                this.buffer.buf[i].Type = BlockType.EOF;
                break;
            }
            if (!this.inLogic) {
                if (c == '{') {
                    this.buffer.buf[i].Type = BlockType.LogicOpen;
                    this.buffer.buf[i].a = this.pos;
                    this.buffer.buf[i].b = this.pos;
                    if (this.peekChar() == '{') {
                        this.warning = new DoubleBraceError(this.pos);
                        this.getChar();
                        this.buffer.buf[i].b = this.pos;
                    }
                    this.inLogic = true;
                    break;
                }
                this.buffer.buf[i].b = this.pos;
                if (this.peekChar() == '{' || this.peekChar() == null) {
                    break;
                }
                continue;
            }
            if (this.inString != null) {
                if (c == this.inString) {
                    this.buffer.buf[i].b = this.pos;
                    this.inString = null;
                    break;
                } else {
                    continue;
                }
            }
            if (c == '}') {
                this.buffer.buf[i].Type = BlockType.LogicClose;
                this.buffer.buf[i].a = this.pos;
                this.buffer.buf[i].b = this.pos;
                if (this.peekChar() == '}') {
                    this.warning = new DoubleBraceError(this.pos);
                    this.getChar();
                    this.buffer.buf[i].b = this.pos;
                }
                this.inLogic = false;
                break;
            }
            if (c == '"' || c == '\'' || c == '`') {
                this.buffer.buf[i].Type = BlockType.String;
                this.buffer.buf[i].a = this.pos;
                this.inString = c;
                continue;
            }
            if (c == ' ' || c == '\t') {
                continue;
            }
            if (this.buffer.buf[i].Type != BlockType.Word) {
                this.buffer.buf[i].Type = BlockType.Word;
                this.buffer.buf[i].a = this.pos;
            }
            if (this.buffer.buf[i].Type == BlockType.Word) {
                this.buffer.buf[i].b = this.pos;
                let next = this.peekChar();
                if (next == ' ' || next == '\t' || next == '}') {
                    break;
                }
            }
        }
    }

    nextFromBuf(): block {
        const out = this.buffer.buf[this.buffer.pos].copy();
        this.next(this.buffer.pos)
        this.buffer.pos = (this.buffer.pos + 1) % seekBufferSize;
        return out;
    }

    peek(): block {
        return this.buffer.buf[this.buffer.pos];
    }

    process(a: block): [string, Error|null] {
        switch (a.Type) {
        case BlockType.PlainText:
            return [a.String(), null]; 
        case BlockType.LogicOpen:
            return this.logicOpen(a);
        // LogicClose and Word/String should only occur within logic blocks and so
        // they should not appear here.
        case BlockType.LogicClose:
        case BlockType.Word:
        case BlockType.String:
            return ["", a.expected(BlockType.LogicOpen, BlockType.PlainText)];
        }
        return null
    }

    processIfBody(ifTrue: boolean): [string, Error|null] {
        let next: block;
        let content: string = "";
        let err: Error;
        while (true) {
            next = this.nextFromBuf();
            if (next.Type == BlockType.EOF) {
                break;
            }
            if (next.Type == BlockType.LogicOpen) {
                const endif = this.peek();
                if (endif.Type == BlockType.Word) {
                    const endifString = endif.String();
                    if (endifString == "endif") {
                        this.nextFromBuf();
                        const shouldBeClose = this.nextFromBuf();
                        if (shouldBeClose.Type == BlockType.LogicClose) {
                            return [content, null];
                        }
                    } else if (endifString == "else") {
                        this.nextFromBuf();
                        const shouldBeClose = this.nextFromBuf();
                        // Invert if condition to decide if we evaluate the next else/else if body.
                        ifTrue = !ifTrue;
                        if (shouldBeClose.Type == BlockType.LogicClose) {
                            // Continue the loop, i.e. print/not print depending on inverted ifTrue.
                            continue;
                        } else if (shouldBeClose.String() == "if") {
                            // Evaluate the if statement, let it calls its own copy of us.
                            let out: string;
                            [out, err] = this.ifStatement(shouldBeClose);
                            if (err != null) return ["", err];
                            if (ifTrue) {
                                content += out;
                            }
                            return [content, err];
                        }
                    }
                }
            }
            // We need to process for the sake of nested if statements.
            let out: string;
            [out, err] = this.process(next);
            if (err != null) return ["", err];
            if (ifTrue) {
                content += out;
            }
        }
        if (next.Type == BlockType.EOF) {
            return ["", next.expectedWord("{endif}")];
        }
        return [content, null];
    }

    logicOpen(open: block): [string, Error|null] {
        const ifWordOrVar = this.nextFromBuf();
        if (ifWordOrVar.Type != BlockType.Word) {
            return ["", ifWordOrVar.expected(BlockType.Word)];
        }

        const closeOrOperand = this.peek();
        if (closeOrOperand.Type == BlockType.LogicClose) {
            return this.templateValue(open, ifWordOrVar);
        }
        return this.ifStatement(ifWordOrVar);
    }

    templateValue(open: block, variable: block): [string, Error|null] {
        const val = this.vals.get(variable.String());
        const close = this.nextFromBuf()
        if (val) {
            return [String(val), null];
        }
        return[open.String() + variable.String() + close.String(), null];
    }

    ifStatement(ifWord: block): [string, Error|null] {
        if (ifWord.String() != "if") {
            return ["", ifWord.expectedWord("\"if\"")];
        }

        const operand = this.nextFromBuf()
        let [val1, err] = this.operand(operand);
        if (err != null) {
            return ["", err];
        }

        const comparisonOrClose = this.nextFromBuf();

        if (comparisonOrClose.Type == BlockType.LogicClose) {
            return this.ifTruthy(operand, val1);
        }
        return this.ifComparison(comparisonOrClose, val1);
    }

    ifTruthy(operand: block, val: any): [string, Error|null] {
        // If Bool(val)
        const positive = this.input[operand.a] != '!';
        return this.processIfBody(positive == truthy(val));
    }

    ifComparison(comparison: block, valA: any): [string, Error|null] {
        // If valA ==/!= valB
        const operandB = this.nextFromBuf();

        const comparisonString = comparison.String();
        if (comparisonString == "=") {
            this.warning = new SingleEqualsError(comparison.a);;
        } else if (comparisonString != "==" && comparisonString != "!=") {
            return ["", comparison.expectedWord("==/=/!=")];
        }

        const [valB, err] = this.operand(operandB);
        if (err != null) {
            return ["", err];
        }

        const shouldBeClose = this.nextFromBuf();
        if (shouldBeClose.Type != BlockType.LogicClose) {
            return ["", shouldBeClose.expected(BlockType.LogicClose)];
        }

        let ifTrue: boolean;
        // No need for else here, the value has already been checked to be valid.
        if (comparisonString == "==" || comparisonString == "=") {
            ifTrue = valA == valB;
        } else if (comparisonString == "!=") {
            ifTrue = valA != valB;
        }
        return this.processIfBody(ifTrue);
    }

    operand(a: block): [any, Error] {
        if (a.Type == BlockType.String) {
            return [a.String(), null];
        } else if (a.Type == BlockType.Word) {
            let name: string;
            if (this.input[a.a] == '!') {
                name = this.input.slice(a.a+1, a.b+1);
            } else {
                name = a.String();
            }
            const val = this.vals.get(name);
            if (val) {
                return [val, null];
            } else {
                return ["", null];
            }
        } else {
            return [null, a.expected(BlockType.Word)];
        }
    }
}
