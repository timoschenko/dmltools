// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex // dml
import "github.com/antlr4-go/antlr/v4"

// BasedmlListener is a complete listener for a parse tree produced by dmlParser.
type BasedmlListener struct{}

var _ dmlListener = &BasedmlListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasedmlListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasedmlListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasedmlListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasedmlListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterDml is called when production dml is entered.
func (s *BasedmlListener) EnterDml(ctx *DmlContext) {}

// ExitDml is called when production dml is exited.
func (s *BasedmlListener) ExitDml(ctx *DmlContext) {}

// EnterDmlCode is called when production dmlCode is entered.
func (s *BasedmlListener) EnterDmlCode(ctx *DmlCodeContext) {}

// ExitDmlCode is called when production dmlCode is exited.
func (s *BasedmlListener) ExitDmlCode(ctx *DmlCodeContext) {}

// EnterDmlTag is called when production dmlTag is entered.
func (s *BasedmlListener) EnterDmlTag(ctx *DmlTagContext) {}

// ExitDmlTag is called when production dmlTag is exited.
func (s *BasedmlListener) ExitDmlTag(ctx *DmlTagContext) {}

// EnterProvisionalTag is called when production provisionalTag is entered.
func (s *BasedmlListener) EnterProvisionalTag(ctx *ProvisionalTagContext) {}

// ExitProvisionalTag is called when production provisionalTag is exited.
func (s *BasedmlListener) ExitProvisionalTag(ctx *ProvisionalTagContext) {}

// EnterDeviceTag is called when production deviceTag is entered.
func (s *BasedmlListener) EnterDeviceTag(ctx *DeviceTagContext) {}

// ExitDeviceTag is called when production deviceTag is exited.
func (s *BasedmlListener) ExitDeviceTag(ctx *DeviceTagContext) {}

// EnterBitorderTag is called when production bitorderTag is entered.
func (s *BasedmlListener) EnterBitorderTag(ctx *BitorderTagContext) {}

// ExitBitorderTag is called when production bitorderTag is exited.
func (s *BasedmlListener) ExitBitorderTag(ctx *BitorderTagContext) {}

// EnterDeviceStatements is called when production deviceStatements is entered.
func (s *BasedmlListener) EnterDeviceStatements(ctx *DeviceStatementsContext) {}

// ExitDeviceStatements is called when production deviceStatements is exited.
func (s *BasedmlListener) ExitDeviceStatements(ctx *DeviceStatementsContext) {}

// EnterExportAsExpression is called when production exportAsExpression is entered.
func (s *BasedmlListener) EnterExportAsExpression(ctx *ExportAsExpressionContext) {}

// ExitExportAsExpression is called when production exportAsExpression is exited.
func (s *BasedmlListener) ExitExportAsExpression(ctx *ExportAsExpressionContext) {}

// EnterExportStatement is called when production exportStatement is entered.
func (s *BasedmlListener) EnterExportStatement(ctx *ExportStatementContext) {}

// ExitExportStatement is called when production exportStatement is exited.
func (s *BasedmlListener) ExitExportStatement(ctx *ExportStatementContext) {}

// EnterImportStatement is called when production importStatement is entered.
func (s *BasedmlListener) EnterImportStatement(ctx *ImportStatementContext) {}

// ExitImportStatement is called when production importStatement is exited.
func (s *BasedmlListener) ExitImportStatement(ctx *ImportStatementContext) {}

// EnterExternStatement is called when production externStatement is entered.
func (s *BasedmlListener) EnterExternStatement(ctx *ExternStatementContext) {}

// ExitExternStatement is called when production externStatement is exited.
func (s *BasedmlListener) ExitExternStatement(ctx *ExternStatementContext) {}

// EnterConstantStatement is called when production constantStatement is entered.
func (s *BasedmlListener) EnterConstantStatement(ctx *ConstantStatementContext) {}

// ExitConstantStatement is called when production constantStatement is exited.
func (s *BasedmlListener) ExitConstantStatement(ctx *ConstantStatementContext) {}

// EnterLogGroupStatement is called when production logGroupStatement is entered.
func (s *BasedmlListener) EnterLogGroupStatement(ctx *LogGroupStatementContext) {}

// ExitLogGroupStatement is called when production logGroupStatement is exited.
func (s *BasedmlListener) ExitLogGroupStatement(ctx *LogGroupStatementContext) {}

// EnterTypedefStatement is called when production typedefStatement is entered.
func (s *BasedmlListener) EnterTypedefStatement(ctx *TypedefStatementContext) {}

// ExitTypedefStatement is called when production typedefStatement is exited.
func (s *BasedmlListener) ExitTypedefStatement(ctx *TypedefStatementContext) {}

// EnterTemplateStatement is called when production templateStatement is entered.
func (s *BasedmlListener) EnterTemplateStatement(ctx *TemplateStatementContext) {}

// ExitTemplateStatement is called when production templateStatement is exited.
func (s *BasedmlListener) ExitTemplateStatement(ctx *TemplateStatementContext) {}

// EnterTemplateBody is called when production templateBody is entered.
func (s *BasedmlListener) EnterTemplateBody(ctx *TemplateBodyContext) {}

// ExitTemplateBody is called when production templateBody is exited.
func (s *BasedmlListener) ExitTemplateBody(ctx *TemplateBodyContext) {}

// EnterDeviceStatement is called when production deviceStatement is entered.
func (s *BasedmlListener) EnterDeviceStatement(ctx *DeviceStatementContext) {}

// ExitDeviceStatement is called when production deviceStatement is exited.
func (s *BasedmlListener) ExitDeviceStatement(ctx *DeviceStatementContext) {}

// EnterToplevelIfStatement is called when production toplevelIfStatement is entered.
func (s *BasedmlListener) EnterToplevelIfStatement(ctx *ToplevelIfStatementContext) {}

// ExitToplevelIfStatement is called when production toplevelIfStatement is exited.
func (s *BasedmlListener) ExitToplevelIfStatement(ctx *ToplevelIfStatementContext) {}

// EnterToplevelElseStatement is called when production toplevelElseStatement is entered.
func (s *BasedmlListener) EnterToplevelElseStatement(ctx *ToplevelElseStatementContext) {}

// ExitToplevelElseStatement is called when production toplevelElseStatement is exited.
func (s *BasedmlListener) ExitToplevelElseStatement(ctx *ToplevelElseStatementContext) {}

// EnterToplevelIf is called when production toplevelIf is entered.
func (s *BasedmlListener) EnterToplevelIf(ctx *ToplevelIfContext) {}

// ExitToplevelIf is called when production toplevelIf is exited.
func (s *BasedmlListener) ExitToplevelIf(ctx *ToplevelIfContext) {}

// EnterArray_list is called when production array_list is entered.
func (s *BasedmlListener) EnterArray_list(ctx *Array_listContext) {}

// ExitArray_list is called when production array_list is exited.
func (s *BasedmlListener) ExitArray_list(ctx *Array_listContext) {}

// EnterObject is called when production object is entered.
func (s *BasedmlListener) EnterObject(ctx *ObjectContext) {}

// ExitObject is called when production object is exited.
func (s *BasedmlListener) ExitObject(ctx *ObjectContext) {}

