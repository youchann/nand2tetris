// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/4/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)
// The algorithm is based on repetitive addition.

// R0,R1の値を読み出す
// どちらかが0ならENDまでジャンプする
    @R1
    D=M
    @R2
    M=D
    @END
    D;JEQ

    @R0
    D=M
    @R2
    M=D
    @END
    D;JEQ

// While文を書いてR2=R2+R0していく
(LOOP)
    // 引き算して0であればENDへ
    @R1
    M=M-1
    D=M
    @END
    D;JEQ
    // 0以上であればR2に追加する
    @R0
    D=M
    @R2
    M=D+M
    // LOOPへ戻る
    @LOOP
    0;JMP
(END)
    @END
    0;JMP