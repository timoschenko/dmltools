// $antlr-format alignTrailingComments true, columnLimit 150, minEmptyLines 1, maxEmptyLinesToKeep 1, reflowComments false, useTab false


// $antlr-format allowShortRulesOnASingleLine false, allowShortBlocksOnASingleLine true, alignSemicolons hanging, alignColons hanging


/*
 * SPDX-License-Identifier: GPL-2.0-only
 * Copyright (C) 2026 Sergei Timoshchenko
 */
grammar dml;

dml
   : (dmlTag)? dmlCode EOF
   ;

dmlCode
   : (provisionalTag)? (deviceTag)? (bitorderTag)? deviceStatements
   ;

dmlTag
   : DML FLOAT_LITERAL SEMI
   {p.checkDmlVersion(localctx)}
   ;

provisionalTag
   : PROVISIONAL identList SEMI
   | PROVISIONAL SEMI
   | PROVISIONAL LPAREN identList RPAREN SEMI
   ;

deviceTag
   : DEVICE objident SEMI
   ;

bitorderTag
   : BITORDER ENDIAN SEMI
   ;

deviceStatements
   : deviceStatement*
   ;

exportAsExpression
   : expression
   ;

exportStatement
   : EXPORT expression AS exportAsExpression SEMI
   ;

importStatement
   : IMPORT STRING_LITERAL SEMI
   ;

externStatement
   : EXTERN (TYPEDEF)? cdecl SEMI
   ;

constantStatement
   : CONSTANT ident ASSIGN expression SEMI
   ;

logGroupStatement
   : LOGGROUP ident SEMI
   ;

typedefStatement
   : TYPEDEF cdecl SEMI
   ;

templateStatement
   : TEMPLATE objident (istemplateExpr)? LBRACE templateBody* RBRACE
   ;

templateBody
   : objectStatement
   | SHARED hook_decl
   ;

deviceStatement
   : templateStatement
   | cCodeBlock
   | logGroupStatement
   | constantStatement
   | typedefStatement
   | externStatement
   | importStatement
   | exportStatement
   | object
   | param
   | methodUnified
   | istemplateStatement
   | toplevelIf
   | errorStatement
   | ineachStatement
   ;

toplevelIfStatement
   : LBRACE deviceStatements RBRACE
   | deviceStatements
   ;

toplevelElseStatement
   : LBRACE deviceStatements RBRACE
   | toplevelIf
   ;

toplevelIf
   : HASH_IF LPAREN conditionExpression RPAREN toplevelIfStatement (HASH_ELSE toplevelElseStatement)?
   ;

array_list
   : (LBRACK arraydef RBRACK)+
   ;

object
   : object_decl
   | sessionOrSavedStatement
   | hook_decl
   ;
   // TODO: Use bitSliceExpr ?

bitrangeDecl
   : AT LBRACK expression (COLON expression)? RBRACK
   ;

object_type
   : ATTRIBUTE
   | BANK
   | CONNECT
   | EVENT
   | GROUP
   | IMPLEMENT
   | INTERFACE
   | PORT
   | SUBDEVICE
   ;

objectDescSpec
   : composed_string_literal
   ;

objectBody
   : (istemplateExpr)? (objectDescSpec)? SEMI
   | (istemplateExpr)? (objectDescSpec)? LBRACE objectStatement* RBRACE
   ;

object_decl
   : REGISTER objident (array_list)? (SIZE sizeSpec)? (offsetSpec)? objectBody
   | FIELD objident (array_list)? (bitrangeDecl)? objectBody
   | (IN)? object_type objident (array_list)? objectBody
   ;

global_decl_kind
   : SESSION
   | SAVED
   ;

local_decl_kind
   : LOCAL
   | SESSION
   | STATIC
   | SAVED
   ;

cvarDecl
   : cdecl (ASSIGN initializer)?
   | LPAREN cdecl_list RPAREN (ASSIGN initializer)?
   ;

