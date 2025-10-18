# Textly

A text parser for generating terminal animations.

## In development
I am trying to get this program to work okay with [espeak](https://www.google.com/url?sa=t&source=web&rct=j&opi=89978449&url=https://espeak.sourceforge.net/&ved=2ahUKEwjunb2Erq6QAxUul-4BHV8qKkoQFnoECBIQAQ&usg=AOvVaw0hrm7DsB6mUHrZ_ecQDWKD). This is the best I was able to do so far:

```bash
./textly examples/simple.txt --flatten -l --delay 0.1s | tee /dev/tty |  espeak --punct=none -g 0
```
