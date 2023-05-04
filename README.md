![](https://i.imgur.com/ORAnuBg.png)
# 2023_Pesak

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/c36accf815a5486a80d747a0db4a3bf0)](https://app.codacy.com/gh/matf-pp/2023_Pesak/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)

Pesak is a 2D [falling-sand](https://en.wikipedia.org/wiki/Falling-sand_game) simulation game written in Go and SDL2. It's designed to simulate interactions between tens of thousands of particles while taking into account their properties such as temperature, density and thermal conductivity. Players can experiment with combining different materials and simulate natural proccesses such as the water cycle. The simulation is also completely modular, making it easy to add new materials to the game and its interface. Pesak supports multiple display modes and many quality-of-life features, providing a fun and engaging experience.

Now with Discord Rich Presence support!

![tuta](./tuta.gif)

## Get yourself a Pesak
The most recent stable and well-tested binaries of Pesak for Linux and Windows are available in the [releases](https://github.com/matf-pp/2023_Pesak/releases) section of this repo. To get started with Pesak, we strongly encourage downloading the executables over compiling them on your machine.

## Build youself a sandcastle
To build Pesak yourself, you will need [Go v1.13+](https://go.dev/dl/) and [SDL2](https://github.com/libsdl-org/SDL/releases) installed on your machine. Detailed instructions for each OS are available [here](https://github.com/veandco/go-sdl2/blob/master/README.md#requirements).

Type the following commands:
```console
$ git clone github.com/matf-pp/2023_Pesak.git
$ cd 2023_Pesak
$ go build
$ ./main
```

## What can I do in Pesak
*    Play in sand
*    Play with fire
*    Heat up the sand
*    Check how hot it is
*    Continue heating it
*    Turn it into lava
*    Observe lava solidification
*    Enjoy the scenery
*    Save the scenery as a png file to share with friends
*    Drag that same (or any other) image back into Pesak

## Controls
*    LMB  -  paint
*    MMB  -  pick pointed material
*    RMB  -  clear
*    1-0  -  pick materials (or from sidebar, or by shift+scroll)
*    P  -  pause/resume simulation
*    T  -  temperature mode
*    D  -  gustina mode
*    N  -  normal mode
*    G  -  change the direction of gravity
*    M  -  mute/unmute sound
*    R  -  restart music from the begining
*    Z/X  -  change volume
*    esc  -  exit

<!-- ![GUI](https://i.imgur.com/JoI7s4I.png) -->
