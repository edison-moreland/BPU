#!/usr/bin/env yosys

source [file join $::env(SYNTH_DIR) "logicworld.tcl"]

LW::loadVerilogSources

LW::generateGateSchematic {
    module inc_dec_rst
    withBuffers yes
    params {
        { N 3 }
        { N 4 }
    }
}

LW::generateGateSchematic {
    module cursor_register
    params {
        { N 3 }
        { N 4 }
    }
}

LW::generateGateSchematic {
    module is_control_char
}