sessionOrSavedDecl
   : global_decl_kind cvarDecl
   ;

sessionOrSavedStatement
   : sessionOrSavedDecl SEMI
   ;

variableDecl
   : local_decl_kind cvarDecl
   ;

variableStatement
   : variableDecl SEMI
   ;

methodOutParams
   : ARROW LPAREN (cdecl_list)? RPAREN
   ;

methodParamsMaybeUntyped
   : cdeclOrInline (COMMA cdeclOrInline)*
   ;

methodQualifiers
   : (SHARED)? INDEPENDENT (STARTUP (MEMOIZED)?)?
   | SHARED
   | INLINE
   ;

methodUnified
   : (methodQualifiers)? METHOD objident LPAREN (methodParamsMaybeUntyped)? RPAREN (methodOutParams)? (THROWS)? SEMI # methodUnifiedDeclaration
   | (methodQualifiers)? METHOD objident LPAREN (methodParamsMaybeUntyped)? RPAREN (methodOutParams)? (THROWS)? (DEFAULT | DEFAULT_NEW)? (LBRACE | LBRACE_NEW) statement* RBRACE # methodUnifiedDefinition
   ;

arraydef
   : ident LT expression
   | ident LT ELLIPSIS
   ;

cCodeBlock
   : (HEADER | HEADER2 | FOOTER) C_CODE_BLOCK
   ;

objectStatement
   : object
   | param
   | methodUnified
   | istemplateStatement
   | objectIf
   | errorStatement
   | ineachStatement
   ;

ineachStatement
   : IN EACH istemplateList LBRACE objectStatement* RBRACE
   ;

objectIfStatement
   : LBRACE objectStatement* RBRACE
   | objectStatement
   | /* nothing */


   ;

objectElseStatement
   : LBRACE objectStatement* RBRACE
   | objectStatement
   | /* nothing */


   ;

objectIf
   : HASH_IF LPAREN conditionExpression RPAREN objectIfStatement (HASH_ELSE objectElseStatement)?
   ;

param
   : PARAM objident (COLON ctypedecl)? ((NEW_ASSIGN | DEFAULT_NEW | ASSIGN | DEFAULT) expression)? SEMI
   | PARAM objident AUTO SEMI
   ;

istemplateList
   : objident
   | LPAREN objident (COMMA objident)* RPAREN
   ;

istemplateExpr
   : IS istemplateList
   ;

istemplateStatement
   : istemplateExpr SEMI
   ;

sizeSpec
   : expression
   ;

offsetSpec
   : AT expression
   ;

cdeclOrInline
   : cdecl
   | INLINE ident
   ;

cdecl
   : (CONST)? basetype cdecl_qualifier_ident_array_func
   ;

signed_or_unsigned_typeident
   : (UNSIGNED | SIGNED)? typeident
   ;

basetype
   : signed_or_unsigned_typeident
   | struct
   | union
   | layout
   | bitfields
   | typeof
   | SEQUENCE LPAREN signed_or_unsigned_typeident RPAREN
   | HOOK LPAREN (cdecl_list)? RPAREN
   ;

typeQualifier
   : CONST
   | MUL
   | VECT
   ;

cdecl_qualifier_ident_array_func
   : (typeQualifier)* cdecl_ident_array_func
   ;

cdecl_ident_array_func
   : ident # cdeclIdent
   | /* empty */

   # cdeclNone
   | cdecl_ident_array_func LBRACK expression RBRACK # cdeclArray
   | cdecl_ident_array_func LPAREN (cdecl_list_variadic)? RPAREN # cdeclFuncArgs
   | LPAREN cdecl_qualifier_ident_array_func RPAREN # cdeclFuncIdent
   ;

cdecl_list
   : cdecl (COMMA cdecl)*
   ;

cdecl_list_variadic
   : ELLIPSIS
   | cdecl_list (COMMA ELLIPSIS)?
   ;

typeof
   : TYPEOF expression
   ;