// EnterBitrangeDecl is called when production bitrangeDecl is entered.
func (s *BasedmlListener) EnterBitrangeDecl(ctx *BitrangeDeclContext) {}

// ExitBitrangeDecl is called when production bitrangeDecl is exited.
func (s *BasedmlListener) ExitBitrangeDecl(ctx *BitrangeDeclContext) {}

// EnterObject_type is called when production object_type is entered.
func (s *BasedmlListener) EnterObject_type(ctx *Object_typeContext) {}

// ExitObject_type is called when production object_type is exited.
func (s *BasedmlListener) ExitObject_type(ctx *Object_typeContext) {}

// EnterObjectDescSpec is called when production objectDescSpec is entered.
func (s *BasedmlListener) EnterObjectDescSpec(ctx *ObjectDescSpecContext) {}

// ExitObjectDescSpec is called when production objectDescSpec is exited.
func (s *BasedmlListener) ExitObjectDescSpec(ctx *ObjectDescSpecContext) {}

// EnterObjectBody is called when production objectBody is entered.
func (s *BasedmlListener) EnterObjectBody(ctx *ObjectBodyContext) {}

// ExitObjectBody is called when production objectBody is exited.
func (s *BasedmlListener) ExitObjectBody(ctx *ObjectBodyContext) {}

// EnterObject_decl is called when production object_decl is entered.
func (s *BasedmlListener) EnterObject_decl(ctx *Object_declContext) {}

// ExitObject_decl is called when production object_decl is exited.
func (s *BasedmlListener) ExitObject_decl(ctx *Object_declContext) {}

// EnterGlobal_decl_kind is called when production global_decl_kind is entered.
func (s *BasedmlListener) EnterGlobal_decl_kind(ctx *Global_decl_kindContext) {}

// ExitGlobal_decl_kind is called when production global_decl_kind is exited.
func (s *BasedmlListener) ExitGlobal_decl_kind(ctx *Global_decl_kindContext) {}

// EnterLocal_decl_kind is called when production local_decl_kind is entered.
func (s *BasedmlListener) EnterLocal_decl_kind(ctx *Local_decl_kindContext) {}

// ExitLocal_decl_kind is called when production local_decl_kind is exited.
func (s *BasedmlListener) ExitLocal_decl_kind(ctx *Local_decl_kindContext) {}

// EnterCvarDecl is called when production cvarDecl is entered.
func (s *BasedmlListener) EnterCvarDecl(ctx *CvarDeclContext) {}

// ExitCvarDecl is called when production cvarDecl is exited.
func (s *BasedmlListener) ExitCvarDecl(ctx *CvarDeclContext) {}

// EnterSessionOrSavedDecl is called when production sessionOrSavedDecl is entered.
func (s *BasedmlListener) EnterSessionOrSavedDecl(ctx *SessionOrSavedDeclContext) {}

// ExitSessionOrSavedDecl is called when production sessionOrSavedDecl is exited.
func (s *BasedmlListener) ExitSessionOrSavedDecl(ctx *SessionOrSavedDeclContext) {}

// EnterSessionOrSavedStatement is called when production sessionOrSavedStatement is entered.
func (s *BasedmlListener) EnterSessionOrSavedStatement(ctx *SessionOrSavedStatementContext) {}

// ExitSessionOrSavedStatement is called when production sessionOrSavedStatement is exited.
func (s *BasedmlListener) ExitSessionOrSavedStatement(ctx *SessionOrSavedStatementContext) {}

// EnterVariableDecl is called when production variableDecl is entered.
func (s *BasedmlListener) EnterVariableDecl(ctx *VariableDeclContext) {}

// ExitVariableDecl is called when production variableDecl is exited.
func (s *BasedmlListener) ExitVariableDecl(ctx *VariableDeclContext) {}

// EnterVariableStatement is called when production variableStatement is entered.
func (s *BasedmlListener) EnterVariableStatement(ctx *VariableStatementContext) {}

// ExitVariableStatement is called when production variableStatement is exited.
func (s *BasedmlListener) ExitVariableStatement(ctx *VariableStatementContext) {}

// EnterMethodOutParams is called when production methodOutParams is entered.
func (s *BasedmlListener) EnterMethodOutParams(ctx *MethodOutParamsContext) {}

// ExitMethodOutParams is called when production methodOutParams is exited.
func (s *BasedmlListener) ExitMethodOutParams(ctx *MethodOutParamsContext) {}

// EnterMethodParamsMaybeUntyped is called when production methodParamsMaybeUntyped is entered.
func (s *BasedmlListener) EnterMethodParamsMaybeUntyped(ctx *MethodParamsMaybeUntypedContext) {}

// ExitMethodParamsMaybeUntyped is called when production methodParamsMaybeUntyped is exited.
func (s *BasedmlListener) ExitMethodParamsMaybeUntyped(ctx *MethodParamsMaybeUntypedContext) {}

// EnterMethodQualifiers is called when production methodQualifiers is entered.
func (s *BasedmlListener) EnterMethodQualifiers(ctx *MethodQualifiersContext) {}

// ExitMethodQualifiers is called when production methodQualifiers is exited.
func (s *BasedmlListener) ExitMethodQualifiers(ctx *MethodQualifiersContext) {}

// EnterMethodUnifiedDeclaration is called when production methodUnifiedDeclaration is entered.
func (s *BasedmlListener) EnterMethodUnifiedDeclaration(ctx *MethodUnifiedDeclarationContext) {}

// ExitMethodUnifiedDeclaration is called when production methodUnifiedDeclaration is exited.
func (s *BasedmlListener) ExitMethodUnifiedDeclaration(ctx *MethodUnifiedDeclarationContext) {}

// EnterMethodUnifiedDefinition is called when production methodUnifiedDefinition is entered.
func (s *BasedmlListener) EnterMethodUnifiedDefinition(ctx *MethodUnifiedDefinitionContext) {}

// ExitMethodUnifiedDefinition is called when production methodUnifiedDefinition is exited.
func (s *BasedmlListener) ExitMethodUnifiedDefinition(ctx *MethodUnifiedDefinitionContext) {}

// EnterArraydef is called when production arraydef is entered.
func (s *BasedmlListener) EnterArraydef(ctx *ArraydefContext) {}

// ExitArraydef is called when production arraydef is exited.
func (s *BasedmlListener) ExitArraydef(ctx *ArraydefContext) {}

// EnterCCodeBlock is called when production cCodeBlock is entered.
func (s *BasedmlListener) EnterCCodeBlock(ctx *CCodeBlockContext) {}

// ExitCCodeBlock is called when production cCodeBlock is exited.
func (s *BasedmlListener) ExitCCodeBlock(ctx *CCodeBlockContext) {}

// EnterObjectStatement is called when production objectStatement is entered.
func (s *BasedmlListener) EnterObjectStatement(ctx *ObjectStatementContext) {}

// ExitObjectStatement is called when production objectStatement is exited.
func (s *BasedmlListener) ExitObjectStatement(ctx *ObjectStatementContext) {}

// EnterIneachStatement is called when production ineachStatement is entered.
func (s *BasedmlListener) EnterIneachStatement(ctx *IneachStatementContext) {}

// ExitIneachStatement is called when production ineachStatement is exited.
func (s *BasedmlListener) ExitIneachStatement(ctx *IneachStatementContext) {}

