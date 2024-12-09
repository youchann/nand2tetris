// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/4/Fill.asm

// Runs an infinite loop that listens to the keyboard input. 
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel. When no key is pressed, 
// the screen should be cleared.

(INIT)
    @SCREEN
    D=A
    // 画像には8K分のレジスタを利用している
    @8192
    D=D+A
    @screen_end
    M=D
(LOOP)
    @SCREEN
    D=A
    // 今塗ろうとしているレジスタを定義
    @screen_ptr
    M=D
    @KBD
    D=M
    @FILL
    D;JNE
    @EMPTY
    0;JMP
// 色を黒に指定
(FILL)
    @color
    M=-1
    @DRAW
    0;JMP
// 色を白に指定
(EMPTY)
    @color
    M=0
    @DRAW
    0;JMP
(DRAW)
    @color
    D=M
    @screen_ptr
    A=M
    M=D

    @screen_ptr
    // 次のレジスタに移動
    M=M+1
    D=M
    @screen_end
    D=M-D
    @DRAW
    D;JGT
    @LOOP
    0;JMP
