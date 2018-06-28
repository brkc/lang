grammar lang;

root
  : statement*
  ;

statement
  : assignment
  | print
  ;

assignment
  : 'let' Id '=' expression '\n'
  ;

print
  : 'print' expression '\n'
  ;

expression
  : term (('+'|'-') term)*
  ;

term
  : atom (('*'|'/') atom)*
  ;

atom
  : Id
  | Number
  | '(' expression ')'
  ;

Id: [a-z]+;
Number: [0-9]+;
Whitespace: [ \t\r]+ -> skip;