structDecls
   : cdecl (COLON expression)? SEMI
   ;

struct
   : STRUCT (ident)? LBRACE structDecls* RBRACE
   | STRUCT ident
   ;

union
   : UNION (ident)? LBRACE structDecls* RBRACE
   | UNION ident
   ;

layoutDecls
   : cdecl SEMI
   ;

layout
   : LAYOUT STRING_LITERAL LBRACE layoutDecls* RBRACE
   ;

bitfields
   : BITFIELDS INTEGER_LITERAL LBRACE bitfields_decls* RBRACE
   ;

bitfields_decls
   : cdecl bitrangeDecl SEMI
   ;

ctypedecl
   : (CONST)? basetype ctypedecl_ptr
   ;

ctypedecl_ptr
   : (MUL CONST?)* ctypedecl_simple
   ;

ctypedecl_simple
   : LPAREN ctypedecl_ptr RPAREN
   | /* empty */


   ;

typeident
   : CHAR
   | DOUBLE
   | FLOAT
   | INT
   | LONG
   | SHORT
   | SIGNED
   | UNSIGNED
   | VOID
   | REGISTER
   | ident
   ;

bitSlice
   : COLON expression (endianflag)?
   | endianflag
   ;
   // assignop is defined as a rule, used in statements

assignop
   : expression ADD_ASSIGN expression
   | expression SUB_ASSIGN expression
   | expression MUL_ASSIGN expression
   | expression DIV_ASSIGN expression
   | expression MOD_ASSIGN expression
   | expression OR_ASSIGN expression
   | expression AND_ASSIGN expression
   | expression XOR_ASSIGN expression
   | expression LSHIFT_ASSIGN expression
   | expression RSHIFT_ASSIGN expression
   ;
   // Expression rule.

expression
   : expression ((DOT | ARROW) expressionAccessIdent)+ # memberAccessExpr
   | expression LBRACK expression bitSlice RBRACK # bitSliceExpr
   | expression LBRACK expression RBRACK # arrayAccessExpr
   | expression (argumentsListEmpty | argumentsList) # functionCallExpr
   | expression op = (INC | DEC) # postIncrementExpr
   | op = (INC | DEC) expression # preIncrementExpr
   | op = (SUB | ADD) expression # unaryExpr
   | BANG expression # notExpr
   | TILDE expression # bitwiseNotExpr
   | BITAND expression # addressofExpr
   | MUL expression # derefExpr
   | CAST LPAREN expression COMMA ctypedecl RPAREN # castExpr
   | SIZEOF expression # sizeofExpr
   | DEFINED expression # definedExpr
   | STRINGIFY LPAREN expression RPAREN # stringifyExpr
   | < assoc = left > expression op = (MUL | DIV | MOD) expression # multiplicativeExpr
   | < assoc = left > expression op = (ADD | SUB) expression # additiveExpr
   | < assoc = left > expression op = (LSHIFT | RSHIFT) expression # shiftExpr
   | < assoc = left > expression op = (LT | GT | LE | GE) expression # relationalExpr
   | < assoc = left > expression op = (EQUAL | NOTEQUAL) expression # equalityExpr
   | < assoc = left > expression op = (BITAND | CARET | BITOR) expression # bitwiseExpr
   | < assoc = left > expression op = (AND | OR) expression # logicalExpr
   | < assoc = right > expression QUESTION expression COLON expression # ternaryIfExpr
   | < assoc = right > expression ALT_QUESTION expression ALT_COLON expression # ternaryConstIfExpr
   | NEW ctypedecl (LBRACK expression RBRACK)? # newExpr
   | EACH objident IN LPAREN expression RPAREN # eachExpr
   | SIZEOFTYPE typeoparg # sizeoftypeExpr
   | LPAREN expression RPAREN # parenExpr
   | LBRACK expression_list RBRACK # arrayLiteralExpr
   | expressionIdent # literalAtom
   ;

