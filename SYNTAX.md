labels need to start at the beginning of a line

line := [ident] [op] [";" comment]


op := ".macro" [ident {"," ident }]
    | ".mend"
    | ".equ" expr
    | ".org" expr
    | ".db" dbOp {"," dbOp }
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

expr := ["-"] factor { "*"|"/"|"%"|"&"|"^" factor } .
factor  := term { "+"|"-"|"|" term } .
term    := "~" term | number | ident | '$'.
number  := digit { digit } 
         | "%" binDigit { binDigit }
         | "&" octDigit { octDigit }
         | "$" hexDigit { hexDigit } .
ident := identChar { identChar | digit }.
identChar := "@" | "." | "_" | alpha .  
 



