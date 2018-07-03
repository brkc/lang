grammar lang;

block
  : statement*
  ;

statement
  : assignment
  | print
  | ifStatement
  ;

assignment
  : 'let' Id '=' expression ';'
  ;

print
  : 'print' expression ';'
  ;

ifStatement
  : 'if' expression ('=='|'!='|'>='|'>'|'<'|'<=') expression '{' block '}'
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