numericLiteral
   : INTEGER_LITERAL
   | HEX_LITERAL
   | BINARY_LITERAL
   | CHAR_LITERAL
   | FLOAT_LITERAL
   ;

stringLiteral
   : STRING_LITERAL
   ;

expressionAccessIdent
   : DEFAULT
   | THIS
   | TEMPLATE
   | objident
   ;

expressionIdent
   : numericLiteral
   | stringLiteral
   | UNDEFINED
   | expressionAccessIdent
   ;

utf8_sconst
   : STRING_LITERAL
   ;

typeoparg
   : ctypedecl
   | LPAREN ctypedecl RPAREN
   ;

endianflag
   : COMMA ident
   ;

expression_list
   : (expression_list_ntc_nonempty)?
   ;

expression_list_ntc_nonempty
   : expression (COMMA expression)* (COMMA)?
   ;

composed_string_literal
   : utf8_sconst (ADD utf8_sconst)*
   ;

bracketed_string_literal
   : composed_string_literal
   | LPAREN composed_string_literal RPAREN
   ;

singleInitializer
   : expression
   | cdeclInitializer
   ;

designatedStructInitializer
   : DOT ident ASSIGN singleInitializer
   ;

designated_struct_initializer_list
   : designatedStructInitializer (COMMA designatedStructInitializer)* (COMMA ELLIPSIS)? (COMMA)?
   ;

cdeclInitializer
   : LBRACE singleInitializerArray RBRACE
   | LBRACE designated_struct_initializer_list RBRACE
   ;

singleInitializerArray
   : singleInitializer (COMMA singleInitializer)* (COMMA)?
   ;

initializer
   : singleInitializer
   | argumentsList
   ;

argumentsList
   : LPAREN singleInitializerArray RPAREN
   ;

argumentsListEmpty
   : LPAREN RPAREN
   ;

statement
   : SEMI // noop
   | afterStatement
   | assertStatement
   | assignop SEMI
   | assignStatement SEMI
   | breakStatement
   | compoundStatement
   | continueStatement
   | deleteStatement
   | doWhileLoopStatement
   | errorStatement
   | expression SEMI
   | foreachStatement
   | forLoopStatement
   | gotoStatement
   | ifStatement
   | logStatement
   | returnStatement
   | selectStatement
   | switchStatement
   | throwStatement
   | tryStatement
   | variableStatement
   | warningStatement
   | whileLoopStatement
   ;

throwStatement
   : THROW SEMI
   ;

continueStatement
   : CONTINUE SEMI
   ;

breakStatement
   : BREAK SEMI
   ;

elseStatement
   : statement
   ;

loopStatement
   : statement
   ;

conditionExpression
   : expression
   ;

assertStatement
   : ASSERT conditionExpression SEMI
   ;

deleteStatement
   : DELETE expression SEMI
   ;

returnStatement
   : RETURN (initializer)? SEMI
   ;

gotoStatement
   : GOTO ident SEMI
   ;
   // log {{{

logKind
   : INFO
   | ERROR
   | SPEC_VIOL
   | WARNING
   | UNIMPL
   | CRITICAL
   | IDENTIFIER
   | utf8_sconst
   ;

logGroup
   : expression
   ;

logFormatString
   : bracketed_string_literal
   ;

thenLevel
   : expression
   ;

logLevel
   : expression (THEN thenLevel)?
   ;

logStatement
   : LOG logKind (COMMA logLevel)? (COMMA logGroup)? COLON logFormatString (COMMA expression_list_ntc_nonempty)? SEMI
   ;
   // }}}

   // Switch-Case {{{

caseOrDefaultStatement
   : CASE expression
   | DEFAULT
   ;

caseStatement
   : caseOrDefaultStatement COLON statement*
   | switchBodyHashIfStatement statement*
   ;

switchBodyStatement
   : caseStatement*
   ;

switchBodyHashElseStatement
   : switchBodyStatement
   ;
   // TOOD: Find a better way?

