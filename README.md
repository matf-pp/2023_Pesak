![](https://i.imgur.com/ORAnuBg.png)
# 2023_Pesak
Pesak is a 2D [falling-sand](https://en.wikipedia.org/wiki/Falling-sand_game) simulation game written in Go and SDL2. It's designed to simulate interactions between tens of thousands of particles while taking into account their properties such as temperature, density and thermal conductivity. Players can experiment with combining different materials and simulate natural proccesses such as the water cycle. The simulation is also completely modular, making it easy to add new materials to the game and its interface. Pesak supports multiple display modes and many quality-of-life features, providing a fun and engaging experience.

Now with Discord Rich Presence support!

![tuta](./tuta.gif)

## Get yourself a Pesak
```
$ go get github.com/matf-pp/2023_Pesak
$ go build
$ ./main
```

## What can I do in Pesak?
* Play in sand
* Play with fire
* Heat up the sand
* Check how hot it is
* Continue heating it
* Turn it into lava
* Observe lava solidification
* Enjoy the scenery
* Save the scenery as a png file to share with friends
* Drag that same (or any other) image back into Pesak 
    
## Controls
- LMB  -  paint
- RMB  -  clear
- 1-0  -  pick materials
- [/]  -  heat up/cool down
- P  -  pause/resume simulation
- T  -  temperature mode
- G  -  gustina mode
- N  -  normal mode
- up/down arrows  -  resize brush
- esc  -  exit

<!-- ![GUI](https://i.imgur.com/JoI7s4I.png) -->
