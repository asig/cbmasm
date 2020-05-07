Labels
======
Labels need to terminate with ":" unless they start at the beginning of the line.

Local labels
------------
A local label is a label that starts with an underscore (`_`).
All local labels are only visible (and need to be resolved) before the next non-local label. 

Syntax
======

line := [ident[":"]] [op] [";" comment]


op := ".macro" [ident {"," ident }]
    | ".mend"
    | ".equ" expr
    | ".org" expr
    | ".byte" dbOp {"," dbOp }
    | ".word" dbOp {"," dbOp }
    | ".reserve" expr ["," dbOp ] 
    | ident [ param {"," param } ].

dbOp := ["<"|">"] expr | string .

string := '"' { stringChar} '"'.

param := "#" ["<"|">"] expr
       | expr``
       | expr "," "X"
       | expr "," "Y"  
       | "(" expr ")"
       | "(" expr "," "X" ")"
       | "(" expr "," "Y" ")"  
       | "(" expr ") ""," "X" 
       | "(" expr ")" "," "Y"       


expr := ["-"] term { "+"|"-"|"|" term } .
term := factor { "*"|"/"|"%"|"&"|"^" factor } . 
factor := "~" factor 
        | number 
        | ident 
        | '*'
        | "(" expr ")" .
number  := digit { digit } 
         | "%" binDigit { binDigit }
         | "&" octDigit { octDigit }
         | "$" hexDigit { hexDigit } .
ident := identChar { identChar | digit }.
identChar := "@" | "." | "_" | alpha .  
 