switchBodyHashIfStatement /* TODO Rework before implementation */

   : HASH_IF LPAREN conditionExpression RPAREN LBRACE switchBodyStatement RBRACE (HASH_ELSE LBRACE switchBodyHashElseStatement RBRACE)?
   | HASH_IF LPAREN conditionExpression RPAREN LBRACE switchBodyStatement RBRACE (HASH_ELSE switchBodyHashElseStatement)?
   | HASH_IF LPAREN conditionExpression RPAREN switchBodyStatement (HASH_ELSE LBRACE switchBodyHashElseStatement RBRACE)?
   | HASH_IF LPAREN conditionExpression RPAREN switchBodyStatement (HASH_ELSE switchBodyHashElseStatement)?
   ;

switchStatement
   : SWITCH LPAREN expression RPAREN LBRACE switchBodyStatement RBRACE
   ;
   // }}}

   // For-loop {{{

forLoopStepStatement
   : innerForLoopStatement (COMMA innerForLoopStatement)*
   ;

innerForLoopStatement
   : assignStatement
   | assignop
   | expression
   ;

forLoopInitStatement
   : variableDecl
   | innerForLoopStatement (COMMA innerForLoopStatement)*
   ;

forLoopStatement
   : FOR LPAREN (forLoopInitStatement)? SEMI (conditionExpression)? SEMI (forLoopStepStatement)? RPAREN loopStatement
   ;
   // }}}

whileLoopStatement
   : WHILE LPAREN conditionExpression RPAREN loopStatement
   ;

doWhileLoopStatement
   : DO loopStatement WHILE LPAREN conditionExpression RPAREN SEMI
   ;

foreachStatement
   : FOREACH ident IN LPAREN expression RPAREN loopStatement
   | HASH_FOREACH ident IN LPAREN expression RPAREN loopStatement
   ;

ifStatement
   : IF LPAREN conditionExpression RPAREN statement (ELSE elseStatement)?
   | HASH_IF LPAREN conditionExpression RPAREN statement (HASH_ELSE elseStatement)?
   ;

catchStatement
   : statement
   ;

tryStatement
   : TRY statement CATCH catchStatement
   ;

selectStatement
   : HASH_SELECT ident IN LPAREN expression RPAREN WHERE LPAREN conditionExpression RPAREN loopStatement (HASH_ELSE elseStatement)?
   | SELECT ident IN LPAREN expression RPAREN WHERE LPAREN conditionExpression RPAREN loopStatement (ELSE elseStatement)?
   ;

assignStatement
   : assign_chain
   | tuple_literal ASSIGN initializer
   ;

assign_chain
   : expression ASSIGN (assign_chain | initializer)
   ;

afterStatement
   : AFTER (afterDelayCondition)? COLON expression SEMI
   ;

afterDelayCondition
   : expression IDENTIFIER
   | expression ARROW LPAREN (identList)? RPAREN
   | expression ARROW ident
   | expression
   ;

tuple_literal
   : LPAREN expression COMMA expression_list_ntc_nonempty RPAREN
   ;

identList
   : ident (COMMA ident)*
   ;

errorStatement
   : ERROR (bracketed_string_literal)? SEMI
   ;

warningStatement
   : WARNING bracketed_string_literal SEMI
   ;

compoundStatement
   : LBRACE statement* RBRACE
   ;

hook_decl
   : HOOK LPAREN (cdecl_list)? RPAREN ident (LBRACK expression RBRACK)* SEMI
   ;

objident
   : REGISTER
   | ident
   ;
   // This rule allows keywords to be used as identifiers,

   // as defined in the original grammar.

