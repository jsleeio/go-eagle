# go-eagle

## overview

This repository contains code and tools for interacting with Autodesk Eagle files.

At present, the below tools are included:

* `panelgen`: create a new blank panel board file in Eurorack, Pulplogic 1U or
  Intellijel 1U formats, at a specified width
* `schroff`: derive a new Eurorack panel board file from the board file for
  your circuit


## panelgen

`panelgen` is used for creating new, blank panels in Eurorack, Pulplogic 1U or
Intellijel 1U formats. An existing Eagle board file is required in order to
derive the desired set of Eagle layer information. This can be any Eagle board
file.

Demonstration usage, creating a 6hp Pulplogic tile:

```
$ ./panelgen -format=pulplogic -reference-board=data/ref.brd -output=mytile.brd -width=6
```


## schroff

`schroff` is used for deriving
[Eurorack module front panels](http://www.doepfer.de/a100_man/a100m_e.htm)
from the Eagle board file for the module's actual circuitry. That is, you
design your module's circuit board in Eagle, and then `schroff` examines the
board file to discover:

* which circuit components (potentiometers, jacks, LEDs, etc) require panel drill holes
* the size of any such drill holes (via the component's `PANEL_DRILL_MM` attribute)
* where the holes should be placed (via the component's origin coordinates)
* header text to be placed in silkscreen at the top of the panel (via the board's `HEADER_TEXT` attribute)
* footer text to be placed in silkscreen at the bottom of the panel (via the board's `FOOTER_TEXT` attribute)

To generate a panel board file:

```
$ ./schroff -copper-pullback=0.5 -hole-stop-radius=2.5 wavolver2-rev1.brd
2019/04/28 17:02:24 board: found HEADER_TEXT attribute with value "Wavolver II"
2019/04/28 17:02:24 board: found FOOTER_TEXT attribute with value "Ian Fritz"
2019/04/28 17:02:24 FOLDMIX: found PANEL_DRILL_MM attribute with value "7"
2019/04/28 17:02:24 FOLDOUT: found PANEL_DRILL_MM attribute with value "6"
2019/04/28 17:02:24 INPUT: found PANEL_DRILL_MM attribute with value "6"
2019/04/28 17:02:24 OFFSET: found PANEL_DRILL_MM attribute with value "7"
2019/04/28 17:02:24 OMOD: found PANEL_DRILL_MM attribute with value "6"
2019/04/28 17:02:24 OMODLEV: found PANEL_DRILL_MM attribute with value "7"
2019/04/28 17:02:24 OUTPUT: found PANEL_DRILL_MM attribute with value "6"
2019/04/28 17:02:24 P2AMP: found PANEL_DRILL_MM attribute with value "7"
2019/04/28 17:02:24 WIDTH: found PANEL_DRILL_MM attribute with value "7"
2019/04/28 17:02:24 WMOD: found PANEL_DRILL_MM attribute with value "6"
2019/04/28 17:02:24 WMODLEV: found PANEL_DRILL_MM attribute with value "7"
```

The output panel file takes the name of the input file and adds the suffix `.panel.brd`:

```
$ ls -l wavolver2-rev1.brd.panel.brd
-rw-r--r--  1 jslee  staff  17912 28 Apr 17:02 wavolver2-rev1.brd.panel.brd
```

### compatibility

At present the generated board files load just fine in Eagle 9.3.2 but are not
accepted by [OSHPark](https://oshpark.com/)'s Eagle board loader. I'm not sure
why this is, but it's most likely *not* OSHPark's fault, so please *don't*
complain to them if you try to use this. Just generate some Gerber files
instead.

### to-do

* import panel text for components from optional component attributes, eg. `PANEL_TEXT`
* exhaustively scan the Eagle DTD and add the various missing items
* BOM generation tool

## copyright

Copyright 2019 John Slee <jslee@jslee.io>.

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

