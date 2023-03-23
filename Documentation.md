# Overview

`cbmasm`'s input is line based. A line typically consists of an optional label, a directive or a mnenonic, parameters,
and an optional comment.

A typical line could look like this:

```
wait:    cmp $d012 ; compare current vertical line pos
```                                                   

# Labels
Labels need to terminate with ":" unless they start at the beginning of the line.

## Local labels
A local label is a label that starts with an underscore (`_`).
All local labels are only visible (and need to be resolved) before the next non-local label.

Local labels in macros are not visible outside the macro. 

# Constants
TODO

# Assembler directives

## Macros

Macros allow you to combine and parametrize often used sequences of operations and instantiate them at a later time. 
The macro's parameters are just text placeholders that will be replaced with argument's "raw" text when instantiated.
Together will conditional assembly (see below), they provide a powerful tool to simplify your assembler sources.

### Labels in macros
Macros can define local labels that are only valid within the macro. Only global labels and local labels that are 
defined in the macro (or passed as an argument) are visible. All local labels that were not passed in need to be 
resolved at the end of the macro.  

### `.macro`
The `.macro` directive starts the macro recording. The line needs to have a label that will be used as the macro's 
name. Optionally, it can have parameters. 

All lines follow the `.macro` line are copied into the macro buffer until a line with `.endm` is reached.

Macros can not be nested.

Example:
```
set16   .macro addr, val ; Set a 16bit value
        lda #<(val)   ; Low byte. Use ( and ) so that complicated expressions can be passed
        sta addr
        lda #>(val)   ; High  byte.
        sta addr+1
        .endm
```

## Conditional assembly
Conditional assembly is controlled by `.ifdef`, `.ifndef`, `.if`, `.else`, `.endif`.

## Other directives

### `.include`
TODO

### `.incbin`
TODO

### `.fail`
TODO

### `.equ`
TODO

### `.org`
TODO

### `.align`
TODO

### `.byte`
TODO

### `.word`
TODO

### `.float`
TODO

### `.reserve`
TODO

### `.cpu`
TODO

### `.platform`
TODO

### `.encoding`
TODO

### `.output`
TODO

# Syntax

```
line := [ident[":"]] [op] [";" comment]

op := ".macro" [ident {"," ident }]
    | ".endm"
    | ".ifdef" ident
    | ".ifndef" ident
    | ".if" expr [relOp expr]
    | ".else"
    | ".endif"
    | ".include" string
    | ".incbin" string [ "," expr ]
    | ".fail" string
    | ".equ" expr
    | ".org" expr
    | ".align" expr
    | ".byte" dbOp {"," dbOp }
    | ".float" expr {"," expr }
    | ".word" expr {"," expr }
    | ".reserve" expr ["," dbOp ]
    | ".cpu" string 
    | ".platform" string 
    | ".encoding" string
    | ".output" string
    | mnemonic [ param {"," param } ]
    | macroname [ actmacroparam {"," actmacroparam } ]
    .
                                         
mnemonic := ident .

macroname := ident .
                    
actmacroparam := ["#" ["<"|">"]] expr .

relOp := ["=" | "!=" | "<=" | "<" | ">=" | >"] .

dbOp := ("<"|">") expr 
      | expr
      .

basicDbOp := expr .

string := '"' { stringChar} '"' .

param := param6502 | paramZ80

param6502 := "#" ["<"|">"] expr
           | expr
           | expr "," "X"
           | expr "," "Y"  
           | "(" expr ")"
           | "(" expr "," "X" ")"
           | "(" expr "," "Y" ")"  
           | "(" expr ") ""," "X" 
           | "(" expr ")" "," "Y"
           .       

paramZ80 := ["<"|">"] expr
          | register
          | cond
          | "(" double-register ")"
          | "(" ["IX"|"IY"] ["+"|"-"] expr ")"
          | "(" expr ")"
          | expr
          .

expr := ["-"] term { "+"|"-"|"|" term } .
term := factor { "*"|"/"|"%"|"&"|"^" factor } . 
factor := "~" factor 
        | number 
        | char-const      
        | string
        | ident 
        | '*'
        | "(" expr ")" 
        | "scr" "(" expr ")" .
        .
number  := digit { digit } 
         | "%" binDigit { binDigit }
         | "&" octDigit { octDigit }
         | "$" hexDigit { hexDigit } .
ident := identChar { identChar | digit }.
identChar := "@" | "." | "_" | alpha .  
```