ident
   : ATTRIBUTE
   | '_'
   | ASYNC
   | AUTO
   | AWAIT
   | BANK
   | BITORDER
   | CALL
   | CLASS
   | CONNECT
   | CONSTANT
   | CRITICAL
   | DATA
   | DEVICE
   | ENDIAN
   | ENUM
   | EVENT
   | FIELD
   | FOOTER
   | GROUP
   | HEADER
   | HEADER2
   | IMPLEMENT
   | IMPORT
   | INDEPENDENT
   | INFO
   | INTERFACE
   | LOGGROUP
   | MEMOIZED
   | METHOD
   | NAMESPACE
   | NOTHROW
   | PARAM
   | PORT
   | PRIVATE
   | PROTECTED
   | PROVISIONAL
   | PUBLIC
   | RESTRICT
   | SAVED
   | SELECT
   | SIZE
   | SPEC_VIOL
   | STARTUP
   | STATIC
   | SUBDEVICE
   | THEN
   | THROWS
   | UNIMPL
   | UNION
   | USING
   | VIRTUAL
   | VOLATILE
   | WITH
   | IDENTIFIER /* must come after all keywords */


   ;
/*
 * =============================================================================
 * LEXER RULES
 * =============================================================================
 */


   // --- Keywords ---

   // All keywords from the 'ident' rule are defined here.

   // This ensures they are tokenized as keywords, not IDENTIFIERs.

AFTER
   : 'after'
   ;

AS
   : 'as'
   ;

ASSERT
   : 'assert'
   ;

ASYNC
   : 'async'
   ;

ATTRIBUTE
   : 'attribute'
   ;

AUTO
   : 'auto'
   ;

AWAIT
   : 'await'
   ;

BANK
   : 'bank'
   ;

BITFIELDS
   : 'bitfields'
   ;

BITORDER
   : 'bitorder'
   ;

BREAK
   : 'break'
   ;

CALL
   : 'call'
   ;

CASE
   : 'case'
   ;

CAST
   : 'cast'
   ;

CATCH
   : 'catch'
   ;

CHAR
   : 'char'
   ;

CLASS
   : 'class'
   ;

CONNECT
   : 'connect'
   ;

CONST
   : 'const'
   ;

CONSTANT
   : 'constant'
   ;

CONTINUE
   : 'continue'
   ;

CRITICAL
   : 'critical'
   ;

DATA
   : 'data'
   ;

DEFAULT
   : 'default'
   ;

DEFAULT_NEW
   : ':default'
   ;

DEFINED
   : 'defined'
   ;

DELETE
   : 'delete'
   ;

DEVICE
   : 'device'
   ;

DML
   : 'dml'
   ;

DO
   : 'do'
   ;

DOUBLE
   : 'double'
   ;

EACH
   : 'each'
   ;

ELSE
   : 'else'
   ;

ENDIAN
   : ENDIAN_BE
   | ENDIAN_LE
   ;

fragment ENDIAN_BE
   : 'be'
   ;

fragment ENDIAN_LE
   : 'le'
   ;

ENUM
   : 'enum'
   ;

ERROR
   : 'error'
   ;

EVENT
   : 'event'
   ;

EXPORT
   : 'export'
   ;

EXTERN
   : 'extern'
   ;

FIELD
   : 'field'
   ;

FLOAT
   : 'float'
   ;

FOOTER
   : 'footer'
   ;

FOR
   : 'for'
   ;

FOREACH
   : 'foreach'
   ;

GOTO
   : 'goto'
   ;

GROUP
   : 'group'
   ;

HASH_ELSE
   : '#else'
   ;

HASH_FOREACH
   : '#foreach'
   ;

HASH_IF
   : '#if'
   ;

HASH_SELECT
   : '#select'
   ;

HEADER
   : 'header'
   ;

HEADER2
   : '_header'
   ;

HOOK
   : 'hook'
   ;

IF
   : 'if'
   ;

IMPLEMENT
   : 'implement'
   ;

IMPORT
   : 'import'
   ;

IN
   : 'in'
   ;

INDEPENDENT
   : 'independent'
   ;

INFO
   : 'info'
   ;

INLINE
   : 'inline'
   ;

