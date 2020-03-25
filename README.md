# overview

This repository contains code and tools for interacting with Autodesk Eagle files.

At present, the below tools are included:

* `panelgen`: create a new blank panel board file in Eurorack, Pulplogic 1U or
  Intellijel 1U formats, at a specified width
* `schroff`: derive a new Eurorack panel board file from the board file for
  your circuit

# installing

On a Mac with Homebrew, you can use my Homebrew tap to install the `panelgen`
and `schroff` commands:

```
brew tap jsleeio/apps
brew install go-eagle
```


# panelgen

`panelgen` is used for creating new, blank panels in Eurorack, Pulplogic 1U or
Intellijel 1U formats. An existing Eagle board file is required in order to
derive the desired set of Eagle layer information. This can be any Eagle board
file.

Demonstration usage, creating a 6hp Pulplogic tile:

```
$ ./panelgen -format=pulplogic -reference-board=data/ref.brd -output=mytile.brd -width=6
```

## commandline options

```
$ ./panelgen --help
Usage of ./panelgen:
  -format string
    	panel format to create (eurorack, pulplogic, intellijel) (default "eurorack")
  -outline-layer string
    	layer to draw board outline in (default "Dimension")
  -output string
    	filename to write new Eagle board file to (default "newpanel.brd")
  -reference-board string
    	reference Eagle board file to read layer information from
  -width int
    	width of the panel, in integer units appropriate for the format (default 4)
```

# schroff

