# Minesweeper

A clone of the good old Minesweeper game, trying to keep the original style as much as possible. Powered by 
[raylib-go](https://github.com/gen2brain/raylib-go).

<img src="img/screenshot.png" width="420">

## How to Play

Just launch the binary and play!

```shell
$ ./minesweeper
```
You can also specify the size of the board and the number of mines.

```shell
$ ./minesweeper -r 30 -c 20 -m 100
````

where, `-r` is the number of rows; `-c` is the number of columns; `-m` is the number of mines.

## How to Build

```shell
$ go build -o minesweeper
```

On Linux, you might want to install the following packages (e.g, Ubuntu):

    - libgl-dev
    - libxi-dev

