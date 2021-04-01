# In development - documentation may not be completely accurate yet.
# Does not work in Windows - all tests have been done in Ubuntu
# arcadeprocesscompanion
Used as a companion to launching non-standard games on RetroPie/EmulationStation.  This application does the following:

- When Steam Proton or Lutris games are launched, the terminal exits immediately, returning the user back to the EmulationStation user interface.  This application can be configured to stay alive until a certain process ends.
- JSON configuration allows for joypad buttons/axis mapping to keyboard buttons.  For instance, some lutris games such as Hollow Knight cannot be exited via any buttons on the joypad.  This allows you to map a button to the Escape key in order to get to the games menu.

## Example commands

### Simulates keyboard events when joypad buttons/axes are pressed, and exits when no processes exist with the name 'gnome-calc'
```
arcadeprocesscompanion gnome-calc ./sample-mappings/test.json
```
