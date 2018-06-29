grammar lang;

root
  : statement*
  ;

statement
  : assignment
  | print
  ;

assignment
  : 'let' Id '=' expression ';'
  ;

print
  : 'print' expression ';'
  ;

expression
  : String
  | mathExpression;

mathExpression
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
String: '"' (~["\r\n] | '\\"')* '"';
Whitespace: [ \t\r\n]+ -> skip;