// EnterObjectIfStatement is called when production objectIfStatement is entered.
func (s *BasedmlListener) EnterObjectIfStatement(ctx *ObjectIfStatementContext) {}

// ExitObjectIfStatement is called when production objectIfStatement is exited.
func (s *BasedmlListener) ExitObjectIfStatement(ctx *ObjectIfStatementContext) {}

// EnterObjectElseStatement is called when production objectElseStatement is entered.
func (s *BasedmlListener) EnterObjectElseStatement(ctx *ObjectElseStatementContext) {}

// ExitObjectElseStatement is called when production objectElseStatement is exited.
func (s *BasedmlListener) ExitObjectElseStatement(ctx *ObjectElseStatementContext) {}

// EnterObjectIf is called when production objectIf is entered.
func (s *BasedmlListener) EnterObjectIf(ctx *ObjectIfContext) {}

// ExitObjectIf is called when production objectIf is exited.
func (s *BasedmlListener) ExitObjectIf(ctx *ObjectIfContext) {}

// EnterParam is called when production param is entered.
func (s *BasedmlListener) EnterParam(ctx *ParamContext) {}

// ExitParam is called when production param is exited.
func (s *BasedmlListener) ExitParam(ctx *ParamContext) {}

// EnterIstemplateList is called when production istemplateList is entered.
func (s *BasedmlListener) EnterIstemplateList(ctx *IstemplateListContext) {}

// ExitIstemplateList is called when production istemplateList is exited.
func (s *BasedmlListener) ExitIstemplateList(ctx *IstemplateListContext) {}

// EnterIstemplateExpr is called when production istemplateExpr is entered.
func (s *BasedmlListener) EnterIstemplateExpr(ctx *IstemplateExprContext) {}

// ExitIstemplateExpr is called when production istemplateExpr is exited.
func (s *BasedmlListener) ExitIstemplateExpr(ctx *IstemplateExprContext) {}

// EnterIstemplateStatement is called when production istemplateStatement is entered.
func (s *BasedmlListener) EnterIstemplateStatement(ctx *IstemplateStatementContext) {}

// ExitIstemplateStatement is called when production istemplateStatement is exited.
func (s *BasedmlListener) ExitIstemplateStatement(ctx *IstemplateStatementContext) {}

// EnterSizeSpec is called when production sizeSpec is entered.
func (s *BasedmlListener) EnterSizeSpec(ctx *SizeSpecContext) {}

// ExitSizeSpec is called when production sizeSpec is exited.
func (s *BasedmlListener) ExitSizeSpec(ctx *SizeSpecContext) {}

// EnterOffsetSpec is called when production offsetSpec is entered.
func (s *BasedmlListener) EnterOffsetSpec(ctx *OffsetSpecContext) {}

// ExitOffsetSpec is called when production offsetSpec is exited.
func (s *BasedmlListener) ExitOffsetSpec(ctx *OffsetSpecContext) {}

// EnterCdeclOrInline is called when production cdeclOrInline is entered.
func (s *BasedmlListener) EnterCdeclOrInline(ctx *CdeclOrInlineContext) {}

// ExitCdeclOrInline is called when production cdeclOrInline is exited.
func (s *BasedmlListener) ExitCdeclOrInline(ctx *CdeclOrInlineContext) {}

// EnterCdecl is called when production cdecl is entered.
func (s *BasedmlListener) EnterCdecl(ctx *CdeclContext) {}

// ExitCdecl is called when production cdecl is exited.
func (s *BasedmlListener) ExitCdecl(ctx *CdeclContext) {}

// EnterSigned_or_unsigned_typeident is called when production signed_or_unsigned_typeident is entered.
func (s *BasedmlListener) EnterSigned_or_unsigned_typeident(ctx *Signed_or_unsigned_typeidentContext) {
}

// ExitSigned_or_unsigned_typeident is called when production signed_or_unsigned_typeident is exited.
func (s *BasedmlListener) ExitSigned_or_unsigned_typeident(ctx *Signed_or_unsigned_typeidentContext) {
}

// EnterBasetype is called when production basetype is entered.
func (s *BasedmlListener) EnterBasetype(ctx *BasetypeContext) {}

// ExitBasetype is called when production basetype is exited.
func (s *BasedmlListener) ExitBasetype(ctx *BasetypeContext) {}

// EnterTypeQualifier is called when production typeQualifier is entered.
func (s *BasedmlListener) EnterTypeQualifier(ctx *TypeQualifierContext) {}

// ExitTypeQualifier is called when production typeQualifier is exited.
func (s *BasedmlListener) ExitTypeQualifier(ctx *TypeQualifierContext) {}

// EnterCdecl_qualifier_ident_array_func is called when production cdecl_qualifier_ident_array_func is entered.
func (s *BasedmlListener) EnterCdecl_qualifier_ident_array_func(ctx *Cdecl_qualifier_ident_array_funcContext) {
}

// ExitCdecl_qualifier_ident_array_func is called when production cdecl_qualifier_ident_array_func is exited.
func (s *BasedmlListener) ExitCdecl_qualifier_ident_array_func(ctx *Cdecl_qualifier_ident_array_funcContext) {
}

// EnterCdeclIdent is called when production cdeclIdent is entered.
func (s *BasedmlListener) EnterCdeclIdent(ctx *CdeclIdentContext) {}

// ExitCdeclIdent is called when production cdeclIdent is exited.
func (s *BasedmlListener) ExitCdeclIdent(ctx *CdeclIdentContext) {}

// EnterCdeclFuncIdent is called when production cdeclFuncIdent is entered.
func (s *BasedmlListener) EnterCdeclFuncIdent(ctx *CdeclFuncIdentContext) {}

// ExitCdeclFuncIdent is called when production cdeclFuncIdent is exited.
func (s *BasedmlListener) ExitCdeclFuncIdent(ctx *CdeclFuncIdentContext) {}

// EnterCdeclArray is called when production cdeclArray is entered.
func (s *BasedmlListener) EnterCdeclArray(ctx *CdeclArrayContext) {}

// ExitCdeclArray is called when production cdeclArray is exited.
func (s *BasedmlListener) ExitCdeclArray(ctx *CdeclArrayContext) {}

// EnterCdeclNone is called when production cdeclNone is entered.
func (s *BasedmlListener) EnterCdeclNone(ctx *CdeclNoneContext) {}

// ExitCdeclNone is called when production cdeclNone is exited.
func (s *BasedmlListener) ExitCdeclNone(ctx *CdeclNoneContext) {}

// EnterCdeclFuncArgs is called when production cdeclFuncArgs is entered.
func (s *BasedmlListener) EnterCdeclFuncArgs(ctx *CdeclFuncArgsContext) {}

// ExitCdeclFuncArgs is called when production cdeclFuncArgs is exited.
func (s *BasedmlListener) ExitCdeclFuncArgs(ctx *CdeclFuncArgsContext) {}

// EnterCdecl_list is called when production cdecl_list is entered.
func (s *BasedmlListener) EnterCdecl_list(ctx *Cdecl_listContext) {}

// ExitCdecl_list is called when production cdecl_list is exited.
func (s *BasedmlListener) ExitCdecl_list(ctx *Cdecl_listContext) {}

