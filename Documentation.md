Overview
========

Labels
======
Labels need to terminate with ":" unless they start at the beginning of the line.

Local labels
------------
A local label is a label that starts with an underscore (`_`).
All local labels are only visible (and need to be resolved) before the next non-local label.

Local labels in macros are not visible outside the macro. 

Assembler directives
====================

Macro directives
--------------------
`.macro` and `.endm` are used to define macros.

Directives for conditional assembly
-----------------------------------
Conditional assembly is controlled by `.ifdef`, `.ifndef`, `.if`, `.else`, `.endif`.

`.include`
----------
TODO

`.incbin`
---------
TODO

`.fail`
-------
TODO

`.equ`
------
TODO

`.org`
------
TODO

`.align`
-------
TODO

`.byte`
-------
TODO

`.word`
-------
TODO

`.reserve`
----------
TODO

`.cpu`
------
TODO

`.platform`
-----------
TODO

Macros
======

- Can define local labels
- Only see the local labels that are defined in the macro
- all local labels need to be resolved at the end of the macro
- can refer to global labels

Conditional assembly
====================

TODO

Syntax
======

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
    | ".incbin" string
    | ".fail" string
    | ".equ" expr
    | ".org" expr
    | ".align" expr
    | ".byte" dbOp {"," dbOp }
    | ".word" expr {"," expr }
    | ".reserve" expr ["," dbOp ]
    | ".cpu" string 
    | ".platform" string 
    | mnemonic [ param {"," param } ]
    | macroname [ actmacroparam {"," actmacroparam } ]
.
                                         
mnemonic := ident .

macroname := ident .
                    
actmacroparam := ["#" ["<"|">"]] expr .

relOp := ["==" | "!=" | "<=" | "<" | ">=" | >"] .

dbOp := ("<"|">") expr 
      | basicDbOp
      | "scr" "(" basicDbOp { "," basicDbOp } ")" .

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

paramZ80 := ["<"|">"] expr
          | register
          | cond
          | "(" double-register ")"
          | "(" ["IX"|"IY"] ["+"|"-"] expr ")"
          | "(" expr ")"
          | expr

expr := ["-"] term { "+"|"-"|"|" term } .
term := factor { "*"|"/"|"%"|"&"|"^" factor } . 
factor := "~" factor 
        | number 
        | char-const      
        | string
        | ident 
        | '*'
        | "(" expr ")" .
number  := digit { digit } 
         | "%" binDigit { binDigit }
         | "&" octDigit { octDigit }
         | "$" hexDigit { hexDigit } .
ident := identChar { identChar | digit }.
identChar := "@" | "." | "_" | alpha .  
```

