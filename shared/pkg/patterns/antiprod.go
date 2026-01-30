// Package patterns provides anti-production pattern detection.
// antiprod.go: Hierarchical P0-P3 detection for production code quality.
// Covers: JS/TS, Rust, Go, Python, Java/Kotlin, Docker, Shell.
// RULE: NEVER suppress warnings — FIX the root cause.
package patterns

import (
	"path/filepath"
	"regexp"
	"strings"
)

// AntiProdLevel represents priority levels for anti-production patterns.
type AntiProdLevel int

const (
	P0MockData   AntiProdLevel = iota // Fake data, placeholders
	P1ProdLeak                        // console.log, TODO, localhost, env, suppression attrs
	P2ErrorBlind                      // Empty catch, non-null !, unwrap, panic
	P3TypeLoose                       // as any, ts-ignore, eslint-disable
)

// AntiProdResult holds the detection result.
type AntiProdResult struct {
	Level   AntiProdLevel
	Code    string
	Match   string
	Message string
}

// === JS/TS patterns ===

var (
	reConsoleLog   = regexp.MustCompile(`\bconsole\.(log|debug|info|warn|error)\s*\(`)
	reTodoComment  = regexp.MustCompile(`(?i)\b(TODO|FIXME|HACK|XXX)\b`)
	reLocalhost    = regexp.MustCompile(`https?://localhost\b`)
	reProcessEnv   = regexp.MustCompile(`process\.env\.\w+`)
	reEmptyCatch   = regexp.MustCompile(`\.catch\s*\(\s*(?:\(\s*\)|_)\s*=>\s*\{\s*\}\s*\)`)
	reNonNullAssert = regexp.MustCompile(`[a-zA-Z_]\w*!\.[a-zA-Z_]`)
	reFetchNoError  = regexp.MustCompile(`fetch\s*\([^)\n]+\)\s*(?:\.then|;)`)
	reAsAny        = regexp.MustCompile(`\bas\s+any\b`)
	reTsIgnore     = regexp.MustCompile(`@ts-ignore|@ts-expect-error`)
	reEslintDis    = regexp.MustCompile(`eslint-disable(?:-next-line)?`)
	reTsNoCheck    = regexp.MustCompile(`@ts-nocheck`)
	// P1: XSS via set:html in Astro
	reAstroSetHtml = regexp.MustCompile(`set:html\s*=`)
	// P1: XSS via dangerouslySetInnerHTML
	reDangerousHTML = regexp.MustCompile(`dangerouslySetInnerHTML`)
	// P1: XSS via innerHTML assignment
	reInnerHTML = regexp.MustCompile(`\.innerHTML\s*=`)
	// P2: double assertion bypasses type safety
	reAsUnknownAs = regexp.MustCompile(`\bas\s+unknown\s+as\b`)
	// P2: bare @ts-expect-error without explanation
	reTsExpectErrorBare = regexp.MustCompile(`@ts-expect-error\s*$`)
	// P2: client:only skips SSR
	reAstroClientOnly = regexp.MustCompile(`client:only`)
	// P1: define:vars with unescaped user input (XSS via JSON serialization)
	reAstroDefineVars = regexp.MustCompile(`define:vars\s*=`)
	// P1: Astro.redirect with dynamic/user input (open redirect)
	reAstroRedirect = regexp.MustCompile(`Astro\.redirect\s*\(`)
	// P2: is:inline script (no bundling, no CSP nonce)
	reAstroIsInline = regexp.MustCompile(`is:inline`)
	// P2: transition:persist without cleanup (memory leak in SPA mode)
	reAstroTransitionPersist = regexp.MustCompile(`transition:persist\b`)
	// P1: Astro.url used in href/src (header injection / XSS via x-forwarded-host)
	reAstroUrlInLink = regexp.MustCompile(`(?:href|src|action)\s*=\s*[{"']?\s*(?:\$\{)?\s*Astro\.url`)
	// P3: any[] type
	reAnyArray = regexp.MustCompile(`:\s*any\[\]`)
	// P3: Record<string, any>
	reRecordAny = regexp.MustCompile(`Record<string,\s*any>`)
	// P3: Object as type (capital O)
	reObjectType = regexp.MustCompile(`:\s*Object\b`)
	// P1: eval() code injection (JS/TS)
	reEval = regexp.MustCompile(`\beval\s*\(`)
	// P1: new Function() is eval equivalent
	reNewFunction = regexp.MustCompile(`new\s+Function\s*\(`)
	// P1: document.write XSS
	reDocumentWrite = regexp.MustCompile(`document\.write\s*\(`)
	// P1: setTimeout/setInterval with string arg (eval equivalent)
	reSetTimeoutString = regexp.MustCompile(`\b(?:setTimeout|setInterval)\s*\(\s*['"]`)
	// P1: JSON.parse without type validation (untyped at runtime)
	reJSONParseAs = regexp.MustCompile(`JSON\.parse\s*\([^)]*\)\s+as\s+`)
	// P1: fetch .json() cast without validation
	reFetchJsonAs = regexp.MustCompile(`\.json\(\)\s*(?:as\s+|then\s*\(\s*\([^)]*\)\s*=>\s*\S+\s+as\s+)`)
	// P2: delete operator (sparse objects, perf hit, type unsafety)
	reDeleteOperator = regexp.MustCompile(`\bdelete\s+\w+[\[.]`)
	// P2: numeric enum (reverse mapping pitfall)
	reNumericEnum = regexp.MustCompile(`\benum\s+\w+\s*\{[^}]*=\s*\d`)
	// P3: String constructor instead of template literal
	reStringConstructor = regexp.MustCompile(`new\s+String\s*\(`)
	// P3: Function type (capital F — use arrow type)
	reFunctionType = regexp.MustCompile(`:\s*Function\b`)
)