// EnterCdecl_list_variadic is called when production cdecl_list_variadic is entered.
func (s *BasedmlListener) EnterCdecl_list_variadic(ctx *Cdecl_list_variadicContext) {}

// ExitCdecl_list_variadic is called when production cdecl_list_variadic is exited.
func (s *BasedmlListener) ExitCdecl_list_variadic(ctx *Cdecl_list_variadicContext) {}

// EnterTypeof is called when production typeof is entered.
func (s *BasedmlListener) EnterTypeof(ctx *TypeofContext) {}

// ExitTypeof is called when production typeof is exited.
func (s *BasedmlListener) ExitTypeof(ctx *TypeofContext) {}

// EnterStructDecls is called when production structDecls is entered.
func (s *BasedmlListener) EnterStructDecls(ctx *StructDeclsContext) {}

// ExitStructDecls is called when production structDecls is exited.
func (s *BasedmlListener) ExitStructDecls(ctx *StructDeclsContext) {}

// EnterStruct is called when production struct is entered.
func (s *BasedmlListener) EnterStruct(ctx *StructContext) {}

// ExitStruct is called when production struct is exited.
func (s *BasedmlListener) ExitStruct(ctx *StructContext) {}

// EnterUnion is called when production union is entered.
func (s *BasedmlListener) EnterUnion(ctx *UnionContext) {}

// ExitUnion is called when production union is exited.
func (s *BasedmlListener) ExitUnion(ctx *UnionContext) {}

// EnterLayoutDecls is called when production layoutDecls is entered.
func (s *BasedmlListener) EnterLayoutDecls(ctx *LayoutDeclsContext) {}

// ExitLayoutDecls is called when production layoutDecls is exited.
func (s *BasedmlListener) ExitLayoutDecls(ctx *LayoutDeclsContext) {}

// EnterLayout is called when production layout is entered.
func (s *BasedmlListener) EnterLayout(ctx *LayoutContext) {}

// ExitLayout is called when production layout is exited.
func (s *BasedmlListener) ExitLayout(ctx *LayoutContext) {}

// EnterBitfields is called when production bitfields is entered.
func (s *BasedmlListener) EnterBitfields(ctx *BitfieldsContext) {}

// ExitBitfields is called when production bitfields is exited.
func (s *BasedmlListener) ExitBitfields(ctx *BitfieldsContext) {}

// EnterBitfields_decls is called when production bitfields_decls is entered.
func (s *BasedmlListener) EnterBitfields_decls(ctx *Bitfields_declsContext) {}

// ExitBitfields_decls is called when production bitfields_decls is exited.
func (s *BasedmlListener) ExitBitfields_decls(ctx *Bitfields_declsContext) {}

// EnterCtypedecl is called when production ctypedecl is entered.
func (s *BasedmlListener) EnterCtypedecl(ctx *CtypedeclContext) {}

// ExitCtypedecl is called when production ctypedecl is exited.
func (s *BasedmlListener) ExitCtypedecl(ctx *CtypedeclContext) {}

// EnterCtypedecl_ptr is called when production ctypedecl_ptr is entered.
func (s *BasedmlListener) EnterCtypedecl_ptr(ctx *Ctypedecl_ptrContext) {}

// ExitCtypedecl_ptr is called when production ctypedecl_ptr is exited.
func (s *BasedmlListener) ExitCtypedecl_ptr(ctx *Ctypedecl_ptrContext) {}

// EnterCtypedecl_simple is called when production ctypedecl_simple is entered.
func (s *BasedmlListener) EnterCtypedecl_simple(ctx *Ctypedecl_simpleContext) {}

// ExitCtypedecl_simple is called when production ctypedecl_simple is exited.
func (s *BasedmlListener) ExitCtypedecl_simple(ctx *Ctypedecl_simpleContext) {}

// EnterTypeident is called when production typeident is entered.
func (s *BasedmlListener) EnterTypeident(ctx *TypeidentContext) {}

// ExitTypeident is called when production typeident is exited.
func (s *BasedmlListener) ExitTypeident(ctx *TypeidentContext) {}

// EnterBitSlice is called when production bitSlice is entered.
func (s *BasedmlListener) EnterBitSlice(ctx *BitSliceContext) {}

// ExitBitSlice is called when production bitSlice is exited.
func (s *BasedmlListener) ExitBitSlice(ctx *BitSliceContext) {}

// EnterAssignop is called when production assignop is entered.
func (s *BasedmlListener) EnterAssignop(ctx *AssignopContext) {}

// ExitAssignop is called when production assignop is exited.
func (s *BasedmlListener) ExitAssignop(ctx *AssignopContext) {}

// EnterNewExpr is called when production newExpr is entered.
func (s *BasedmlListener) EnterNewExpr(ctx *NewExprContext) {}

// ExitNewExpr is called when production newExpr is exited.
func (s *BasedmlListener) ExitNewExpr(ctx *NewExprContext) {}

// EnterTernaryConstIfExpr is called when production ternaryConstIfExpr is entered.
func (s *BasedmlListener) EnterTernaryConstIfExpr(ctx *TernaryConstIfExprContext) {}

// ExitTernaryConstIfExpr is called when production ternaryConstIfExpr is exited.
func (s *BasedmlListener) ExitTernaryConstIfExpr(ctx *TernaryConstIfExprContext) {}

// EnterCastExpr is called when production castExpr is entered.
func (s *BasedmlListener) EnterCastExpr(ctx *CastExprContext) {}

// ExitCastExpr is called when production castExpr is exited.
func (s *BasedmlListener) ExitCastExpr(ctx *CastExprContext) {}

// EnterDefinedExpr is called when production definedExpr is entered.
func (s *BasedmlListener) EnterDefinedExpr(ctx *DefinedExprContext) {}

// ExitDefinedExpr is called when production definedExpr is exited.
func (s *BasedmlListener) ExitDefinedExpr(ctx *DefinedExprContext) {}

// EnterLogicalExpr is called when production logicalExpr is entered.
func (s *BasedmlListener) EnterLogicalExpr(ctx *LogicalExprContext) {}

// ExitLogicalExpr is called when production logicalExpr is exited.
func (s *BasedmlListener) ExitLogicalExpr(ctx *LogicalExprContext) {}

// EnterBitSliceExpr is called when production bitSliceExpr is entered.
func (s *BasedmlListener) EnterBitSliceExpr(ctx *BitSliceExprContext) {}

// ExitBitSliceExpr is called when production bitSliceExpr is exited.
func (s *BasedmlListener) ExitBitSliceExpr(ctx *BitSliceExprContext) {}

// EnterEachExpr is called when production eachExpr is entered.
func (s *BasedmlListener) EnterEachExpr(ctx *EachExprContext) {}

// ExitEachExpr is called when production eachExpr is exited.
func (s *BasedmlListener) ExitEachExpr(ctx *EachExprContext) {}

// EnterAdditiveExpr is called when production additiveExpr is entered.
func (s *BasedmlListener) EnterAdditiveExpr(ctx *AdditiveExprContext) {}

// ExitAdditiveExpr is called when production additiveExpr is exited.
func (s *BasedmlListener) ExitAdditiveExpr(ctx *AdditiveExprContext) {}

