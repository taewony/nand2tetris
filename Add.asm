// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.

// Computes R0 = 2 + 3  (R0 refers to RAM[0])

@2
D=A // D=M[2]
@3
D=D+A // D=D+M[3]
@0
M=D // m[0] = D