INT
   : 'int'
   ;

INTERFACE
   : 'interface'
   ;

IS
   : 'is'
   ;

LAYOUT
   : 'layout'
   ;

LOCAL
   : 'local'
   ;

LOG
   : 'log'
   ;

LOGGROUP
   : 'loggroup'
   ;

LONG
   : 'long'
   ;

MEMOIZED
   : 'memoized'
   ;

METHOD
   : 'method'
   ;

NAMESPACE
   : 'namespace'
   ;

NEW
   : 'new'
   ;

NOTHROW
   : 'nothrow'
   ;

PARAM
   : 'param'
   ;

PORT
   : 'port'
   ;

PRIVATE
   : 'private'
   ;

PROTECTED
   : 'protected'
   ;

PROVISIONAL
   : 'provisional'
   ;

PUBLIC
   : 'public'
   ;

REGISTER
   : 'register'
   ;

RESTRICT
   : 'restrict'
   ;

RETURN
   : 'return'
   ;

SAVED
   : 'saved'
   ;

SELECT
   : 'select'
   ;

SEQUENCE
   : 'sequence'
   ;

SESSION
   : 'session'
   ;

SHARED
   : 'shared'
   ;

SHORT
   : 'short'
   ;

SIGNED
   : 'signed'
   ;

SIZE
   : 'size'
   ;

SIZEOF
   : 'sizeof'
   ;

SIZEOFTYPE
   : 'sizeoftype'
   ;

SPEC_VIOL
   : 'spec_viol'
   ;

STARTUP
   : 'startup'
   ;

STATIC
   : 'static'
   ;

STRINGIFY
   : 'stringify'
   ;

STRUCT
   : 'struct'
   ;

SUBDEVICE
   : 'subdevice'
   ;

SWITCH
   : 'switch'
   ;

TEMPLATE
   : 'template'
   ;

THEN
   : 'then'
   ;

THIS
   : 'this'
   ;

THROW
   : 'throw'
   ;

THROWS
   : 'throws'
   ;

TRY
   : 'try'
   ;

TYPEDEF
   : 'typedef'
   ;

TYPEOF
   : 'typeof'
   ;

UNDEFINED
   : 'undefined'
   ;

UNIMPL
   : 'unimpl'
   ;

UNION
   : 'union'
   ;

UNSIGNED
   : 'unsigned'
   ;

USING
   : 'using'
   ;

VECT
   : 'vect'
   ;

VIRTUAL
   : 'virtual'
   ;

VOID
   : 'void'
   ;

VOLATILE
   : 'volatile'
   ;

WARNING
   : WARNINGA
   | WARNINGB
   ;

fragment WARNINGA
   : 'warning'
   ;

fragment WARNINGB
   : '_warning'
   ;

WHERE
   : 'where'
   ;

WHILE
   : 'while'
   ;

WITH
   : 'with'
   ;
   // --- Literals ---

   // These are simplified regexes. You may need to refine them,

   // especially for float and char/string escapes.

INTEGER_LITERAL
   : [0-9] [0-9_]*
   ;

HEX_LITERAL
   : '0' [xX] [_0-9a-fA-F]+
   ;

BINARY_LITERAL
   : '0' [bB] [_01]+
   ;

FLOAT_LITERAL
   : [0-9]+ DOT [0-9]+ ([eE] [+-]? [0-9]+)?
   | DOT [0-9]+ ([eE] [+-]? [0-9]+)?
   | [0-9]+ [eE] [+-]? [0-9]+
   ;

CHAR_LITERAL
   : '\'' (~ ('\'' | '\\') | ('\\' .)) '\''
   ;

STRING_LITERAL
   : '"' (~ ('"' | '\\') | ('\\' .))* '"'
   ;
   // This is an assumption for the `"%{ ... %}"` blocks.

   // It non-greedily matches everything between '%{' and '%}'.

   // This will fail if the content contains nested '%}'

   // A more robust way uses lexer modes, but is more complex.