// EnterBitwiseNotExpr is called when production bitwiseNotExpr is entered.
func (s *BasedmlListener) EnterBitwiseNotExpr(ctx *BitwiseNotExprContext) {}

// ExitBitwiseNotExpr is called when production bitwiseNotExpr is exited.
func (s *BasedmlListener) ExitBitwiseNotExpr(ctx *BitwiseNotExprContext) {}

// EnterParenExpr is called when production parenExpr is entered.
func (s *BasedmlListener) EnterParenExpr(ctx *ParenExprContext) {}

// ExitParenExpr is called when production parenExpr is exited.
func (s *BasedmlListener) ExitParenExpr(ctx *ParenExprContext) {}

// EnterUnaryExpr is called when production unaryExpr is entered.
func (s *BasedmlListener) EnterUnaryExpr(ctx *UnaryExprContext) {}

// ExitUnaryExpr is called when production unaryExpr is exited.
func (s *BasedmlListener) ExitUnaryExpr(ctx *UnaryExprContext) {}

// EnterArrayAccessExpr is called when production arrayAccessExpr is entered.
func (s *BasedmlListener) EnterArrayAccessExpr(ctx *ArrayAccessExprContext) {}

// ExitArrayAccessExpr is called when production arrayAccessExpr is exited.
func (s *BasedmlListener) ExitArrayAccessExpr(ctx *ArrayAccessExprContext) {}

// EnterFunctionCallExpr is called when production functionCallExpr is entered.
func (s *BasedmlListener) EnterFunctionCallExpr(ctx *FunctionCallExprContext) {}

// ExitFunctionCallExpr is called when production functionCallExpr is exited.
func (s *BasedmlListener) ExitFunctionCallExpr(ctx *FunctionCallExprContext) {}

// EnterMemberAccessExpr is called when production memberAccessExpr is entered.
func (s *BasedmlListener) EnterMemberAccessExpr(ctx *MemberAccessExprContext) {}

// ExitMemberAccessExpr is called when production memberAccessExpr is exited.
func (s *BasedmlListener) ExitMemberAccessExpr(ctx *MemberAccessExprContext) {}

// EnterTernaryIfExpr is called when production ternaryIfExpr is entered.
func (s *BasedmlListener) EnterTernaryIfExpr(ctx *TernaryIfExprContext) {}

// ExitTernaryIfExpr is called when production ternaryIfExpr is exited.
func (s *BasedmlListener) ExitTernaryIfExpr(ctx *TernaryIfExprContext) {}

// EnterPostIncrementExpr is called when production postIncrementExpr is entered.
func (s *BasedmlListener) EnterPostIncrementExpr(ctx *PostIncrementExprContext) {}

// ExitPostIncrementExpr is called when production postIncrementExpr is exited.
func (s *BasedmlListener) ExitPostIncrementExpr(ctx *PostIncrementExprContext) {}

// EnterRelationalExpr is called when production relationalExpr is entered.
func (s *BasedmlListener) EnterRelationalExpr(ctx *RelationalExprContext) {}

// ExitRelationalExpr is called when production relationalExpr is exited.
func (s *BasedmlListener) ExitRelationalExpr(ctx *RelationalExprContext) {}

// EnterPreIncrementExpr is called when production preIncrementExpr is entered.
func (s *BasedmlListener) EnterPreIncrementExpr(ctx *PreIncrementExprContext) {}

// ExitPreIncrementExpr is called when production preIncrementExpr is exited.
func (s *BasedmlListener) ExitPreIncrementExpr(ctx *PreIncrementExprContext) {}

// EnterSizeoftypeExpr is called when production sizeoftypeExpr is entered.
func (s *BasedmlListener) EnterSizeoftypeExpr(ctx *SizeoftypeExprContext) {}

// ExitSizeoftypeExpr is called when production sizeoftypeExpr is exited.
func (s *BasedmlListener) ExitSizeoftypeExpr(ctx *SizeoftypeExprContext) {}

// EnterShiftExpr is called when production shiftExpr is entered.
func (s *BasedmlListener) EnterShiftExpr(ctx *ShiftExprContext) {}

// ExitShiftExpr is called when production shiftExpr is exited.
func (s *BasedmlListener) ExitShiftExpr(ctx *ShiftExprContext) {}

// EnterAddressofExpr is called when production addressofExpr is entered.
func (s *BasedmlListener) EnterAddressofExpr(ctx *AddressofExprContext) {}

// ExitAddressofExpr is called when production addressofExpr is exited.
func (s *BasedmlListener) ExitAddressofExpr(ctx *AddressofExprContext) {}

// EnterBitwiseExpr is called when production bitwiseExpr is entered.
func (s *BasedmlListener) EnterBitwiseExpr(ctx *BitwiseExprContext) {}

// ExitBitwiseExpr is called when production bitwiseExpr is exited.
func (s *BasedmlListener) ExitBitwiseExpr(ctx *BitwiseExprContext) {}

// EnterSizeofExpr is called when production sizeofExpr is entered.
func (s *BasedmlListener) EnterSizeofExpr(ctx *SizeofExprContext) {}

// ExitSizeofExpr is called when production sizeofExpr is exited.
func (s *BasedmlListener) ExitSizeofExpr(ctx *SizeofExprContext) {}

// EnterLiteralAtom is called when production literalAtom is entered.
func (s *BasedmlListener) EnterLiteralAtom(ctx *LiteralAtomContext) {}

// ExitLiteralAtom is called when production literalAtom is exited.
func (s *BasedmlListener) ExitLiteralAtom(ctx *LiteralAtomContext) {}

// EnterNotExpr is called when production notExpr is entered.
func (s *BasedmlListener) EnterNotExpr(ctx *NotExprContext) {}

// ExitNotExpr is called when production notExpr is exited.
func (s *BasedmlListener) ExitNotExpr(ctx *NotExprContext) {}

// EnterDerefExpr is called when production derefExpr is entered.
func (s *BasedmlListener) EnterDerefExpr(ctx *DerefExprContext) {}

// ExitDerefExpr is called when production derefExpr is exited.
func (s *BasedmlListener) ExitDerefExpr(ctx *DerefExprContext) {}

// EnterStringifyExpr is called when production stringifyExpr is entered.
func (s *BasedmlListener) EnterStringifyExpr(ctx *StringifyExprContext) {}

// ExitStringifyExpr is called when production stringifyExpr is exited.
func (s *BasedmlListener) ExitStringifyExpr(ctx *StringifyExprContext) {}

// EnterArrayLiteralExpr is called when production arrayLiteralExpr is entered.
func (s *BasedmlListener) EnterArrayLiteralExpr(ctx *ArrayLiteralExprContext) {}

// ExitArrayLiteralExpr is called when production arrayLiteralExpr is exited.
func (s *BasedmlListener) ExitArrayLiteralExpr(ctx *ArrayLiteralExprContext) {}

// EnterMultiplicativeExpr is called when production multiplicativeExpr is entered.
func (s *BasedmlListener) EnterMultiplicativeExpr(ctx *MultiplicativeExprContext) {}

// ExitMultiplicativeExpr is called when production multiplicativeExpr is exited.
func (s *BasedmlListener) ExitMultiplicativeExpr(ctx *MultiplicativeExprContext) {}

// EnterEqualityExpr is called when production equalityExpr is entered.
func (s *BasedmlListener) EnterEqualityExpr(ctx *EqualityExprContext) {}

