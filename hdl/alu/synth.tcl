#!/usr/bin/env yosys

source [file join $::env(SYNTH_DIR) "logicworld.tcl"]

LW::loadVerilogSources

set modules {
    "alu_adder_cs_half"
    "alu_adder_half"
    "alu_comparator_half"
    "alu_logic_slice"
}

foreach module $modules {
    LW::generateGateSchematic $LW::withBuffers $module
}