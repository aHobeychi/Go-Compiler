# Go Compiler

A Compiler Written from scratch in Go.

## Language Features

* Procedural Programming Language
* Syntax Similar to C
* Contains Arrays
* Contains Loops, Functions
* Contains Read, Write operations
* Integer Data Types
* Contains arithmetic operators
* Contains logical operators
* Allows for Ternary expressions

## Compiler Features

* Type Checking with error messaging
* Variable Declaration Checks
* Translates Code to Conform to a simplified Risc instruction set.

## Compiler Phases

1. File Lexer
2. Parser, receives tokens from the lexer.
3. Semantic Checks: type checking, inheritance checks, ...
4. Translation

## Usage

1. Compile the project using: go build
2. Run the executable using -i [source_path]
3. Run the Moon Machine on the .m output file

The grammar specifications are contained within the specifications folder and
sample programs within the samples folder.
