#!/usr/bin/env yosys

source [file join $::env(SYNTH_DIR) "logicworld.tcl"]

LW::loadVerilogSources

set nModules {
    "full_adder"
    "comparator"
}
foreach nModule $nModules {
    LW::generateGateSchematic [dict create \
        module $nModule \
        withBuffers yes \
        params {
            { N 4 }
            { N 8 }
        }
    ]
}

foreach module $modules {
    LW::generateGateSchematic $LW::withBuffers $module
}