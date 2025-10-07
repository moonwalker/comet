# AGENTS.md - Guidelines for AI Agents Working on Comet

This document provides clear instructions for AI agents (like GitHub Copilot, Claude, ChatGPT, etc.) working on the Comet codebase.

## 🎯 Core Principles

### 1. **Code Over Documentation**
- Focus on **writing code** that solves problems
- Keep documentation **minimal and essential**
- Don't create documentation for the sake of creating documentation
- Only update docs when functionality actually changes

### 2. **No Random Markdown Generation**
**❌ DO NOT:**
- Create summary files after every change
- Generate "IMPLEMENTATION_SUMMARY.md", "FINAL_SUMMARY.md", "SUCCESS.md", etc.
- Create progress tracking documents
- Make duplicate documentation
- Create meta-documentation about what you just did

**✅ DO:**
- Update existing relevant documentation (README.md, CHANGELOG.md)
- Create documentation only for NEW user-facing features
- Keep documentation in appropriate places (see structure below)

### 3. **Documentation Structure**

```
Root Level:
├── README.md              (User-facing overview and quick start)
├── CHANGELOG.md           (Version history)
├── CONTRIBUTING.md        (If we have contribution guidelines)
└── AGENTS.md              (This file - for AI agents)

docs/:
├── architecture.md        (System design)
├── best-practices.md      (Usage recommendations)
├── dsl-improvements.md    (DSL features documentation)
├── dsl-quick-reference.md (Quick syntax reference)
├── userland-patterns.md   (User-created patterns guide)
├── its-just-javascript.md (JavaScript extensibility guide)
└── cross-stack-references.md (Feature-specific docs)

stacks/_examples/:
└── *.stack.js             (Working code examples)
```

**Do not create files like:**
- ❌ IMPLEMENTATION_SUMMARY.md
- ❌ FINAL_SUMMARY.md
- ❌ SUCCESS.md
- ❌ IMPLEMENTATION_COMPLETE.md
- ❌ EMPHASIS_COMPLETE.md
- ❌ Or any other meta-summary files

## 🛠️ Development Workflow

### When Implementing a Feature:

1. **Write the code** in appropriate Go files
2. **Add tests** if applicable
3. **Update CHANGELOG.md** with the change
4. **Update README.md** only if it's a user-facing feature
5. **Add documentation** in `docs/` only if it's a complex feature needing explanation
6. **Add examples** in `stacks/_examples/` if it helps users understand

**Stop there.** Don't create summary files.

### When Fixing a Bug:

1. **Fix the code**
2. **Update CHANGELOG.md** in the "Fixed" section
3. **Done.** No documentation needed for most bug fixes.

### When Refactoring:

1. **Refactor the code**
2. **Update CHANGELOG.md** if it affects users
3. **Done.** Internal refactoring rarely needs documentation.

## 📝 Documentation Guidelines

### What to Document:

✅ **User-facing features** - How users interact with new functionality
✅ **Breaking changes** - Migration guides when necessary
✅ **Complex concepts** - Architecture decisions, design patterns
✅ **API changes** - New functions, changed signatures
✅ **Configuration options** - New settings in comet.yaml

### What NOT to Document:

❌ **Implementation details** - Code should be self-documenting
❌ **Obvious functionality** - Don't over-explain simple code
❌ **Work-in-progress** - No "TODO" documents
❌ **Summary of what you just did** - No meta-documentation
❌ **Multiple versions of the same info** - No duplication

### Documentation Style:

- **Concise** - Get to the point quickly
- **Practical** - Show code examples
- **Accurate** - Test examples before documenting
- **Up-to-date** - Remove outdated info

## 🏗️ Code Guidelines

### Go Code:

- Follow existing code style in the project
- Add comments for non-obvious logic
- Keep functions small and focused
- Use meaningful variable names
- Handle errors appropriately

### File Structure:

```
internal/
├── parser/
│   └── js/
│       └── js.go       (JavaScript parser/interpreter)
├── schema/             (Data structures)
├── secrets/            (Secrets management)
├── exec/               (Terraform/OpenTofu execution)
└── cli/                (CLI output formatting)

cmd/                    (CLI commands)
├── root.go
├── plan.go
├── apply.go
└── ...
```

### Adding New DSL Functions:

When adding new JavaScript functions available in stack files:

1. **Add the function** in `internal/parser/js/js.go`
2. **Register it** in the `Parse()` method with `vm.rt.Set()`
3. **Test it** with a real stack file
4. **Document it** in `docs/dsl-quick-reference.md` (one place only!)
5. **Add example** in `stacks/_examples/` if helpful

### Testing:

- Write tests for new functionality
- Test with real stack files in `stacks/`
- Don't create test-specific markdown documentation

## 🚫 Common Mistakes to Avoid

### ❌ Mistake 1: Documentation Proliferation
```
Agent creates:
- IMPLEMENTATION_SUMMARY.md
- FINAL_SUMMARY.md  
- SUCCESS.md
- IMPLEMENTATION_COMPLETE.md
```
**Why it's bad:** Clutters the repository, duplicates information, confuses users.

**✅ Instead:** Update CHANGELOG.md and README.md only.

### ❌ Mistake 2: Over-documenting Obvious Code
```markdown
## The envs() Function
This function sets environment variables...
(20 more paragraphs explaining what environment variables are)
```
**Why it's bad:** Users don't need computer science lectures.

**✅ Instead:** Show a code example and move on.

### ❌ Mistake 3: Creating Duplicate Docs
```
Same information in:
- README.md
- docs/guide.md
- docs/quickstart.md
- docs/introduction.md
```
**Why it's bad:** Hard to maintain, versions get out of sync.

**✅ Instead:** One piece of information in one logical place.

## ✅ Good Examples

### Example 1: Adding a Feature

```
Files changed:
✓ internal/parser/js/js.go          (code implementation)
✓ CHANGELOG.md                       (added to unreleased)
✓ docs/dsl-quick-reference.md        (syntax added)
✓ stacks/_examples/new-feature.js    (working example)

Files NOT created:
✗ IMPLEMENTATION_SUMMARY.md
✗ FEATURE_COMPLETE.md
```

### Example 2: Fixing a Bug

```
Files changed:
✓ internal/exec/tf/terraform.go      (bug fix)
✓ CHANGELOG.md                       (noted in "Fixed" section)

Files NOT created:
✗ BUG_FIX_SUMMARY.md
✗ Any other documentation
```

### Example 3: Refactoring

```
Files changed:
✓ internal/parser/parser.go          (refactored code)
✓ CHANGELOG.md                       (only if user-visible)

Files NOT created:
✗ REFACTORING_NOTES.md
✗ CODE_IMPROVEMENTS.md
```

## 📋 Checklist for Agents

Before completing a task, verify:

- [ ] Code is written and tested
- [ ] CHANGELOG.md is updated (if user-visible change)
- [ ] README.md is updated (if major feature)
- [ ] Existing docs are updated (if behavior changed)
- [ ] Examples added (if helpful for understanding)
- [ ] **Did NOT create summary/meta documentation files**
- [ ] **Did NOT duplicate information across files**
- [ ] **Did NOT over-document obvious functionality**

## 🎓 Philosophy

**Comet's philosophy extends to documentation:**

> Provide **minimal, essential documentation**. Users are smart - they can read code and examples. Don't patronize them with verbose explanations.

**Code > Docs**

- Good code is self-documenting
- Examples are better than prose
- Less is more

## 🤝 When in Doubt

**Ask yourself:**

1. Does this documentation help users accomplish a task?
2. Is this information already documented elsewhere?
3. Will this document need to be maintained?
4. Is this creating noise in the repository?

**If unsure, default to:**
- ✅ Update CHANGELOG.md
- ✅ Update README.md if truly necessary
- ❌ Don't create new markdown files

## 📚 Required Reading

Before working on Comet, understand:

1. **Philosophy:** Minimal, unopinionated tooling
2. **DSL Design:** JavaScript superset, users build their own abstractions
3. **Documentation:** Essential only, no fluff

## 🔄 This File

**AGENTS.md itself should be:**
- Updated when development guidelines change
- Kept concise and practical
- The **single source of truth** for AI agent guidelines

**AGENTS.md should NOT be:**
- A changelog of agent activities
- A list of every task completed
- Documentation about documentation

---

## Summary

**TL;DR for AI Agents:**

1. **Write code, not documentation**
2. **Update CHANGELOG.md for changes**
3. **Update README.md for major features**
4. **Don't create random markdown files**
5. **Keep it simple**

**When you finish a task, stop.** Don't create summary files. The git commit message is your summary.

---

*Last updated: 2025-10-07*