// === Rust patterns (STRICT — no suppression, no shortcuts) ===

var (
	reEnvNoDefault     = regexp.MustCompile(`std::env::var\(\s*"[^"]+"\s*\)\.unwrap\(\)`)
	reRustUnwrap       = regexp.MustCompile(`\.unwrap\(\)`)
	reRustExpect       = regexp.MustCompile(`\.expect\(\s*"`)
	reRustAllowDead    = regexp.MustCompile(`#\[allow\(dead_code\)\]`)
	reRustAllowUnused  = regexp.MustCompile(`#\[allow\(unused`)
	reRustAllowClippy  = regexp.MustCompile(`#\[allow\(clippy::`)
	reRustCfgAllowDead = regexp.MustCompile(`#\[cfg_attr\([^)]*allow\(dead_code\)`)
	// Detects renaming vars to _var to suppress unused warnings
	// Matches: let _foo = ..., _bar: Type, mut _baz
	// Does NOT match: _ = (wildcard discard) or _single underscore
	reRustUnderscoreVar = regexp.MustCompile(`\blet\s+(?:mut\s+)?_[a-zA-Z]\w*\s*[=:]`)
	// P1: debug macro left in code
	reRustDbg = regexp.MustCompile(`\bdbg!\s*\(`)
	// P1: println!/eprintln! instead of structured logging
	reRustPrintln = regexp.MustCompile(`\b(?:println!|eprintln!)\s*\(`)
	// P1: todo!/unimplemented! panics at runtime
	reRustTodoMacro = regexp.MustCompile(`\b(?:todo!|unimplemented!)\s*\(`)
	// P1: unsafe block without SAFETY justification
	reRustUnsafe = regexp.MustCompile(`\bunsafe\s*\{`)
	// P2: .expect() with short/generic message (<=20 chars)
	reRustShortExpect = regexp.MustCompile(`\.expect\(\s*"[^"]{0,20}"\s*\)`)
	// P2: excessive .clone() (heuristic)
	reRustClone = regexp.MustCompile(`\.clone\(\)`)
	// P1: panic!() in non-main Rust files
	reRustPanic = regexp.MustCompile(`\bpanic!\s*\(`)
	// P2: mem::forget leaks resources
	reRustMemForget = regexp.MustCompile(`mem::forget\s*\(`)
	// P1: process::exit skips destructors
	reRustProcessExit = regexp.MustCompile(`(?:std::)?process::exit\s*\(`)
	// P1: Box::leak intentional memory leak
	reRustBoxLeak = regexp.MustCompile(`Box::leak\s*\(`)
	// P1: mem::transmute dangerous type coercion
	reRustTransmute = regexp.MustCompile(`(?:mem::)?transmute\s*[:<(]`)
	// P1: unreachable_unchecked — UB if ever reached
	reRustUnreachableUnchecked = regexp.MustCompile(`unreachable_unchecked\s*\(`)
	// P2: #[allow(unused_must_use)] silences important return values
	reRustAllowMustUse = regexp.MustCompile(`#\[allow\(unused_must_use\)\]`)
	// P2: lossy `as` casts (u32 as u8, i64 as i32, etc.)
	reRustLossyCast = regexp.MustCompile(`\bas\s+(?:u8|i8|u16|i16|u32|i32|f32)\b`)
	// P3: String::from("...").as_str() — unnecessary allocation
	reRustStringFromAsStr = regexp.MustCompile(`String::from\s*\([^)]+\)\.as_str\(\)`)
)

