grammar lang;

block
  : statement*
  ;

statement
  : declaration
  | print
  | ifStatement
  | whileStatement
  | functionStatement
  | returnStatement
  | assignment
  | callExpression
  ;

declaration
  : 'var' Id '=' booleanExpression ';'
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

functionStatement
  : 'def' Id '(' (Id (',' Id)?)? ')' '{' block '}'
  ;

returnStatement
  : 'return' booleanExpression ';'
  ;

assignment
  : Id '=' booleanExpression ';'
  ;

callExpression
  : Id '(' (booleanExpression (',' booleanExpression)?)? ')'
  ;

booleanExpression
  : andExpression ('or' andExpression)*
  ;

andExpression
  : condition ('and' condition)*
  ;

condition
  : logicalOperand (('=='|'!='|'>='|'>'|'<'|'<=') logicalOperand)?
  ;

logicalOperand
  : term (('+'|'-') term)*
  ;

term
  : logicalNotExpression (('*'|'/') logicalNotExpression)*
  ;

logicalNotExpression
  : 'not' logicalNotExpression
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
