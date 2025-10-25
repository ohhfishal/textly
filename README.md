# Textly

A text parser for generating terminal animations.

## In development

### New Features
- [X] Comments
- [X] Clearing the screen `{clear}`
- [ ] Better control of whitespace
    ```
    Here is some text {Something to denote ignore newline]
    This is on the same line
    ```
- [ ] Decorators
    ```
    @fast { Here is some fast text }
    @slow { Here is some slowwww text }
    @delete {
    Way for deleting across lines 
    and also alternative to [text]
    }
    @(fast, delete, tabindented) {
        Combine multiple decorators for once block.
        This text is typed fast and all deleted when done.
        Also the tabs at the start of the line are ignored!
    }
    ```
- [ ] Color support 
    ```
    Normal Text
    @red{ Here is some red text }
    ```
- [ ] Header to set options and macros
    ```
    ---
    shell: enable
    macros:
        - name: example
          color: 0xFFFFFF
          speed: fast
    ---
    @example {
    Text using custom macro
    }
    ```
- [ ] espeak integration?
    ```
    @espeak{ Say this aloud }
    ```

### Espeak
I am trying to get this program to work okay with [espeak](https://www.google.com/url?sa=t&source=web&rct=j&opi=89978449&url=https://espeak.sourceforge.net/&ved=2ahUKEwjunb2Erq6QAxUul-4BHV8qKkoQFnoECBIQAQ&usg=AOvVaw0hrm7DsB6mUHrZ_ecQDWKD). This is the best I was able to do so far:

```bash
./textly examples/simple.txt --flatten -l --delay 0.1s | tee /dev/tty |  espeak --punct=none -g 0
```
