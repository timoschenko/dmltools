// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex // dml
import "github.com/antlr4-go/antlr/v4"

// dmlListener is a complete listener for a parse tree produced by dmlParser.
type dmlListener interface {
	antlr.ParseTreeListener

	// EnterDml is called when entering the dml production.
	EnterDml(c *DmlContext)

	// EnterDmlCode is called when entering the dmlCode production.
	EnterDmlCode(c *DmlCodeContext)

	// EnterDmlTag is called when entering the dmlTag production.
	EnterDmlTag(c *DmlTagContext)

	// EnterProvisionalTag is called when entering the provisionalTag production.
	EnterProvisionalTag(c *ProvisionalTagContext)

	// EnterDeviceTag is called when entering the deviceTag production.
	EnterDeviceTag(c *DeviceTagContext)

	// EnterBitorderTag is called when entering the bitorderTag production.
	EnterBitorderTag(c *BitorderTagContext)

	// EnterDeviceStatements is called when entering the deviceStatements production.
	EnterDeviceStatements(c *DeviceStatementsContext)

	// EnterExportAsExpression is called when entering the exportAsExpression production.
	EnterExportAsExpression(c *ExportAsExpressionContext)

	// EnterExportStatement is called when entering the exportStatement production.
	EnterExportStatement(c *ExportStatementContext)

	// EnterImportStatement is called when entering the importStatement production.
	EnterImportStatement(c *ImportStatementContext)

	// EnterExternStatement is called when entering the externStatement production.
	EnterExternStatement(c *ExternStatementContext)

	// EnterConstantStatement is called when entering the constantStatement production.
	EnterConstantStatement(c *ConstantStatementContext)

	// EnterLogGroupStatement is called when entering the logGroupStatement production.
	EnterLogGroupStatement(c *LogGroupStatementContext)

	// EnterTypedefStatement is called when entering the typedefStatement production.
	EnterTypedefStatement(c *TypedefStatementContext)

	// EnterTemplateStatement is called when entering the templateStatement production.
	EnterTemplateStatement(c *TemplateStatementContext)

	// EnterTemplateBody is called when entering the templateBody production.
	EnterTemplateBody(c *TemplateBodyContext)

	// EnterDeviceStatement is called when entering the deviceStatement production.
	EnterDeviceStatement(c *DeviceStatementContext)

	// EnterToplevelIfStatement is called when entering the toplevelIfStatement production.
	EnterToplevelIfStatement(c *ToplevelIfStatementContext)

	// EnterToplevelElseStatement is called when entering the toplevelElseStatement production.
	EnterToplevelElseStatement(c *ToplevelElseStatementContext)

	// EnterToplevelIf is called when entering the toplevelIf production.
	EnterToplevelIf(c *ToplevelIfContext)

	// EnterArray_list is called when entering the array_list production.
	EnterArray_list(c *Array_listContext)

	// EnterObject is called when entering the object production.
	EnterObject(c *ObjectContext)

	// EnterBitrangeDecl is called when entering the bitrangeDecl production.
	EnterBitrangeDecl(c *BitrangeDeclContext)

	// EnterObject_type is called when entering the object_type production.
	EnterObject_type(c *Object_typeContext)

	// EnterObjectDescSpec is called when entering the objectDescSpec production.
	EnterObjectDescSpec(c *ObjectDescSpecContext)

	// EnterObjectBody is called when entering the objectBody production.
	EnterObjectBody(c *ObjectBodyContext)

	// EnterObject_decl is called when entering the object_decl production.
	EnterObject_decl(c *Object_declContext)

	// EnterGlobal_decl_kind is called when entering the global_decl_kind production.
	EnterGlobal_decl_kind(c *Global_decl_kindContext)

	// EnterLocal_decl_kind is called when entering the local_decl_kind production.
	EnterLocal_decl_kind(c *Local_decl_kindContext)

	// EnterCvarDecl is called when entering the cvarDecl production.
	EnterCvarDecl(c *CvarDeclContext)

	// EnterSessionOrSavedDecl is called when entering the sessionOrSavedDecl production.
	EnterSessionOrSavedDecl(c *SessionOrSavedDeclContext)

	// EnterSessionOrSavedStatement is called when entering the sessionOrSavedStatement production.
	EnterSessionOrSavedStatement(c *SessionOrSavedStatementContext)

	// EnterVariableDecl is called when entering the variableDecl production.
	EnterVariableDecl(c *VariableDeclContext)

	// EnterVariableStatement is called when entering the variableStatement production.
	EnterVariableStatement(c *VariableStatementContext)

	// EnterMethodOutParams is called when entering the methodOutParams production.
	EnterMethodOutParams(c *MethodOutParamsContext)

	// EnterMethodParamsMaybeUntyped is called when entering the methodParamsMaybeUntyped production.
	EnterMethodParamsMaybeUntyped(c *MethodParamsMaybeUntypedContext)

	// EnterMethodQualifiers is called when entering the methodQualifiers production.
	EnterMethodQualifiers(c *MethodQualifiersContext)

	// EnterMethodUnifiedDeclaration is called when entering the methodUnifiedDeclaration production.
	EnterMethodUnifiedDeclaration(c *MethodUnifiedDeclarationContext)

	// EnterMethodUnifiedDefinition is called when entering the methodUnifiedDefinition production.
	EnterMethodUnifiedDefinition(c *MethodUnifiedDefinitionContext)

	// EnterArraydef is called when entering the arraydef production.
	EnterArraydef(c *ArraydefContext)

	// EnterCCodeBlock is called when entering the cCodeBlock production.
	EnterCCodeBlock(c *CCodeBlockContext)

	// EnterObjectStatement is called when entering the objectStatement production.
	EnterObjectStatement(c *ObjectStatementContext)

	// EnterIneachStatement is called when entering the ineachStatement production.
	EnterIneachStatement(c *IneachStatementContext)

	// EnterObjectIfStatement is called when entering the objectIfStatement production.
	EnterObjectIfStatement(c *ObjectIfStatementContext)

	// EnterObjectElseStatement is called when entering the objectElseStatement production.
	EnterObjectElseStatement(c *ObjectElseStatementContext)

	// EnterObjectIf is called when entering the objectIf production.
	EnterObjectIf(c *ObjectIfContext)

	// EnterParam is called when entering the param production.
	EnterParam(c *ParamContext)

	// EnterIstemplateList is called when entering the istemplateList production.
	EnterIstemplateList(c *IstemplateListContext)

	// EnterIstemplateExpr is called when entering the istemplateExpr production.
	EnterIstemplateExpr(c *IstemplateExprContext)

	// EnterIstemplateStatement is called when entering the istemplateStatement production.
	EnterIstemplateStatement(c *IstemplateStatementContext)

	// EnterSizeSpec is called when entering the sizeSpec production.
	EnterSizeSpec(c *SizeSpecContext)

	// EnterOffsetSpec is called when entering the offsetSpec production.
	EnterOffsetSpec(c *OffsetSpecContext)

	// EnterCdeclOrInline is called when entering the cdeclOrInline production.
	EnterCdeclOrInline(c *CdeclOrInlineContext)

	// EnterCdecl is called when entering the cdecl production.
	EnterCdecl(c *CdeclContext)

	// EnterSigned_or_unsigned_typeident is called when entering the signed_or_unsigned_typeident production.
	EnterSigned_or_unsigned_typeident(c *Signed_or_unsigned_typeidentContext)

	// EnterBasetype is called when entering the basetype production.
	EnterBasetype(c *BasetypeContext)

	// EnterTypeQualifier is called when entering the typeQualifier production.
	EnterTypeQualifier(c *TypeQualifierContext)

	// EnterCdecl_qualifier_ident_array_func is called when entering the cdecl_qualifier_ident_array_func production.
	EnterCdecl_qualifier_ident_array_func(c *Cdecl_qualifier_ident_array_funcContext)

	// EnterCdeclIdent is called when entering the cdeclIdent production.
	EnterCdeclIdent(c *CdeclIdentContext)

	// EnterCdeclFuncIdent is called when entering the cdeclFuncIdent production.
	EnterCdeclFuncIdent(c *CdeclFuncIdentContext)

	// EnterCdeclArray is called when entering the cdeclArray production.
	EnterCdeclArray(c *CdeclArrayContext)

	// EnterCdeclNone is called when entering the cdeclNone production.
	EnterCdeclNone(c *CdeclNoneContext)

	// EnterCdeclFuncArgs is called when entering the cdeclFuncArgs production.
	EnterCdeclFuncArgs(c *CdeclFuncArgsContext)

	// EnterCdecl_list is called when entering the cdecl_list production.
	EnterCdecl_list(c *Cdecl_listContext)

	// EnterCdecl_list_variadic is called when entering the cdecl_list_variadic production.
	EnterCdecl_list_variadic(c *Cdecl_list_variadicContext)

	// EnterTypeof is called when entering the typeof production.
	EnterTypeof(c *TypeofContext)

	// EnterStructDecls is called when entering the structDecls production.
	EnterStructDecls(c *StructDeclsContext)

	// EnterStruct is called when entering the struct production.
	EnterStruct(c *StructContext)

	// EnterUnion is called when entering the union production.
	EnterUnion(c *UnionContext)

	// EnterLayoutDecls is called when entering the layoutDecls production.
	EnterLayoutDecls(c *LayoutDeclsContext)

	// EnterLayout is called when entering the layout production.
	EnterLayout(c *LayoutContext)

	// EnterBitfields is called when entering the bitfields production.
	EnterBitfields(c *BitfieldsContext)

	// EnterBitfields_decls is called when entering the bitfields_decls production.
	EnterBitfields_decls(c *Bitfields_declsContext)

	// EnterCtypedecl is called when entering the ctypedecl production.
	EnterCtypedecl(c *CtypedeclContext)

	// EnterCtypedecl_ptr is called when entering the ctypedecl_ptr production.
	EnterCtypedecl_ptr(c *Ctypedecl_ptrContext)

	// EnterCtypedecl_simple is called when entering the ctypedecl_simple production.
	EnterCtypedecl_simple(c *Ctypedecl_simpleContext)

	// EnterTypeident is called when entering the typeident production.
	EnterTypeident(c *TypeidentContext)

	// EnterBitSlice is called when entering the bitSlice production.
	EnterBitSlice(c *BitSliceContext)

	// EnterAssignop is called when entering the assignop production.
	EnterAssignop(c *AssignopContext)

	// EnterNewExpr is called when entering the newExpr production.
	EnterNewExpr(c *NewExprContext)

	// EnterTernaryConstIfExpr is called when entering the ternaryConstIfExpr production.
	EnterTernaryConstIfExpr(c *TernaryConstIfExprContext)

	// EnterCastExpr is called when entering the castExpr production.
	EnterCastExpr(c *CastExprContext)

	// EnterDefinedExpr is called when entering the definedExpr production.
	EnterDefinedExpr(c *DefinedExprContext)

	// EnterLogicalExpr is called when entering the logicalExpr production.
	EnterLogicalExpr(c *LogicalExprContext)

	// EnterBitSliceExpr is called when entering the bitSliceExpr production.
	EnterBitSliceExpr(c *BitSliceExprContext)

	// EnterEachExpr is called when entering the eachExpr production.
	EnterEachExpr(c *EachExprContext)

	// EnterAdditiveExpr is called when entering the additiveExpr production.
	EnterAdditiveExpr(c *AdditiveExprContext)

	// EnterBitwiseNotExpr is called when entering the bitwiseNotExpr production.
	EnterBitwiseNotExpr(c *BitwiseNotExprContext)

	// EnterParenExpr is called when entering the parenExpr production.
	EnterParenExpr(c *ParenExprContext)

	// EnterUnaryExpr is called when entering the unaryExpr production.
	EnterUnaryExpr(c *UnaryExprContext)

	// EnterArrayAccessExpr is called when entering the arrayAccessExpr production.
	EnterArrayAccessExpr(c *ArrayAccessExprContext)

	// EnterFunctionCallExpr is called when entering the functionCallExpr production.
	EnterFunctionCallExpr(c *FunctionCallExprContext)

	// EnterMemberAccessExpr is called when entering the memberAccessExpr production.
	EnterMemberAccessExpr(c *MemberAccessExprContext)

	// EnterTernaryIfExpr is called when entering the ternaryIfExpr production.
	EnterTernaryIfExpr(c *TernaryIfExprContext)

	// EnterPostIncrementExpr is called when entering the postIncrementExpr production.
	EnterPostIncrementExpr(c *PostIncrementExprContext)

	// EnterRelationalExpr is called when entering the relationalExpr production.
	EnterRelationalExpr(c *RelationalExprContext)

	// EnterPreIncrementExpr is called when entering the preIncrementExpr production.
	EnterPreIncrementExpr(c *PreIncrementExprContext)

	// EnterSizeoftypeExpr is called when entering the sizeoftypeExpr production.
	EnterSizeoftypeExpr(c *SizeoftypeExprContext)

	// EnterShiftExpr is called when entering the shiftExpr production.
	EnterShiftExpr(c *ShiftExprContext)

	// EnterAddressofExpr is called when entering the addressofExpr production.
	EnterAddressofExpr(c *AddressofExprContext)

	// EnterBitwiseExpr is called when entering the bitwiseExpr production.
	EnterBitwiseExpr(c *BitwiseExprContext)

	// EnterSizeofExpr is called when entering the sizeofExpr production.
	EnterSizeofExpr(c *SizeofExprContext)

	// EnterLiteralAtom is called when entering the literalAtom production.
	EnterLiteralAtom(c *LiteralAtomContext)

	// EnterNotExpr is called when entering the notExpr production.
	EnterNotExpr(c *NotExprContext)

	// EnterDerefExpr is called when entering the derefExpr production.
	EnterDerefExpr(c *DerefExprContext)

	// EnterStringifyExpr is called when entering the stringifyExpr production.
	EnterStringifyExpr(c *StringifyExprContext)

	// EnterArrayLiteralExpr is called when entering the arrayLiteralExpr production.
	EnterArrayLiteralExpr(c *ArrayLiteralExprContext)

	// EnterMultiplicativeExpr is called when entering the multiplicativeExpr production.
	EnterMultiplicativeExpr(c *MultiplicativeExprContext)

	// EnterEqualityExpr is called when entering the equalityExpr production.
	EnterEqualityExpr(c *EqualityExprContext)

	// EnterNumericLiteral is called when entering the numericLiteral production.
	EnterNumericLiteral(c *NumericLiteralContext)

	// EnterStringLiteral is called when entering the stringLiteral production.
	EnterStringLiteral(c *StringLiteralContext)

	// EnterExpressionAccessIdent is called when entering the expressionAccessIdent production.
	EnterExpressionAccessIdent(c *ExpressionAccessIdentContext)

	// EnterExpressionIdent is called when entering the expressionIdent production.
	EnterExpressionIdent(c *ExpressionIdentContext)

	// EnterUtf8_sconst is called when entering the utf8_sconst production.
	EnterUtf8_sconst(c *Utf8_sconstContext)

	// EnterTypeoparg is called when entering the typeoparg production.
	EnterTypeoparg(c *TypeopargContext)

	// EnterEndianflag is called when entering the endianflag production.
	EnterEndianflag(c *EndianflagContext)

	// EnterExpression_list is called when entering the expression_list production.
	EnterExpression_list(c *Expression_listContext)

	// EnterExpression_list_ntc_nonempty is called when entering the expression_list_ntc_nonempty production.
	EnterExpression_list_ntc_nonempty(c *Expression_list_ntc_nonemptyContext)

	// EnterComposed_string_literal is called when entering the composed_string_literal production.
	EnterComposed_string_literal(c *Composed_string_literalContext)

	// EnterBracketed_string_literal is called when entering the bracketed_string_literal production.
	EnterBracketed_string_literal(c *Bracketed_string_literalContext)

	// EnterSingleInitializer is called when entering the singleInitializer production.
	EnterSingleInitializer(c *SingleInitializerContext)

	// EnterDesignatedStructInitializer is called when entering the designatedStructInitializer production.
	EnterDesignatedStructInitializer(c *DesignatedStructInitializerContext)

	// EnterDesignated_struct_initializer_list is called when entering the designated_struct_initializer_list production.
	EnterDesignated_struct_initializer_list(c *Designated_struct_initializer_listContext)

	// EnterCdeclInitializer is called when entering the cdeclInitializer production.
	EnterCdeclInitializer(c *CdeclInitializerContext)

	// EnterSingleInitializerArray is called when entering the singleInitializerArray production.
	EnterSingleInitializerArray(c *SingleInitializerArrayContext)

	// EnterInitializer is called when entering the initializer production.
	EnterInitializer(c *InitializerContext)

	// EnterArgumentsList is called when entering the argumentsList production.
	EnterArgumentsList(c *ArgumentsListContext)

	// EnterArgumentsListEmpty is called when entering the argumentsListEmpty production.
	EnterArgumentsListEmpty(c *ArgumentsListEmptyContext)

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterThrowStatement is called when entering the throwStatement production.
	EnterThrowStatement(c *ThrowStatementContext)

	// EnterContinueStatement is called when entering the continueStatement production.
	EnterContinueStatement(c *ContinueStatementContext)

	// EnterBreakStatement is called when entering the breakStatement production.
	EnterBreakStatement(c *BreakStatementContext)

	// EnterElseStatement is called when entering the elseStatement production.
	EnterElseStatement(c *ElseStatementContext)

	// EnterLoopStatement is called when entering the loopStatement production.
	EnterLoopStatement(c *LoopStatementContext)

	// EnterConditionExpression is called when entering the conditionExpression production.
	EnterConditionExpression(c *ConditionExpressionContext)

	// EnterAssertStatement is called when entering the assertStatement production.
	EnterAssertStatement(c *AssertStatementContext)

	// EnterDeleteStatement is called when entering the deleteStatement production.
	EnterDeleteStatement(c *DeleteStatementContext)

	// EnterReturnStatement is called when entering the returnStatement production.
	EnterReturnStatement(c *ReturnStatementContext)

	// EnterGotoStatement is called when entering the gotoStatement production.
	EnterGotoStatement(c *GotoStatementContext)

	// EnterLogKind is called when entering the logKind production.
	EnterLogKind(c *LogKindContext)

	// EnterLogGroup is called when entering the logGroup production.
	EnterLogGroup(c *LogGroupContext)

	// EnterLogFormatString is called when entering the logFormatString production.
	EnterLogFormatString(c *LogFormatStringContext)

	// EnterThenLevel is called when entering the thenLevel production.
	EnterThenLevel(c *ThenLevelContext)

	// EnterLogLevel is called when entering the logLevel production.
	EnterLogLevel(c *LogLevelContext)

	// EnterLogStatement is called when entering the logStatement production.
	EnterLogStatement(c *LogStatementContext)

	// EnterCaseOrDefaultStatement is called when entering the caseOrDefaultStatement production.
	EnterCaseOrDefaultStatement(c *CaseOrDefaultStatementContext)

	// EnterCaseStatement is called when entering the caseStatement production.
	EnterCaseStatement(c *CaseStatementContext)

	// EnterSwitchBodyStatement is called when entering the switchBodyStatement production.
	EnterSwitchBodyStatement(c *SwitchBodyStatementContext)

	// EnterSwitchBodyHashElseStatement is called when entering the switchBodyHashElseStatement production.
	EnterSwitchBodyHashElseStatement(c *SwitchBodyHashElseStatementContext)

	// EnterSwitchBodyHashIfStatement is called when entering the switchBodyHashIfStatement production.
	EnterSwitchBodyHashIfStatement(c *SwitchBodyHashIfStatementContext)

	// EnterSwitchStatement is called when entering the switchStatement production.
	EnterSwitchStatement(c *SwitchStatementContext)

	// EnterForLoopStepStatement is called when entering the forLoopStepStatement production.
	EnterForLoopStepStatement(c *ForLoopStepStatementContext)

	// EnterInnerForLoopStatement is called when entering the innerForLoopStatement production.
	EnterInnerForLoopStatement(c *InnerForLoopStatementContext)

	// EnterForLoopInitStatement is called when entering the forLoopInitStatement production.
	EnterForLoopInitStatement(c *ForLoopInitStatementContext)

	// EnterForLoopStatement is called when entering the forLoopStatement production.
	EnterForLoopStatement(c *ForLoopStatementContext)

	// EnterWhileLoopStatement is called when entering the whileLoopStatement production.
	EnterWhileLoopStatement(c *WhileLoopStatementContext)

	// EnterDoWhileLoopStatement is called when entering the doWhileLoopStatement production.
	EnterDoWhileLoopStatement(c *DoWhileLoopStatementContext)

	// EnterForeachStatement is called when entering the foreachStatement production.
	EnterForeachStatement(c *ForeachStatementContext)

	// EnterIfStatement is called when entering the ifStatement production.
	EnterIfStatement(c *IfStatementContext)

	// EnterCatchStatement is called when entering the catchStatement production.
	EnterCatchStatement(c *CatchStatementContext)

	// EnterTryStatement is called when entering the tryStatement production.
	EnterTryStatement(c *TryStatementContext)

	// EnterSelectStatement is called when entering the selectStatement production.
	EnterSelectStatement(c *SelectStatementContext)

	// EnterAssignStatement is called when entering the assignStatement production.
	EnterAssignStatement(c *AssignStatementContext)

	// EnterAssign_chain is called when entering the assign_chain production.
	EnterAssign_chain(c *Assign_chainContext)

	// EnterAfterStatement is called when entering the afterStatement production.
	EnterAfterStatement(c *AfterStatementContext)

	// EnterAfterDelayCondition is called when entering the afterDelayCondition production.
	EnterAfterDelayCondition(c *AfterDelayConditionContext)

	// EnterTuple_literal is called when entering the tuple_literal production.
	EnterTuple_literal(c *Tuple_literalContext)

	// EnterIdentList is called when entering the identList production.
	EnterIdentList(c *IdentListContext)

	// EnterErrorStatement is called when entering the errorStatement production.
	EnterErrorStatement(c *ErrorStatementContext)

	// EnterWarningStatement is called when entering the warningStatement production.
	EnterWarningStatement(c *WarningStatementContext)

	// EnterCompoundStatement is called when entering the compoundStatement production.
	EnterCompoundStatement(c *CompoundStatementContext)

	// EnterHook_decl is called when entering the hook_decl production.
	EnterHook_decl(c *Hook_declContext)

	// EnterObjident is called when entering the objident production.
	EnterObjident(c *ObjidentContext)

	// EnterIdent is called when entering the ident production.
	EnterIdent(c *IdentContext)

	// ExitDml is called when exiting the dml production.
	ExitDml(c *DmlContext)

	// ExitDmlCode is called when exiting the dmlCode production.
	ExitDmlCode(c *DmlCodeContext)

	// ExitDmlTag is called when exiting the dmlTag production.
	ExitDmlTag(c *DmlTagContext)

	// ExitProvisionalTag is called when exiting the provisionalTag production.
	ExitProvisionalTag(c *ProvisionalTagContext)

	// ExitDeviceTag is called when exiting the deviceTag production.
	ExitDeviceTag(c *DeviceTagContext)

	// ExitBitorderTag is called when exiting the bitorderTag production.
	ExitBitorderTag(c *BitorderTagContext)

	// ExitDeviceStatements is called when exiting the deviceStatements production.
	ExitDeviceStatements(c *DeviceStatementsContext)

	// ExitExportAsExpression is called when exiting the exportAsExpression production.
	ExitExportAsExpression(c *ExportAsExpressionContext)

	// ExitExportStatement is called when exiting the exportStatement production.
	ExitExportStatement(c *ExportStatementContext)

	// ExitImportStatement is called when exiting the importStatement production.
	ExitImportStatement(c *ImportStatementContext)

	// ExitExternStatement is called when exiting the externStatement production.
	ExitExternStatement(c *ExternStatementContext)

	// ExitConstantStatement is called when exiting the constantStatement production.
	ExitConstantStatement(c *ConstantStatementContext)

	// ExitLogGroupStatement is called when exiting the logGroupStatement production.
	ExitLogGroupStatement(c *LogGroupStatementContext)

	// ExitTypedefStatement is called when exiting the typedefStatement production.
	ExitTypedefStatement(c *TypedefStatementContext)

	// ExitTemplateStatement is called when exiting the templateStatement production.
	ExitTemplateStatement(c *TemplateStatementContext)

	// ExitTemplateBody is called when exiting the templateBody production.
	ExitTemplateBody(c *TemplateBodyContext)

	// ExitDeviceStatement is called when exiting the deviceStatement production.
	ExitDeviceStatement(c *DeviceStatementContext)

	// ExitToplevelIfStatement is called when exiting the toplevelIfStatement production.
	ExitToplevelIfStatement(c *ToplevelIfStatementContext)

	// ExitToplevelElseStatement is called when exiting the toplevelElseStatement production.
	ExitToplevelElseStatement(c *ToplevelElseStatementContext)

	// ExitToplevelIf is called when exiting the toplevelIf production.
	ExitToplevelIf(c *ToplevelIfContext)

	// ExitArray_list is called when exiting the array_list production.
	ExitArray_list(c *Array_listContext)

	// ExitObject is called when exiting the object production.
	ExitObject(c *ObjectContext)

	// ExitBitrangeDecl is called when exiting the bitrangeDecl production.
	ExitBitrangeDecl(c *BitrangeDeclContext)

	// ExitObject_type is called when exiting the object_type production.
	ExitObject_type(c *Object_typeContext)

	// ExitObjectDescSpec is called when exiting the objectDescSpec production.
	ExitObjectDescSpec(c *ObjectDescSpecContext)

	// ExitObjectBody is called when exiting the objectBody production.
	ExitObjectBody(c *ObjectBodyContext)

	// ExitObject_decl is called when exiting the object_decl production.
	ExitObject_decl(c *Object_declContext)

	// ExitGlobal_decl_kind is called when exiting the global_decl_kind production.
	ExitGlobal_decl_kind(c *Global_decl_kindContext)

	// ExitLocal_decl_kind is called when exiting the local_decl_kind production.
	ExitLocal_decl_kind(c *Local_decl_kindContext)

	// ExitCvarDecl is called when exiting the cvarDecl production.
	ExitCvarDecl(c *CvarDeclContext)

	// ExitSessionOrSavedDecl is called when exiting the sessionOrSavedDecl production.
	ExitSessionOrSavedDecl(c *SessionOrSavedDeclContext)

	// ExitSessionOrSavedStatement is called when exiting the sessionOrSavedStatement production.
	ExitSessionOrSavedStatement(c *SessionOrSavedStatementContext)

	// ExitVariableDecl is called when exiting the variableDecl production.
	ExitVariableDecl(c *VariableDeclContext)

	// ExitVariableStatement is called when exiting the variableStatement production.
	ExitVariableStatement(c *VariableStatementContext)

	// ExitMethodOutParams is called when exiting the methodOutParams production.
	ExitMethodOutParams(c *MethodOutParamsContext)

	// ExitMethodParamsMaybeUntyped is called when exiting the methodParamsMaybeUntyped production.
	ExitMethodParamsMaybeUntyped(c *MethodParamsMaybeUntypedContext)

	// ExitMethodQualifiers is called when exiting the methodQualifiers production.
	ExitMethodQualifiers(c *MethodQualifiersContext)

	// ExitMethodUnifiedDeclaration is called when exiting the methodUnifiedDeclaration production.
	ExitMethodUnifiedDeclaration(c *MethodUnifiedDeclarationContext)

	// ExitMethodUnifiedDefinition is called when exiting the methodUnifiedDefinition production.
	ExitMethodUnifiedDefinition(c *MethodUnifiedDefinitionContext)

	// ExitArraydef is called when exiting the arraydef production.
	ExitArraydef(c *ArraydefContext)

	// ExitCCodeBlock is called when exiting the cCodeBlock production.
	ExitCCodeBlock(c *CCodeBlockContext)

	// ExitObjectStatement is called when exiting the objectStatement production.
	ExitObjectStatement(c *ObjectStatementContext)

	// ExitIneachStatement is called when exiting the ineachStatement production.
	ExitIneachStatement(c *IneachStatementContext)

	// ExitObjectIfStatement is called when exiting the objectIfStatement production.
	ExitObjectIfStatement(c *ObjectIfStatementContext)

	// ExitObjectElseStatement is called when exiting the objectElseStatement production.
	ExitObjectElseStatement(c *ObjectElseStatementContext)

	// ExitObjectIf is called when exiting the objectIf production.
	ExitObjectIf(c *ObjectIfContext)

	// ExitParam is called when exiting the param production.
	ExitParam(c *ParamContext)

	// ExitIstemplateList is called when exiting the istemplateList production.
	ExitIstemplateList(c *IstemplateListContext)

	// ExitIstemplateExpr is called when exiting the istemplateExpr production.
	ExitIstemplateExpr(c *IstemplateExprContext)

	// ExitIstemplateStatement is called when exiting the istemplateStatement production.
	ExitIstemplateStatement(c *IstemplateStatementContext)

	// ExitSizeSpec is called when exiting the sizeSpec production.
	ExitSizeSpec(c *SizeSpecContext)

	// ExitOffsetSpec is called when exiting the offsetSpec production.
	ExitOffsetSpec(c *OffsetSpecContext)

	// ExitCdeclOrInline is called when exiting the cdeclOrInline production.
	ExitCdeclOrInline(c *CdeclOrInlineContext)

	// ExitCdecl is called when exiting the cdecl production.
	ExitCdecl(c *CdeclContext)

	// ExitSigned_or_unsigned_typeident is called when exiting the signed_or_unsigned_typeident production.
	ExitSigned_or_unsigned_typeident(c *Signed_or_unsigned_typeidentContext)

	// ExitBasetype is called when exiting the basetype production.
	ExitBasetype(c *BasetypeContext)

	// ExitTypeQualifier is called when exiting the typeQualifier production.
	ExitTypeQualifier(c *TypeQualifierContext)

	// ExitCdecl_qualifier_ident_array_func is called when exiting the cdecl_qualifier_ident_array_func production.
	ExitCdecl_qualifier_ident_array_func(c *Cdecl_qualifier_ident_array_funcContext)

	// ExitCdeclIdent is called when exiting the cdeclIdent production.
	ExitCdeclIdent(c *CdeclIdentContext)

	// ExitCdeclFuncIdent is called when exiting the cdeclFuncIdent production.
	ExitCdeclFuncIdent(c *CdeclFuncIdentContext)

	// ExitCdeclArray is called when exiting the cdeclArray production.
	ExitCdeclArray(c *CdeclArrayContext)

	// ExitCdeclNone is called when exiting the cdeclNone production.
	ExitCdeclNone(c *CdeclNoneContext)

	// ExitCdeclFuncArgs is called when exiting the cdeclFuncArgs production.
	ExitCdeclFuncArgs(c *CdeclFuncArgsContext)

	// ExitCdecl_list is called when exiting the cdecl_list production.
	ExitCdecl_list(c *Cdecl_listContext)

	// ExitCdecl_list_variadic is called when exiting the cdecl_list_variadic production.
	ExitCdecl_list_variadic(c *Cdecl_list_variadicContext)

	// ExitTypeof is called when exiting the typeof production.
	ExitTypeof(c *TypeofContext)

	// ExitStructDecls is called when exiting the structDecls production.
	ExitStructDecls(c *StructDeclsContext)

	// ExitStruct is called when exiting the struct production.
	ExitStruct(c *StructContext)

	// ExitUnion is called when exiting the union production.
	ExitUnion(c *UnionContext)

	// ExitLayoutDecls is called when exiting the layoutDecls production.
	ExitLayoutDecls(c *LayoutDeclsContext)

	// ExitLayout is called when exiting the layout production.
	ExitLayout(c *LayoutContext)

	// ExitBitfields is called when exiting the bitfields production.
	ExitBitfields(c *BitfieldsContext)

	// ExitBitfields_decls is called when exiting the bitfields_decls production.
	ExitBitfields_decls(c *Bitfields_declsContext)

	// ExitCtypedecl is called when exiting the ctypedecl production.
	ExitCtypedecl(c *CtypedeclContext)

	// ExitCtypedecl_ptr is called when exiting the ctypedecl_ptr production.
	ExitCtypedecl_ptr(c *Ctypedecl_ptrContext)

	// ExitCtypedecl_simple is called when exiting the ctypedecl_simple production.
	ExitCtypedecl_simple(c *Ctypedecl_simpleContext)

	// ExitTypeident is called when exiting the typeident production.
	ExitTypeident(c *TypeidentContext)

	// ExitBitSlice is called when exiting the bitSlice production.
	ExitBitSlice(c *BitSliceContext)

	// ExitAssignop is called when exiting the assignop production.
	ExitAssignop(c *AssignopContext)

	// ExitNewExpr is called when exiting the newExpr production.
	ExitNewExpr(c *NewExprContext)

	// ExitTernaryConstIfExpr is called when exiting the ternaryConstIfExpr production.
	ExitTernaryConstIfExpr(c *TernaryConstIfExprContext)

	// ExitCastExpr is called when exiting the castExpr production.
	ExitCastExpr(c *CastExprContext)

	// ExitDefinedExpr is called when exiting the definedExpr production.
	ExitDefinedExpr(c *DefinedExprContext)

	// ExitLogicalExpr is called when exiting the logicalExpr production.
	ExitLogicalExpr(c *LogicalExprContext)

	// ExitBitSliceExpr is called when exiting the bitSliceExpr production.
	ExitBitSliceExpr(c *BitSliceExprContext)

	// ExitEachExpr is called when exiting the eachExpr production.
	ExitEachExpr(c *EachExprContext)

	// ExitAdditiveExpr is called when exiting the additiveExpr production.
	ExitAdditiveExpr(c *AdditiveExprContext)

	// ExitBitwiseNotExpr is called when exiting the bitwiseNotExpr production.
	ExitBitwiseNotExpr(c *BitwiseNotExprContext)

	// ExitParenExpr is called when exiting the parenExpr production.
	ExitParenExpr(c *ParenExprContext)

	// ExitUnaryExpr is called when exiting the unaryExpr production.
	ExitUnaryExpr(c *UnaryExprContext)

	// ExitArrayAccessExpr is called when exiting the arrayAccessExpr production.
	ExitArrayAccessExpr(c *ArrayAccessExprContext)

	// ExitFunctionCallExpr is called when exiting the functionCallExpr production.
	ExitFunctionCallExpr(c *FunctionCallExprContext)

	// ExitMemberAccessExpr is called when exiting the memberAccessExpr production.
	ExitMemberAccessExpr(c *MemberAccessExprContext)

	// ExitTernaryIfExpr is called when exiting the ternaryIfExpr production.
	ExitTernaryIfExpr(c *TernaryIfExprContext)

	// ExitPostIncrementExpr is called when exiting the postIncrementExpr production.
	ExitPostIncrementExpr(c *PostIncrementExprContext)

	// ExitRelationalExpr is called when exiting the relationalExpr production.
	ExitRelationalExpr(c *RelationalExprContext)

	// ExitPreIncrementExpr is called when exiting the preIncrementExpr production.
	ExitPreIncrementExpr(c *PreIncrementExprContext)

	// ExitSizeoftypeExpr is called when exiting the sizeoftypeExpr production.
	ExitSizeoftypeExpr(c *SizeoftypeExprContext)

	// ExitShiftExpr is called when exiting the shiftExpr production.
	ExitShiftExpr(c *ShiftExprContext)

	// ExitAddressofExpr is called when exiting the addressofExpr production.
	ExitAddressofExpr(c *AddressofExprContext)

	// ExitBitwiseExpr is called when exiting the bitwiseExpr production.
	ExitBitwiseExpr(c *BitwiseExprContext)

	// ExitSizeofExpr is called when exiting the sizeofExpr production.
	ExitSizeofExpr(c *SizeofExprContext)

	// ExitLiteralAtom is called when exiting the literalAtom production.
	ExitLiteralAtom(c *LiteralAtomContext)

	// ExitNotExpr is called when exiting the notExpr production.
	ExitNotExpr(c *NotExprContext)

	// ExitDerefExpr is called when exiting the derefExpr production.
	ExitDerefExpr(c *DerefExprContext)

	// ExitStringifyExpr is called when exiting the stringifyExpr production.
	ExitStringifyExpr(c *StringifyExprContext)

	// ExitArrayLiteralExpr is called when exiting the arrayLiteralExpr production.
	ExitArrayLiteralExpr(c *ArrayLiteralExprContext)

	// ExitMultiplicativeExpr is called when exiting the multiplicativeExpr production.
	ExitMultiplicativeExpr(c *MultiplicativeExprContext)

	// ExitEqualityExpr is called when exiting the equalityExpr production.
	ExitEqualityExpr(c *EqualityExprContext)

	// ExitNumericLiteral is called when exiting the numericLiteral production.
	ExitNumericLiteral(c *NumericLiteralContext)

	// ExitStringLiteral is called when exiting the stringLiteral production.
	ExitStringLiteral(c *StringLiteralContext)

	// ExitExpressionAccessIdent is called when exiting the expressionAccessIdent production.
	ExitExpressionAccessIdent(c *ExpressionAccessIdentContext)

	// ExitExpressionIdent is called when exiting the expressionIdent production.
	ExitExpressionIdent(c *ExpressionIdentContext)

	// ExitUtf8_sconst is called when exiting the utf8_sconst production.
	ExitUtf8_sconst(c *Utf8_sconstContext)

	// ExitTypeoparg is called when exiting the typeoparg production.
	ExitTypeoparg(c *TypeopargContext)

	// ExitEndianflag is called when exiting the endianflag production.
	ExitEndianflag(c *EndianflagContext)

	// ExitExpression_list is called when exiting the expression_list production.
	ExitExpression_list(c *Expression_listContext)

	// ExitExpression_list_ntc_nonempty is called when exiting the expression_list_ntc_nonempty production.
	ExitExpression_list_ntc_nonempty(c *Expression_list_ntc_nonemptyContext)

	// ExitComposed_string_literal is called when exiting the composed_string_literal production.
	ExitComposed_string_literal(c *Composed_string_literalContext)

	// ExitBracketed_string_literal is called when exiting the bracketed_string_literal production.
	ExitBracketed_string_literal(c *Bracketed_string_literalContext)

	// ExitSingleInitializer is called when exiting the singleInitializer production.
	ExitSingleInitializer(c *SingleInitializerContext)

	// ExitDesignatedStructInitializer is called when exiting the designatedStructInitializer production.
	ExitDesignatedStructInitializer(c *DesignatedStructInitializerContext)

	// ExitDesignated_struct_initializer_list is called when exiting the designated_struct_initializer_list production.
	ExitDesignated_struct_initializer_list(c *Designated_struct_initializer_listContext)

	// ExitCdeclInitializer is called when exiting the cdeclInitializer production.
	ExitCdeclInitializer(c *CdeclInitializerContext)

	// ExitSingleInitializerArray is called when exiting the singleInitializerArray production.
	ExitSingleInitializerArray(c *SingleInitializerArrayContext)

	// ExitInitializer is called when exiting the initializer production.
	ExitInitializer(c *InitializerContext)

	// ExitArgumentsList is called when exiting the argumentsList production.
	ExitArgumentsList(c *ArgumentsListContext)

	// ExitArgumentsListEmpty is called when exiting the argumentsListEmpty production.
	ExitArgumentsListEmpty(c *ArgumentsListEmptyContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitThrowStatement is called when exiting the throwStatement production.
	ExitThrowStatement(c *ThrowStatementContext)

	// ExitContinueStatement is called when exiting the continueStatement production.
	ExitContinueStatement(c *ContinueStatementContext)

	// ExitBreakStatement is called when exiting the breakStatement production.
	ExitBreakStatement(c *BreakStatementContext)

	// ExitElseStatement is called when exiting the elseStatement production.
	ExitElseStatement(c *ElseStatementContext)

	// ExitLoopStatement is called when exiting the loopStatement production.
	ExitLoopStatement(c *LoopStatementContext)

	// ExitConditionExpression is called when exiting the conditionExpression production.
	ExitConditionExpression(c *ConditionExpressionContext)

	// ExitAssertStatement is called when exiting the assertStatement production.
	ExitAssertStatement(c *AssertStatementContext)

	// ExitDeleteStatement is called when exiting the deleteStatement production.
	ExitDeleteStatement(c *DeleteStatementContext)

	// ExitReturnStatement is called when exiting the returnStatement production.
	ExitReturnStatement(c *ReturnStatementContext)

	// ExitGotoStatement is called when exiting the gotoStatement production.
	ExitGotoStatement(c *GotoStatementContext)

	// ExitLogKind is called when exiting the logKind production.
	ExitLogKind(c *LogKindContext)

	// ExitLogGroup is called when exiting the logGroup production.
	ExitLogGroup(c *LogGroupContext)

	// ExitLogFormatString is called when exiting the logFormatString production.
	ExitLogFormatString(c *LogFormatStringContext)

	// ExitThenLevel is called when exiting the thenLevel production.
	ExitThenLevel(c *ThenLevelContext)

	// ExitLogLevel is called when exiting the logLevel production.
	ExitLogLevel(c *LogLevelContext)

	// ExitLogStatement is called when exiting the logStatement production.
	ExitLogStatement(c *LogStatementContext)

	// ExitCaseOrDefaultStatement is called when exiting the caseOrDefaultStatement production.
	ExitCaseOrDefaultStatement(c *CaseOrDefaultStatementContext)

	// ExitCaseStatement is called when exiting the caseStatement production.
	ExitCaseStatement(c *CaseStatementContext)

	// ExitSwitchBodyStatement is called when exiting the switchBodyStatement production.
	ExitSwitchBodyStatement(c *SwitchBodyStatementContext)

	// ExitSwitchBodyHashElseStatement is called when exiting the switchBodyHashElseStatement production.
	ExitSwitchBodyHashElseStatement(c *SwitchBodyHashElseStatementContext)

	// ExitSwitchBodyHashIfStatement is called when exiting the switchBodyHashIfStatement production.
	ExitSwitchBodyHashIfStatement(c *SwitchBodyHashIfStatementContext)

	// ExitSwitchStatement is called when exiting the switchStatement production.
	ExitSwitchStatement(c *SwitchStatementContext)

	// ExitForLoopStepStatement is called when exiting the forLoopStepStatement production.
	ExitForLoopStepStatement(c *ForLoopStepStatementContext)

	// ExitInnerForLoopStatement is called when exiting the innerForLoopStatement production.
	ExitInnerForLoopStatement(c *InnerForLoopStatementContext)

	// ExitForLoopInitStatement is called when exiting the forLoopInitStatement production.
	ExitForLoopInitStatement(c *ForLoopInitStatementContext)

	// ExitForLoopStatement is called when exiting the forLoopStatement production.
	ExitForLoopStatement(c *ForLoopStatementContext)

	// ExitWhileLoopStatement is called when exiting the whileLoopStatement production.
	ExitWhileLoopStatement(c *WhileLoopStatementContext)

	// ExitDoWhileLoopStatement is called when exiting the doWhileLoopStatement production.
	ExitDoWhileLoopStatement(c *DoWhileLoopStatementContext)

	// ExitForeachStatement is called when exiting the foreachStatement production.
	ExitForeachStatement(c *ForeachStatementContext)

	// ExitIfStatement is called when exiting the ifStatement production.
	ExitIfStatement(c *IfStatementContext)

	// ExitCatchStatement is called when exiting the catchStatement production.
	ExitCatchStatement(c *CatchStatementContext)

	// ExitTryStatement is called when exiting the tryStatement production.
	ExitTryStatement(c *TryStatementContext)

	// ExitSelectStatement is called when exiting the selectStatement production.
	ExitSelectStatement(c *SelectStatementContext)

	// ExitAssignStatement is called when exiting the assignStatement production.
	ExitAssignStatement(c *AssignStatementContext)

	// ExitAssign_chain is called when exiting the assign_chain production.
	ExitAssign_chain(c *Assign_chainContext)

	// ExitAfterStatement is called when exiting the afterStatement production.
	ExitAfterStatement(c *AfterStatementContext)

	// ExitAfterDelayCondition is called when exiting the afterDelayCondition production.
	ExitAfterDelayCondition(c *AfterDelayConditionContext)

	// ExitTuple_literal is called when exiting the tuple_literal production.
	ExitTuple_literal(c *Tuple_literalContext)

	// ExitIdentList is called when exiting the identList production.
	ExitIdentList(c *IdentListContext)

	// ExitErrorStatement is called when exiting the errorStatement production.
	ExitErrorStatement(c *ErrorStatementContext)

	// ExitWarningStatement is called when exiting the warningStatement production.
	ExitWarningStatement(c *WarningStatementContext)

	// ExitCompoundStatement is called when exiting the compoundStatement production.
	ExitCompoundStatement(c *CompoundStatementContext)

	// ExitHook_decl is called when exiting the hook_decl production.
	ExitHook_decl(c *Hook_declContext)

	// ExitObjident is called when exiting the objident production.
	ExitObjident(c *ObjidentContext)

	// ExitIdent is called when exiting the ident production.
	ExitIdent(c *IdentContext)
}