// ExitEqualityExpr is called when production equalityExpr is exited.
func (s *BasedmlListener) ExitEqualityExpr(ctx *EqualityExprContext) {}

// EnterNumericLiteral is called when production numericLiteral is entered.
func (s *BasedmlListener) EnterNumericLiteral(ctx *NumericLiteralContext) {}

// ExitNumericLiteral is called when production numericLiteral is exited.
func (s *BasedmlListener) ExitNumericLiteral(ctx *NumericLiteralContext) {}

// EnterStringLiteral is called when production stringLiteral is entered.
func (s *BasedmlListener) EnterStringLiteral(ctx *StringLiteralContext) {}

// ExitStringLiteral is called when production stringLiteral is exited.
func (s *BasedmlListener) ExitStringLiteral(ctx *StringLiteralContext) {}

// EnterExpressionAccessIdent is called when production expressionAccessIdent is entered.
func (s *BasedmlListener) EnterExpressionAccessIdent(ctx *ExpressionAccessIdentContext) {}

// ExitExpressionAccessIdent is called when production expressionAccessIdent is exited.
func (s *BasedmlListener) ExitExpressionAccessIdent(ctx *ExpressionAccessIdentContext) {}

// EnterExpressionIdent is called when production expressionIdent is entered.
func (s *BasedmlListener) EnterExpressionIdent(ctx *ExpressionIdentContext) {}

// ExitExpressionIdent is called when production expressionIdent is exited.
func (s *BasedmlListener) ExitExpressionIdent(ctx *ExpressionIdentContext) {}

// EnterUtf8_sconst is called when production utf8_sconst is entered.
func (s *BasedmlListener) EnterUtf8_sconst(ctx *Utf8_sconstContext) {}

// ExitUtf8_sconst is called when production utf8_sconst is exited.
func (s *BasedmlListener) ExitUtf8_sconst(ctx *Utf8_sconstContext) {}

// EnterTypeoparg is called when production typeoparg is entered.
func (s *BasedmlListener) EnterTypeoparg(ctx *TypeopargContext) {}

// ExitTypeoparg is called when production typeoparg is exited.
func (s *BasedmlListener) ExitTypeoparg(ctx *TypeopargContext) {}

// EnterEndianflag is called when production endianflag is entered.
func (s *BasedmlListener) EnterEndianflag(ctx *EndianflagContext) {}

// ExitEndianflag is called when production endianflag is exited.
func (s *BasedmlListener) ExitEndianflag(ctx *EndianflagContext) {}

// EnterExpression_list is called when production expression_list is entered.
func (s *BasedmlListener) EnterExpression_list(ctx *Expression_listContext) {}

// ExitExpression_list is called when production expression_list is exited.
func (s *BasedmlListener) ExitExpression_list(ctx *Expression_listContext) {}

// EnterExpression_list_ntc_nonempty is called when production expression_list_ntc_nonempty is entered.
func (s *BasedmlListener) EnterExpression_list_ntc_nonempty(ctx *Expression_list_ntc_nonemptyContext) {
}

// ExitExpression_list_ntc_nonempty is called when production expression_list_ntc_nonempty is exited.
func (s *BasedmlListener) ExitExpression_list_ntc_nonempty(ctx *Expression_list_ntc_nonemptyContext) {
}

// EnterComposed_string_literal is called when production composed_string_literal is entered.
func (s *BasedmlListener) EnterComposed_string_literal(ctx *Composed_string_literalContext) {}

// ExitComposed_string_literal is called when production composed_string_literal is exited.
func (s *BasedmlListener) ExitComposed_string_literal(ctx *Composed_string_literalContext) {}

// EnterBracketed_string_literal is called when production bracketed_string_literal is entered.
func (s *BasedmlListener) EnterBracketed_string_literal(ctx *Bracketed_string_literalContext) {}

// ExitBracketed_string_literal is called when production bracketed_string_literal is exited.
func (s *BasedmlListener) ExitBracketed_string_literal(ctx *Bracketed_string_literalContext) {}

// EnterSingleInitializer is called when production singleInitializer is entered.
func (s *BasedmlListener) EnterSingleInitializer(ctx *SingleInitializerContext) {}

// ExitSingleInitializer is called when production singleInitializer is exited.
func (s *BasedmlListener) ExitSingleInitializer(ctx *SingleInitializerContext) {}

// EnterDesignatedStructInitializer is called when production designatedStructInitializer is entered.
func (s *BasedmlListener) EnterDesignatedStructInitializer(ctx *DesignatedStructInitializerContext) {}

// ExitDesignatedStructInitializer is called when production designatedStructInitializer is exited.
func (s *BasedmlListener) ExitDesignatedStructInitializer(ctx *DesignatedStructInitializerContext) {}

// EnterDesignated_struct_initializer_list is called when production designated_struct_initializer_list is entered.
func (s *BasedmlListener) EnterDesignated_struct_initializer_list(ctx *Designated_struct_initializer_listContext) {
}

// ExitDesignated_struct_initializer_list is called when production designated_struct_initializer_list is exited.
func (s *BasedmlListener) ExitDesignated_struct_initializer_list(ctx *Designated_struct_initializer_listContext) {
}

// EnterCdeclInitializer is called when production cdeclInitializer is entered.
func (s *BasedmlListener) EnterCdeclInitializer(ctx *CdeclInitializerContext) {}

// ExitCdeclInitializer is called when production cdeclInitializer is exited.
func (s *BasedmlListener) ExitCdeclInitializer(ctx *CdeclInitializerContext) {}

// EnterSingleInitializerArray is called when production singleInitializerArray is entered.
func (s *BasedmlListener) EnterSingleInitializerArray(ctx *SingleInitializerArrayContext) {}

// ExitSingleInitializerArray is called when production singleInitializerArray is exited.
func (s *BasedmlListener) ExitSingleInitializerArray(ctx *SingleInitializerArrayContext) {}

// EnterInitializer is called when production initializer is entered.
func (s *BasedmlListener) EnterInitializer(ctx *InitializerContext) {}

// ExitInitializer is called when production initializer is exited.
func (s *BasedmlListener) ExitInitializer(ctx *InitializerContext) {}

// EnterArgumentsList is called when production argumentsList is entered.
func (s *BasedmlListener) EnterArgumentsList(ctx *ArgumentsListContext) {}

// ExitArgumentsList is called when production argumentsList is exited.
func (s *BasedmlListener) ExitArgumentsList(ctx *ArgumentsListContext) {}

// EnterArgumentsListEmpty is called when production argumentsListEmpty is entered.
func (s *BasedmlListener) EnterArgumentsListEmpty(ctx *ArgumentsListEmptyContext) {}

// ExitArgumentsListEmpty is called when production argumentsListEmpty is exited.
func (s *BasedmlListener) ExitArgumentsListEmpty(ctx *ArgumentsListEmptyContext) {}

// EnterStatement is called when production statement is entered.
func (s *BasedmlListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BasedmlListener) ExitStatement(ctx *StatementContext) {}

// EnterThrowStatement is called when production throwStatement is entered.
func (s *BasedmlListener) EnterThrowStatement(ctx *ThrowStatementContext) {}

