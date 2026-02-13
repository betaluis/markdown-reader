# mdview Test Document

This document tests various markdown features supported by **mdview**.

## Headers

### Level 3 Header
#### Level 4 Header
##### Level 5 Header
###### Level 6 Header

## Text Formatting

**Bold text** and *italic text* and ***bold italic***.

~~Strikethrough text~~ (GitHub Flavored Markdown).

## Lists

### Unordered List
- First item
- Second item
  - Nested item 1
  - Nested item 2
- Third item

### Ordered List
1. First step
2. Second step
3. Third step
   1. Sub-step A
   2. Sub-step B

## Links and Images

[Visit GitHub](https://github.com)

Auto-linked URL: https://github.com

## Code

Inline `code` with backticks.

### Code Blocks

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, mdview!")
    
    // Syntax highlighting test
    for i := 0; i < 10; i++ {
        fmt.Printf("Count: %d\n", i)
    }
}
```

```python
def fibonacci(n):
    """Generate Fibonacci sequence"""
    a, b = 0, 1
    for _ in range(n):
        yield a
        a, b = b, a + b

# Print first 10 Fibonacci numbers
for num in fibonacci(10):
    print(num)
```

```javascript
// JavaScript example
const greet = (name) => {
    console.log(`Hello, ${name}!`);
};

greet('mdview');
```

## Blockquotes

> This is a blockquote.
> 
> It can span multiple lines.
> 
> > Nested blockquotes are also supported.

## Tables

| Feature | Status | Notes |
|---------|--------|-------|
| Markdown rendering | ✅ Done | Using Goldmark |
| Syntax highlighting | ✅ Done | Using Chroma |
| Live reload | ✅ Done | WebSocket |
| GitHub CSS | ✅ Done | Clean styling |

## Horizontal Rule

---

## Task Lists

- [x] Create project structure
- [x] Implement markdown renderer
- [x] Add live reload
- [ ] Add more features
- [ ] Write documentation

## Math (if supported)

Inline math: E = mc²

## Emoji (if supported)

:rocket: :star: :heart:

---

**Test Instructions:**

1. Open this file with `mdview examples/test.md`
2. Edit this file and save
3. Watch the browser automatically reload
4. Try different markdown features
5. Check syntax highlighting for different languages

*Last updated: 2026-02-12*
