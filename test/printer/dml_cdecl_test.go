// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintTypes_Basic(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "101-basic-int",
			code:     "typedef int my_int;\n",
			expected: "typedef int my_int;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "102-messy-whitespace",
			code:     "typedef   int    my_int  ;\n",
			expected: "typedef int my_int;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "103-new-lines",
			code:     "typedef   \nint  \n  my_int \n ;\n",
			expected: "typedef int my_int;\n",
			pc:       func(f func()) { f() },
		},

		{
			name:     "104-very-long-line-not-wrapped",
			code:     "typedef    int   this_is_very_long_line;\n",
			expected: "typedef int this_is_very_long_line;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "201-const-int",
			code:     "typedef const int const_intv;\n",
			expected: "typedef const int const_intv;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "202-messy-whitespace",
			code:     "typedef   const   int    const_intv  ;\n",
			expected: "typedef const int const_intv;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "203-new-lines",
			code:     "typedef   \nconst  \n  int  \n const_intv \n ;\n",
			expected: "typedef const int const_intv;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "204-very-long-line-not-wrapped",
			code:     "typedef    const   int   this_is_very_long_line_for_an_const_long;\n",
			expected: "typedef const int this_is_very_long_line_for_an_const_long;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "301-unsigned-int",
			code:     "typedef unsigned    int usig_intv;\n",
			expected: "typedef unsigned int usig_intv;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "302-messy-whitespace",
			code:     "typedef   unsigned   int    usig_intv  ;\n",
			expected: "typedef unsigned int usig_intv;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "303-new-lines",
			code:     "typedef   \nunsigned  \n  int  \n usig_intv \n ;\n",
			expected: "typedef unsigned int usig_intv;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "304-very-long-line-not-wrapped",
			code:     "typedef    unsigned   int   this_is_very_long_line_for_an_unsigned_long;\n",
			expected: "typedef unsigned int this_is_very_long_line_for_an_unsigned_long;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "401-signed-int",
			code:     "typedef signed int sig_intv;\n",
			expected: "typedef signed int sig_intv;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "402-messy-whitespace",
			code:     "typedef   signed   int    sig_intv  ;\n",
			expected: "typedef signed int sig_intv;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "403-new-lines",
			code:     "typedef   \nsigned  \n  int  \n sig_intv \n ;\n",
			expected: "typedef signed int sig_intv;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "404-very-long-line-not-wrapped",
			code:     "typedef    signed   int   this_is_very_long_line_for_an_signed_long;\n",
			expected: "typedef signed int this_is_very_long_line_for_an_signed_long;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "501-pointer",
			code:     "typedef int * my_int_ptr;\n",
			expected: "typedef int *my_int_ptr;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "502-messy-whitespace",
			code:     "typedef   int    *   my_int_ptr  ;\n",
			expected: "typedef int *my_int_ptr;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "503-new-lines",
			code:     "typedef   \nint  \n *  \n my_int_ptr \n ;\n",
			expected: "typedef int *my_int_ptr;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "504-very-long-line-not-wrapped",
			code:     "typedef    int   *   this_is_very_long_line_for_a_pointer;\n",
			expected: "typedef int *this_is_very_long_line_for_a_pointer;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "601-const-pointer",
			code:     "typedef const int * const my_int_ptr;\n",
			expected: "typedef const int *const my_int_ptr;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "602-messy-whitespace",
			code:     "typedef   const    int   *   const    my_int_ptr  ;\n",
			expected: "typedef const int *const my_int_ptr;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "603-new-lines",
			code:     "typedef   \nconst  \n int \n* \nconst \n my_int_ptr \n ;\n",
			expected: "typedef const int *const my_int_ptr;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "604-very-long-line-not-wrapped",
			code:     "typedef    const   int   *   const   this_is_very_long_line_for_a_const_pointer;\n",
			expected: "typedef const int *const this_is_very_long_line_for_a_const_pointer;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "701-const-pointer",
			code:     "typedef const foo_t ***const c_foo_ptr;\n",
			expected: "typedef const foo_t ***const c_foo_ptr;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintTypes_Array(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "101-array",
			code:     "typedef int my_array[10];\n",
			expected: "typedef int my_array[10];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "102-messy-whitespace",
			code:     "typedef   int    my_array  [  10  ]  ;\n",
			expected: "typedef int my_array[10];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "103-new-lines",
			code:     "typedef   \nint  \n  my_array \n [ \n 10 \n ] \n ;\n",
			expected: "typedef int my_array[10];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "104-very-long-line-not-wrapped",
			code:     "typedef    int   this_is_very_long_line_for_an_array_type[10];\n",
			expected: "typedef int this_is_very_long_line_for_an_array_type[\n    10\n];\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintTypes_Function(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "101-function",
			code:     "typedef void my_func(int, char);\n",
			expected: "typedef void my_func(int, char);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "102-messy-whitespace",
			code:     "typedef   void    my_func  (  int  ,   char  )  ;\n",
			expected: "typedef void my_func(int, char);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "103-new-lines",
			code:     "typedef   \nvoid  \n  my_func \n ( \n int \n , \n char \n ) \n ;\n",
			expected: "typedef void my_func(int, char);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "104-very-long-line-not-wrapped",
			code:     "typedef    void   this_is_very_long_line_for_a_function_type(int, char);\n",
			expected: "typedef void this_is_very_long_line_for_a_function_type(\n    int, char\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "201-function-pointer",
			code:     "typedef void (*my_func_name_ptr)(int, char);\n",
			expected: "typedef void (*my_func_name_ptr)(\n    int, char\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "202-messy-whitespace",
			code:     "typedef   void    (  *  my_func_name_ptr  )  (  int  ,   char  )  ;\n",
			expected: "typedef void (*my_func_name_ptr)(\n    int, char\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "203-new-lines",
			code:     "typedef   \nvoid  \n  ( \n * \n my_func_name_ptr \n ) \n ( \n int \n , \n char \n ) \n ;\n",
			expected: "typedef void (*my_func_name_ptr)(\n    int, char\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "204-very-long-line-not-wrapped",
			code:     "typedef    void   (*this_is_very_long_line_for_a_function_pointer)(int, char);\n",
			expected: "typedef void (*this_is_very_long_line_for_a_function_pointer)(\n    int, char\n);\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintTypes_Struct(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-empty-struct",
			code:     "typedef    struct {   }    my_struct_t;\n",
			expected: "typedef struct {\n} my_struct_t;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-single-field",
			code:     "typedef    struct    {  int   a;  } \nmy_struct_t;\n",
			expected: "typedef struct {\n    int a;\n} my_struct_t;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-multi-field",
			code:     "typedef   struct  {  int a;    float \n b;} my_struct_t;\n",
			expected: "typedef struct {\n    int a;\n    float b;\n} my_struct_t;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-messy-whitespace",
			code:     "typedef    struct {int  a  ;float b;} my_struct_t;\n",
			expected: "typedef struct {\n    int a;\n    float b;\n} my_struct_t;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintTypes_Union(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-empty-union",
			code:     "typedef    union {   }    my_union_t;\n",
			expected: "typedef union {\n} my_union_t;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-single-field",
			code:     "typedef    union    {  int   a;  } \nmy_union_t;\n",
			expected: "typedef union {\n    int a;\n} my_union_t;\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "03-multi-field",
			code: "typedef   union  {  int a;    float \n b;} my_union_t;\n",
			expected: "typedef union {\n" +
				"    int a;\n" +
				"    float b;\n" +
				"} my_union_t;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "04-messy-whitespace",
			code: "typedef    union {int  a  ;float b;} my_union_t;\n",
			expected: "typedef union {\n" +
				"    int a;\n" +
				"    float b;\n" +
				"} my_union_t;\n",
			pc: func(f func()) { f() },
		},
		{
			name:     "05-tagged-union",
			code:     "typedef    union   my_tag   {  int   a;  } \nmy_union_t;\n",
			expected: "typedef union my_tag {\n    int a;\n} my_union_t;\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "06-c-bitfields-in-struct",
			code: "typedef struct {\n" +
				"    uint32   a  :8;\n" +
				"    uint32   b  :16;\n" +
				"    uint32      :8;\n" +
				"} my_struct_t;\n",
			expected: "typedef struct {\n" +
				"    uint32 a : 8;\n" +
				"    uint32 b : 16;\n" +
				"    uint32 : 8;\n" +
				"} my_struct_t;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "07-anonymous-nested-struct-in-union",
			code: "typedef union {\n" +
				"    struct {\n" +
				"        uint32 a   :8;\n" +
				"        uint32 b:   8;\n" +
				"    };\n" +
				"    uint32 value;\n" +
				"} my_union_t;\n",
			expected: "typedef union {\n" +
				"    struct {\n" +
				"        uint32 a : 8;\n" +
				"        uint32 b : 8;\n" +
				"    };\n" +
				"    uint32 value;\n" +
				"} my_union_t;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "08-pointers-and-arrays",
			code: "typedef union {\n" +
				"    int *ptr;\n" +
				"    int arr[4];\n" +
				"} my_union_t;\n",
			expected: "typedef union {\n" +
				"    int *ptr;\n" +
				"    int arr[4];\n" +
				"} my_union_t;\n",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintTypes_Layout(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-empty-layout",
			code:     "typedef    layout    `my_layout`    {     }    l_t;\n",
			expected: "typedef layout `my_layout` {\n} l_t;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-single-field",
			code:     "typedef layout    `my_layout`    {    int    a;    }   l_t;\n",
			expected: "typedef layout `my_layout` {\n    int a;\n} l_t;\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "03-placeholder-names-issue-76",
			code: "typedef layout `little-endian` {\n" +
				"    uint8 _;\n" +
				"    uint8 data;\n" +
				"    uint16 __;\n" +
				"} foo_t;\n",
			expected: "typedef layout `little-endian` {\n" +
				"    uint8 _;\n" +
				"    uint8 data;\n" +
				"    uint16 __;\n" +
				"} foo_t;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "04-placeholder-triple-underscore-issue-76",
			code: "typedef layout `little-endian` {\n" +
				"    uint16 ___;\n" +
				"} foo_t;\n",
			expected: "typedef layout `little-endian` {\n" +
				"    uint16 ___;\n" +
				"} foo_t;\n",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintTypes_Bitfields(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-empty-bitfields",
			code:     "typedef     bitfields   32{    }bf_t;\n",
			expected: "typedef bitfields 32 {\n} bf_t;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-single-field-range",
			code:     "typedef   bitfields  32   {  int \n a   @  [0 : 4] ; } bf_t;\n",
			expected: "typedef bitfields 32 {\n    int a @ [0:4];\n} bf_t;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-multi-field",
			code:     "typedef     bitfields   32   {  int   a   @   [0    : 4] ;  int  b@[5];}  bf_t;\n",
			expected: "typedef bitfields 32 {\n    int a @ [0:4];\n    int b @ [5];\n} bf_t;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintTypes_Typeof(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-basic",
			code:     "  typedef    typeof   ( a )   t_t;\n",
			expected: "typedef typeof(a) t_t;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-expr",
			code:     "typedef   typeof(     a +     b)t_t;\n",
			expected: "typedef typeof(a + b) t_t;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintTypes_Sequence(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-basic",
			code:     "   typedef    sequence   (  int)seq_t   ; \n",
			expected: "typedef sequence(int) seq_t;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintTypes_Hook(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-empty",
			code:     "   typedef    hook  (  )   h_t  ;\n",
			expected: "typedef hook() h_t;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-args",
			code:     "  typedef    hook  (  int , char  )h_t ;\n",
			expected: "typedef hook(int, char) h_t;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