`schroff` is used for deriving
[Eurorack module front panels](http://www.doepfer.de/a100_man/a100m_e.htm)
from the Eagle board file for the module's actual circuitry. That is, you
design your module's circuit board in Eagle, and then `schroff` examines the
board file to discover:

* which circuit components (potentiometers, jacks, LEDs, etc) require panel drill holes
* the size of any such drill holes (via the component's `PANEL_DRILL_MM` attribute)
* where the holes should be placed (via the component's origin coordinates)
* where the legend text should be placed (via the component's origin coordinates, and optional offset)
* header text to be placed in silkscreen at the top of the panel (via the board's `PANEL_HEADER_TEXT` attribute)
* footer text to be placed in silkscreen at the bottom of the panel (via the board's `PANEL_FOOTER_TEXT` attribute)

Components that need panel holes must have a `PANEL_DRILL_MM` attribute.

## list of global and component attributes

attribute name                    | type      | default value    | purpose
--------------------------------- | --------- | ---------------- | --------------------------------------------------------------------
`PANEL_HEADER_LAYER`              | global    | `tStop`          | layer to place header text on
`PANEL_HEADER_OFFSET_X`           | global    | `0.0`            | nudge panel header text left or right (millimetres)
`PANEL_HEADER_OFFSET_Y`           | global    | `0.0`            | nudge panel header text up or down (millimetres)
`PANEL_HEADER_TEXT`               | global    | `<HEADER_TEXT>`  | text for header section of panel
`PANEL_FOOTER_LAYER`              | global    | `tStop`          | layer to place footer text on
`PANEL_FOOTER_OFFSET_X`           | global    | `0.0`            | nudge panel footer text left or right (millimetres)
`PANEL_FOOTER_OFFSET_Y`           | global    | `0.0`            | nudge panel footer text up or down (millimetres)
`PANEL_FOOTER_TEXT`               | global    | `<FOOTER_TEXT>`  | text for footer section of panel
`PANEL_LEGEND_LAYER`              | global    | `tStop`          | layer to place panel legend text on
`PANEL_LEGEND_SKIP_RE`            | global    | _none_           | [RE2](https://github.com/google/re2/wiki/Syntax) expression; if a component name matches, legend text is skipped
`PANEL_DRILL_MM`                  | component | _none_           | panel drill size to create for a component. Required for drill holes.
`PANEL_HOLE_STOP_WIDTH`           | component | `2.0`            | override the width of the stop-mask ring around the component hole
`PANEL_LEGEND_LOCATION`           | component | `above`          | set to `below` to place the legend text `below` the component instead of `above`
`PANEL_LEGEND_OFFSET_X`           | component | `0.0`            | nudge panel legend text left or right (millimetres)
`PANEL_LEGEND_OFFSET_Y`           | component | `0.0`            | nudge panel legend text up or down (millimetres)
`PANEL_LEGEND_TICKS`              | component | `no`             | set to `yes` to add tick marks around component hole, eg. for potentiometers
`PANEL_LEGEND_TICKS_COUNT`        | component | `11`             | number of ticks to draw
`PANEL_LEGEND_TICKS_END_ANGLE`    | component | `225.0`          | ending polar angle to which to draw ticks, in degrees. Zero degrees is at 9 o'clock
`PANEL_LEGEND_TICKS_LENGTH`       | component | `2.0`            | length of ticks
`PANEL_LEGEND_TICKS_LABELS`       | component | `no`             | set to `yes` to add text labels next to tick marks
`PANEL_LEGEND_TICKS_LABELS_TEXTS` | component | _none_           | labels for tick marks, separated with `,`. Quantity must match `PANEL_LEGEND_TICKS_COUNT`
`PANEL_LEGEND_TICKS_START_ANGLE`  | component | `-45.0`          | starting polar angle from which to draw ticks, in degrees. Zero degrees is at 9 o'clock
`PANEL_LEGEND_TICKS_WIDTH`        | component | `0.5`            | width of ticks
`PANEL_LEGEND`                    | component | _component name_ | override panel legend text for a component

## commandline options

```
$ ./schroff --help
Usage of ./schroff:
  -format string
    	panel format to create (eurorack, pulplogic, intellijel) (default "eurorack")
  -hole-stop-radius float
    	Radius to pull back soldermask around a hole (default 2)
  -text-size float
    	label text size (default 2.25)
  -text-spacing float
    	spacing between a hole and its related label (default 3.5)
```

## generating board files

To generate a panel board file:

```
$ schroff morphlag-rev2.brd
2019/06/02 17:48:17 FALL: found PANEL_DRILL_MM attribute with value 7
2019/06/02 17:48:17 IN: found PANEL_DRILL_MM attribute with value 6
2019/06/02 17:48:17 OUT: found PANEL_DRILL_MM attribute with value 6
2019/06/02 17:48:17 OUTPOL: found PANEL_DRILL_MM attribute with value 6
2019/06/02 17:48:17 RISE: found PANEL_DRILL_MM attribute with value 7
2019/06/02 17:48:17 SHAPE: found PANEL_DRILL_MM attribute with value 7
2019/06/02 17:48:17 SW1: found PANEL_DRILL_MM attribute with value 4.5
2019/06/02 17:48:17 POLARIZE: found PANEL_DRILL_MM attribute with value 7
2019/06/02 17:48:17 OUTINV: found PANEL_DRILL_MM attribute with value 6
2019/06/02 17:48:17 MANUAL: found PANEL_DRILL_MM attribute with value 7.5
```

The output panel file takes the name of the input file and adds the suffix `.panel.brd`:

```
$ ls -l wavolver2-rev1.brd.panel.brd
-rw-r--r--  1 jslee  staff  17912 28 Apr 17:02 wavolver2-rev1.brd.panel.brd
```

# compatibility

At present the generated board files load just fine in Eagle 9.3.2+ (probably
many earlier versions also!) but are _not_ accepted by
[OSHPark](https://oshpark.com/)'s Eagle board loader.  I'm not sure why this
is, but it's most likely *not* OSHPark's fault, so please *don't* complain to
them if you try to use this. Just generate some Gerber files instead.

# to-do

* exhaustively scan the Eagle DTD and add the various missing items (libraries!)
* ability to define custom panel formats, eg. to fit a specific custom enclosure
* BOM generation tool

# copyright

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