// === Go patterns (STRICT) ===

var (
	reGoFmtPrint     = regexp.MustCompile(`\bfmt\.(Print|Printf|Println)\s*\(`)
	reGoPanic        = regexp.MustCompile(`\bpanic\s*\(`)
	reGoBlankDiscard = regexp.MustCompile(`\b_\s*=\s*\w+\.\w+\(`) // _ = foo.Bar() discarding error
	reGoNolint       = regexp.MustCompile(`//\s*nolint`)
	// P1: os.Exit outside main
	reGoOsExit = regexp.MustCompile(`\bos\.Exit\s*\(`)
	// P2: defer inside loop (resource leak)
	reGoDeferInLoop = regexp.MustCompile(`for\s.*\{[^}]*defer\s`)
	// P2: time.Sleep in handler code (usually a hack)
	reGoTimeSleep = regexp.MustCompile(`\btime\.Sleep\s*\(`)
)

// === Python patterns ===

var (
	rePyPrint      = regexp.MustCompile(`\bprint\s*\(`)
	rePyBareExcept = regexp.MustCompile(`except\s*:\s*$`)
	rePyExceptPass = regexp.MustCompile(`except\s+\w+.*:\s*\n\s*pass\s*$`)
	rePyTypeIgnore = regexp.MustCompile(`#\s*type:\s*ignore`)
	rePyNoqa       = regexp.MustCompile(`#\s*noqa`)
	// P1: eval/exec code injection
	rePyEval = regexp.MustCompile(`\b(?:eval|exec)\s*\(`)
	// P1: debug statements left in code
	rePyDebug = regexp.MustCompile(`\b(?:import\s+pdb|pdb\.set_trace|breakpoint)\s*\(?`)
	// P1: os.system shell injection risk
	rePyOsSystem = regexp.MustCompile(`\bos\.system\s*\(`)
	// P1: pickle.load on untrusted data
	rePyPickleLoad = regexp.MustCompile(`\bpickle\.loads?\s*\(`)
	// P1: yaml.load without SafeLoader
	rePyYamlLoad = regexp.MustCompile(`\byaml\.load\s*\(`)
	// P2: assert in non-test code (disabled with -O)
	rePyAssert = regexp.MustCompile(`\bassert\s+`)
)

// === Java/Kotlin patterns ===

var (
	reJavaSysOut      = regexp.MustCompile(`System\.out\.print`)
	reJavaEmptyCatch  = regexp.MustCompile(`catch\s*\([^)]*\)\s*\{\s*\}`)
	reJavaSuppress    = regexp.MustCompile(`@SuppressWarnings`)
)

// === Docker/Infra patterns ===

var (
	reDockerLatest = regexp.MustCompile(`(?i)^FROM\s+\S+:latest`)
	reChmod777     = regexp.MustCompile(`chmod\s+777`)
	rePipeToBash   = regexp.MustCompile(`(curl|wget)\s+[^|]*\|\s*(ba)?sh`)
	// P2: ADD instead of COPY (ADD can fetch URLs, extract archives)
	reDockerADD = regexp.MustCompile(`(?im)^ADD\s+`)
	// P1: no USER directive (running as root)
	reDockerUser = regexp.MustCompile(`(?im)^USER\s+`)
)

