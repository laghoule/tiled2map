# Tiled2map

Tiled can be used to build a simple 2 layer map, be exported in JSON format, and be repacked in files format used by project writtend in assembly language.

Support a limited subset of fonctionality of Tiled. The goal of this project is to:

- Create a map in raw binary format
- Create a tileset in raw binary format
- Create a tileset in PNG format
- Create an offset table in assembly, of each tile of the tileset for used in assembly language
- Custom properties of object for use with `tiles attributs bits` in the game engine

## Supported format

2 layer:

- Name: bg (map object)
- Name: fg (object layer)
