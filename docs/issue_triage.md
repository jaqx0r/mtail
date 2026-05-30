# google/mtail open issue triage

Date: 2026-05-27  
Source: GitHub API + manual inspection of comments and PRs in `jaqx0r/mtail`

## Methodology

- Fetched all 57 open items from `github.com/google/mtail`
- Filtered out PRs, dependabot bots, and OSS-Fuzz issues → 37 real issues
- Inspected every issue body + comment thread for status
- Searched `jaqx0r/mtail` for any fixing PRs, commits, or branches referencing each issue

## Key: status classification

| Label | Meaning |
|---|---|
| **Done** | Fixed by a PR/commit in `jaqx0r/mtail` — no action needed |
| **Valid Bug** | Confirmed reproducible defect with no known fix in this fork |
| **Answered** | Already answered by maintainer, no further action required |
| **Feature Request** | Enhancement welcome via PR, no code written yet |

---

## Done (fixed in jaqx0r/mtail)

| # | Title | Fix |
|---|---|---|
| 531 | Log rotation fails under WSL/Docker | PR #453 in jaqx0r/mtail: `FILE_SHARE_DELETE` on Windows |
| 903 | /metrics fails with >70k metrics | PR #454: reproduction benchmark added; root cause may still need fixing |
| 340 | Regex concatenation doesn't work as expected | Grammar: `const_pattern` and `concat_start` non-terminals allow `id_expr` at start of concatenation |

---

## Valid Bugs (need addressing)

| # | Title | Notes |
|---|---|---|
| 837 | `limit` keyword doesn't work with `hidden` keyword | GC never evicts hidden metrics; no fix known |
| 763 | `limit` keyword not releasing memory after eviction | Go map `map[string]*LabelValue` never shrinks (golang/go#20135) |
| 807 | Windows: mtail locks watched files, prevents deletion/rotation | Partially overlaps #531 but Windows-specific locking may remain |
| 725 | File with initially-incorrect permissions later fixed — tailer loops | PR #456 in jaqx0r/mtail: close fs.lines on rotation reopen failure |
| 636 | Duplicate metric collected after SIGHUP reload | Metric declaration moves line numbers, dupe check compares source location |
| 608 | Index out of range `[-1]` in vm.go:342 | `{neg <nil> 2}` instruction; known to jaqx0r, "tricky to fix without reengineering compiler" |
| 590 | VM parses IP address `10.0.0.1` as float with histogram | Type inference sees single `.` as float pattern, fails on multiple dots; only with histograms |
| 581 | Checker warns about numeric capture group refs | Internally regex assigns numeric indices even for named groups; checker should ignore them |
| 505 | NFS stale file handle causes 100% CPU loop | jaqx0r pushed ESTALE handling to HEAD but issue may recur; AWS EFS users heavily affected |
| 191 | Const patterns can't be used alone as match condition | Need `// + A` prefix to avoid "capref not defined"; milestone 3.0.0 |
| 190 | Can't chain `=~` match with `/regex/` in same expression | Both define capref 0 → compile error; workaround: nested blocks |
| 267 | Need way to test if capture group defined in current scope | When group matches `-` instead of digits, `$response_size` is undefined → runtime error |
| 263 | Timestamps lose sub-second precision | `timestamp()` calls `Time.Unix()` (1s resolution); suggestion: use float64 of nanoseconds |
| 230 | Type inference converts `"05"` to `"5"` — breaks strptime | `$month` from `\d{2}` becomes int `5`, string conversion yields `"5"` not `"05"` |
| 180 | No `BEGIN` block to initialize metrics at startup | Milestone 3.0.0; assigned to jaqx0r |
| 221 | Runtime errors logged at INFO, not visible to users | Users don't see errors unless they watch stderr at INFO level |

---

## Already Answered / No Further Action

| # | Title | Verdict |
|---|---|---|
| 972 | How is parser/runtime implemented? | jaqx0r: "read the code, it's well commented" |
| 971 | Why was fsnotify removed for polling? | jaqx0r: fsnotify was "complex, slow, unperformant" |
| 970 | Multi-line text matching? | mtail is line-oriented; not planned |
| 965 | Docker container logs as input? | No specialist client code; tee to socket instead |
| 955 | Integrate mtail into own Go program? | "What parts?" |
| 929 | This repo is obsolete, use jaqx0r/mtail | The redirect notice itself |
| 900 | v3.0.6 tag created but no binaries released | Workflow broken; check releases tab |

---

## Feature Requests (open for contribution)

| # | Title | Notes |
|---|---|---|
| 684 | Exemplar support (OpenMetrics) | jaqx0r interested; needs spec for trace ID capture |
| 679 | Push metrics to Pushgateway | HTTP PUT of Prometheus text format |
| 678 | Helm chart for Kubernetes | User `jmatis` shared one at github.com/jmatis/mtail |
| 748 | Custom annotations for Alertmanager | Blocked on exemplar support |
| 724 | Write unparsed log lines to another file | jaqx0r not keen; existed in Google-internal v2.0 |
| 542 | Windows service registration flag | `golang.org/x/sys/windows/svc` |
| 333 | Temporary variables (non-state) | Workaround: `hidden` + `del` |
| 314 | Official systemd service file | Debian bug #886894 |
| 259 | Option to disable Go runtime metrics | Question about `go_gc_duration_seconds` etc. |
| 234 | Better docs on `hidden string` variables | Covers OOMK log linkage pattern |
| 199 | Publish tailer as standalone library | Internal now; user maintains own copy |

---

## Suggested priorities

### Immediate (production-impacting bugs)

1. **#837 / #763** — `limit` with `hidden` broken; memory never released
2. **#725** — file permission fix not detected (infinite loop)
3. **#505** — NFS stale file handle (100% CPU, affects cloud deployments)
4. **#636** — duplicate metrics after SIGHUP reload
5. **#230** — type inference strips leading zeroes (breaks `strptime`)

### Reproducing test cases to write first

These can be added as table-driven tests in `internal/runtime/runtime_integration_test.go`:

- `TestLimitWithHidden` (#837)
- `TestStrptimeLeadingZeros` (#230)
- `TestIPAddressAsFloat` (#590) — regex `[\d\.]+` + histogram
- `TestConstPatternAsMatch` (#191) — const without `//` prefix
- `TestChainedMatchCaprefConflict` (#190) — `getfilename() =~ "x" && /pattern/`
- `TestUndefinedCaptureGroup` (#267) — optional group matching `-`
