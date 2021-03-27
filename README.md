# SnakeGame_Linux  I develop this solution on Ubuntu20.04 PC, using VSCode IDE 

Since go does not support user input nicely, I deal with keyboard input using third party library--"github.com/gdamore/tcell". In fact, the hard part of this 
task is finding a tool to process user input, and I have dig in google for hours to find a usable one.

I assume the game's boader size is within range[10~100], both height and width is in that range. Also, snake should be auto moving until user press keyboard 
to change the direction. I've also implemented the manually moving snake version, but there's no fun in doing that, so I do not upload the code in that version.

I've been trying hard to avoid sharing global state, so a game state structure include everything to pass to the update game function.

The program work flow is: start/run => user input boarder size => print game initial state => user input direction => change snake head direction 
=> update game state and son on, until game over => user press esc key to exit the program and return back to the terminal.