// ExitThrowStatement is called when production throwStatement is exited.
func (s *BasedmlListener) ExitThrowStatement(ctx *ThrowStatementContext) {}

// EnterContinueStatement is called when production continueStatement is entered.
func (s *BasedmlListener) EnterContinueStatement(ctx *ContinueStatementContext) {}

// ExitContinueStatement is called when production continueStatement is exited.
func (s *BasedmlListener) ExitContinueStatement(ctx *ContinueStatementContext) {}

// EnterBreakStatement is called when production breakStatement is entered.
func (s *BasedmlListener) EnterBreakStatement(ctx *BreakStatementContext) {}

// ExitBreakStatement is called when production breakStatement is exited.
func (s *BasedmlListener) ExitBreakStatement(ctx *BreakStatementContext) {}

// EnterElseStatement is called when production elseStatement is entered.
func (s *BasedmlListener) EnterElseStatement(ctx *ElseStatementContext) {}

// ExitElseStatement is called when production elseStatement is exited.
func (s *BasedmlListener) ExitElseStatement(ctx *ElseStatementContext) {}

// EnterLoopStatement is called when production loopStatement is entered.
func (s *BasedmlListener) EnterLoopStatement(ctx *LoopStatementContext) {}

// ExitLoopStatement is called when production loopStatement is exited.
func (s *BasedmlListener) ExitLoopStatement(ctx *LoopStatementContext) {}

// EnterConditionExpression is called when production conditionExpression is entered.
func (s *BasedmlListener) EnterConditionExpression(ctx *ConditionExpressionContext) {}

// ExitConditionExpression is called when production conditionExpression is exited.
func (s *BasedmlListener) ExitConditionExpression(ctx *ConditionExpressionContext) {}

// EnterAssertStatement is called when production assertStatement is entered.
func (s *BasedmlListener) EnterAssertStatement(ctx *AssertStatementContext) {}

// ExitAssertStatement is called when production assertStatement is exited.
func (s *BasedmlListener) ExitAssertStatement(ctx *AssertStatementContext) {}

// EnterDeleteStatement is called when production deleteStatement is entered.
func (s *BasedmlListener) EnterDeleteStatement(ctx *DeleteStatementContext) {}

// ExitDeleteStatement is called when production deleteStatement is exited.
func (s *BasedmlListener) ExitDeleteStatement(ctx *DeleteStatementContext) {}

// EnterReturnStatement is called when production returnStatement is entered.
func (s *BasedmlListener) EnterReturnStatement(ctx *ReturnStatementContext) {}

// ExitReturnStatement is called when production returnStatement is exited.
func (s *BasedmlListener) ExitReturnStatement(ctx *ReturnStatementContext) {}

// EnterGotoStatement is called when production gotoStatement is entered.
func (s *BasedmlListener) EnterGotoStatement(ctx *GotoStatementContext) {}

// ExitGotoStatement is called when production gotoStatement is exited.
func (s *BasedmlListener) ExitGotoStatement(ctx *GotoStatementContext) {}

// EnterLogKind is called when production logKind is entered.
func (s *BasedmlListener) EnterLogKind(ctx *LogKindContext) {}

// ExitLogKind is called when production logKind is exited.
func (s *BasedmlListener) ExitLogKind(ctx *LogKindContext) {}

// EnterLogGroup is called when production logGroup is entered.
func (s *BasedmlListener) EnterLogGroup(ctx *LogGroupContext) {}

// ExitLogGroup is called when production logGroup is exited.
func (s *BasedmlListener) ExitLogGroup(ctx *LogGroupContext) {}

// EnterLogFormatString is called when production logFormatString is entered.
func (s *BasedmlListener) EnterLogFormatString(ctx *LogFormatStringContext) {}

// ExitLogFormatString is called when production logFormatString is exited.
func (s *BasedmlListener) ExitLogFormatString(ctx *LogFormatStringContext) {}

// EnterThenLevel is called when production thenLevel is entered.
func (s *BasedmlListener) EnterThenLevel(ctx *ThenLevelContext) {}

// ExitThenLevel is called when production thenLevel is exited.
func (s *BasedmlListener) ExitThenLevel(ctx *ThenLevelContext) {}

// EnterLogLevel is called when production logLevel is entered.
func (s *BasedmlListener) EnterLogLevel(ctx *LogLevelContext) {}

// ExitLogLevel is called when production logLevel is exited.
func (s *BasedmlListener) ExitLogLevel(ctx *LogLevelContext) {}

// EnterLogStatement is called when production logStatement is entered.
func (s *BasedmlListener) EnterLogStatement(ctx *LogStatementContext) {}

// ExitLogStatement is called when production logStatement is exited.
func (s *BasedmlListener) ExitLogStatement(ctx *LogStatementContext) {}

// EnterCaseOrDefaultStatement is called when production caseOrDefaultStatement is entered.
func (s *BasedmlListener) EnterCaseOrDefaultStatement(ctx *CaseOrDefaultStatementContext) {}

// ExitCaseOrDefaultStatement is called when production caseOrDefaultStatement is exited.
func (s *BasedmlListener) ExitCaseOrDefaultStatement(ctx *CaseOrDefaultStatementContext) {}

// EnterCaseStatement is called when production caseStatement is entered.
func (s *BasedmlListener) EnterCaseStatement(ctx *CaseStatementContext) {}

// ExitCaseStatement is called when production caseStatement is exited.
func (s *BasedmlListener) ExitCaseStatement(ctx *CaseStatementContext) {}

// EnterSwitchBodyStatement is called when production switchBodyStatement is entered.
func (s *BasedmlListener) EnterSwitchBodyStatement(ctx *SwitchBodyStatementContext) {}

// ExitSwitchBodyStatement is called when production switchBodyStatement is exited.
func (s *BasedmlListener) ExitSwitchBodyStatement(ctx *SwitchBodyStatementContext) {}

// EnterSwitchBodyHashElseStatement is called when production switchBodyHashElseStatement is entered.
func (s *BasedmlListener) EnterSwitchBodyHashElseStatement(ctx *SwitchBodyHashElseStatementContext) {}

// ExitSwitchBodyHashElseStatement is called when production switchBodyHashElseStatement is exited.
func (s *BasedmlListener) ExitSwitchBodyHashElseStatement(ctx *SwitchBodyHashElseStatementContext) {}

// EnterSwitchBodyHashIfStatement is called when production switchBodyHashIfStatement is entered.
func (s *BasedmlListener) EnterSwitchBodyHashIfStatement(ctx *SwitchBodyHashIfStatementContext) {}

// ExitSwitchBodyHashIfStatement is called when production switchBodyHashIfStatement is exited.
func (s *BasedmlListener) ExitSwitchBodyHashIfStatement(ctx *SwitchBodyHashIfStatementContext) {}

// EnterSwitchStatement is called when production switchStatement is entered.
func (s *BasedmlListener) EnterSwitchStatement(ctx *SwitchStatementContext) {}

// ExitSwitchStatement is called when production switchStatement is exited.
func (s *BasedmlListener) ExitSwitchStatement(ctx *SwitchStatementContext) {}

// EnterForLoopStepStatement is called when production forLoopStepStatement is entered.
func (s *BasedmlListener) EnterForLoopStepStatement(ctx *ForLoopStepStatementContext) {}

