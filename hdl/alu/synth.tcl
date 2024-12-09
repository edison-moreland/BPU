#!/usr/bin/env yosys

source [file join $::env(SYNTH_DIR) "logicworld.tcl"]

LW::loadVerilogSources


LW::simulate {
    module "alu_tb"
}

LW::generateGateSchematic {
    module "alu"
    withBuffers yes
    params {
        { N 8 }
    }
}

LW::generateGateSchematic {
    module "adder"
    withBuffers yes
    params {
        { N 4 }
        { N 8 }
        { N 16 }
    }
}

LW::generateGateSchematic {
    module "multiplier"
    withBuffers yes
    params {
        { N 4 }
        { N 8 }
    }
}