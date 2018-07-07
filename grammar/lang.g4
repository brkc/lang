grammar lang;

block
  : statement*
  ;

statement
  : declaration
  | assignment
  | print
  | ifStatement
  | whileStatement
  ;

declaration
  : 'var' Id '=' booleanExpression ';'
  ;

assignment
  : Id '=' booleanExpression ';'
  ;

print
  : 'print' booleanExpression ';'
  ;

ifStatement
  : 'if' booleanExpression '{' block '}'
  ;

whileStatement
  : 'while' booleanExpression '{' block '}'
  ;

booleanExpression
  : andExpression ('||' andExpression)*
  ;

andExpression
  : condition ('&&' condition)*
  ;

condition
  : expression (('=='|'!='|'>='|'>'|'<'|'<=') expression)?
  ;

expression
  : term (('+'|'-') term)*
  ;

term
  : atom (('*'|'/') atom)*
  ;

logicalNotExpression
  : '!' logicalNotExpression
  | atom
  ;

atom
  : Id
  | Number
  | String
  | ('true'|'false')
  | '(' booleanExpression ')'
  ;

Id: [a-z]+;
Number: [0-9]+;
String: '"' (~["\r\n] | '\\"')* '"';
Whitespace: [ \t\r\n]+ -> skip;