C_CODE_BLOCK
   : '%{' .*? '%}'
   ;
   // --- Identifier ---

   // Must come *after* all keywords

IDENTIFIER
   : [a-zA-Z_] [a-zA-Z0-9_]*
   ;
   // --- Symbols ---

   // Operators and Punctuation

LPAREN
   : '('
   ;

RPAREN
   : ')'
   ;

LBRACE
   : '{'
   ;

LBRACE_NEW
   : ':{'
   ;

RBRACE
   : '}'
   ;

LBRACK
   : '['
   ;

RBRACK
   : ']'
   ;

SEMI
   : ';'
   ;

COMMA
   : ','
   ;

DOT
   : '.'
   ;

ELLIPSIS
   : '...'
   ;

ASSIGN
   : '='
   ;

NEW_ASSIGN
   : ':='
   ;

GT
   : '>'
   ;

LT
   : '<'
   ;

BANG
   : '!'
   ;

TILDE
   : '~'
   ;

QUESTION
   : '?'
   ;

COLON
   : ':'
   ;

EQUAL
   : '=='
   ;

LE
   : '<='
   ;

GE
   : '>='
   ;

NOTEQUAL
   : '!='
   ;

AND
   : '&&'
   ;

OR
   : '||'
   ;

INC
   : '++'
   ;

DEC
   : '--'
   ;

ADD
   : '+'
   ;

SUB
   : '-'
   ;

MUL
   : '*'
   ;

DIV
   : '/'
   ;

BITAND
   : '&'
   ;

BITOR
   : '|'
   ;

CARET
   : '^'
   ;

MOD
   : '%'
   ;

ADD_ASSIGN
   : '+='
   ;

SUB_ASSIGN
   : '-='
   ;

MUL_ASSIGN
   : '*='
   ;

DIV_ASSIGN
   : '/='
   ;

MOD_ASSIGN
   : '%='
   ;

AND_ASSIGN
   : '&='
   ;

OR_ASSIGN
   : '|='
   ;

XOR_ASSIGN
   : '^='
   ;

LSHIFT_ASSIGN
   : '<<='
   ;

RSHIFT_ASSIGN
   : '>>='
   ;

LSHIFT
   : '<<'
   ;

RSHIFT
   : '>>'
   ;

ARROW
   : '->'
   ;

AT
   : '@'
   ;

ALT_QUESTION
   : '#?'
   ;

ALT_COLON
   : '#:'
   ;
   // --- Whitespace and Comments ---

   // These are skipped and do not produce tokens for the parser.

WS
   : [ \t\r\n]+ -> skip
   ;

fragment COMMENT_BEGIN
   : '//'
   ;

fragment COMMENT_INLINE_BEGIN
   : '/*'
   ;

fragment COMMENT_INLINE_END
   : '*/'
   ;

COMMENT
   : COMMENT_BEGIN .*? ('\n' | EOF) -> channel (HIDDEN)
   ;

COMMENT_BLOCK_EOL
   : COMMENT_INLINE_BEGIN (~ [\r\n*] | MUL+ ~ [\r\n/*])* MUL* COMMENT_INLINE_END [ \t]* '\r'? '\n' -> channel (HIDDEN)
   ;

COMMENT_BLOCK
   : COMMENT_INLINE_BEGIN (~ [\r\n*] | MUL+ ~ [\r\n/*])* MUL* COMMENT_INLINE_END -> channel (HIDDEN)
   ;

COMMENT_INLINE_EOL
   : COMMENT_INLINE_BEGIN (~ [*] | MUL+ ~ [/*])* MUL* COMMENT_INLINE_END [ \t]* '\r'? '\n' -> channel (HIDDEN)
   ;

COMMENT_INLINE
   : COMMENT_INLINE_BEGIN (~ [*] | MUL+ ~ [/*])* MUL* COMMENT_INLINE_END -> channel (HIDDEN)
   ;
