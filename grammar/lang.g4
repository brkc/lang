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
  : 'let' Id '=' booleanExpression ';'
  ;

print
  : 'print' booleanExpression ';'
  ;

ifStatement
  : 'if' booleanExpression '{' block '}'
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
  : String
  | mathExpression
  | ('true'|'false')
  | Id
  ;

mathExpression
  : term (('+'|'-') term)*
  ;

term
  : atom (('*'|'/') atom)*
  ;

atom
  : Id
  | Number
  | '(' booleanExpression ')'
  | logicalNotExpression
  ;

logicalNotExpression
  : '!' booleanExpression
  ;

Id: [a-z]+;
Number: [0-9]+;
String: '"' (~["\r\n] | '\\"')* '"';
Whitespace: [ \t\r\n]+ -> skip;
