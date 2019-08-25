# docopt Grammar

## start parsing with participle

It seems that I will need to Capture some pattern
Aloas looking Basic parsing example using Line delimiter.

## grammar

Comment
Usage: Usage_Program
LongBlank String Usage_Expr EOL

Options:


Usage_Expr := LongBlank String Usage_def EOL
Usage_def  := 
  // expr ::= seq ( '|' seq )* ;
  // seq ::= ( atom [ '...' ] )* ;
  // atom ::= '(' expr ')' | '[' expr ']' | 'options' | long | shorts | argument | command ;
  // long ::= '--' chars [ ( ' ' | '=' ) chars ] ;

## print docopt input (usage) AST - DONE

use `--print-ast` for printing AST parsed from Usage given with `-h`