// ExitForLoopStepStatement is called when production forLoopStepStatement is exited.
func (s *BasedmlListener) ExitForLoopStepStatement(ctx *ForLoopStepStatementContext) {}

// EnterInnerForLoopStatement is called when production innerForLoopStatement is entered.
func (s *BasedmlListener) EnterInnerForLoopStatement(ctx *InnerForLoopStatementContext) {}

// ExitInnerForLoopStatement is called when production innerForLoopStatement is exited.
func (s *BasedmlListener) ExitInnerForLoopStatement(ctx *InnerForLoopStatementContext) {}

// EnterForLoopInitStatement is called when production forLoopInitStatement is entered.
func (s *BasedmlListener) EnterForLoopInitStatement(ctx *ForLoopInitStatementContext) {}

// ExitForLoopInitStatement is called when production forLoopInitStatement is exited.
func (s *BasedmlListener) ExitForLoopInitStatement(ctx *ForLoopInitStatementContext) {}

// EnterForLoopStatement is called when production forLoopStatement is entered.
func (s *BasedmlListener) EnterForLoopStatement(ctx *ForLoopStatementContext) {}

// ExitForLoopStatement is called when production forLoopStatement is exited.
func (s *BasedmlListener) ExitForLoopStatement(ctx *ForLoopStatementContext) {}

// EnterWhileLoopStatement is called when production whileLoopStatement is entered.
func (s *BasedmlListener) EnterWhileLoopStatement(ctx *WhileLoopStatementContext) {}

// ExitWhileLoopStatement is called when production whileLoopStatement is exited.
func (s *BasedmlListener) ExitWhileLoopStatement(ctx *WhileLoopStatementContext) {}

// EnterDoWhileLoopStatement is called when production doWhileLoopStatement is entered.
func (s *BasedmlListener) EnterDoWhileLoopStatement(ctx *DoWhileLoopStatementContext) {}

// ExitDoWhileLoopStatement is called when production doWhileLoopStatement is exited.
func (s *BasedmlListener) ExitDoWhileLoopStatement(ctx *DoWhileLoopStatementContext) {}

// EnterForeachStatement is called when production foreachStatement is entered.
func (s *BasedmlListener) EnterForeachStatement(ctx *ForeachStatementContext) {}

// ExitForeachStatement is called when production foreachStatement is exited.
func (s *BasedmlListener) ExitForeachStatement(ctx *ForeachStatementContext) {}

// EnterIfStatement is called when production ifStatement is entered.
func (s *BasedmlListener) EnterIfStatement(ctx *IfStatementContext) {}

// ExitIfStatement is called when production ifStatement is exited.
func (s *BasedmlListener) ExitIfStatement(ctx *IfStatementContext) {}

// EnterCatchStatement is called when production catchStatement is entered.
func (s *BasedmlListener) EnterCatchStatement(ctx *CatchStatementContext) {}

// ExitCatchStatement is called when production catchStatement is exited.
func (s *BasedmlListener) ExitCatchStatement(ctx *CatchStatementContext) {}

// EnterTryStatement is called when production tryStatement is entered.
func (s *BasedmlListener) EnterTryStatement(ctx *TryStatementContext) {}

// ExitTryStatement is called when production tryStatement is exited.
func (s *BasedmlListener) ExitTryStatement(ctx *TryStatementContext) {}

// EnterSelectStatement is called when production selectStatement is entered.
func (s *BasedmlListener) EnterSelectStatement(ctx *SelectStatementContext) {}

// ExitSelectStatement is called when production selectStatement is exited.
func (s *BasedmlListener) ExitSelectStatement(ctx *SelectStatementContext) {}

// EnterAssignStatement is called when production assignStatement is entered.
func (s *BasedmlListener) EnterAssignStatement(ctx *AssignStatementContext) {}

// ExitAssignStatement is called when production assignStatement is exited.
func (s *BasedmlListener) ExitAssignStatement(ctx *AssignStatementContext) {}

// EnterAssign_chain is called when production assign_chain is entered.
func (s *BasedmlListener) EnterAssign_chain(ctx *Assign_chainContext) {}

// ExitAssign_chain is called when production assign_chain is exited.
func (s *BasedmlListener) ExitAssign_chain(ctx *Assign_chainContext) {}

// EnterAfterStatement is called when production afterStatement is entered.
func (s *BasedmlListener) EnterAfterStatement(ctx *AfterStatementContext) {}

// ExitAfterStatement is called when production afterStatement is exited.
func (s *BasedmlListener) ExitAfterStatement(ctx *AfterStatementContext) {}

// EnterAfterDelayCondition is called when production afterDelayCondition is entered.
func (s *BasedmlListener) EnterAfterDelayCondition(ctx *AfterDelayConditionContext) {}

// ExitAfterDelayCondition is called when production afterDelayCondition is exited.
func (s *BasedmlListener) ExitAfterDelayCondition(ctx *AfterDelayConditionContext) {}

// EnterTuple_literal is called when production tuple_literal is entered.
func (s *BasedmlListener) EnterTuple_literal(ctx *Tuple_literalContext) {}

// ExitTuple_literal is called when production tuple_literal is exited.
func (s *BasedmlListener) ExitTuple_literal(ctx *Tuple_literalContext) {}

// EnterIdentList is called when production identList is entered.
func (s *BasedmlListener) EnterIdentList(ctx *IdentListContext) {}

// ExitIdentList is called when production identList is exited.
func (s *BasedmlListener) ExitIdentList(ctx *IdentListContext) {}

// EnterErrorStatement is called when production errorStatement is entered.
func (s *BasedmlListener) EnterErrorStatement(ctx *ErrorStatementContext) {}

// ExitErrorStatement is called when production errorStatement is exited.
func (s *BasedmlListener) ExitErrorStatement(ctx *ErrorStatementContext) {}

// EnterWarningStatement is called when production warningStatement is entered.
func (s *BasedmlListener) EnterWarningStatement(ctx *WarningStatementContext) {}

// ExitWarningStatement is called when production warningStatement is exited.
func (s *BasedmlListener) ExitWarningStatement(ctx *WarningStatementContext) {}

// EnterCompoundStatement is called when production compoundStatement is entered.
func (s *BasedmlListener) EnterCompoundStatement(ctx *CompoundStatementContext) {}

// ExitCompoundStatement is called when production compoundStatement is exited.
func (s *BasedmlListener) ExitCompoundStatement(ctx *CompoundStatementContext) {}

// EnterHook_decl is called when production hook_decl is entered.
func (s *BasedmlListener) EnterHook_decl(ctx *Hook_declContext) {}

// ExitHook_decl is called when production hook_decl is exited.
func (s *BasedmlListener) ExitHook_decl(ctx *Hook_declContext) {}

// EnterObjident is called when production objident is entered.
func (s *BasedmlListener) EnterObjident(ctx *ObjidentContext) {}

// ExitObjident is called when production objident is exited.
func (s *BasedmlListener) ExitObjident(ctx *ObjidentContext) {}

// EnterIdent is called when production ident is entered.
func (s *BasedmlListener) EnterIdent(ctx *IdentContext) {}

// ExitIdent is called when production ident is exited.
func (s *BasedmlListener) ExitIdent(ctx *IdentContext) {}
