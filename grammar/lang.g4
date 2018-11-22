grammar lang;

block
  : statement*
  ;

statement
  : declaration
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

ifStatement
  : 'if' booleanExpression '{' block '}'
  ;

whileStatement
  : 'while' booleanExpression '{' block '}'
  ;

functionStatement
  : 'fn' Id '(' (Id (',' Id)?)? ')' '{' block '}'
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

Id: [a-zA-Z_][a-zA-Z_0-9]*;
Number: [0-9]+;
String: '"' (~["\r\n] | '\\"')* '"';
Whitespace: [ \t\r\n]+ -> skip;
