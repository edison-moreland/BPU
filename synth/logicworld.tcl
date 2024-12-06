package require fileutil
package require textutil

yosys -import

namespace eval ::LW {
    namespace eval v {
        variable targetDir $::env(TARGET_DIR)
        variable outputDir $::env(OUTPUT_DIR)
        variable synthDir $::env(SYNTH_DIR)

        variable sourceFiles {tes test}

        variable logicworldLib [file join $synthDir "logicworld.lib"]
        variable memoryTechamp [file join $synthDir "flipflop2latch_techmap.v"]
        variable abcTechmap [file join $synthDir "abc_techmap.m4"]

        variable legalDffs {
            $_DLATCH_P_
            $_DFF_P_
            $_ALDFF_PP_
            $_SR_PP_
            $_DFF_PP0_
        }

        variable netlistsvgSkin [file join $synthDir "netlistsvg_skin.svg"]
    }

    proc loadVerilogSources {} {
        set verilogOutput [file join $v::outputDir "verilog"]

        if {[file exists $verilogOutput] != 1} {
            file mkdir $verilogOutput
        }

        # Load verilog files first
        foreach vfile [fileutil::findByPattern $v::targetDir *.v] {
            set vfilename [file rootname [file tail $vfile]]
            set vdir [file join $verilogOutput ./[file dirname [textutil::trimPrefix $vfile $v::targetDir]]]
            set vfileOutput [file normalize [file join "$vdir" "./${vfilename}.v"]]

            if {[file exists $vdir] != 1} {
                file mkdir $vdir
            }

            file copy -force $vfile $vfileOutput
            dict set v::sourceFiles $vfilename $vfileOutput
        }

        # Now load system verilog files
        foreach svfile [fileutil::findByPattern $v::targetDir *.sv] {
            set svfilename [file rootname [file tail $svfile]]
            set svdir [file join $verilogOutput ./[file dirname [textutil::trimPrefix $svfile $v::targetDir]]]
            set svfileOutput [file normalize [file join "$svdir" "./${svfilename}.v"]]

            if {[file exists $svdir] != 1} {
                file mkdir $svdir
            }

            exec sv2v -w $svfileOutput $svfile
            dict set v::sourceFiles $svfilename $svfileOutput
        }
    }

    set withoutBuffers 0
    set withBuffers 1

    proc generateModuleSchematic {moduleName { moduleParameters {{}} }} {
        yosys read -vlog2k [dict get $v::sourceFiles $moduleName]
        yosys design -stash $moduleName


        foreach parameterSet $moduleParameters {
            yosys design -load $moduleName

            set outputName $moduleName
            dict for { key value } $parameterSet {
                chparam -set $key $value $moduleName
                set outputName "${outputName}_$key$value"
            }

            # The real magic
            # ::LW::synthToGates $moduleName $insertBuffers $netlistOutput

            # set netlistOutput [::LW::outputPath "netlists" "${outputName}.json"]
            # write_json $netlistOutput
            hierarchy -top $outputName
            yosys proc
            clean
            opt

            set schematicOutput [::LW::outputPath "schematics" "${outputName}"]
            show -format svg -prefix $schematicOutput $outputName

            # exec netlistsvg $netlistOutput -o $schematicOutput --skin $v::netlistsvgSkin
        }
    }

    proc generateGateSchematic {  { moduleParameters {} } } {
        set moduleName [dict get $moduleParameters module]
        set moduleParams [::LW::dictGetDefault $moduleParameters params {}]

        set insertBuffers [::LW::dictGetDefault $moduleParameters withBuffers 0]

        yosys read -vlog2k [dict get $v::sourceFiles $moduleName]
        yosys design -stash $moduleName

        set abcScript [::LW::prepareABCScript $insertBuffers]


        foreach parameterSet $moduleParams {
            yosys design -load $moduleName

            set outputName $moduleName
            dict for { key value } $parameterSet {
                chparam -set $key $value $moduleName
                set outputName "${outputName}_$key$value"
            }

            # The real magic
            ::LW::synthToGates $moduleName $abcScript

            set netlistOutput [::LW::outputPath "netlists" "${outputName}.json"]
            write_json $netlistOutput

            set schematicOutput [::LW::outputPath "schematics" "${outputName}.svg"]
            exec netlistsvg $netlistOutput -o $schematicOutput --skin $v::netlistsvgSkin
        }
    }

    proc outputPath { folder filename } {
        set outputDir [file join $v::outputDir $folder]
        if {[file exists $outputDir] != 1} {
            file mkdir $outputDir
        }
        return [file join $outputDir $filename]
    }

    proc synthToGates { moduleName abcScript } {
        read_liberty -lib $v::logicworldLib

        # Generic asic synth
        synth -flatten -top $moduleName

        # Memory techmap
        dfflegalize {*}[::LW::dfflegalizeArgs]
        techmap -autoproc -map $v::memoryTechamp
        opt_merge
        freduce -inv

        # Final opt and techmap
        opt -full
        abc -script $abcScript -liberty $v::logicworldLib
        opt_clean
    }

    proc dictGetDefault { dict key default } {
        if [dict exists $dict $key] {
            return [dict get $dict $key]
        } else {
            return $default
        }
    }

    proc prepareABCScript {insertBuffers} {
        set abcTechmapNew [file join $v::outputDir "../abc_techmap"]
        set defines {}
        if {$insertBuffers} {
            set abcTechmapNew "${abcTechmapNew}_with_buffers"
            dict set defines insert-buffers 1
        }
        m4 $v::abcTechmap $abcTechmapNew $defines

        return $abcTechmapNew
    }

    proc m4 {inputFile outputFile defines} {
        set defineArgs {}
        dict for {key value} $defines {
            lappend defineArgs "--define=$key=$value"
        }

        exec m4 <$inputFile >$outputFile {*}$defineArgs
    }

    proc dfflegalizeArgs {} {
        set da {}
        foreach dff $v::legalDffs {
            lappend da -cell $dff 01
        }

        return $da
    }

}