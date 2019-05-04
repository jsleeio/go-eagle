#!/bin/sh

for format in eurorack intellijel pulplogic ; do
  for width in 6 12 ; do
    ./panelgen \
      -width="${width}" \
      -reference-board=data/ref.brd \
      -output="${format}-${width}hp.brd" \
      -format="${format}"
  done
done
