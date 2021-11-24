# Dynamic Entity Creator
In Rogue Champions new game elements will be able to created using the web editor.
Also the balance data will be loaded from backend.

## Requirements
- Upload current game object data to backend
- Create/Update game objects using backend data
- Dynamically load balance data from backend and update values while playing the game
    - At first a console command might be used e.g /load_version 'v1'


## Brainstorm
- Should we skip the gameobject representation entirely?
I don't know why do we need gameobject when we can pull all the data from backend and create prefabs from that.


## To-do
- There should be a conversion beta
