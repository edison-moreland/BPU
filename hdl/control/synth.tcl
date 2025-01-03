#!/usr/bin/env yosys

source [file join $::env(SYNTH_DIR) "logicworld.tcl"]

LW::loadVerilogSources

set modules {
    "microop_counter"
    "inc_rst"
}

foreach module $modules {
    LW::generateGateSchematic $LW::withBuffers $module
}