// DetectAntiProd runs hierarchical P0→P3 checks.
// Returns on first detection at each level (P0 first, then P1, etc.).
// All levels are checked — caller decides whether to stop at first hit.
func DetectAntiProd(filePath, content string) []AntiProdResult {
	if content == "" || isAllowlisted(filePath) {
		return nil
	}

	var results []AntiProdResult

	// P0: Mock data (delegate to existing DetectMockData)
	if detected, reason := DetectMockData(filePath, content); detected {
		results = append(results, AntiProdResult{
			Level:   P0MockData,
			Code:    "MOCK_DATA",
			Match:   reason,
			Message: "Replace hardcoded data with real API fetch. Use: fetch('/api/...') or sqlx::query_as.",
		})
	}

	// P1: Production leaks (frontend + backend)
	if isFrontendFile(filePath) {
		if m := reConsoleLog.FindString(content); m != "" {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "console.log",
				Message: "Remove debug output or use structured logger.",
			})
		}
	}
	if m := reTodoComment.FindString(content); m != "" {
		results = append(results, AntiProdResult{
			Level:   P1ProdLeak,
			Code:    "PROD_LEAK",
			Match:   "TODO/FIXME",
			Message: "Implement or create ticket. Do not ship TODO comments.",
		})
	}
	if isNonConfigFile(filePath) && reLocalhost.MatchString(content) {
		results = append(results, AntiProdResult{
			Level:   P1ProdLeak,
			Code:    "PROD_LEAK",
			Match:   "http://localhost",
			Message: "Use config/environment variable for URLs.",
		})
	}
	if isFrontendFile(filePath) {
		for _, line := range strings.Split(content, "\n") {
			if reProcessEnv.MatchString(line) && !hasEnvFallback(line) {
				results = append(results, AntiProdResult{
					Level:   P1ProdLeak,
					Code:    "PROD_LEAK",
					Match:   "process.env without fallback",
					Message: "Add fallback: process.env.X ?? 'default' or process.env.X || 'default'.",
				})
				break
			}
		}
	}
	if isBackendFile(filePath) && reEnvNoDefault.MatchString(content) {
		results = append(results, AntiProdResult{
			Level:   P1ProdLeak,
			Code:    "PROD_LEAK",
			Match:   "env::var().unwrap()",
			Message: "Use env::var().unwrap_or_else() or env::var().ok().",
		})
	}

	// P1: Go — fmt.Print in non-test .go
	if isGoFile(filePath) {
		if reGoFmtPrint.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "fmt.Print",
				Message: "Use structured logger (slog/zerolog) instead of fmt.Print.",
			})
		}
		base := strings.ToLower(filepath.Base(filePath))
		if reGoOsExit.MatchString(content) && base != "main.go" {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "os.Exit outside main",
				Message: "Return error instead of os.Exit. Only main() should call os.Exit.",
			})
		}
		if reGoNolint.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P3TypeLoose,
				Code:    "TYPE_LOOSE",
				Match:   "nolint",
				Message: "Fix the lint issue instead of suppressing with nolint.",
			})
		}
	}

	// P1: Python — print() in non-test .py
	if isPythonFile(filePath) {
		if rePyPrint.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "print()",
				Message: "Use logging module instead of print().",
			})
		}
		if rePyEval.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "eval()/exec()",
				Message: "Code injection risk. Use ast.literal_eval() or safer alternatives.",
			})
		}
		if rePyDebug.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "pdb/breakpoint",
				Message: "Remove debug statements before shipping.",
			})
		}
		if rePyOsSystem.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "os.system()",
				Message: "Shell injection risk. Use subprocess.run() with shell=False.",
			})
		}
		if rePyPickleLoad.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "pickle.load()",
				Message: "Arbitrary code execution risk. Use JSON or a safe serializer.",
			})
		}
		if rePyYamlLoad.MatchString(content) && !strings.Contains(content, "SafeLoader") && !strings.Contains(content, "safe_load") {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "yaml.load() without SafeLoader",
				Message: "Use yaml.safe_load() or yaml.load(Loader=SafeLoader).",
			})
		}
		if rePyAssert.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "assert in production code",
				Message: "Assert is disabled with python -O. Use if/raise for runtime checks.",
			})
		}
		if rePyNoqa.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P3TypeLoose,
				Code:    "TYPE_LOOSE",
				Match:   "noqa",
				Message: "Fix the lint issue instead of suppressing with noqa.",
			})
		}
	}

	// P1: Java/Kotlin — System.out.print in non-test
	if isJavaFile(filePath) {
		if reJavaSysOut.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "System.out.print",
				Message: "Use SLF4J/Log4j instead of System.out.",
			})
		}
	}

	// P1: Docker/Infra
	if isDockerfile(filePath) {
		if reDockerLatest.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "FROM :latest",
				Message: "Pin image version instead of using :latest.",
			})
		}
	}
	if isDockerfile(filePath) {
		if reDockerADD.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "ADD instead of COPY",
				Message: "Use COPY instead of ADD. ADD can fetch URLs and extract archives unexpectedly.",
			})
		}
		if !reDockerUser.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "no USER directive",
				Message: "Container runs as root. Add USER directive for least-privilege.",
			})
		}
	}
	if reChmod777.MatchString(content) {
		results = append(results, AntiProdResult{
			Level:   P1ProdLeak,
			Code:    "PROD_LEAK",
			Match:   "chmod 777",
			Message: "Use least-privilege permissions instead of chmod 777.",
		})
	}
	if (isDockerfile(filePath) || isShellFile(filePath)) && rePipeToBash.MatchString(content) {
		results = append(results, AntiProdResult{
			Level:   P1ProdLeak,
			Code:    "PROD_LEAK",
			Match:   "curl|bash",
			Message: "Download and verify before piping to shell.",
		})
	}

	// P2: Error blindness
	if isGoFile(filePath) {
		goBase := strings.ToLower(filepath.Base(filePath))
		if reGoPanic.MatchString(content) && goBase != "main.go" {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "panic()",
				Message: "Return error instead of panic.",
			})
		}
		if reGoBlankDiscard.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "_ = (error discarded)",
				Message: "Handle the error instead of discarding with blank identifier.",
			})
		}
		if reGoDeferInLoop.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "defer in loop",
				Message: "Defer in loop won't execute until function returns. Extract to separate function.",
			})
		}
		if reGoTimeSleep.MatchString(content) && isHandlerFile(filePath) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "time.Sleep in handler",
				Message: "Use context-aware waiting (select on ctx.Done) instead of time.Sleep.",
			})
		}
	}
	if isPythonFile(filePath) {
		if rePyBareExcept.MatchString(content) || rePyExceptPass.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "bare except / except pass",
				Message: "Handle errors explicitly instead of bare except or pass.",
			})
		}
		if rePyTypeIgnore.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "type: ignore",
				Message: "Add proper type annotation instead of type: ignore.",
			})
		}
	}
	if isJavaFile(filePath) {
		if reJavaEmptyCatch.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "empty catch block",
				Message: "Handle exception instead of empty catch block.",
			})
		}
	}

	// P3: Java suppression
	if isJavaFile(filePath) && reJavaSuppress.MatchString(content) {
		results = append(results, AntiProdResult{
			Level:   P3TypeLoose,
			Code:    "TYPE_LOOSE",
			Match:   "@SuppressWarnings",
			Message: "Fix generic types instead of @SuppressWarnings.",
		})
	}

	// P2: Error blindness (frontend)
	if isFrontendFile(filePath) {
		if reEmptyCatch.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   ".catch(() => {})",
				Message: "Handle errors: log, show user feedback, or re-throw.",
			})
		}
		if reNonNullAssert.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "non-null assertion (!.)",
				Message: "Use optional chaining (?.) instead of non-null assertion.",
			})
		}
		if reFetchNoError.MatchString(content) && !strings.Contains(content, ".catch") && !strings.Contains(content, "try") {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "fetch() without error handling",
				Message: "Wrap fetch in try/catch or add .catch() handler.",
			})
		}
	}
	if isBackendFile(filePath) && strings.HasSuffix(strings.ToLower(filePath), ".rs") {
		// Only flag .unwrap() in handler files, not tests
		if reRustUnwrap.MatchString(content) && isHandlerFile(filePath) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   ".unwrap() in handler",
				Message: "Use ? operator instead of .unwrap() in request handlers.",
			})
		}
	}

	// Rust-specific patterns (expanded)
	if isRustFile(filePath) {
		if reRustDbg.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "dbg!()",
				Message: "Remove dbg!() — runs in release builds. Use tracing::debug!().",
			})
		}
		if reRustPrintln.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "println!/eprintln!",
				Message: "Use tracing/log crate instead of println!/eprintln!.",
			})
		}
		if reRustTodoMacro.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "todo!/unimplemented!",
				Message: "Implement before shipping. todo!/unimplemented! panics at runtime.",
			})
		}
		rsBase := strings.ToLower(filepath.Base(filePath))
		if reRustPanic.MatchString(content) && rsBase != "main.rs" {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "panic!()",
				Message: "Return Result/Option instead of panic!. Only main/tests should panic.",
			})
		}
		if reRustUnsafe.MatchString(content) && !strings.Contains(content, "// SAFETY:") {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "unsafe block",
				Message: "Justify unsafe block with // SAFETY: comment or remove.",
			})
		}
		if reRustShortExpect.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   ".expect() with generic message",
				Message: "Use ? operator or provide a context-rich expect message.",
			})
		}
		if reRustMemForget.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "mem::forget()",
				Message: "mem::forget leaks resources. Use ManuallyDrop or drop() explicitly.",
			})
		}
		if matches := reRustClone.FindAllString(content, -1); len(matches) > 2 {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "excessive .clone()",
				Message: "Excessive .clone() — consider borrowing or Arc/Rc for shared ownership.",
			})
		}
		if reRustProcessExit.MatchString(content) && rsBase != "main.rs" {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "process::exit()",
				Message: "process::exit skips destructors. Return from main or use anyhow::bail!.",
			})
		}
		if reRustBoxLeak.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "Box::leak()",
				Message: "Box::leak intentionally leaks memory. Use Arc or lifetime-bound references.",
			})
		}
		if reRustTransmute.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "mem::transmute",
				Message: "transmute is extremely unsafe. Use TryFrom, bytemuck, or zerocopy instead.",
			})
		}
		if reRustUnreachableUnchecked.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "unreachable_unchecked()",
				Message: "UB if ever reached. Use unreachable!() which panics safely.",
			})
		}
		if reRustAllowMustUse.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "#[allow(unused_must_use)]",
				Message: "Handle the return value. #[must_use] exists for a reason.",
			})
		}
		if reRustLossyCast.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "lossy as cast",
				Message: "Lossy cast may truncate/overflow. Use TryFrom or .try_into() with error handling.",
			})
		}
		if reRustStringFromAsStr.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P3TypeLoose,
				Code:    "TYPE_LOOSE",
				Match:   "String::from().as_str()",
				Message: "Unnecessary allocation. Use string literal directly.",
			})
		}
	}

	// Astro-specific patterns
	if isAstroFile(filePath) {
		if reAstroSetHtml.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "set:html",
				Message: "XSS risk: set:html injects unescaped HTML. Sanitize input or use {text}.",
			})
		}
		if reAstroDefineVars.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "define:vars",
				Message: "define:vars inlines via JSON.stringify — XSS if user input reaches it. Sanitize first.",
			})
		}
		if reAstroUrlInLink.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "Astro.url in link",
				Message: "Astro.url reflects x-forwarded-host header. Validate or allowlist before using in links (CVE-2025-61925).",
			})
		}
		if reAstroRedirect.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "Astro.redirect()",
				Message: "Open redirect risk if input is user-controlled. Validate redirect target against allowlist.",
			})
		}
		if reAstroClientOnly.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "client:only",
				Message: "client:only skips SSR — ensure no-JS fallback exists.",
			})
		}
		if reAstroIsInline.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "is:inline",
				Message: "is:inline skips bundling and CSP nonce injection. Use bundled scripts where possible.",
			})
		}
		if reAstroTransitionPersist.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "transition:persist",
				Message: "Persisted elements survive navigation — add cleanup logic to prevent memory leaks.",
			})
		}
	}

	// JS/TS security patterns (frontend)
	if isFrontendFile(filePath) {
		if reEval.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "eval()",
				Message: "Code injection risk. Never use eval() — refactor to avoid dynamic code execution.",
			})
		}
		if reNewFunction.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "new Function()",
				Message: "new Function() is eval equivalent. Refactor to avoid dynamic code execution.",
			})
		}
		if reDocumentWrite.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "document.write()",
				Message: "XSS risk and blocks parsing. Use DOM manipulation instead.",
			})
		}
		if reSetTimeoutString.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "setTimeout/setInterval with string",
				Message: "String arg to setTimeout/setInterval is eval. Pass a function instead.",
			})
		}
		if reJSONParseAs.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "JSON.parse() as Type",
				Message: "JSON.parse returns unknown at runtime. Validate with Zod/io-ts instead of casting.",
			})
		}
		if reFetchJsonAs.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   ".json() as Type",
				Message: "External data is untyped at runtime. Validate with schema (Zod) before casting.",
			})
		}
		if reDangerousHTML.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "dangerouslySetInnerHTML",
				Message: "XSS risk: sanitize HTML before injecting.",
			})
		}
		if reInnerHTML.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P1ProdLeak,
				Code:    "PROD_LEAK",
				Match:   "innerHTML =",
				Message: "XSS risk: use textContent or sanitize before innerHTML.",
			})
		}
		if reAsUnknownAs.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "as unknown as",
				Message: "Double assertion bypasses type safety. Use type guard or refactor types.",
			})
		}
		if reTsExpectErrorBare.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "@ts-expect-error (bare)",
				Message: "Add explanation after @ts-expect-error.",
			})
		}
		if reDeleteOperator.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "delete operator",
				Message: "delete creates sparse objects and breaks V8 optimizations. Use destructuring or omit.",
			})
		}
		if reNumericEnum.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P2ErrorBlind,
				Code:    "ERROR_BLIND",
				Match:   "numeric enum",
				Message: "Numeric enums have reverse mapping pitfalls. Use string enums or as const objects.",
			})
		}
	}

	// P3: Type looseness (frontend only)
	if isFrontendFile(filePath) {
		if reAsAny.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P3TypeLoose,
				Code:    "TYPE_LOOSE",
				Match:   "as any",
				Message: "Use proper type narrowing instead of 'as any'.",
			})
		}
		if reTsIgnore.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P3TypeLoose,
				Code:    "TYPE_LOOSE",
				Match:   "@ts-ignore",
				Message: "Fix the type error instead of suppressing it.",
			})
		}
		if reEslintDis.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P3TypeLoose,
				Code:    "TYPE_LOOSE",
				Match:   "eslint-disable",
				Message: "Fix the lint error instead of disabling the rule.",
			})
		}
		if reTsNoCheck.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P3TypeLoose,
				Code:    "TYPE_LOOSE",
				Match:   "@ts-nocheck",
				Message: "Remove @ts-nocheck and fix type errors.",
			})
		}
		if reAnyArray.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P3TypeLoose,
				Code:    "TYPE_LOOSE",
				Match:   "any[]",
				Message: "Use typed array instead of any[].",
			})
		}
		if reRecordAny.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P3TypeLoose,
				Code:    "TYPE_LOOSE",
				Match:   "Record<string, any>",
				Message: "Use a specific type instead of Record<string, any>.",
			})
		}
		if reObjectType.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P3TypeLoose,
				Code:    "TYPE_LOOSE",
				Match:   "Object type",
				Message: "Use object (lowercase) or a specific interface.",
			})
		}
		if reStringConstructor.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P3TypeLoose,
				Code:    "TYPE_LOOSE",
				Match:   "new String()",
				Message: "new String() creates wrapper object, not primitive. Use String() or template literal.",
			})
		}
		if reFunctionType.MatchString(content) {
			results = append(results, AntiProdResult{
				Level:   P3TypeLoose,
				Code:    "TYPE_LOOSE",
				Match:   "Function type",
				Message: "Use specific function signature (e.g. () => void) instead of Function.",
			})
		}
	}

	return results
}

// isNonConfigFile returns true for files that should not contain localhost.
func isNonConfigFile(path string) bool {
	p := strings.ToLower(path)
	configPatterns := []string{
		"config", ".env", "astro.config", "vite.config", "next.config",
		"wrangler.toml", "docker-compose", "Caddyfile", ".toml",
		"dev_ports", "constants",
	}
	for _, pat := range configPatterns {
		if strings.Contains(p, pat) {
			return false
		}
	}
	return true
}

// isHandlerFile returns true for Rust handler/route files.
func isHandlerFile(path string) bool {
	p := strings.ToLower(path)
	return strings.Contains(p, "handler") || strings.Contains(p, "routes") ||
		strings.Contains(p, "lib.rs") || strings.Contains(p, "main.rs")
}

// hasEnvFallback checks if a line with process.env already has a fallback.
func hasEnvFallback(line string) bool {
	return strings.Contains(line, "??") || strings.Contains(line, "||") ||
		strings.Contains(line, "import.meta.env")
}
