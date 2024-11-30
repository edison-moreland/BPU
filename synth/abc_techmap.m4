# First pass
&get -n;
&st;
&dch -f -e -r -m;
&nf;

# Optimize for area (taken from SiliconCompiler scripts)
define(`area_opt', `&st;&syn2;&if -g -K 6;&synch2;&nf -a {D};')
area_opt
area_opt
area_opt
area_opt
area_opt

&put;

mfs3 -v -es;
topo;

ifdef(`insert-buffers', `addbuffs;', `')

print_stats -m;

